-- Graph versions CRUD operations

-- name: CreateGraphVersion :one
INSERT INTO graph_versions (id, project_id, parent_version_id, name, description, is_working_set)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetGraphVersion :one
SELECT * FROM graph_versions
WHERE id = ?;

-- name: ListGraphVersionsByProject :many
SELECT * FROM graph_versions
WHERE project_id = ?
ORDER BY created_at DESC;

-- name: GetWorkingSetVersion :one
SELECT * FROM graph_versions
WHERE project_id = ? AND is_working_set = TRUE;

-- name: UpdateGraphVersion :one
UPDATE graph_versions
SET name = ?, description = ?
WHERE id = ?
RETURNING *;

-- name: SetWorkingSet :exec
UPDATE graph_versions
SET is_working_set = CASE WHEN id = ? THEN TRUE ELSE FALSE END
WHERE project_id = ?;

-- name: DeleteGraphVersion :exec
DELETE FROM graph_versions
WHERE id = ?;