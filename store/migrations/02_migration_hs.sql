-- +goose Up

INSERT INTO Project (name, desc) VALUES ("sports", "i wanna play football");
INSERT INTO Project (name, desc) VALUES ("gym", "i wanna lift weights");

-- +goose Down
DELETE * FROM Project WHERE name IN ("sports", "gym");
