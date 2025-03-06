package postgres

import (
	"context"

	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/models"
)

// insert into Subscription
func (p *PostgresImplementation) InsertSubscription(ctx context.Context, subscription models.Subscription) (*models.Subscription, error) {
	query := `
		INSERT INTO Subscription (id, user_id, current_state, success_order_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, user_id, current_state, success_order_id,created_at, updated_at;`

	err := p.Pool.QueryRow(ctx, query,
		subscription.ID,
		subscription.UserID,
		subscription.CurrentState,
		subscription.SuccessOrderID,
	).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.CurrentState,
		&subscription.SuccessOrderID,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

// GetSubscriptionByUserID
func (p *PostgresImplementation) GetSubscriptionByUserID(ctx context.Context, userID string) (*models.Subscription, error) {
	var subscription models.Subscription

	query := `
		SELECT id, user_id, success_order_id, current_state, created_at, updated_at
		FROM Subscription
		WHERE user_id = $1;`

	err := p.Pool.QueryRow(ctx, query, userID).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.SuccessOrderID,
		&subscription.CurrentState,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

// is allowed to use
// if current_state == PAID or if created_at < 14 days, return true
func (p *PostgresImplementation) IsAllowedToUse(ctx context.Context, userID string) (bool, error) {
	var subscription models.Subscription

	query := `
		SELECT id, user_id, success_order_id, current_state, created_at, updated_at
		FROM Subscription
		WHERE user_id = $1 and (current_state = $2 or created_at > NOW() - INTERVAL '14 days');`

	err := p.Pool.QueryRow(ctx, query, userID, constants.SubscriptionTypeActive).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.SuccessOrderID,
		&subscription.CurrentState,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

// update subscription status to active given user_id, order_id
func (p *PostgresImplementation) UpdateSubscriptionStatus(ctx context.Context, subID, orderID string) error {
	query := `
		UPDATE Subscription
		SET current_state = $1, success_order_id = $2, updated_at = NOW()
		WHERE id = $3;`

	_, err := p.Pool.Exec(ctx, query, constants.SubscriptionTypeActive, orderID, subID)
	if err != nil {
		return err
	}

	return nil
}
