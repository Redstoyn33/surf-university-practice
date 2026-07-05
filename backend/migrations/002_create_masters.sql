-- +goose Up
CREATE TABLE masters (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT   NOT NULL,
    photo_url  TEXT   NOT NULL,
    rating     REAL   NOT NULL DEFAULT 0.0 CHECK (rating = 0.0 OR (rating >= 1.0 AND rating <= 5.0)),
    level      TEXT   NOT NULL CHECK (level IN ('новичок', 'опытный'))
);

-- +goose Down
DROP TABLE masters;
