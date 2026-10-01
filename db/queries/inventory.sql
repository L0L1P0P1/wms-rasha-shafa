-- name: InsertMovement :one
INSERT INTO inventory_movements (
    sku_id,
    pku_id,
    lot_id,
    package_count,
    quantity,
    source_node_id,
    destination_node_id,
    reason,
    reference_id,
    allocation_id,
    is_sealed,
    operator_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: ListMovementsBySKU :many
SELECT * FROM inventory_movements
WHERE sku_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListMovementsByNode :many
SELECT * FROM inventory_movements
WHERE source_node_id = $1 OR destination_node_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: LockBalanceForUpdate :one
SELECT id, on_hand, package_count, allocated, version 
FROM inventory_balances 
WHERE node_id = $1 AND pku_id = $2 
  AND lot_id IS NOT DISTINCT FROM $3 AND is_sealed = $4
FOR UPDATE;

-- name: ProjectSourceDeduction :execrows
UPDATE inventory_balances
SET on_hand = on_hand - $1, 
    package_count = package_count - $2, 
    version = version + 1,
    updated_at = now()
WHERE id = $3 AND version = $4;

-- name: ProjectDestinationAddition :exec
INSERT INTO inventory_balances (
    node_id, sku_id, pku_id, lot_id, is_sealed, on_hand, package_count
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT (node_id, pku_id, lot_id, is_sealed) 
DO UPDATE SET 
    on_hand = inventory_balances.on_hand + EXCLUDED.on_hand,
    package_count = inventory_balances.package_count + EXCLUDED.package_count,
    version = inventory_balances.version + 1,
    updated_at = now();

-- name: FindReservableBalances :many
SELECT b.*
FROM inventory_balances b
LEFT JOIN lots l ON b.lot_id = l.id
WHERE b.sku_id = $1
  AND b.pku_id = $2
  AND (b.on_hand - b.allocated) >= $3
  AND (l.status IS NULL OR l.status = 'AVAILABLE')
ORDER BY l.expires_at ASC NULLS LAST, b.on_hand DESC;

-- name: ListBalancesByNode :many
SELECT b.*, s.id AS sku_uuid, s.name AS sku_name, p.unit_name
FROM inventory_balances b
JOIN stock_keeping_units s ON b.sku_id = s.id
JOIN sku_packaging_units p ON b.pku_id = p.id
WHERE b.node_id = $1;

-- name: CreateAllocation :one
INSERT INTO inventory_allocations (
    order_line_id,
    balance_id,
    sku_id,
    allocated_quantity
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetAllocationsByOrderLine :many
SELECT * FROM inventory_allocations
WHERE order_line_id = $1;

-- name: IncrementAllocationPicked :one
UPDATE inventory_allocations
SET picked_quantity = picked_quantity + $2
WHERE id = $1 AND (picked_quantity + $2) <= allocated_quantity
RETURNING *;

-- name: DeleteAllocation :execrows
DELETE FROM inventory_allocations
WHERE id = $1;

-- name: UpdateBalanceAllocation :execrows
UPDATE inventory_balances
SET allocated = allocated + $1,
    version = version + 1,
    updated_at = now()
WHERE id = $2 AND version = $3 AND (on_hand - allocated - $1) >= 0;
