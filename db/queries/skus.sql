-- name: CreateSKU :one
INSERT INTO stock_keeping_units (
    name,
    base_uom,
    is_discrete,
    requires_lot_tracking,
    attributes
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetSKUByID :one
SELECT * FROM stock_keeping_units
WHERE id = $1 LIMIT 1;

-- name: ListSKUs :many
SELECT * FROM stock_keeping_units
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: SearchSKUsByName :many
SELECT * FROM stock_keeping_units
WHERE name % $1
ORDER BY similarity(name, $1) DESC
LIMIT $2;

-- name: UpdateSKU :one
UPDATE stock_keeping_units
SET
    name = COALESCE(sqlc.narg('name'), name),
    is_discrete = COALESCE(sqlc.narg('is_discrete'), is_discrete),
    requires_lot_tracking = COALESCE(sqlc.narg('requires_lot_tracking'), requires_lot_tracking),
    attributes = COALESCE(sqlc.narg('attributes'), attributes)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteSKU :execrows
DELETE FROM stock_keeping_units
WHERE id = $1;

-- name: CreatePackagingUnit :one
INSERT INTO sku_packaging_units (
    sku_id,
    unit_name,
    conversion_factor,
    parent_packaging_unit_id,
    is_base_unit,
    allows_break_bulk,
    barcode,
    tare_weight_kg,
    gross_volume_cm3
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetPackagingUnitByID :one
SELECT * FROM sku_packaging_units
WHERE id = $1 LIMIT 1;

-- name: GetPackagingUnitByBarcode :one
SELECT * FROM sku_packaging_units
WHERE barcode = $1 LIMIT 1;

-- name: ListPackagingUnitsBySKU :many
SELECT * FROM sku_packaging_units
WHERE sku_id = $1
ORDER BY conversion_factor ASC;

-- name: UpdatePackagingUnit :one
UPDATE sku_packaging_units
SET
    unit_name = COALESCE(sqlc.narg('unit_name'), unit_name),
    conversion_factor = COALESCE(sqlc.narg('conversion_factor'), conversion_factor),
    allows_break_bulk = COALESCE(sqlc.narg('allows_break_bulk'), allows_break_bulk),
    barcode = COALESCE(sqlc.narg('barcode'), barcode),
    tare_weight_kg = COALESCE(sqlc.narg('tare_weight_kg'), tare_weight_kg),
    gross_volume_cm3 = COALESCE(sqlc.narg('gross_volume_cm3'), gross_volume_cm3)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeletePackagingUnit :execrows
DELETE FROM sku_packaging_units
WHERE id = $1;

-- name: CreateLot :one
INSERT INTO lots (
    sku_id,
    lot_number,
    status,
    manufactured_at,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetLotByID :one
SELECT * FROM lots
WHERE id = $1 LIMIT 1;

-- name: GetLotBySKUAndNumber :one
SELECT * FROM lots
WHERE sku_id = $1 AND lot_number = $2 LIMIT 1;

-- name: ListLotsBySKU :many
SELECT * FROM lots
WHERE sku_id = $1
ORDER BY expires_at ASC NULLS LAST;

-- name: UpdateLotStatus :one
UPDATE lots
SET status = $2
WHERE id = $1
RETURNING *;
