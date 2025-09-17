-- Annotations CRUD operations

-- name: CreateAnnotation :one
INSERT INTO annotations (id, entity_id, annotation_type, content, metadata, agent_name)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetAnnotation :one
SELECT * FROM annotations
WHERE id = ?;

-- name: ListAnnotationsByEntity :many
SELECT * FROM annotations
WHERE entity_id = ?
ORDER BY created_at DESC;

-- name: ListAnnotationsByType :many
SELECT * FROM annotations
WHERE entity_id = ? AND annotation_type = ?
ORDER BY created_at DESC;

-- name: ListAnnotationsByAgent :many
SELECT * FROM annotations
WHERE agent_name = ?
ORDER BY created_at DESC;

-- name: UpdateAnnotation :one
UPDATE annotations
SET content = ?, metadata = ?
WHERE id = ?
RETURNING *;

-- name: DeleteAnnotation :exec
DELETE FROM annotations
WHERE id = ?;

-- name: DeleteAnnotationsByEntity :exec
DELETE FROM annotations
WHERE entity_id = ?;