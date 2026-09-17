-- name: CreateWarehouseTask :one
INSERT INTO warehouse_tasks (
    task_type,
    status,
    priority,
    sku_id,
    pku_id,
    lot_id,
    package_count,
    quantity,
    target_node_id,
    source_node_id,
    destination_node_id,
    order_line_id,
    allocation_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING *;

-- name: GetWarehouseTaskByID :one
SELECT * FROM warehouse_tasks
WHERE id = $1 LIMIT 1;

-- name: ListPendingTasks :many
SELECT * FROM warehouse_tasks
WHERE status = 'PENDING'
ORDER BY priority DESC, created_at ASC
LIMIT $1;

-- name: AssignTask :one
UPDATE warehouse_tasks
SET
    status = 'ASSIGNED',
    assigned_operator_id = $2
WHERE id = $1 AND status = 'PENDING'
RETURNING *;

-- name: CompleteTask :one
UPDATE warehouse_tasks
SET
    status = 'COMPLETED',
    completed_at = now()
WHERE id = $1 AND status IN ('ASSIGNED', 'IN_PROGRESS')
RETURNING *;

-- name: CancelTask :one
UPDATE warehouse_tasks
SET status = 'CANCELLED'
WHERE id = $1 AND status != 'COMPLETED'
RETURNING *;
