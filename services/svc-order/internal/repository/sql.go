package repository

import (
	"context"
	"fmt"
	"ops-monorepo/services/svc-order/internal/model"
	sql "ops-monorepo/shared-libs/storage/postgres"
	"time"

	"github.com/google/uuid"
)

type (
	IOrderSQLRepository interface {
		BeginTransaction(ctx context.Context) (sql.PgxTx, error)
		RollbackTransaction(ctx context.Context, tx sql.PgxTx) error
		CommitTransaction(ctx context.Context, tx sql.PgxTx) error

		// insert order
		InsertOrderWithTx(ctx context.Context, tx sql.PgxTx, order *model.Order) error
		UpdateOrderWithTx(ctx context.Context, tx sql.PgxTx, order *model.Order) error
		UpdateOrderStatus(ctx context.Context, orderId uuid.UUID, status string) error
		UpdateOrderStatusWithTx(ctx context.Context, tx sql.PgxTx, orderId uuid.UUID, status string) error

		// insert item order
		InsertItemOrderWithTx(ctx context.Context, tx sql.PgxTx, itemOrder model.ItemOrder) error
		UpdateItemOrderWithTx(ctx context.Context, tx sql.PgxTx, itemOrder model.ItemOrder) error

		// insert order with items
		InsertOrderWithItems(ctx context.Context, order *model.Order, items []model.ItemOrder) error

		// get order
		GetOrderById(ctx context.Context, orderId uuid.UUID) (*model.Order, error)
		GetOrderItemsByOrderId(ctx context.Context, orderId uuid.UUID) ([]model.ItemOrder, error)
		GetOrderWithItems(ctx context.Context, orderId uuid.UUID) (*model.Order, []model.ItemOrder, error)
	}

	OrderSQLRepository struct {
		Pgx *sql.PostgresPgx
	}
)

func NewOrderRepository(pgx *sql.PostgresPgx) *OrderSQLRepository {
	return &OrderSQLRepository{
		Pgx: pgx,
	}
}

func (o *OrderSQLRepository) BeginTransaction(ctx context.Context) (sql.PgxTx, error) {
	return o.Pgx.Pool().Begin(ctx)
}

func (o *OrderSQLRepository) RollbackTransaction(ctx context.Context, tx sql.PgxTx) error {
	return tx.Rollback(ctx)
}

func (o *OrderSQLRepository) CommitTransaction(ctx context.Context, tx sql.PgxTx) error {
	return tx.Commit(ctx)
}

func (o *OrderSQLRepository) InsertOrderWithTx(ctx context.Context, tx sql.PgxTx, order *model.Order) error {
	query := `
		INSERT INTO order_service.orders (id, user_id, user_email, status, total_amount, currency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	if order.Id == uuid.Nil {
		order.Id = uuid.New()
	}

	now := time.Now()
	if order.CreatedAt.IsZero() {
		order.CreatedAt = now
	}
	if order.UpdateAt.IsZero() {
		order.UpdateAt = now
	}

	_, err := tx.Exec(ctx, query,
		order.Id,
		order.UserId,
		order.UserEmail,
		order.Status,
		order.TotalAmount,
		order.Currency,
		order.CreatedAt,
		order.UpdateAt,
	)

	return err
}

func (o *OrderSQLRepository) UpdateOrderWithTx(ctx context.Context, tx sql.PgxTx, order *model.Order) error {
	query := `
		UPDATE order_service.orders 
		SET user_id = $2, user_email = $3, status = $4, total_amount = $5, currency = $6, updated_at = $7
		WHERE id = $1
	`

	order.UpdateAt = time.Now()

	_, err := tx.Exec(ctx, query,
		order.Id,
		order.UserId,
		order.UserEmail,
		order.Status,
		order.TotalAmount,
		order.Currency,
		order.UpdateAt,
	)

	return err
}

func (o *OrderSQLRepository) UpdateOrderStatus(ctx context.Context, orderId uuid.UUID, status string) error {
	query := `
		UPDATE order_service.orders 
		SET status = $2, updated_at = $3
		WHERE id = $1
	`

	_, err := o.Pgx.Pool().Exec(ctx, query, orderId, status, time.Now())
	return err
}

func (o *OrderSQLRepository) UpdateOrderStatusWithTx(ctx context.Context, tx sql.PgxTx, orderId uuid.UUID, status string) error {
	query := `
		UPDATE order_service.orders 
		SET status = $2, updated_at = $3
		WHERE id = $1
	`

	_, err := tx.Exec(ctx, query, orderId, status, time.Now())
	return err
}

func (o *OrderSQLRepository) InsertItemOrderWithTx(ctx context.Context, tx sql.PgxTx, itemOrder model.ItemOrder) error {
	query := `
		INSERT INTO order_service.order_items (id, order_id, sku, quantity_per_uom,  price_per_uom, uom_code)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	if itemOrder.Id == uuid.Nil {
		itemOrder.Id = uuid.New()
	}

	_, err := tx.Exec(ctx, query,
		itemOrder.Id,
		itemOrder.OrderId,
		itemOrder.Sku,
		itemOrder.QuantityPerUom,
		itemOrder.PricePerUom,
		itemOrder.UomCode,
	)

	return err
}

