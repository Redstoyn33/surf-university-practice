.PHONY: run build migrate-up migrate-down migrate-status test

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

migrate-up:
	goose -dir migrations sqlite3 ./data/glini.db up

migrate-down:
	goose -dir migrations sqlite3 ./data/glini.db down

migrate-status:
	goose -dir migrations sqlite3 ./data/glini.db status

test:
	go test ./... -v -count=1
