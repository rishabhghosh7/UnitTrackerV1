-- name: GetProject :many
SELECT name, desc, created_ts, updated_ts FROM Project WHERE id IN (sqlc.slice('ids'));

-- name: 
