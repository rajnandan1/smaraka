package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rajnandan1/smaraka/models"
)

func (p *PostgresImplementation) InsertNewUserOrganization(ctx context.Context, userOrganization models.UserOrganizations) error {
	query := `
		INSERT INTO user_organizations (id, user_id, organization_id, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := p.Pool.Exec(ctx, query,
		userOrganization.ID,
		userOrganization.UserID,
		userOrganization.OrganizationID,
		userOrganization.Role,
		time.Now(),
		time.Now())

	if err != nil {
		return fmt.Errorf("failed to insert user organization: %v", err)
	}

	return nil
}

func (p *PostgresImplementation) GetLastUserOrganizationByUserID(ctx context.Context, userID string) (*models.UserOrganizations, error) {
	var userOrganization models.UserOrganizations

	query := `
		SELECT id, user_id, organization_id, role, created_at, updated_at
		FROM user_organizations
		WHERE user_id = $1
		ORDER BY updated_at DESC
		LIMIT 1;`

	err := p.Pool.QueryRow(ctx, query, userID).Scan(
		&userOrganization.ID,
		&userOrganization.UserID,
		&userOrganization.OrganizationID,
		&userOrganization.Role,
		&userOrganization.CreatedAt,
		&userOrganization.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err // No rows found
		}
		return nil, fmt.Errorf("failed to retrieve last user organization: %v", err)
	}

	return &userOrganization, nil
}
func (p *PostgresImplementation) UpdateUserOrganizationUpdatedAt(ctx context.Context, id string) error {
	query := `
	UPDATE user_organizations 
	SET updated_at = $1 
	WHERE id = $2;`

	_, err := p.Pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update user organization updated_at: %v", err)
	}

	return nil
}

// giver org id and user id, return user organization
func (p *PostgresImplementation) GetUserOrganizationByOrgIDAndUserID(ctx context.Context, orgID, userID string) (*models.UserOrganizations, error) {
	var userOrganization models.UserOrganizations

	query := `
		SELECT id, user_id, organization_id, role, created_at, updated_at
		FROM user_organizations
		WHERE organization_id = $1 AND user_id = $2;`

	err := p.Pool.QueryRow(ctx, query, orgID, userID).Scan(
		&userOrganization.ID,
		&userOrganization.UserID,
		&userOrganization.OrganizationID,
		&userOrganization.Role,
		&userOrganization.CreatedAt,
		&userOrganization.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err // No rows found
		}
		return nil, fmt.Errorf("failed to retrieve user organization by org id and user id: %v", err)
	}

	return &userOrganization, nil
}
