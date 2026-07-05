-- +goose Up
CREATE TABLE masters_programs (
    master_id  INTEGER NOT NULL REFERENCES masters(id) ON DELETE RESTRICT,
    program_id INTEGER NOT NULL REFERENCES programs(id) ON DELETE RESTRICT,
    PRIMARY KEY (master_id, program_id)
);

-- +goose Down
DROP TABLE masters_programs;
