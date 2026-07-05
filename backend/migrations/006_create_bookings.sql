-- +goose Up
CREATE TABLE bookings (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id           INTEGER NOT NULL REFERENCES clients(id) ON DELETE RESTRICT,
    slot_id             INTEGER NOT NULL REFERENCES slots(id) ON DELETE RESTRICT,
    status              TEXT    NOT NULL DEFAULT 'активна'
                            CHECK (status IN ('активна', 'отменена клиентом', 'отменена мастерской')),
    rental_selected     INTEGER NOT NULL DEFAULT 0 CHECK (rental_selected IN (0, 1)),
    created_at          TEXT    NOT NULL DEFAULT (datetime('now')),
    cancellation_reason TEXT,
    CHECK (cancellation_reason IS NULL OR status = 'отменена мастерской')
);

CREATE INDEX idx_bookings_client_id ON bookings(client_id);
CREATE INDEX idx_bookings_slot_id ON bookings(slot_id);

-- гарантия: один клиент — одна активная бронь на слот
CREATE UNIQUE INDEX idx_active_booking
    ON bookings(client_id, slot_id)
    WHERE status = 'активна';

-- +goose Down
DROP TABLE bookings;
