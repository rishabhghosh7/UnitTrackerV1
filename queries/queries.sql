-- name: GetProject :many
SELECT id, name, desc, created_ts, updated_ts FROM Project WHERE id IN (sqlc.slice('ids'));

-- name: GetProjectByName :one
SELECT id, name, desc, created_ts, updated_ts FROM Project WHERE name=?;

-- name: CreateProject :exec
INSERT INTO Project(name, desc, created_ts, updated_ts) VALUES(?, ?, ?, ?);

-- name: ListProjects :many
SELECT id, name, desc, created_ts, updated_ts FROM Project;

-- name: UpdateProject :exec
UPDATE Project SET desc = ? WHERE id = ?;

-- name: AddUnit :exec
INSERT INTO Unit(project_id, created_ts, updated_ts) VALUES(?, ?, ?);

-- name: GetUnits :many
SELECT id, project_id, created_ts, updated_ts FROM Unit WHERE project_id IN (sqlc.slice('project_ids'));


