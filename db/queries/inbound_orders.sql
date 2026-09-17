-- name: CreateInboundOrder :one
INSERT INTO inbound_orders (
    po_number,
    vendor_name,
    status,
    expected_delivery
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetInboundOrderByID :one
SELECT * FROM inbound_orders
WHERE id = $1 LIMIT 1;

-- name: GetInboundOrderByPONumber :one
SELECT * FROM inbound_orders
WHERE po_number = $1 LIMIT 1;

-- name: UpdateInboundOrderStatus :one
UPDATE inbound_orders
SET
    status = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: CreateInboundOrderLine :one
INSERT INTO inbound_order_lines (
    inbound_order_id,
    sku_id,
    pku_id,
    expected_quantity,
    received_quantity
) VALUES (
    $1, $2, $3, $4, 0
)
RETURNING *;

-- name: GetInboundOrderLines :many
SELECT il.*, s.id, s.name AS sku_name, p.unit_name
FROM inbound_order_lines il
JOIN stock_keeping_units s ON il.sku_id = s.id
JOIN sku_packaging_units p ON il.pku_id = p.id
WHERE il.inbound_order_id = $1;

-- name: IncrementInboundLineReceived :one
UPDATE inbound_order_lines
SET received_quantity = received_quantity + $2
WHERE id = $1
RETURNING *;

-- name: DeleteInboundOrder :execrows
DELETE FROM inbound_orders
WHERE id = $1;
