-- +goose Up
CREATE TABLE IF NOT EXISTS Project (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   name TEXT NOT NULL UNIQUE,
   desc TEXT
);
INSERT INTO Project (name, desc) VALUES ("work", "tracking pomodoros for work");
INSERT INTO Project (name, desc) VALUES ("guitar", "i wanna play neon!");

-- +goose Down
DROP TABLE Project;
