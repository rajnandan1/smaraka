package postgres

import (
	"context"
	"fmt"

	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/models"
)

// Insert a new secret
func (p *PostgresImplementation) InsertNewSecret(ctx context.Context, secret models.DbSecret) error {
	query := `
		INSERT INTO secrets (id, organization_id, secret_type, secret_value, current_state, creator_id, created_at, last_used_at, updated_at, secret_name)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW(), NOW(), $7)
		RETURNING id, organization_id, secret_type, secret_value, current_state, creator_id, created_at, last_used_at, updated_at;`

	err := p.Pool.QueryRow(ctx, query,
		secret.ID,
		secret.OrganizationID,
		secret.SecretType,
		secret.SecretValue,
		secret.CurrentState,
		secret.CreatorID,
		secret.SecretName,
	).Scan(
		&secret.ID,
		&secret.OrganizationID,
		&secret.SecretType,
		&secret.SecretValue,
		&secret.CurrentState,
		&secret.CreatorID,
		&secret.CreatedAt,
		&secret.LastUsedAt,
		&secret.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert secret: %v", err)
	}

	return nil
}

// Get a secret by Organization ID and Secret Type and Secret Value and status ACTIVE
func (p *PostgresImplementation) GetSecretByOrgAndValue(ctx context.Context, organizationID, secretType, secretValue string) (*models.DbSecret, error) {
	var secret models.DbSecret

	query := `
		SELECT id, organization_id, secret_type, secret_value, current_state, creator_id, created_at, last_used_at, updated_at, secret_name
		FROM secrets
		WHERE organization_id = $1 AND secret_type = $2 AND secret_value = $3 AND current_state = $4;`

	err := p.Pool.QueryRow(ctx, query, organizationID, secretType, secretValue, constants.SecretStatusActive).Scan(
		&secret.ID,
		&secret.OrganizationID,
		&secret.SecretType,
		&secret.SecretValue,
		&secret.CurrentState,
		&secret.CreatorID,
		&secret.CreatedAt,
		&secret.LastUsedAt,
		&secret.UpdatedAt,
		&secret.SecretName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret: %v", err)
	}

	return &secret, nil
}

// get by id
func (p *PostgresImplementation) GetSecretByID(ctx context.Context, id string) (*models.DbSecret, error) {
	var secret models.DbSecret

	query := `
		SELECT id, organization_id, secret_type, secret_value, current_state, creator_id, created_at, last_used_at, updated_at, secret_name
		FROM secrets
		WHERE id = $1;`

	err := p.Pool.QueryRow(ctx, query, id).Scan(
		&secret.ID,
		&secret.OrganizationID,
		&secret.SecretType,
		&secret.SecretValue,
		&secret.CurrentState,
		&secret.CreatorID,
		&secret.CreatedAt,
		&secret.LastUsedAt,
		&secret.UpdatedAt,
		&secret.SecretName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret: %v", err)
	}

	return &secret, nil
}

// return active secrets for an organization
func (p *PostgresImplementation) GetSecretsByOrganizationID(ctx context.Context, organizationID, secretType string) ([]*models.DbSecret, error) {
	var secrets []*models.DbSecret

	query := `
		SELECT id, organization_id, secret_type, secret_value, current_state, creator_id, created_at, last_used_at, updated_at, secret_name
		FROM secrets
		WHERE organization_id = $1 AND current_state = $2 AND secret_type = $3 order by created_at DESC;`

	rows, err := p.Pool.Query(ctx, query, organizationID, constants.SecretStatusActive, secretType)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secrets: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var secret models.DbSecret
		err := rows.Scan(
			&secret.ID,
			&secret.OrganizationID,
			&secret.SecretType,
			&secret.SecretValue,
			&secret.CurrentState,
			&secret.CreatorID,
			&secret.CreatedAt,
			&secret.LastUsedAt,
			&secret.UpdatedAt,
			&secret.SecretName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve secrets: %v", err)
		}

		secrets = append(secrets, &secret)
	}

	return secrets, nil

}

// deactivate a secret given id and organization id
func (p *PostgresImplementation) DeactivateSecret(ctx context.Context, id, organizationID string) error {
	query := `
		UPDATE secrets
		SET current_state = $1, updated_at = NOW()
		WHERE id = $2 AND organization_id = $3
		RETURNING id, organization_id, secret_type, secret_value, current_state, creator_id, created_at, last_used_at, updated_at, secret_name;`

	var secret models.DbSecret
	err := p.Pool.QueryRow(ctx, query, constants.SecretStatusInactive, id, organizationID).Scan(
		&secret.ID,
		&secret.OrganizationID,
		&secret.SecretType,
		&secret.SecretValue,
		&secret.CurrentState,
		&secret.CreatorID,
		&secret.CreatedAt,
		&secret.LastUsedAt,
		&secret.UpdatedAt,
		&secret.SecretName,
	)
	if err != nil {
		return fmt.Errorf("failed to deactivate secret: %v", err)
	}

	return nil
}

// get secret by secret value
func (p *PostgresImplementation) GetSecretByValue(ctx context.Context, secretValue string) (*models.DbSecret, error) {
	var secret models.DbSecret

	query := `
		SELECT id, organization_id, secret_type, secret_value, current_state, creator_id, created_at, last_used_at, updated_at, secret_name
		FROM secrets
		WHERE secret_value = $1 AND current_state = $2;`

	err := p.Pool.QueryRow(ctx, query, secretValue, constants.SecretStatusActive).Scan(
		&secret.ID,
		&secret.OrganizationID,
		&secret.SecretType,
		&secret.SecretValue,
		&secret.CurrentState,
		&secret.CreatorID,
		&secret.CreatedAt,
		&secret.LastUsedAt,
		&secret.UpdatedAt,
		&secret.SecretName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret: %v", err)
	}

	return &secret, nil
}
