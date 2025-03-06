package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/rajnandan1/smaraka/models"
)

func (p *PostgresImplementation) InsertNewOrganization(ctx context.Context, organization models.Organizations) error {
	query := `
		INSERT INTO organizations (id, name, creator_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5);`

	_, err := p.Pool.Exec(ctx, query,
		organization.ID,
		organization.Name,
		organization.CreatorID,
		time.Now(),
		time.Now())

	if err != nil {
		return fmt.Errorf("failed to insert organization: %v", err)
	}

	return nil
}

func (p *PostgresImplementation) GetOrganizationsByCreatorID(creatorID string) ([]models.Organizations, error) {
	var organizations []models.Organizations

	query := `
		SELECT id, name, creator_id, created_at, updated_at
		FROM organizations
		WHERE creator_id = $1;`

	rows, err := p.Pool.Query(context.Background(), query, creatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve organizations: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var org models.Organizations
		if err := rows.Scan(&org.ID, &org.Name, &org.CreatorID, &org.CreatedAt, &org.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan organization: %v", err)
		}
		organizations = append(organizations, org)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %v", err)
	}

	return organizations, nil
}

// get organization by id
func (p *PostgresImplementation) GetOrganizationByID(ctx context.Context, id string) (*models.Organizations, error) {
	var organization models.Organizations

	query := `
		SELECT id, name, creator_id, created_at, updated_at
		FROM organizations
		WHERE id = $1;`

	err := p.Pool.QueryRow(ctx, query, id).Scan(&organization.ID, &organization.Name, &organization.CreatorID, &organization.CreatedAt, &organization.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve organization: %v", err)
	}

	return &organization, nil
}
