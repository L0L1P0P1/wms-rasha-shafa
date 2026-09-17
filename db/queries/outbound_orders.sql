-- name: CreateOutboundOrder :one
INSERT INTO outbound_orders (
    order_number,
    customer_name,
    status,
    priority
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetOutboundOrderByID :one
SELECT * FROM outbound_orders
WHERE id = $1 LIMIT 1;

-- name: GetOutboundOrderByNumber :one
SELECT * FROM outbound_orders
WHERE order_number = $1 LIMIT 1;

-- name: ListOrdersByStatus :many
SELECT * FROM outbound_orders
WHERE status = $1
ORDER BY priority DESC, created_at ASC
LIMIT $2 OFFSET $3;

-- name: UpdateOutboundOrderStatus :one
UPDATE outbound_orders
SET
    status = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: CreateOutboundOrderLine :one
INSERT INTO outbound_order_lines (
    order_id,
    sku_id,
    pku_id,
    requested_quantity,
    fulfilled_quantity
) VALUES (
    $1, $2, $3, $4, 0
)
RETURNING *;

-- name: GetOutboundOrderLines :many
SELECT ol.*, s.id, s.name AS sku_name, p.unit_name
FROM outbound_order_lines ol
JOIN stock_keeping_units s ON ol.sku_id = s.id
JOIN sku_packaging_units p ON ol.pku_id = p.id
WHERE ol.order_id = $1;

-- name: IncrementOrderLineFulfillment :one
UPDATE outbound_order_lines
SET fulfilled_quantity = fulfilled_quantity + $2
WHERE id = $1 AND (fulfilled_quantity + $2) <= requested_quantity
RETURNING *;

-- name: DeleteOutboundOrder :execrows
DELETE FROM outbound_orders
WHERE id = $1;
