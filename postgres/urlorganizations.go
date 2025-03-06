package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/models"
)

func (p *PostgresImplementation) InsertNewURLOrganization(ctx context.Context, urlOrganization models.URLOrganizations) (*models.URLOrganizations, error) {
	query := `
		INSERT INTO url_organizations (id, url_id, organization_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())`

	_, err := p.Pool.Exec(ctx, query,
		urlOrganization.ID,
		urlOrganization.URLID,
		urlOrganization.OrganizationID,
		urlOrganization.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert url organization: %v", err)
	}

	return p.GetURLOrganizationByID(ctx, urlOrganization.ID)
}

func (p *PostgresImplementation) GetURLOrganizationByID(ctx context.Context, id string) (*models.URLOrganizations, error) {
	var urlOrganization models.URLOrganizations

	query := `
		SELECT id, url_id, organization_id, created_at, updated_at
		FROM url_organizations
		WHERE id = $1;`

	err := p.Pool.QueryRow(ctx, query, id).Scan(
		&urlOrganization.ID,
		&urlOrganization.URLID,
		&urlOrganization.OrganizationID,
		&urlOrganization.CreatedAt,
		&urlOrganization.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve url organization: %v", err)
	}

	return &urlOrganization, nil
}

// GetURLCountForOrganization
func (p *PostgresImplementation) GetURLCountForOrganization(ctx context.Context, organizationID string) (int, error) {
	query := `
		SELECT count(*)
		FROM url_organizations
		WHERE organization_id = $1 and status = $2;`

	var count int
	err := p.Pool.QueryRow(ctx, query, organizationID, constants.URLStatusActive).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve url count for organization: %v", err)
	}

	return count, nil
}

// GetURLOrganizationsByURLIDOrgID
func (p *PostgresImplementation) GetURLOrganizationsByURLIDOrgID(ctx context.Context, urlID string, organizationID string) (*models.URLOrganizations, error) {
	var urlOrganization models.URLOrganizations

	query := `
		SELECT id, url_id, organization_id, created_at, updated_at
		FROM url_organizations
		WHERE url_id = $1 AND organization_id = $2;`

	err := p.Pool.QueryRow(ctx, query, urlID, organizationID).Scan(
		&urlOrganization.ID,
		&urlOrganization.URLID,
		&urlOrganization.OrganizationID,
		&urlOrganization.CreatedAt,
		&urlOrganization.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve url organization: %v", err)
	}

	return &urlOrganization, nil
}

// UpdateURLStatusByID
func (p *PostgresImplementation) UpdateURLStatusByID(ctx context.Context, id string, status string) error {
	query := `
		UPDATE url_organizations
		SET status = $1, updated_at = $2
		WHERE id = $3;`

	_, err := p.Pool.Exec(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update url organization status: %v", err)
	}

	return nil
}

// DeleteURLByID
func (p *PostgresImplementation) DeleteURLByID(ctx context.Context, id string) error {
	query := `
		DELETE FROM url_organizations
		WHERE id = $1;`

	_, err := p.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete url organization: %v", err)
	}

	return nil
}

// GetNewURLsForOrganization(organizationID string, firstId string, pageSize int) ([]models.URLOrganizations, error)
func (p *PostgresImplementation) GetNewURLsForOrganization(ctx context.Context, organizationID string, firstId string, pageSize int) ([]models.URLOrganizations, error) {
	var urlOrganizations []models.URLOrganizations

	query := `
		SELECT id, url_id, organization_id, status, created_at, updated_at
		FROM url_organizations
		WHERE organization_id = $1 AND id > $2 and status = $3
		ORDER BY id ASC
		LIMIT $4;`

	rows, err := p.Pool.Query(ctx, query, organizationID, firstId, constants.URLStatusActive, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve new urls for organization: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var urlOrganization models.URLOrganizations
		err := rows.Scan(
			&urlOrganization.ID,
			&urlOrganization.URLID,
			&urlOrganization.OrganizationID,
			&urlOrganization.Status,
			&urlOrganization.CreatedAt,
			&urlOrganization.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan url organization: %v", err)
		}
		urlOrganizations = append(urlOrganizations, urlOrganization)
	}

	return urlOrganizations, nil
}

// delete urls by ids
func (p *PostgresImplementation) DeleteURLsByIDs(ctx context.Context, ids []string, orgID string) error {
	query := `
		UPDATE url_organizations
		SET status = $1
		WHERE organization_id = $2 AND id = ANY($3);`

	_, err := p.Pool.Exec(ctx, query, constants.URLStatusDeleted, orgID, ids)
	if err != nil {
		return fmt.Errorf("failed to delete urls by ids: %v", err)
	}

	return nil
}
