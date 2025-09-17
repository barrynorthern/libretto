-- Relationships CRUD operations

-- name: CreateRelationship :one
INSERT INTO relationships (id, version_id, from_entity_id, to_entity_id, relationship_type, properties)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetRelationship :one
SELECT * FROM relationships
WHERE id = ?;

-- name: ListRelationshipsByVersion :many
SELECT * FROM relationships
WHERE version_id = ?
ORDER BY created_at DESC;

-- name: ListRelationshipsByEntity :many
SELECT * FROM relationships
WHERE (from_entity_id = ? OR to_entity_id = ?)
ORDER BY created_at DESC;

-- name: ListRelationshipsByType :many
SELECT * FROM relationships
WHERE version_id = ? AND relationship_type = ?
ORDER BY created_at DESC;

-- name: GetRelationshipsBetweenEntities :many
SELECT * FROM relationships
WHERE from_entity_id = ? AND to_entity_id = ?;

-- name: UpdateRelationship :one
UPDATE relationships
SET properties = ?
WHERE id = ?
RETURNING *;

-- name: DeleteRelationship :exec
DELETE FROM relationships
WHERE id = ?;

-- name: DeleteRelationshipsByEntity :exec
DELETE FROM relationships
WHERE from_entity_id = ? OR to_entity_id = ?;