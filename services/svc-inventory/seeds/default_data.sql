-- Seed UOMs (unchanged as these use codes, not UUIDs)
INSERT INTO inventory_service.uom (code, name, description) VALUES
('EA', 'Each', 'Individual unit'),
('PK', 'Pack', 'Pack of items'),
('KG', 'Kilogram', 'Weight measurement'),
('G', 'Gram', 'Weight measurement'),
('L', 'Liter', 'Volume measurement'),
('M', 'Meter', 'Length measurement');

-- Seed Categories
INSERT INTO inventory_service.product_categories (id, name, description) VALUES
('9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d', 'Electronics', 'Electronic devices and components'),
('1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed', 'Clothing', 'Apparel and accessories'),
('6ec0bd7f-11c0-43da-975e-2a8ad9ebae0b', 'Groceries', 'Food and household items'),
('f47ac10b-58cc-4372-a567-0e02b2c3d479', 'Furniture', 'Home and office furniture'),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Books', 'Books and publications');

-- Seed Products
INSERT INTO inventory_service.products (id, name, description, category_id) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'Smartphone X', 'Latest smartphone model', '9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d'),
('550e8400-e29b-41d4-a716-446655440001', 'Wireless Headphones', 'Noise cancelling headphones', '9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d'),
('550e8400-e29b-41d4-a716-446655440002', 'T-Shirt', 'Cotton t-shirt', '1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed'),
('550e8400-e29b-41d4-a716-446655440003', 'Jeans', 'Denim jeans', '1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed'),
('550e8400-e29b-41d4-a716-446655440004', 'Rice 5kg', 'Premium quality rice', '6ec0bd7f-11c0-43da-975e-2a8ad9ebae0b'),
('550e8400-e29b-41d4-a716-446655440005', 'Olive Oil', 'Extra virgin olive oil', '6ec0bd7f-11c0-43da-975e-2a8ad9ebae0b'),
('550e8400-e29b-41d4-a716-446655440006', 'Office Chair', 'Ergonomic office chair', 'f47ac10b-58cc-4372-a567-0e02b2c3d479'),
('550e8400-e29b-41d4-a716-446655440007', 'Coffee Table', 'Modern coffee table', 'f47ac10b-58cc-4372-a567-0e02b2c3d479'),
('550e8400-e29b-41d4-a716-446655440008', 'Programming Book', 'Guide to Go programming', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
('550e8400-e29b-41d4-a716-446655440009', 'Cookbook', 'International recipes', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11');

-- Seed SKUs (updated with new product_id references)
INSERT INTO inventory_service.skus (sku, product_id, variant_attributes, default_uom) VALUES
('SMARTPHONE-X-BLACK', '550e8400-e29b-41d4-a716-446655440000', '{"color": "black", "storage": "128GB"}', 'EA'),
('SMARTPHONE-X-WHITE', '550e8400-e29b-41d4-a716-446655440000', '{"color": "white", "storage": "256GB"}', 'EA'),
('HEADPHONES-BLACK', '550e8400-e29b-41d4-a716-446655440001', '{"color": "black"}', 'EA'),
('TSHIRT-M-WHITE', '550e8400-e29b-41d4-a716-446655440002', '{"color": "white", "size": "M"}', 'EA'),
('TSHIRT-L-BLUE', '550e8400-e29b-41d4-a716-446655440002', '{"color": "blue", "size": "L"}', 'EA'),
('JEANS-32-BLUE', '550e8400-e29b-41d4-a716-446655440003', '{"color": "blue", "waist": 32, "length": 34}', 'EA'),
('RICE-5KG', '550e8400-e29b-41d4-a716-446655440004', '{"weight": 5}', 'PK'),
('OLIVE-OIL-1L', '550e8400-e29b-41d4-a716-446655440005', '{"volume": 1}', 'L'),
('CHAIR-BLACK', '550e8400-e29b-41d4-a716-446655440006', '{"color": "black"}', 'EA'),
('TABLE-WALNUT', '550e8400-e29b-41d4-a716-446655440007', '{"material": "walnut"}', 'EA'),
('GO-BOOK', '550e8400-e29b-41d4-a716-446655440008', '{"author": "John Doe"}', 'EA'),
('COOKBOOK-INTL', '550e8400-e29b-41d4-a716-446655440009', '{"cuisine": "international"}', 'EA'),
('SMARTPHONE-X-BLUE', '550e8400-e29b-41d4-a716-446655440000', '{"color": "blue", "storage": "64GB"}', 'EA'),
('HEADPHONES-WHITE', '550e8400-e29b-41d4-a716-446655440001', '{"color": "white"}', 'EA'),
('JEANS-30-BLACK', '550e8400-e29b-41d4-a716-446655440003', '{"color": "black", "waist": 30, "length": 32}', 'EA');

-- Seed SKU Inventory (unchanged as these reference SKU codes)
INSERT INTO inventory_service.sku_inventory (sku, current_stock, reserved_stock, min_stock_level, max_stock_level) VALUES
('SMARTPHONE-X-BLACK', 50, 5, 10, 100),
('SMARTPHONE-X-WHITE', 30, 2, 5, 50),
('HEADPHONES-BLACK', 100, 15, 20, 200),
('TSHIRT-M-WHITE', 200, 30, 50, 500),
('TSHIRT-L-BLUE', 150, 20, 30, 300),
('JEANS-32-BLUE', 75, 10, 15, 150),
('RICE-5KG', 500, 50, 100, 1000),
('OLIVE-OIL-1L', 300, 25, 50, 600),
('CHAIR-BLACK', 40, 3, 5, 80),
('TABLE-WALNUT', 25, 1, 3, 50),
('GO-BOOK', 80, 5, 10, 150),
('COOKBOOK-INTL', 60, 4, 5, 120),
('SMARTPHONE-X-BLUE', 20, 1, 5, 40),
('HEADPHONES-WHITE', 80, 8, 10, 150),
('JEANS-30-BLACK', 60, 7, 10, 120);

-- Seed SKU Prices
INSERT INTO inventory_service.sku_prices (sku, uom_code, currency, unit_price, valid_from, valid_to) VALUES
('SMARTPHONE-X-BLACK', 'EA', 'USD', 799.99, '2023-01-01', NULL),
('SMARTPHONE-X-WHITE', 'EA','USD', 899.99, '2023-01-01', NULL),
('HEADPHONES-BLACK', 'EA','USD', 199.99, '2023-01-01', NULL),
('TSHIRT-M-WHITE', 'EA','USD', 19.99, '2023-01-01', NULL),
('TSHIRT-L-BLUE', 'EA','USD', 19.99, '2023-01-01', NULL),
('JEANS-32-BLUE', 'EA','USD', 49.99, '2023-01-01', NULL),
('RICE-5KG', 'PK','USD', 12.99, '2023-01-01', NULL),
('OLIVE-OIL-1L', 'L', 'USD', 9.99, '2023-01-01', NULL),
('CHAIR-BLACK', 'EA','USD', 149.99, '2023-01-01', NULL),
('TABLE-WALNUT', 'EA','USD', 199.99, '2023-01-01', NULL),
('GO-BOOK', 'EA','USD', 39.99, '2023-01-01', NULL),
('COOKBOOK-INTL', 'EA','USD', 29.99, '2023-01-01', NULL),
('SMARTPHONE-X-BLUE', 'EA','USD', 699.99, '2023-01-01', NULL),
('HEADPHONES-WHITE', 'EA','USD', 199.99, '2023-01-01', NULL),
('JEANS-30-BLACK', 'EA','USD', 49.99, '2023-01-01', NULL);