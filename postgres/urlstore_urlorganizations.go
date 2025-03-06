package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/models"
)

func (p *PostgresImplementation) GetAllURLsForORG(ctx context.Context, orgID string) (*[]models.URLStore, *[]models.URLOrganizations, error) {
	// Prepare SQL statement
	query := `
		SELECT uo.id, uo.created_at, us.title, us.url
		FROM url_organizations uo
		JOIN url_store us ON uo.url_id = us.id
		WHERE uo.organization_id = $1
		ORDER BY uo.created_at ASC;`

	rows, err := p.Pool.Query(ctx, query, orgID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve urls for organization: %v", err)
	}
	defer rows.Close()

	var urlStores []models.URLStore
	var urlOrganizations []models.URLOrganizations
	for rows.Next() {
		var urlStore models.URLStore
		var urlOrganization models.URLOrganizations
		err := rows.Scan(
			&urlOrganization.ID,
			&urlOrganization.CreatedAt,
			&urlStore.Title,
			&urlStore.URL,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan url store: %v", err)
		}
		urlStores = append(urlStores, urlStore)
		urlOrganizations = append(urlOrganizations, urlOrganization)
	}

	return &urlStores, &urlOrganizations, nil
}

func (p *PostgresImplementation) SearchURLs(ctx context.Context, orgID, query, domain string) ([]*models.URLResponses, error) {
	// Prepare SQL statement

	terms := strings.Fields(query) // Split the query by whitespace
	boostedQuery := query + "^2"

	queryStr := `
        SELECT us.id as url_id, us.title, us.url, us.excerpt, us.image_sm, us.image_lg, us.color, 
        uo.id as organization_relation_id, uo.status as organization_url_status, paradedb.score(us.id)
        FROM url_organizations uo
        JOIN url_store us ON uo.url_id = us.id 
        WHERE uo.organization_id = $1 AND (
            us.id @@@ paradedb.phrase('full_content', $3::text[], slop => 10)
            or 
            us.id @@@ paradedb.phrase_prefix('excerpt', $3::text[])
            or 
						us.id @@@ paradedb.term('domain', $4)
            or 
            us.title @@@ $5
        ) and uo.status = $2`

	queryStr += " ORDER BY paradedb.score(us.id) DESC limit 100;"

	var rows pgx.Rows
	rows, err := p.Pool.Query(ctx, queryStr, orgID, constants.URLStatusActive, terms, query, boostedQuery)

	if err != nil {
		return nil, fmt.Errorf("failed to search urls: %v", err)
	}
	defer rows.Close()

	var urlStores []*models.URLResponses
	for rows.Next() {
		var urlStore models.URLResponses
		err := rows.Scan(
			&urlStore.URLID,
			&urlStore.Title,
			&urlStore.URL,
			&urlStore.Excerpt,
			&urlStore.ImageSmall,
			&urlStore.ImageLarge,
			&urlStore.AccentColor,
			&urlStore.OrganizationRelationID,
			&urlStore.OrganizationURLStatus,
			&urlStore.Score,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan url store: %v", err)
		}
		urlStores = append(urlStores, &urlStore)
	}

	return urlStores, nil
}

