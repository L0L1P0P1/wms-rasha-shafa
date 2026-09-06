-- +goose Up
CREATE TYPE movement_reason AS ENUM (
    'RECEIVING',        -- Dock receipt into temporary staging
    'PUTAWAY',          -- Moving from staging to bulk/active bin
    'PICK',             -- Order fulfillment pick from bin to cart/tote
    'INTERNAL_TRANSFER',-- Re-balancing stock between bins/aisles
    'CYCLE_COUNT',      -- Correcting discrepancy found during counting
    'DAMAGED',          -- Moving out of circulation to scrap/quarantine
    'RETURN'            -- Customer or RMA return processing
);

CREATE TABLE inventory_movements (
    id bigint 
        GENERATED ALWAYS AS IDENTITY 
        PRIMARY KEY,
        
    sku_id bigint NOT NULL REFERENCES stock_keeping_units(id),
    
    quantity integer NOT NULL CHECK (quantity > 0),

    -- NULL source_location_id implies creation from external origin (e.g., Receiving dock arrival)
    source_location_id bigint REFERENCES locations(id) ON DELETE RESTRICT,

    -- NULL destination_location_id implies exit from warehouse (e.g., Shipped order, Damaged write-off)
    destination_location_id bigint REFERENCES locations(id) ON DELETE RESTRICT,

    reason movement_reason NOT NULL,

    -- Reference to business entities driving this physical movement
    reference_id text COLLATE "C", -- e.g., 'PO-98214', 'PICK-00412', 'ORDER-5512'

    -- Who physically moved it (system user ID or device/terminal identifier)
    operator_id text COLLATE "C" NOT NULL,

    -- Traceability / batch tracking
    lot_number text COLLATE "C",

    created_at timestamptz NOT NULL DEFAULT now(),

    -- Prevent no-op transactions where source and destination are identical
    CONSTRAINT chk_source_destination_differ CHECK (
        source_location_id IS DISTINCT FROM destination_location_id
    ),

    -- Enforce that a movement must touch at least one physical location
    CONSTRAINT chk_has_at_least_one_location CHECK (
        source_location_id IS NOT NULL OR destination_location_id IS NOT NULL
    )
);

CREATE INDEX idx_movements_sku_history 
ON inventory_movements (sku_id, created_at DESC);

CREATE INDEX idx_movements_source_loc 
ON inventory_movements (source_location_id, created_at DESC) 
WHERE source_location_id IS NOT NULL;

CREATE INDEX idx_movements_dest_loc 
ON inventory_movements (destination_location_id, created_at DESC) 
WHERE destination_location_id IS NOT NULL;

CREATE INDEX idx_movements_reference 
ON inventory_movements (reference_id) 
WHERE reference_id IS NOT NULL;

-- +goose Down
