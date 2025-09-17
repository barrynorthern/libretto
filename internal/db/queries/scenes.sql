-- name: CreateScene :one
INSERT INTO scenes (id, title, summary, content)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetScene :one
SELECT * FROM scenes
WHERE id = ?;

-- name: ListScenes :many
SELECT * FROM scenes
ORDER BY created_at DESC;

-- name: UpdateScene :one
UPDATE scenes
SET title = ?, summary = ?, content = ?
WHERE id = ?
RETURNING *;

-- name: DeleteScene :exec
DELETE FROM scenes
WHERE id = ?;
