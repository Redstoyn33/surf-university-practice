package repository

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"github.com/glini/backend/migrations"
	"github.com/pressly/goose/v3"
)

func NewDB(dbPath string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, fmt.Errorf("enable WAL: %w", err)
	}

	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, fmt.Errorf("set goose dialect: %w", err)
	}

	if err := goose.Up(db.DB, "."); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return db, nil
}

func NewTestDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("open test db: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	goose.SetBaseFS(migrations.FS)
	goose.SetLogger(goose.NopLogger())

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}

func MustRunMigrations(db *sql.DB) {
	if err := RunMigrations(db); err != nil {
		slog.Error("migrations failed", "error", err)
		panic(err)
	}
}

func NewTestDBWithMigrations() *sqlx.DB {
	db, err := NewTestDB()
	if err != nil {
		panic(err)
	}
	MustRunMigrations(db.DB)
	return db
}
