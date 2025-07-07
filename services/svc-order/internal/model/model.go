package model

import (
	inventoryv1 "pb_schemas/inventory/v1"
	"time"

	"github.com/google/uuid"
	"github.com/robaho/fixed"
)

const (
	ORDER_STATUS_PENDING            = "PENDING"
	ORDER_STATUS_CONFIRMED          = "CONFIRMED"
	ORDER_STATUS_FAILED_RESERVATION = "FAILED_RESERVATION"
	ORDER_STATUS_CANCELLED          = "CANCELLED"
)

type (
	Order struct {
		Id          uuid.UUID   `json:"uuid"`
		UserId      string      `json:"user_id"`
		UserEmail   string      `json:"user_email"`
		Status      string      `json:"status"`
		TotalAmount fixed.Fixed `json:"total_amount"`
		Currency    string      `json:"currency"`
		CreatedAt   time.Time   `json:"created_at"`
		UpdateAt    time.Time   `json:"updated_at"`
	}

	ItemOrder struct {
		Id             uuid.UUID   `json:"id"`
		OrderId        uuid.UUID   `json:"order_id"`
		Sku            string      `json:"sku"`
		QuantityPerUom fixed.Fixed `json:"quantity_per_uom"`
		PricePerUom    fixed.Fixed `json:"price_per_uom"`
		UomCode        string      `json:"uom_code"`
	}

	OrderWithItems struct {
		Order
		Items []ItemOrder `json:"items"`
	}

	OrderResponse struct {
		Order                OrderWithItems `json:"order"`
		FailedProcessedStock *inventoryv1.FailedProcessedItems
	}

	OrderedItemStockStatus = inventoryv1.InventoryStatus
)
