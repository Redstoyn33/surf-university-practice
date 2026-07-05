package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/glini/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type ClientRepo struct {
	db *sqlx.DB
}

func NewClientRepo(db *sqlx.DB) *ClientRepo {
	return &ClientRepo{db: db}
}

func (r *ClientRepo) InsertClient(ctx context.Context, login, passwordHash string) (domain.Client, error) {
	var c domain.Client
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO clients (login, password_hash) VALUES ($1, $2)
		 RETURNING id, login, password_hash, created_at`,
		login, passwordHash,
	).StructScan(&c)
	if err != nil {
		if isUniqueConstraintErr(err) {
			return domain.Client{}, fmt.Errorf("login already exists: %w", domain.ErrDuplicate)
		}
		return domain.Client{}, fmt.Errorf("insert client: %w", err)
	}
	return c, nil
}

func (r *ClientRepo) GetClientByLogin(ctx context.Context, login string) (domain.Client, error) {
	var c domain.Client
	err := r.db.GetContext(ctx, &c,
		`SELECT id, login, password_hash, created_at FROM clients WHERE login = $1`, login)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Client{}, fmt.Errorf("client not found: %w", domain.ErrNotFound)
		}
		return domain.Client{}, fmt.Errorf("get client by login: %w", err)
	}
	return c, nil
}

func (r *ClientRepo) GetClientByID(ctx context.Context, id int64) (domain.Client, error) {
	var c domain.Client
	err := r.db.GetContext(ctx, &c,
		`SELECT id, login, password_hash, created_at FROM clients WHERE id = $1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Client{}, fmt.Errorf("client not found: %w", domain.ErrNotFound)
		}
		return domain.Client{}, fmt.Errorf("get client by id: %w", err)
	}
	return c, nil
}
