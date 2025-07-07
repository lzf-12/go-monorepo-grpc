package model

import "time"

const (
	ReservedStatus = "RESERVED"
	ReleasedStatus = "RELEASED"
)

// StockStatus represents the inventory status of a single SKU
type StockStatus struct {
	SKU               string  `json:"sku"`
	AvailableQuantity float64 `json:"available_quantity"`
	ReservedQuantity  float64 `json:"reserved_quantity"`
	TotalQuantity     float64 `json:"total_quantity"`
	SKU_UOM           string  `json:"sku_uom"`
	SKUPrice          float64 `json:"sku_price"`
	SKUCurrency       string  `json:"sku_currency"`
}

type ReservationHistory struct {
	Id         string     `json:"id"`
	OrderId    string     `json:"order_id"`
	Sku        string     `json:"sku"`
	Quantity   float64    `json:"quantity"`
	Uom        string     `json:"uom"`
	Status     string     `json:"status"`
	ReservedAt time.Time  `json:"reserved_at"`
	ReleasedAt *time.Time `json:"released_at"`
}
