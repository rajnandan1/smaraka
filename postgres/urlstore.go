package postgres

import (
	"context"
	"fmt"

	"github.com/rajnandan1/smaraka/models"
)

func (p *PostgresImplementation) InsertNewURLStore(ctx context.Context, urlStore models.URLStore) (*models.URLStore, error) {
	query := `
		INSERT INTO url_store (id, url, domain, title, image_sm, image_lg, excerpt, color, status, full_content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW());`

	_, err := p.Pool.Exec(ctx, query,
		urlStore.ID,
		urlStore.URL,
		urlStore.Domain,
		urlStore.Title,
		urlStore.ImageSmall,
		urlStore.ImageLarge,
		urlStore.Excerpt,
		urlStore.AccentColor,
		urlStore.Status,
		urlStore.FullText,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert URL store record: %v", err)
	}

	return p.GetURLStoreByID(ctx, urlStore.ID)
}

func (p *PostgresImplementation) GetURLStoreByURL(ctx context.Context, url string) (*models.URLStore, error) {
	var urlStore models.URLStore

	query := `
		SELECT id, url, domain, title, image_sm, image_lg, excerpt, color, status, full_content, created_at, updated_at
		FROM url_store
		WHERE url = $1;`

	err := p.Pool.QueryRow(ctx, query, url).Scan(
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

func (p *PostgresImplementation) GetURLStoreByID(ctx context.Context, id string) (*models.URLStore, error) {
	var urlStore models.URLStore

	query := `
		SELECT id, url, domain, title, image_sm, image_lg, excerpt, color, status, full_content, created_at, updated_at
		FROM url_store
		WHERE id = $1;`

	err := p.Pool.QueryRow(ctx, query, id).Scan(
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

func (p *PostgresImplementation) UpdateURLStoreByID(ctx context.Context, id string, urlData models.URLStore) (*models.URLStore, error) {
	query := `
		UPDATE url_store
		SET title = $1, image_sm = $2, image_lg = $3, excerpt = $4, color = $5, status = $6, full_content = $7, updated_at = $8
		WHERE id = $9;`

	_, err := p.Pool.Exec(ctx, query, urlData.Title, urlData.ImageSmall, urlData.ImageLarge, urlData.Excerpt, urlData.AccentColor, urlData.Status, urlData.FullText, urlData.UpdatedAt, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update url store: %v", err)
	}

	return p.GetURLStoreByID(ctx, urlData.ID)
}

func (p *PostgresImplementation) UpdateURLStoreByURL(ctx context.Context, url string, urlData models.URLStore) (*models.URLStore, error) {
	query := `
		UPDATE url_store
		SET title = $1, image_sm = $2, image_lg = $3, excerpt = $4, color = $5, status = $6, full_content = $7, updated_at = $8
		WHERE url = $9;`

	_, err := p.Pool.Exec(ctx, query, urlData.Title, urlData.ImageSmall, urlData.ImageLarge, urlData.Excerpt, urlData.AccentColor, urlData.Status, urlData.FullText, urlData.UpdatedAt, url)
	if err != nil {
		return nil, fmt.Errorf("failed to update url store: %v", err)
	}

	return p.GetURLStoreByURL(ctx, urlData.URL)
}

func (p *PostgresImplementation) GetURLStoreByIDs(ctx context.Context, ids []string) ([]models.URLStore, error) {
	var urlStores []models.URLStore

	query := `
		SELECT id, url, domain, title, image_sm, image_lg, excerpt, color, status, created_at, updated_at
		FROM url_store
		WHERE id = ANY($1)
		ORDER BY array_position($1, id);`

	rows, err := p.Pool.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve url store: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var urlStore models.URLStore
		err := rows.Scan(
			&urlStore.ID,
			&urlStore.URL,
			&urlStore.Domain,
			&urlStore.Title,
			&urlStore.ImageSmall,
			&urlStore.ImageLarge,
			&urlStore.Excerpt,
			&urlStore.AccentColor,
			&urlStore.Status,
			&urlStore.CreatedAt,
			&urlStore.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve url store: %v", err)
		}
		urlStores = append(urlStores, urlStore)
	}

	return urlStores, nil
}
func (p *PostgresImplementation) GetURLsByIDs(ctx context.Context, ids []string) ([]models.URLStore, error) {
	var urls []models.URLStore

	// Define SQL query with array_position to maintain ordering based on ids
	query := `
		SELECT id, url, domain, title, image_sm, image_lg, excerpt, color, status, full_content, created_at, updated_at
		FROM url_store
		WHERE id = ANY($1)
		ORDER BY array_position($1, id);`

	// Execute the query with ids as an array argument
	rows, err := p.Pool.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve URL store records: %v", err)
	}
	defer rows.Close()

	// Iterate over rows and scan each into a url_store model
	for rows.Next() {
		var url models.URLStore
		err = rows.Scan(
			&url.ID,
			&url.URL,
			&url.Domain,
			&url.Title,
			&url.ImageSmall,
			&url.ImageLarge,
			&url.Excerpt,
			&url.AccentColor,
			&url.Status,
			&url.FullText,
			&url.CreatedAt,
			&url.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan URL store record: %v", err)
		}
		urls = append(urls, url)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over URL store records: %v", err)
	}

	return urls, nil
}
