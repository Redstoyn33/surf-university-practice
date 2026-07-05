package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/glini/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type SlotRepo struct {
	db          *sqlx.DB
	masterRepo  *MasterRepo
	programRepo *ProgramRepo
}

func NewSlotRepo(db *sqlx.DB, masterRepo *MasterRepo, programRepo *ProgramRepo) *SlotRepo {
	return &SlotRepo{db: db, masterRepo: masterRepo, programRepo: programRepo}
}

func (r *SlotRepo) QuerySlots(ctx context.Context, filter domain.SlotFilter) ([]domain.Slot, error) {
	query := `SELECT id, date_time, end_time, program_id, master_id,
		total_spots, available_spots, rental_available, rental_price
		FROM slots WHERE 1=1`
	args := []any{}
	argIdx := 1

	if filter.DateFrom != "" {
		query += fmt.Sprintf(" AND date_time >= $%d", argIdx)
		args = append(args, filter.DateFrom)
		argIdx++
	}
	if filter.DateTo != "" {
		query += fmt.Sprintf(" AND date_time <= $%d", argIdx)
		args = append(args, filter.DateTo+"T23:59:59Z")
		argIdx++
	}
	if filter.MasterID != nil {
		query += fmt.Sprintf(" AND master_id = $%d", argIdx)
		args = append(args, *filter.MasterID)
		argIdx++
	}
	if filter.ProgramID != nil {
		query += fmt.Sprintf(" AND program_id = $%d", argIdx)
		args = append(args, *filter.ProgramID)
		argIdx++
	}
	query += " ORDER BY date_time"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query slots: %w", err)
	}
	defer rows.Close()

	var slots []domain.Slot
	for rows.Next() {
		var s domain.Slot
		if err := rows.StructScan(&s); err != nil {
			return nil, fmt.Errorf("scan slot: %w", err)
		}
		slots = append(slots, s)
	}

	for i := range slots {
		if err := r.enrichSlot(ctx, &slots[i]); err != nil {
			return nil, err
		}
	}

	return slots, nil
}

func (r *SlotRepo) GetSlotByID(ctx context.Context, id int64) (domain.Slot, error) {
	var s domain.Slot
	err := r.db.GetContext(ctx, &s,
		`SELECT id, date_time, end_time, program_id, master_id,
		 total_spots, available_spots, rental_available, rental_price
		 FROM slots WHERE id = $1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Slot{}, fmt.Errorf("slot not found: %w", domain.ErrNotFound)
		}
		return domain.Slot{}, fmt.Errorf("get slot by id: %w", err)
	}

	if err := r.enrichSlot(ctx, &s); err != nil {
		return domain.Slot{}, err
	}

	return s, nil
}

func (r *SlotRepo) DecrementSpots(ctx context.Context, slotID int64) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE slots SET available_spots = available_spots - 1
		 WHERE id = $1 AND available_spots > 0`, slotID)
	if err != nil {
		return fmt.Errorf("decrement spots: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("no available spots: %w", domain.ErrConflict)
	}
	return nil
}

func (r *SlotRepo) DecrementSpotsTx(tx *sqlx.Tx, slotID int64) error {
	res, err := tx.Exec(
		`UPDATE slots SET available_spots = available_spots - 1
		 WHERE id = $1 AND available_spots > 0`, slotID)
	if err != nil {
		return fmt.Errorf("decrement spots: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("no available spots: %w", domain.ErrConflict)
	}
	return nil
}

func (r *SlotRepo) IncrementSpots(ctx context.Context, slotID int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE slots SET available_spots = available_spots + 1
		 WHERE id = $1`, slotID)
	if err != nil {
		return fmt.Errorf("increment spots: %w", err)
	}
	return nil
}

func (r *SlotRepo) IncrementSpotsTx(tx *sqlx.Tx, slotID int64) error {
	_, err := tx.Exec(
		`UPDATE slots SET available_spots = available_spots + 1
		 WHERE id = $1`, slotID)
	if err != nil {
		return fmt.Errorf("increment spots: %w", err)
	}
	return nil
}

func (r *SlotRepo) enrichSlot(ctx context.Context, s *domain.Slot) error {
	if s.ProgramID > 0 {
		program, err := r.programRepo.GetProgramByID(ctx, s.ProgramID)
		if err != nil {
			return err
		}
		s.Program = &program
	}
	if s.MasterID > 0 {
		master, err := r.masterRepo.GetMasterByID(ctx, s.MasterID)
		if err != nil {
			return err
		}
		s.Master = &master
	}
	return nil
}

func buildSlotFilter(values map[string]any) (string, []any) {
	where := []string{}
	args := []any{}
	idx := 1
	for k, v := range values {
		where = append(where, fmt.Sprintf("%s = $%d", k, idx))
		args = append(args, v)
		idx++
	}
	return strings.Join(where, " AND "), args
}
