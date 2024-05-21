-- +goose Up
CREATE TABLE IF NOT EXISTS projects (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   name TEXT NOT NULL UNIQUE,
   desc TEXT
);

-- +goose Down
