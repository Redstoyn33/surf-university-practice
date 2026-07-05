package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/glini/backend/internal/domain"
	"github.com/glini/backend/internal/repository"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

type bookingTestEnv struct {
	svc     *BookingService
	db      *sqlx.DB
	ctx     context.Context
	slotID  int64
}

func newBookingTestEnv(t *testing.T, hoursUntilSlot int) *bookingTestEnv {
	t.Helper()
	db := repository.NewTestDBWithMigrations()
	t.Cleanup(func() { db.Close() })

	masterRepo := repository.NewMasterRepo(db)
	programRepo := repository.NewProgramRepo(db)
	slotRepo := repository.NewSlotRepo(db, masterRepo, programRepo)
	bookingRepo := repository.NewBookingRepo(db, slotRepo, masterRepo, programRepo)
	svc := NewBookingService(bookingRepo, slotRepo)
	ctx := context.Background()

	_, err := db.Exec(`INSERT INTO clients (id, login, password_hash) VALUES (1, 'test@test.com', 'hash')`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO masters (id, name, photo_url, level) VALUES (1, 'Test Master', 'http://test.com', 'опытный')`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO programs (id, name, max_capacity) VALUES (1, 'Test Program', 10)`)
	require.NoError(t, err)

	future := time.Now().UTC().Add(time.Duration(hoursUntilSlot) * time.Hour)
	dt := future.Format(time.RFC3339)
	end := future.Add(2 * time.Hour).Format(time.RFC3339)

	_, err = db.Exec(
		`INSERT INTO slots (id, date_time, end_time, program_id, master_id, total_spots, available_spots)
		 VALUES (1, $1, $2, 1, 1, 5, 5)`, dt, end)
	require.NoError(t, err)

	if hoursUntilSlot > 0 {
		_, err = db.Exec(
			`INSERT INTO slots (id, date_time, end_time, program_id, master_id, total_spots, available_spots)
			 VALUES (2, $1, $2, 1, 1, 5, 0)`, dt, end)
		require.NoError(t, err)
	}

	return &bookingTestEnv{svc: svc, db: db, ctx: ctx, slotID: 1}
}

func TestBookingService_CreateBooking_Success(t *testing.T) {
	env := newBookingTestEnv(t, 24)

	booking, err := env.svc.CreateBooking(env.ctx, 1, env.slotID, false)
	require.NoError(t, err)
	require.Greater(t, booking.ID, int64(0))
	require.Equal(t, "активна", booking.Status)
}

func TestBookingService_CreateBooking_NoSpots(t *testing.T) {
	env := newBookingTestEnv(t, 24)

	_, err := env.svc.CreateBooking(env.ctx, 1, 2, false)
	require.Error(t, err)
	require.True(t, errors.Is(err, domain.ErrConflict))
}

func TestBookingService_CreateBooking_DoubleBooking(t *testing.T) {
	env := newBookingTestEnv(t, 24)

	_, err := env.svc.CreateBooking(env.ctx, 1, env.slotID, false)
	require.NoError(t, err)

	_, err = env.svc.CreateBooking(env.ctx, 1, env.slotID, false)
	require.Error(t, err)
	require.True(t, errors.Is(err, domain.ErrConflict))
}

func TestBookingService_CancelBooking_Success(t *testing.T) {
	env := newBookingTestEnv(t, 24)

	booking, err := env.svc.CreateBooking(env.ctx, 1, env.slotID, false)
	require.NoError(t, err)

	cancelled, err := env.svc.CancelBooking(env.ctx, booking.ID, 1)
	require.NoError(t, err)
	require.Equal(t, "отменена клиентом", cancelled.Status)

	var avail int
	err = env.db.Get(&avail, "SELECT available_spots FROM slots WHERE id = 1")
	require.NoError(t, err)
	require.Equal(t, 5, avail)
}

func TestBookingService_CancelBooking_TooSoon(t *testing.T) {
	env := newBookingTestEnv(t, 2)

	booking, err := env.svc.CreateBooking(env.ctx, 1, env.slotID, false)
	require.NoError(t, err)

	_, err = env.svc.CancelBooking(env.ctx, booking.ID, 1)
	require.Error(t, err)
	require.True(t, errors.Is(err, domain.ErrValidation))
}

func TestBookingService_CancelBooking_NotOwned(t *testing.T) {
	env := newBookingTestEnv(t, 24)

	_, err := env.db.Exec(`INSERT INTO clients (id, login, password_hash) VALUES (2, 'other@test.com', 'hash')`)
	require.NoError(t, err)

	booking, err := env.svc.CreateBooking(env.ctx, 1, env.slotID, false)
	require.NoError(t, err)

	_, err = env.svc.CancelBooking(env.ctx, booking.ID, 2)
	require.Error(t, err)
}
