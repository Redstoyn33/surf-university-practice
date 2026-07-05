-- +goose Up
CREATE TABLE programs (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    name          TEXT    NOT NULL,
    description   TEXT    NOT NULL DEFAULT '',
    max_capacity  INTEGER NOT NULL CHECK (max_capacity > 0)
);

-- +goose Down
DROP TABLE programs;
