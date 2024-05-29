
-- +goose Up

CREATE TABLE IF NOT EXISTS Project (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  desc TEXT,
  updated_ts INTEGER NOT NULL,
  created_ts INTEGER NOT NULL
);

INSERT INTO Project (name, desc, created_ts, updated_ts) VALUES ("work", "tracking pomodoros for work", CAST(strftime("%s", "now") AS INTEGER), CAST(strftime("%s", "now") AS INTEGER));
INSERT INTO Project (name, desc, created_ts, updated_ts) VALUES ("guitar", "i wanna play neon!", CAST(strftime("%s", "now") AS INTEGER), CAST(strftime("%s", "now") AS INTEGER));
INSERT INTO Project (name, desc, created_ts, updated_ts) VALUES ("sports", "i wanna play football", CAST(strftime("%s", "now") AS INTEGER), CAST(strftime("%s", "now") AS INTEGER));
INSERT INTO Project (name, desc, created_ts, updated_ts) VALUES ("gym", "i wanna lift weights", CAST(strftime("%s", "now") AS INTEGER), CAST(strftime("%s", "now") AS INTEGER));

CREATE TABLE IF NOT EXISTS Unit (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  project_id INTEGER NOT NULL REFERENCES Project(id),
  updated_ts INTEGER NOT NULL,
  created_ts INTEGER NOT NULL
);
