-- Projects CRUD operations

-- name: CreateProject :one
INSERT INTO projects (id, name, theme, genre, description)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetProject :one
SELECT * FROM projects
WHERE id = ?;

-- name: ListProjects :many
SELECT * FROM projects
ORDER BY created_at DESC;

-- name: UpdateProject :one
UPDATE projects
SET name = ?, theme = ?, genre = ?, description = ?
WHERE id = ?
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = ?;