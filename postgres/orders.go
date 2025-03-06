package postgres

import (
	"context"

	"github.com/rajnandan1/smaraka/constants"
	"github.com/rajnandan1/smaraka/models"
)

// create order
func (p *PostgresImplementation) InsertOrder(ctx context.Context, order models.Orders) (*models.Orders, error) {
	query := `
		INSERT INTO Orders (id, user_id, customer_name, customer_email,subscription_id, variant_id, store_id, current_state, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, user_id, customer_name, customer_email, subscription_id, variant_id, store_id, current_state, created_at, updated_at;`

	err := p.Pool.QueryRow(ctx, query,
		order.ID,
		order.UserID,
		order.CustomerName,
		order.CustomerEmail,
		order.SubscriptionID,
		order.VariantID,
		order.StoreID,
		order.CurrentState,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.CustomerName,
		&order.CustomerEmail,
		&order.SubscriptionID,
		&order.VariantID,
		&order.StoreID,
		&order.CurrentState,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// Update order to PAID along with order object
func (p *PostgresImplementation) UpdateOrder(ctx context.Context, order_id, order_object string) (*models.Orders, error) {
	query := `
		UPDATE Orders
		SET current_state = $1, order_object = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, user_id, customer_name, customer_email, subscription_id, variant_id, store_id, current_state, order_object, created_at, updated_at;`

	var order models.Orders
	err := p.Pool.QueryRow(ctx, query, constants.OrderPaid, order_object, order_id).Scan(
		&order.ID,
		&order.UserID,
		&order.CustomerName,
		&order.CustomerEmail,
		&order.SubscriptionID,
		&order.VariantID,
		&order.StoreID,
		&order.CurrentState,
		&order.OrderObject,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
