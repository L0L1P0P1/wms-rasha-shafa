-- name: CreateStorageNode :one
INSERT INTO storage_nodes (
    parent_id,
    path,
    code,
    node_type,
    is_movable,
    max_weight_kg,
    max_volume_cm3,
    is_active
) VALUES (
    $1,
    ''::ltree, -- Trigger set_storage_node_path computes the real path
    $2, $3, $4, $5, $6, $7
)
RETURNING id, parent_id, path::text AS path, code, node_type, is_movable, max_weight_kg, max_volume_cm3, is_active, created_at;

-- name: GetStorageNodeByID :one
SELECT id, parent_id, path::text AS path, code, node_type, is_movable, max_weight_kg, max_volume_cm3, is_active, created_at
FROM storage_nodes
WHERE id = $1 LIMIT 1;

-- name: GetStorageNodeByCode :one
SELECT id, parent_id, path::text AS path, code, node_type, is_movable, max_weight_kg, max_volume_cm3, is_active, created_at
FROM storage_nodes
WHERE code = $1 LIMIT 1;

-- name: ListSubtreeNodes :many
SELECT 
    node.id, 
    node.parent_id, 
    node.path::text AS path, 
    node.code, 
    node.node_type, 
    node.is_movable, 
    node.max_weight_kg, 
    node.max_volume_cm3, 
    node.is_active, 
    node.created_at
FROM storage_nodes AS node
WHERE node.path <@ (
    SELECT root.path 
    FROM storage_nodes AS root 
    WHERE root.id = $1
)
ORDER BY node.path;

-- name: ListDirectChildNodes :many
SELECT id, parent_id, path::text AS path, code, node_type, is_movable, max_weight_kg, max_volume_cm3, is_active, created_at
FROM storage_nodes
WHERE parent_id = $1
ORDER BY code ASC;

-- name: ReparentStorageNode :one
UPDATE storage_nodes
SET parent_id = $2
WHERE id = $1
RETURNING id, parent_id, path::text AS path, code, node_type, is_movable, max_weight_kg, max_volume_cm3, is_active, created_at;

-- name: UpdateStorageNodeDetails :one
UPDATE storage_nodes
SET
    code = COALESCE(sqlc.narg('code'), code),
    is_movable = COALESCE(sqlc.narg('is_movable'), is_movable),
    max_weight_kg = COALESCE(sqlc.narg('max_weight_kg'), max_weight_kg),
    max_volume_cm3 = COALESCE(sqlc.narg('max_volume_cm3'), max_volume_cm3),
    is_active = COALESCE(sqlc.narg('is_active'), is_active)
WHERE id = sqlc.arg('id')
RETURNING id, parent_id, path::text AS path, code, node_type, is_movable, max_weight_kg, max_volume_cm3, is_active, created_at;

-- name: DeleteStorageNode :execrows
DELETE FROM storage_nodes
WHERE id = $1;
