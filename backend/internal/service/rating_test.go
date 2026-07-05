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

type ratingTestEnv struct {
	svc      *RatingService
	db       *sqlx.DB
	ctx      context.Context
	clientID int64
	slotID   int64
	masterID int64
}

func newRatingTestEnv(t *testing.T, hoursSinceEnd int) *ratingTestEnv {
	t.Helper()
	db := repository.NewTestDBWithMigrations()
	t.Cleanup(func() { db.Close() })

	masterRepo := repository.NewMasterRepo(db)
	programRepo := repository.NewProgramRepo(db)
	slotRepo := repository.NewSlotRepo(db, masterRepo, programRepo)
	bookingRepo := repository.NewBookingRepo(db, slotRepo, masterRepo, programRepo)
	clientRepo := repository.NewClientRepo(db)
	ratingRepo := repository.NewRatingRepo(db)

	svc := NewRatingService(ratingRepo, bookingRepo, slotRepo)
	_ = clientRepo
	ctx := context.Background()

	_, err := db.Exec(`INSERT INTO clients (id, login, password_hash) VALUES (1, 'test@test.com', 'hash')`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO masters (id, name, photo_url, level) VALUES (1, 'Test Master', 'http://test.com', 'опытный')`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO programs (id, name, max_capacity) VALUES (1, 'Test Program', 10)`)
	require.NoError(t, err)

	slotStart := time.Now().UTC().Add(-2 * time.Hour).Add(-time.Duration(hoursSinceEnd) * time.Hour)
	slotEnd := time.Now().UTC().Add(-time.Duration(hoursSinceEnd) * time.Hour)
	slotID := int64(1)

	_, err = db.Exec(
		`INSERT INTO slots (id, date_time, end_time, program_id, master_id, total_spots, available_spots)
		 VALUES ($1, $2, $3, 1, 1, 5, 0)`,
		slotID, slotStart.Format(time.RFC3339), slotEnd.Format(time.RFC3339))
	require.NoError(t, err)

	_, err = db.Exec(
		`INSERT INTO bookings (id, client_id, slot_id, status, created_at)
		 VALUES (1, 1, $1, 'активна', datetime('now'))`, slotID)
	require.NoError(t, err)

	return &ratingTestEnv{
		svc:      svc,
		db:       db,
		ctx:      ctx,
		clientID: 1,
		slotID:   slotID,
		masterID: 1,
	}
}

func TestRatingService_CreateRating_Success(t *testing.T) {
	env := newRatingTestEnv(t, 3)

	rating, err := env.svc.CreateRating(env.ctx, env.clientID, env.masterID, env.slotID, 5)
	require.NoError(t, err)
	require.Greater(t, rating.ID, int64(0))
	require.Equal(t, 5, rating.Score)
}

func TestRatingService_CreateRating_TooEarly(t *testing.T) {
	env := newRatingTestEnv(t, 0)

	_, err := env.svc.CreateRating(env.ctx, env.clientID, env.masterID, env.slotID, 5)
	require.Error(t, err)
	require.True(t, errors.Is(err, domain.ErrValidation))
}

func TestRatingService_CreateRating_TooLate(t *testing.T) {
	env := newRatingTestEnv(t, 49)

	_, err := env.svc.CreateRating(env.ctx, env.clientID, env.masterID, env.slotID, 5)
	require.Error(t, err)
	require.True(t, errors.Is(err, domain.ErrValidation))
}

func TestRatingService_CreateRating_Duplicate(t *testing.T) {
	env := newRatingTestEnv(t, 3)

	_, err := env.svc.CreateRating(env.ctx, env.clientID, env.masterID, env.slotID, 5)
	require.NoError(t, err)

	_, err = env.svc.CreateRating(env.ctx, env.clientID, env.masterID, env.slotID, 4)
	require.Error(t, err)
	require.True(t, errors.Is(err, domain.ErrConflict))
}

func TestRatingService_CreateRating_NoBooking(t *testing.T) {
	env := newRatingTestEnv(t, 3)

	_, err := env.svc.CreateRating(env.ctx, env.clientID, env.masterID, 999, 5)
	require.Error(t, err)
}
