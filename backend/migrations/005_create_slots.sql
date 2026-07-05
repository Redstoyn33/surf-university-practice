-- +goose Up
CREATE TABLE slots (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    date_time        TEXT    NOT NULL,
    end_time         TEXT    NOT NULL CHECK (end_time > date_time),
    program_id       INTEGER NOT NULL REFERENCES programs(id) ON DELETE RESTRICT,
    master_id        INTEGER NOT NULL REFERENCES masters(id) ON DELETE RESTRICT,
    total_spots      INTEGER NOT NULL CHECK (total_spots > 0),
    available_spots  INTEGER NOT NULL CHECK (available_spots >= 0 AND available_spots <= total_spots),
    rental_available INTEGER NOT NULL DEFAULT 0 CHECK (rental_available IN (0, 1)),
    rental_price     REAL    NOT NULL DEFAULT 0 CHECK (rental_price >= 0),
    CHECK (rental_available = 1 OR rental_price = 0)
);

CREATE INDEX idx_slots_date_time ON slots(date_time);
CREATE INDEX idx_slots_program_id ON slots(program_id);
CREATE INDEX idx_slots_master_id ON slots(master_id);

-- +goose Down
DROP TABLE slots;
