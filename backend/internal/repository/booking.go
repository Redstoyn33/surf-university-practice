package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/glini/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type BookingRepo struct {
	db         *sqlx.DB
	slotRepo   *SlotRepo
	masterRepo *MasterRepo
	programRepo *ProgramRepo
}

func NewBookingRepo(db *sqlx.DB, slotRepo *SlotRepo, masterRepo *MasterRepo, programRepo *ProgramRepo) *BookingRepo {
	return &BookingRepo{
		db:          db,
		slotRepo:    slotRepo,
		masterRepo:  masterRepo,
		programRepo: programRepo,
	}
}

func (r *BookingRepo) InsertBooking(ctx context.Context, clientID, slotID int64, rentalSelected bool) (domain.Booking, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return domain.Booking{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if err := r.slotRepo.DecrementSpotsTx(tx, slotID); err != nil {
		return domain.Booking{}, err
	}

	var b domain.Booking
	rental := 0
	if rentalSelected {
		rental = 1
	}
	err = tx.QueryRowxContext(ctx,
		`INSERT INTO bookings (client_id, slot_id, rental_selected, created_at)
		 VALUES ($1, $2, $3, datetime('now'))
		 RETURNING id, client_id, slot_id, status, rental_selected, created_at, cancellation_reason`,
		clientID, slotID, rental,
	).StructScan(&b)
	if err != nil {
		if isUniqueConstraintErr(err) {
			return domain.Booking{}, fmt.Errorf("duplicate booking: %w", domain.ErrConflict)
		}
		return domain.Booking{}, fmt.Errorf("insert booking: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return domain.Booking{}, fmt.Errorf("commit tx: %w", err)
	}

	if err := r.enrichBooking(ctx, &b); err != nil {
		return domain.Booking{}, err
	}

	return b, nil
}

func (r *BookingRepo) QueryBookingsByClient(ctx context.Context, clientID int64, statusFilter string) ([]domain.Booking, error) {
	query := `SELECT id, client_id, slot_id, status, rental_selected, created_at, cancellation_reason
		FROM bookings WHERE client_id = $1`
	args := []any{clientID}
	if statusFilter != "" {
		query += " AND status = $2"
		args = append(args, statusFilter)
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query bookings: %w", err)
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		var b domain.Booking
		if err := rows.StructScan(&b); err != nil {
			return nil, fmt.Errorf("scan booking: %w", err)
		}
		bookings = append(bookings, b)
	}

	for i := range bookings {
		if err := r.enrichBooking(ctx, &bookings[i]); err != nil {
			return nil, err
		}
	}

	return bookings, nil
}

func (r *BookingRepo) GetBookingByID(ctx context.Context, id int64) (domain.Booking, error) {
	var b domain.Booking
	err := r.db.GetContext(ctx, &b,
		`SELECT id, client_id, slot_id, status, rental_selected, created_at, cancellation_reason
		 FROM bookings WHERE id = $1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Booking{}, fmt.Errorf("booking not found: %w", domain.ErrNotFound)
		}
		return domain.Booking{}, fmt.Errorf("get booking: %w", err)
	}

	if err := r.enrichBooking(ctx, &b); err != nil {
		return domain.Booking{}, err
	}

	return b, nil
}

func (r *BookingRepo) UpdateBookingStatus(ctx context.Context, id int64, status string, reason *string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE bookings SET status = $1, cancellation_reason = $2 WHERE id = $3`,
		status, reason, id)
	if err != nil {
		return fmt.Errorf("update booking status: %w", err)
	}
	return nil
}

func (r *BookingRepo) CancelBookingTx(ctx context.Context, bookingID int64, slotID int64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`UPDATE bookings SET status = 'отменена клиентом', cancellation_reason = NULL WHERE id = $1`,
		bookingID)
	if err != nil {
		return fmt.Errorf("update booking: %w", err)
	}

	if err := r.slotRepo.IncrementSpotsTx(tx, slotID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *BookingRepo) enrichBooking(ctx context.Context, b *domain.Booking) error {
	if b.SlotID > 0 {
		slot, err := r.slotRepo.GetSlotByID(ctx, b.SlotID)
		if err != nil {
			return err
		}
		b.Slot = &slot
	}
	return nil
}

type slotRepoForBooking interface {
	GetSlotByID(ctx context.Context, id int64) (domain.Slot, error)
	DecrementSpotsTx(tx *sqlx.Tx, slotID int64) error
	IncrementSpotsTx(tx *sqlx.Tx, slotID int64) error
}

var _ slotRepoForBooking = (*SlotRepo)(nil)

func init() {
	// Ensure time import is used
	_ = time.Time{}
}
