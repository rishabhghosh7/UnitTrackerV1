-- +goose Up
CREATE TABLE IF NOT EXISTS Unit (
   project_id INTEGER NOT NULL REFERENCES Project(id),
   create_ts INTEGER NOT NULL
);
