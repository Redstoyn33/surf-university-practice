package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/glini/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type ProgramRepo struct {
	db *sqlx.DB
}

func NewProgramRepo(db *sqlx.DB) *ProgramRepo {
	return &ProgramRepo{db: db}
}

func (r *ProgramRepo) QueryPrograms(ctx context.Context) ([]domain.Program, error) {
	rows, err := r.db.QueryxContext(ctx,
		`SELECT id, name, description, max_capacity FROM programs ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("query programs: %w", err)
	}
	defer rows.Close()

	var programs []domain.Program
	for rows.Next() {
		var p domain.Program
		if err := rows.StructScan(&p); err != nil {
			return nil, fmt.Errorf("scan program: %w", err)
		}
		programs = append(programs, p)
	}

	for i := range programs {
		masterIDs, err := r.getMasterIDs(ctx, programs[i].ID)
		if err != nil {
			return nil, fmt.Errorf("get masters for program %d: %w", programs[i].ID, err)
		}
		programs[i].MasterIDs = masterIDs
	}

	return programs, nil
}

func (r *ProgramRepo) GetProgramByID(ctx context.Context, id int64) (domain.Program, error) {
	var p domain.Program
	err := r.db.GetContext(ctx, &p,
		`SELECT id, name, description, max_capacity FROM programs WHERE id = $1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Program{}, fmt.Errorf("program not found: %w", domain.ErrNotFound)
		}
		return domain.Program{}, fmt.Errorf("get program by id: %w", err)
	}

	masterIDs, err := r.getMasterIDs(ctx, p.ID)
	if err != nil {
		return domain.Program{}, fmt.Errorf("get masters for program %d: %w", p.ID, err)
	}
	p.MasterIDs = masterIDs

	return p, nil
}

func (r *ProgramRepo) getMasterIDs(ctx context.Context, programID int64) ([]int64, error) {
	var ids []int64
	err := r.db.SelectContext(ctx, &ids,
		`SELECT master_id FROM masters_programs WHERE program_id = $1 ORDER BY master_id`, programID)
	if err != nil {
		return nil, err
	}
	if ids == nil {
		return []int64{}, nil
	}
	return ids, nil
}
