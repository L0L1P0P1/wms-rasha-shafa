-- 1. Physical Location Coordinates
CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    code VARCHAR(30) UNIQUE NOT NULL, -- e.g., 'A-01-2-A'
    zone VARCHAR(20) DEFAULT 'MAIN',  -- 'MAIN', 'COLD_STORAGE', 'RECEIVING', 'PACKING'
    is_active BOOLEAN DEFAULT TRUE
);

-- 2. Master Product Definitions
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL,
    unit VARCHAR(20) DEFAULT 'EA',
    min_stock_threshold INT DEFAULT 5,
    is_expirable BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 3. Lots and Expirations (Crucial for dental resins/anesthetics)
CREATE TABLE product_lots (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id),
    lot_number VARCHAR(50) NOT NULL,
    expiration_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(product_id, lot_number)
);

-- 4. Immutable Stock Movement Ledger
CREATE TYPE movement_reason AS ENUM (
    'PURCHASE_RECEIPT', 
    'PICK_ORDER', 
    'PUTAWAY', 
    'INTERNAL_TRANSFER', 
    'CYCLE_COUNT_ADJUSTMENT', 
    'SCRAP_EXPIRED'
);

CREATE TABLE stock_movements (
    id BIGSERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id),
    lot_id INT REFERENCES product_lots(id),
    from_location_id INT REFERENCES locations(id), -- NULL for vendor inbound
    to_location_id INT REFERENCES locations(id),   -- NULL for customer outbound
    quantity INT NOT NULL CHECK (quantity > 0),
    reason movement_reason NOT NULL,
    reference_id VARCHAR(100), -- Order ID, Purchase Order ID, or Invoice ID
    performed_by VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- View to get real-time on-hand stock by location and lot
CREATE VIEW v_stock_balances AS
SELECT 
    p.sku,
    p.name,
    l.code AS location_code,
    pl.lot_number,
    pl.expiration_date,
    SUM(
        CASE 
            WHEN sm.to_location_id = l.id THEN sm.quantity 
            WHEN sm.from_location_id = l.id THEN -sm.quantity 
            ELSE 0 
        END
    ) AS on_hand_qty
FROM stock_movements sm
JOIN products p ON sm.product_id = p.id
LEFT JOIN product_lots pl ON sm.lot_id = pl.id
JOIN locations l ON l.id IN (sm.from_location_id, sm.to_location_id)
GROUP BY p.sku, p.name, l.code, pl.lot_number, pl.expiration_date
HAVING SUM(CASE WHEN sm.to_location_id = l.id THEN sm.quantity WHEN sm.from_location_id = l.id THEN -sm.quantity ELSE 0 END) > 0;
