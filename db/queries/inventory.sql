-- name: UpsertInventoryBalance :one
INSERT INTO inventory_balances (
    node_id,
    sku_id,
    pku_id,
    lot_id,
    is_sealed,
    package_count,
    on_hand,
    allocated
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, 0
)
ON CONFLICT (node_id, pku_id, lot_id, is_sealed)
DO UPDATE SET
    package_count = inventory_balances.package_count + EXCLUDED.package_count,
    on_hand = inventory_balances.on_hand + EXCLUDED.on_hand,
    updated_at = now()
RETURNING *;

-- name: GetInventoryBalanceByID :one
SELECT * FROM inventory_balances
WHERE id = $1 LIMIT 1;

-- name: GetBalanceForUpdate :one
SELECT * FROM inventory_balances
WHERE id = $1
FOR UPDATE;

-- name: FindReservableBalances :many
SELECT b.*
FROM inventory_balances b
LEFT JOIN lots l ON b.lot_id = l.id
WHERE b.sku_id = $1
  AND b.pku_id = $2
  AND (b.on_hand - b.allocated) >= $3
  AND (l.status IS NULL OR l.status = 'AVAILABLE')
ORDER BY l.expires_at ASC NULLS LAST, b.on_hand DESC;

-- name: ApplyBalanceAllocation :one
UPDATE inventory_balances
SET
    allocated = allocated + $2,
    updated_at = now()
WHERE id = $1 AND (on_hand - allocated) >= $2
RETURNING *;

-- name: DeductBalanceFulfillment :one
UPDATE inventory_balances
SET
    on_hand = on_hand - $2,
    allocated = allocated - $2,
    package_count = GREATEST(0, package_count - $3),
    updated_at = now()
WHERE id = $1 AND allocated >= $2 AND on_hand >= $2
RETURNING *;

-- name: ReleaseBalanceAllocation :one
UPDATE inventory_balances
SET
    allocated = allocated - $2,
    updated_at = now()
WHERE id = $1 AND allocated >= $2
RETURNING *;

-- name: ListBalancesByNode :many
SELECT b.*, s.id, s.name AS sku_name, p.unit_name
FROM inventory_balances b
JOIN stock_keeping_units s ON b.sku_id = s.id
JOIN sku_packaging_units p ON b.pku_id = p.id
WHERE b.node_id = $1;

-- name: CreateAllocation :one
INSERT INTO inventory_allocations (
    order_line_id,
    balance_id,
    allocated_quantity
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetAllocationsByOrderLine :many
SELECT * FROM inventory_allocations
WHERE order_line_id = $1;

-- name: DeleteAllocation :execrows
DELETE FROM inventory_allocations
WHERE id = $1;

-- name: RecordMovement :one
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
