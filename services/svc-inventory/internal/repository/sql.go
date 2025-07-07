package repository

import (
	"context"
	"errors"
	"fmt"
	"ops-monorepo/services/svc-inventory/internal/model"
	rg "ops-monorepo/shared-libs/regexp"
	sql "ops-monorepo/shared-libs/storage/postgres"

	"github.com/robaho/fixed"
)

type IInventorySQLRepository interface {
	BeginTransaction(ctx context.Context) (sql.PgxTx, error)
	RollbackTransaction(ctx context.Context, tx sql.PgxTx) error
	CommitTransaction(ctx context.Context, tx sql.PgxTx) error

	CheckStockWithMultipleSkus(ctx context.Context, skus []string) (data []model.StockStatus, missingSkus []string, err error)
	GetStockStatus(ctx context.Context, sku string) (*model.StockStatus, error)

	ReserveStock(ctx context.Context, orderId, sku string, quantity float64) error
	ReleaseStock(ctx context.Context, sku string, quantity float64) error

	GetReservationHistoryByOrderIdAndstatus(ctx context.Context, orderId string, status string) ([]model.ReservationHistory, error)
}

type InventorySQLRepository struct {
	Pgx *sql.PostgresPgx
}

func NewInventoryRepository(pgx *sql.PostgresPgx) IInventorySQLRepository {
	return &InventorySQLRepository{
		Pgx: pgx,
	}
}

func (o *InventorySQLRepository) BeginTransaction(ctx context.Context) (sql.PgxTx, error) {
	return o.Pgx.Pool().Begin(ctx)
}

func (o *InventorySQLRepository) RollbackTransaction(ctx context.Context, tx sql.PgxTx) error {
	return tx.Rollback(ctx)
}

func (o *InventorySQLRepository) CommitTransaction(ctx context.Context, tx sql.PgxTx) error {
	return tx.Commit(ctx)
}