func (p *PostgresImplementation) GetURLsForOrganization(ctx context.Context, organizationID string, lastID string, pageSize int) ([]*models.URLResponses, error) {
	var urlOrganizations []*models.URLResponses

	query := `
        SELECT us.id as url_id, us.title, us.url, us.excerpt, us.image_sm, us.image_lg, us.color, 
        uo.id as organization_relation_id, uo.status as organization_url_status
        FROM url_organizations uo
        JOIN url_store us ON uo.url_id = us.id
        WHERE uo.organization_id = $1 AND uo.id < $2 and uo.status = $3
        ORDER BY uo.id DESC
        LIMIT $4`

	rows, err := p.Pool.Query(ctx, query, organizationID, lastID, constants.URLStatusActive, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve urls for organization: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		// Initialize a new struct for each row
		urlOrganization := &models.URLResponses{}

		err := rows.Scan(
			&urlOrganization.URLID,
			&urlOrganization.Title,
			&urlOrganization.URL,
			&urlOrganization.Excerpt,
			&urlOrganization.ImageSmall,
			&urlOrganization.ImageLarge,
			&urlOrganization.AccentColor,
			&urlOrganization.OrganizationRelationID,
			&urlOrganization.OrganizationURLStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan url organization: %v", err)
		}

		urlOrganizations = append(urlOrganizations, urlOrganization)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return urlOrganizations, nil
}

// GetURLStoreByURLOrgIDOrgID
func (p *PostgresImplementation) GetURLStoreByURLOrgIDOrgID(ctx context.Context, urlOrgID, orgID string) (*models.URLStore, error) {
	query := `
		SELECT us.id, us.url, us.domain, us.title, us.image_sm, us.image_lg, us.excerpt, us.color, us.status, us.full_content, us.created_at, us.updated_at
		FROM url_organizations uo
		JOIN url_store us ON uo.url_id = us.id
		WHERE uo.id = $1 AND uo.organization_id = $2;`

	row := p.Pool.QueryRow(ctx, query, urlOrgID, orgID)

	var urlStore models.URLStore
	err := row.Scan(
		&urlStore.ID,
		&urlStore.URL,
		&urlStore.Domain,
		&urlStore.Title,
		&urlStore.ImageSmall,
		&urlStore.ImageLarge,
		&urlStore.Excerpt,
		&urlStore.AccentColor,
		&urlStore.Status,
		&urlStore.FullText,
		&urlStore.CreatedAt,
		&urlStore.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve url store: %v", err)
	}

	return &urlStore, nil
}

func (p *PostgresImplementation) GetSingleURLForOrganization(ctx context.Context, organizationID string, urlOrgID string) (*models.URLResponses, error) {
	query := `
		SELECT us.id as url_id, us.title, us.url, us.excerpt, us.image_sm, us.image_lg, us.color, 
		uo.id as organization_relation_id, uo.status as organization_url_status
		FROM url_organizations uo
		JOIN url_store us ON uo.url_id = us.id
		WHERE uo.organization_id = $1 AND uo.id = $2`

	row := p.Pool.QueryRow(ctx, query, organizationID, urlOrgID)

	var urlOrganization models.URLResponses
	err := row.Scan(
		&urlOrganization.URLID,
		&urlOrganization.Title,
		&urlOrganization.URL,
		&urlOrganization.Excerpt,
		&urlOrganization.ImageSmall,
		&urlOrganization.ImageLarge,
		&urlOrganization.AccentColor,
		&urlOrganization.OrganizationRelationID,
		&urlOrganization.OrganizationURLStatus,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve url organization: %v", err)
	}

	return &urlOrganization, nil
}

func (p *PostgresImplementation) GetSingleURLForOrganizationURL(ctx context.Context, organizationID string, url string) (*models.URLResponses, error) {
	query := `
		SELECT us.id as url_id, us.title, us.url, us.excerpt, us.image_sm, us.image_lg, us.color, 
		uo.id as organization_relation_id, uo.status as organization_url_status
		FROM url_organizations uo
		JOIN url_store us ON uo.url_id = us.id
		WHERE uo.organization_id = $1 AND us.url = $2`

	row := p.Pool.QueryRow(ctx, query, organizationID, url)

	var urlOrganization models.URLResponses
	err := row.Scan(
		&urlOrganization.URLID,
		&urlOrganization.Title,
		&urlOrganization.URL,
		&urlOrganization.Excerpt,
		&urlOrganization.ImageSmall,
		&urlOrganization.ImageLarge,
		&urlOrganization.AccentColor,
		&urlOrganization.OrganizationRelationID,
		&urlOrganization.OrganizationURLStatus,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve url organization: %v", err)
	}

	return &urlOrganization, nil
}
