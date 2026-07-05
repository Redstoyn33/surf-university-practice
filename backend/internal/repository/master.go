package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/glini/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type MasterRepo struct {
	db *sqlx.DB
}

func NewMasterRepo(db *sqlx.DB) *MasterRepo {
	return &MasterRepo{db: db}
}

func (r *MasterRepo) QueryMasters(ctx context.Context) ([]domain.Master, error) {
	rows, err := r.db.QueryxContext(ctx,
		`SELECT id, name, photo_url, rating, level FROM masters ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("query masters: %w", err)
	}
	defer rows.Close()

	var masters []domain.Master
	for rows.Next() {
		var m domain.Master
		if err := rows.StructScan(&m); err != nil {
			return nil, fmt.Errorf("scan master: %w", err)
		}
		masters = append(masters, m)
	}

	for i := range masters {
		programIDs, err := r.getProgramIDs(ctx, masters[i].ID)
		if err != nil {
			return nil, fmt.Errorf("get programs for master %d: %w", masters[i].ID, err)
		}
		masters[i].ProgramIDs = programIDs
	}

	return masters, nil
}

func (r *MasterRepo) GetMasterByID(ctx context.Context, id int64) (domain.Master, error) {
	var m domain.Master
	err := r.db.GetContext(ctx, &m,
		`SELECT id, name, photo_url, rating, level FROM masters WHERE id = $1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Master{}, fmt.Errorf("master not found: %w", domain.ErrNotFound)
		}
		return domain.Master{}, fmt.Errorf("get master by id: %w", err)
	}

	programIDs, err := r.getProgramIDs(ctx, m.ID)
	if err != nil {
		return domain.Master{}, fmt.Errorf("get programs for master %d: %w", m.ID, err)
	}
	m.ProgramIDs = programIDs

	return m, nil
}

func (r *MasterRepo) getProgramIDs(ctx context.Context, masterID int64) ([]int64, error) {
	var ids []int64
	err := r.db.SelectContext(ctx, &ids,
		`SELECT program_id FROM masters_programs WHERE master_id = $1 ORDER BY program_id`, masterID)
	if err != nil {
		return nil, err
	}
	if ids == nil {
		return []int64{}, nil
	}
	return ids, nil
}
