-- +goose Up
CREATE TABLE clients (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    login      TEXT    NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE UNIQUE INDEX idx_clients_login ON clients(login);

-- +goose Down
DROP TABLE clients;
