CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA IF NOT EXISTS inventory_service;

CREATE TABLE IF NOT exists inventory_service.uom (
    code VARCHAR(20) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT exists inventory_service.product_categories (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    parent_id UUID REFERENCES inventory_service.product_categories(id),
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT exists inventory_service.products (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category_id UUID REFERENCES inventory_service.product_categories(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    discontinued BOOLEAN NOT NULL DEFAULT FALSE
);


CREATE TABLE IF NOT exists inventory_service.skus (
    sku VARCHAR(50) PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES inventory_service.products(id),
    variant_attributes JSONB, -- color, size, etc.
    default_uom VARCHAR(20) NOT NULL REFERENCES inventory_service.uom(code),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT exists inventory_service.sku_inventory (
    sku VARCHAR(50) PRIMARY KEY REFERENCES inventory_service.skus(sku),
    current_stock DECIMAL(12, 3) NOT NULL DEFAULT 0 CHECK (current_stock >= 0),
    reserved_stock DECIMAL(12, 3) NOT NULL DEFAULT 0 CHECK (reserved_stock >= 0),
    min_stock_level DECIMAL(12, 3) NOT NULL DEFAULT 0,
    max_stock_level DECIMAL(12, 3),
    last_stock_update TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT exists inventory_service.sku_prices (
    sku VARCHAR(50) NOT NULL REFERENCES inventory_service.skus(sku),
    uom_code VARCHAR(20) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL,
    valid_from TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    valid_to TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (sku, currency, valid_from)
);

CREATE TABLE IF NOT exists inventory_service.reservation_history (
    id UUID PRIMARY KEY NOT NULL,
    order_id UUID, -- References order_service.orders(id)
    sku VARCHAR(50) REFERENCES inventory_service.skus(sku),
    quantity DECIMAL(10, 2) NOT NULL,
    uom VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('RESERVED', 'RELEASED')),
    reserved_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    released_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_skus_product ON inventory_service.skus(product_id);
CREATE INDEX idx_sku_prices_active ON inventory_service.sku_prices(sku, is_active, valid_from, valid_to);