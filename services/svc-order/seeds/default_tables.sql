CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA IF NOT EXISTS order_service;

CREATE TABLE IF NOT EXISTS order_service.orders (
    id UUID PRIMARY KEY not null DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    user_email VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('PENDING', 'CONFIRMED', 'CANCELLED')),
    total_amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_service.order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES order_service.orders(id) ON DELETE CASCADE,
    sku VARCHAR(50) NOT NULL, -- References inventory_service.skus(sku)
    quantity_per_uom DECIMAL(10, 2),
    price_per_uom DECIMAL(10, 2) NOT NULL,
    uom_code VARCHAR(20) NOT NULL, -- References inventory_service.units_of_measure(code)
    CONSTRAINT unique_order_sku UNIQUE (order_id, sku)
);  

CREATE INDEX IF NOT EXISTS idx_order_user ON order_service.orders(user_id);
CREATE INDEX IF NOT EXISTS idx_order_status ON order_service.orders(status);
CREATE INDEX IF NOT EXISTS idx_order_created ON order_service.orders(created_at);
CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_service.order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_sku ON order_service.order_items(sku);