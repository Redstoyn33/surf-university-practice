package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func setupBookingTest(t *testing.T) (*BookingRepo, *SlotRepo, context.Context) {
	t.Helper()
	db := NewTestDBWithMigrations()
	t.Cleanup(func() { db.Close() })

	masterRepo := NewMasterRepo(db)
	programRepo := NewProgramRepo(db)
	slotRepo := NewSlotRepo(db, masterRepo, programRepo)
	bookingRepo := NewBookingRepo(db, slotRepo, masterRepo, programRepo)
	ctx := context.Background()

	_, err := db.Exec(`INSERT INTO clients (id, login, password_hash) VALUES (1, 'test@test.com', 'hash')`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO masters (id, name, photo_url, level) VALUES (1, 'Test Master', 'http://test.com', 'опытный')`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO programs (id, name, max_capacity) VALUES (1, 'Test Program', 10)`)
	require.NoError(t, err)

	return bookingRepo, slotRepo, ctx
}

func TestBookingRepo_InsertBooking_Success(t *testing.T) {
	bookingRepo, slotRepo, ctx := setupBookingTest(t)

	now := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	end := time.Now().UTC().Add(26 * time.Hour).Format(time.RFC3339)

	_, err := bookingRepo.db.Exec(
		`INSERT INTO slots (id, date_time, end_time, program_id, master_id, total_spots, available_spots)
		 VALUES (1, $1, $2, 1, 1, 5, 5)`, now, end)
	require.NoError(t, err)

	booking, err := bookingRepo.InsertBooking(ctx, 1, 1, false)
	require.NoError(t, err)
	require.Greater(t, booking.ID, int64(0))
	require.Equal(t, "активна", booking.Status)

	var avail int
	err = bookingRepo.db.Get(&avail, "SELECT available_spots FROM slots WHERE id = 1")
	require.NoError(t, err)
	require.Equal(t, 4, avail)

	_ = slotRepo
}

func TestBookingRepo_InsertBooking_Duplicate(t *testing.T) {
	bookingRepo, _, ctx := setupBookingTest(t)

	now := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	end := time.Now().UTC().Add(26 * time.Hour).Format(time.RFC3339)

	_, err := bookingRepo.db.Exec(
		`INSERT INTO slots (id, date_time, end_time, program_id, master_id, total_spots, available_spots)
		 VALUES (1, $1, $2, 1, 1, 5, 5)`, now, end)
	require.NoError(t, err)

	_, err = bookingRepo.InsertBooking(ctx, 1, 1, false)
	require.NoError(t, err)

	_, err = bookingRepo.InsertBooking(ctx, 1, 1, false)
	require.Error(t, err)
}

func TestBookingRepo_QueryBookingsByClient(t *testing.T) {
	bookingRepo, _, ctx := setupBookingTest(t)

	now := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	end := time.Now().UTC().Add(26 * time.Hour).Format(time.RFC3339)

	_, err := bookingRepo.db.Exec(
		`INSERT INTO slots (id, date_time, end_time, program_id, master_id, total_spots, available_spots)
		 VALUES (1, $1, $2, 1, 1, 5, 5)`, now, end)
	require.NoError(t, err)

	_, err = bookingRepo.InsertBooking(ctx, 1, 1, false)
	require.NoError(t, err)

	bookings, err := bookingRepo.QueryBookingsByClient(ctx, 1, "")
	require.NoError(t, err)
	require.Len(t, bookings, 1)
	require.NotNil(t, bookings[0].Slot)
}

func TestBookingRepo_GetBookingByID(t *testing.T) {
	bookingRepo, _, ctx := setupBookingTest(t)

	now := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	end := time.Now().UTC().Add(26 * time.Hour).Format(time.RFC3339)

	_, err := bookingRepo.db.Exec(
		`INSERT INTO slots (id, date_time, end_time, program_id, master_id, total_spots, available_spots)
		 VALUES (1, $1, $2, 1, 1, 5, 5)`, now, end)
	require.NoError(t, err)

	created, err := bookingRepo.InsertBooking(ctx, 1, 1, false)
	require.NoError(t, err)

	booking, err := bookingRepo.GetBookingByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, booking.ID)
	require.NotNil(t, booking.Slot)
}

func TestBookingRepo_DecrementSpots_NoSpots(t *testing.T) {
	_, slotRepo, ctx := setupBookingTest(t)

	now := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	end := time.Now().UTC().Add(26 * time.Hour).Format(time.RFC3339)

	_, err := slotRepo.db.Exec(
		`INSERT INTO slots (id, date_time, end_time, program_id, master_id, total_spots, available_spots)
		 VALUES (1, $1, $2, 1, 1, 2, 0)`, now, end)
	require.NoError(t, err)

	err = slotRepo.DecrementSpots(ctx, 1)
	require.Error(t, err)
}
