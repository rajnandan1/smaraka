package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rajnandan1/smaraka/models"
)

// InsertNewUser inserts a new user into theusers table
func (p *PostgresImplementation) InsertNewUser(ctx context.Context, user models.Users) (*models.Users, error) {
	query := `
		INSERT INTO users (id, email, name, password_hash, created_at, seen_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`

	_, err := p.Pool.Exec(ctx, query,
		user.ID,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.CreatedAt,
		user.SeenAt,
		user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}

	return p.GetUserByID(ctx, user.ID)
}

// GetUserByEmail retrieves a user by their email from theusers table
func (p *PostgresImplementation) GetUserByEmail(ctx context.Context, email string) (*models.Users, error) {
	var user models.Users

	query := `
		SELECT id, email, name, created_at, seen_at, updated_at
		FROM users
		WHERE email = $1;`

	err := p.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.CreatedAt,
		&user.SeenAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err // No user found with that email
		}
		return nil, fmt.Errorf("failed to retrieve user: %v", err)
	}

	return &user, nil
}

// GetUserByID ...
func (p *PostgresImplementation) GetUserByID(ctx context.Context, id string) (*models.Users, error) {
	var user models.Users

	query := `
		SELECT id, email, name, created_at, seen_at, updated_at
		FROM users
		WHERE id = $1;`

	err := p.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.CreatedAt,
		&user.SeenAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err // No user found with that ID
		}
		return nil, fmt.Errorf("failed to retrieve user: %v", err)
	}

	return &user, nil
}

// UpdateUserSeenAtByID
func (p *PostgresImplementation) UpdateUserSeenAtByID(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET seen_at = now()
		WHERE id = $1;`

	_, err := p.Pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to update user seen at: %v", err)
	}

	return nil
}

// get password has by email id
func (p *PostgresImplementation) GetPasswordHashByEmail(ctx context.Context, email string) (string, error) {
	var passwordHash string

	query := `
		SELECT password_hash
		FROM users
		WHERE email = $1;`

	err := p.Pool.QueryRow(ctx, query, email).Scan(&passwordHash)

	if err != nil {
		if err == pgx.ErrNoRows {
			return "", err // No user found with that email
		}
		return "", fmt.Errorf("failed to retrieve password hash: %v", err)
	}

	return passwordHash, nil
}
