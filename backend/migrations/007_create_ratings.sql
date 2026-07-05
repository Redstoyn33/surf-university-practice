-- +goose Up
CREATE TABLE ratings (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id  INTEGER NOT NULL REFERENCES clients(id) ON DELETE RESTRICT,
    master_id  INTEGER NOT NULL REFERENCES masters(id) ON DELETE RESTRICT,
    slot_id    INTEGER NOT NULL REFERENCES slots(id) ON DELETE RESTRICT,
    score      INTEGER NOT NULL CHECK (score >= 1 AND score <= 5),
    created_at TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_ratings_master_id ON ratings(master_id);

-- один клиент — одна оценка за слот
CREATE UNIQUE INDEX idx_ratings_client_slot
    ON ratings(client_id, slot_id);

-- +goose Down
DROP TABLE ratings;
