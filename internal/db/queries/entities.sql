-- Entities CRUD operations

-- name: CreateEntity :one
INSERT INTO entities (id, version_id, entity_type, name, data)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetEntity :one
SELECT * FROM entities
WHERE id = ?;

-- name: ListEntitiesByVersion :many
SELECT * FROM entities
WHERE version_id = ?
ORDER BY created_at DESC;

-- name: ListEntitiesByType :many
SELECT * FROM entities
WHERE version_id = ? AND entity_type = ?
ORDER BY created_at DESC;

-- name: UpdateEntity :one
UPDATE entities
SET name = ?, data = ?
WHERE id = ?
RETURNING *;

-- name: DeleteEntity :exec
DELETE FROM entities
WHERE id = ?;

-- name: CountEntitiesByType :one
SELECT COUNT(*) FROM entities
WHERE version_id = ? AND entity_type = ?;