func (o *OrderSQLRepository) UpdateItemOrderWithTx(ctx context.Context, tx sql.PgxTx, itemOrder model.ItemOrder) error {
	query := `
		UPDATE order_service.order_items 
		SET order_id = $2, sku = $3, quantity_per_uom = $4, price_per_uom = $5, uom_code = $6
		WHERE id = $1
	`

	_, err := tx.Exec(ctx, query,
		itemOrder.Id,
		itemOrder.OrderId,
		itemOrder.Sku,
		itemOrder.QuantityPerUom,
		itemOrder.PricePerUom,
		itemOrder.UomCode,
	)

	return err
}

// batch operations
func (o *OrderSQLRepository) InsertOrderWithItems(ctx context.Context, order *model.Order, items []model.ItemOrder) error {
	tx, err := o.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			o.RollbackTransaction(ctx, tx)
		}
	}()

	if err = o.InsertOrderWithTx(ctx, tx, order); err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	for _, item := range items {
		item.OrderId = order.Id
		if err = o.InsertItemOrderWithTx(ctx, tx, item); err != nil {
			return fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	if err = o.CommitTransaction(ctx, tx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetOrderById
func (o *OrderSQLRepository) GetOrderById(ctx context.Context, orderId uuid.UUID) (*model.Order, error) {
	query := `
		SELECT id, user_id, user_email, status, total_amount, currency, created_at, updated_at
		FROM order_service.orders 
		WHERE id = $1
	`

	var order model.Order
	row := o.Pgx.Pool().QueryRow(ctx, query, orderId)

	err := row.Scan(
		&order.Id,
		&order.UserId,
		&order.UserEmail,
		&order.Status,
		&order.TotalAmount,
		&order.Currency,
		&order.CreatedAt,
		&order.UpdateAt,
	)

	if err != nil {
		if err == sql.PgxErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &order, nil
}

func (o *OrderSQLRepository) GetOrderItemsByOrderId(ctx context.Context, orderId uuid.UUID) ([]model.ItemOrder, error) {
	query := `
		SELECT id, order_id, sku, quantity_per_uom,  price_per_uom, uom_code
		FROM order_service.order_items 
		WHERE order_id = $1
	`

	rows, err := o.Pgx.Pool().Query(ctx, query, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.ItemOrder
	for rows.Next() {
		var item model.ItemOrder
		err := rows.Scan(
			&item.Id,
			&item.OrderId,
			&item.Sku,
			&item.QuantityPerUom,
			&item.PricePerUom,
			&item.UomCode,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (o *OrderSQLRepository) GetOrderWithItems(ctx context.Context, orderId uuid.UUID) (*model.Order, []model.ItemOrder, error) {
	order, err := o.GetOrderById(ctx, orderId)
	if err != nil {
		return nil, nil, err
	}

	if order == nil {
		return nil, nil, nil
	}

	items, err := o.GetOrderItemsByOrderId(ctx, orderId)
	if err != nil {
		return nil, nil, err
	}

	return order, items, nil
}
