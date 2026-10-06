-- name: CreateStorageNode :one
WITH new_node AS (
  INSERT INTO storage_nodes (
      parent_id,
      code,
      node_type,
      is_movable,
      max_weight_kg,
      max_volume_cm3,
      is_active
  ) VALUES (
      $1, $2, $3, $4, $5, $6, $7
  )
  RETURNING id, parent_id, path::text AS path, code, node_type, is_movable, max_weight_kg, max_volume_cm3, is_active, created_at
)
UPDATE storage_nodes n
SET path = CASE
  WHEN i.parent_id IS NULL THEN i.id::text::ltree
  ELSE (SELECT path FROM storage_nodes WHERE id = i.parent_id) || i.id::text::ltree
END
FROM new_node i
WHERE n.id = i.id
RETURNING n.*;


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

-- name: ReparentStorageNode :exec
UPDATE storage_nodes
SET
  parent_id = CASE 
    WHEN storage_nodes.id = t.t_id THEN p.p_id 
    ELSE storage_nodes.parent_id 
  END,
  path = p.p_path || subpath(storage_nodes.path, nlevel(t.t_path) - 1)
FROM 
  (SELECT sn1.id AS t_id, sn1.path AS t_path FROM storage_nodes sn1 WHERE sn1.id = $1) AS t
CROSS JOIN 
  (SELECT sn2.id AS p_id, sn2.path AS p_path FROM storage_nodes sn2 WHERE sn2.id = $2) AS p
WHERE storage_nodes.path <@ t.t_path
  AND NOT (p.p_path <@ t.t_path);

-- name: UpdateStorageNodePath :exec
UPDATE storage_nodes SET path = $2 WHERE id = $1;

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

-- name: RecordNodeRelocation :one
INSERT INTO node_relocations (
    node_id,
    from_parent_id,
    to_parent_id,
    from_path,
    to_path,
    operator_id
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;