// checks the availability of multiple SKUs
func (r *InventorySQLRepository) CheckStockWithMultipleSkus(ctx context.Context, skus []string) (data []model.StockStatus, missingSkus []string, err error) {
	if len(skus) == 0 {
		return nil, []string{}, errors.New("no SKUs provided")
	}

	query := `
		SELECT 
			s.sku,
			si.current_stock,
			si.reserved_stock,
			(si.current_stock - si.reserved_stock) as available_quantity,
			s.default_uom,
			sp.unit_price,
			sp.currency
		FROM 
			inventory_service.skus s
		JOIN 
			inventory_service.sku_inventory si ON s.sku = si.sku
		JOIN 
			inventory_service.sku_prices sp ON s.sku = sp.sku
		WHERE 
			s.sku = ANY($1)
			AND sp.is_active = true
			AND (sp.valid_to IS NULL OR sp.valid_to > NOW())
	`

	// clean hidden whitespaces
	q := rg.ReplaceWhitesWithSingleSpace(query)
	rows, err := r.Pgx.Pool().Query(ctx, q, skus)
	if err != nil {
		return nil, []string{}, fmt.Errorf("failed to query inventory: %w", err)
	}
	defer rows.Close()

	var results []model.StockStatus
	for rows.Next() {
		var item model.StockStatus
		err := rows.Scan(
			&item.SKU,
			&item.TotalQuantity,
			&item.ReservedQuantity,
			&item.AvailableQuantity,
			&item.SKU_UOM,
			&item.SKUPrice,
			&item.SKUCurrency,
		)
		if err != nil {
			return nil, []string{}, fmt.Errorf("failed to scan inventory row: %w", err)
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, []string{}, fmt.Errorf("error after scanning inventory rows: %w", err)
	}

	// if we didn't find all requested SKUs, identify which ones are missing
	if len(results) != len(skus) {
		foundSKUs := make(map[string]bool)
		for _, item := range results {
			foundSKUs[item.SKU] = true
		}

		var missingSKUs []string
		for _, sku := range skus {
			if !foundSKUs[sku] {
				missingSKUs = append(missingSKUs, sku)
			}
		}

		return results, missingSKUs, fmt.Errorf("some SKUs not found: %v", missingSKUs)
	}
	return results, []string{}, nil
}

// returns the stock status for a single SKU
func (r *InventorySQLRepository) GetStockStatus(ctx context.Context, sku string) (*model.StockStatus, error) {
	query := `
		SELECT 
			s.sku,
			si.current_stock,
			si.reserved_stock,
			(si.current_stock - si.reserved_stock) as available_quantity,
			s.default_uom,
			sp.unit_price
		FROM 
			inventory_service.skus s
		JOIN 
			inventory_service.sku_inventory si ON s.sku = si.sku
		JOIN 
			inventory_service.sku_prices sp ON s.sku = sp.sku
		WHERE 
			s.sku = $1
			AND sp.is_active = true
			AND (sp.valid_to IS NULL OR sp.valid_to > NOW())
		LIMIT 1
	`

	var item model.StockStatus
	err := r.Pgx.Pool().QueryRow(ctx, query, sku).Scan(
		&item.SKU,
		&item.TotalQuantity,
		&item.ReservedQuantity,
		&item.AvailableQuantity,
		&item.SKU_UOM,
		&item.SKUPrice,
	)

	if err != nil {
		if errors.Is(err, sql.PgxErrNoRows) {
			return nil, fmt.Errorf("SKU not found: %s", sku)
		}
		return nil, fmt.Errorf("failed to get inventory status: %w", err)
	}

	return &item, nil
}

// ReserveInventory reserves inventory for a single SKU
func (r *InventorySQLRepository) ReserveStock(ctx context.Context, orderId, sku string, quantity float64) error {
	pgx := r.Pgx.Pool()
	tx, err := pgx.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// get available quantity first
	var available float64
	err = tx.QueryRow(ctx,
		"SELECT (current_stock - reserved_stock) FROM inventory_service.sku_inventory WHERE sku = $1",
		sku,
	).Scan(&available)

	if err != nil {
		return fmt.Errorf("failed to check available quantity: %w", err)
	}

	if available < quantity {
		return fmt.Errorf("insufficient available quantity for SKU %s: requested %.2f, available %.2f",
			sku, quantity, available)
	}

	// increment reserved the inventory
	_, err = tx.Exec(ctx,
		"UPDATE inventory_service.sku_inventory SET reserved_stock = reserved_stock + $1 WHERE sku = $2",
		quantity, sku,
	)

	if err != nil {
		return fmt.Errorf("failed to reserve inventory: %w", err)
	}

	// insert reservation history
	_, err = tx.Exec(ctx,
		`INSERT INTO inventory_service.reservation_history 
		(id, order_id, sku, quantity, uom, status, reserved_at, released_at) 
		VALUES (gen_random_uuid(), $1, $2, $3, 'EA', $4, NOW(), NULL)`,
		orderId, sku, quantity, model.ReservedStatus,
	)

	if err != nil {
		return fmt.Errorf("failed to insert reservation history: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *InventorySQLRepository) GetReservationHistoryByOrderIdAndstatus(ctx context.Context, orderId string, status string) ([]model.ReservationHistory, error) {
	pgx := r.Pgx.Pool()

	query := `
		SELECT 
			id,
			order_id,
			sku,
			quantity,
			uom,
			status,
			reserved_at,
			released_at
		FROM inventory_service.reservation_history 
		WHERE order_id = $1 AND status = $2
		ORDER BY reserved_at DESC
	`

	rows, err := pgx.Query(ctx, query, orderId, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query reservation history: %w", err)
	}
	defer rows.Close()

	var histories []model.ReservationHistory
	for rows.Next() {
		var history model.ReservationHistory
		err := rows.Scan(
			&history.Id,
			&history.OrderId,
			&history.Sku,
			&history.Quantity,
			&history.Uom,
			&history.Status,
			&history.ReservedAt,
			&history.ReleasedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reservation history row: %w", err)
		}
		histories = append(histories, history)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", err)
	}

	return histories, nil
}

// ReleaseInventory releases reserved inventory for a SKU
func (r *InventorySQLRepository) ReleaseStock(ctx context.Context, sku string, quantity float64) error {
	tx, err := r.Pgx.Pool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// check reserved quantity first
	var reserved float64
	err = tx.QueryRow(ctx,
		"SELECT reserved_stock FROM inventory_service.sku_inventory WHERE sku = $1",
		sku,
	).Scan(&reserved)

	if err != nil {
		return fmt.Errorf("failed to check reserved quantity: %w", err)
	}

	if reserved < quantity {
		return fmt.Errorf("insufficient reserved quantity for SKU %s: requested to release %.2f, reserved %.2f",
			sku, quantity, reserved)
	}

	// release the inventory
	_, err = tx.Exec(ctx,
		"UPDATE inventory_service.sku_inventory SET reserved_stock = reserved_stock - $1 WHERE sku = $2",
		quantity, sku,
	)

	if err != nil {
		return fmt.Errorf("failed to release inventory: %w", err)
	}

	// update reservation history status by order_id to RELEASED

	return tx.Commit(ctx)
}

// updates the inventory levels for a SKU
func (r *InventorySQLRepository) IncrementInventory(ctx context.Context, sku string, adjustment fixed.Fixed) error {
	_, err := r.Pgx.Pool().Exec(ctx,
		`UPDATE inventory_service.sku_inventory 
		SET current_stock = current_stock + $1, 
			last_stock_update = NOW() 
		WHERE sku = $2`,
		adjustment, sku,
	)

	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	return nil
}
