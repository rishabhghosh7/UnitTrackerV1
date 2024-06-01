-- name: GetProject :many
SELECT id, name, desc, created_ts, updated_ts FROM Project WHERE id IN (sqlc.slice('ids'));

-- name: ListProjects :many
SELECT id, name, desc, created_ts, updated_ts FROM Project;

-- name: GetProjectByName :one
SELECT id, name, desc, created_ts, updated_ts FROM Project WHERE name = ?;

-- name: CreateProject :exec
INSERT INTO Project(name, desc, created_ts, updated_ts) VALUES(?, ?, ?, ?);

-- name: GetUnits :many
SELECT id, project_id, created_ts, updated_ts FROM Unit WHERE project_id IN (sqlc.slice('project_ids'));

-- name: AddUnit :exec
INSERT INTO Unit(project_id, created_ts, updated_ts) VALUES(?, ?, ?);


