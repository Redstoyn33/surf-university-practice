package repository

import (
	"context"
	"testing"
	"time"

	"github.com/glini/backend/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestSlotRepo_QuerySlots(t *testing.T) {
	db := NewTestDBWithMigrations()
	defer db.Close()

	masterRepo := NewMasterRepo(db)
	programRepo := NewProgramRepo(db)
	repo := NewSlotRepo(db, masterRepo, programRepo)
	ctx := context.Background()

	err := SeedData(db)
	require.NoError(t, err)

	slots, err := repo.QuerySlots(ctx, domain.SlotFilter{})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(slots), 1)
	require.NotNil(t, slots[0].Program)
	require.NotNil(t, slots[0].Master)
}

func TestSlotRepo_QuerySlots_FilterByMaster(t *testing.T) {
	db := NewTestDBWithMigrations()
	defer db.Close()

	masterRepo := NewMasterRepo(db)
	programRepo := NewProgramRepo(db)
	repo := NewSlotRepo(db, masterRepo, programRepo)
	ctx := context.Background()

	err := SeedData(db)
	require.NoError(t, err)

	masterID := int64(3)
	slots, err := repo.QuerySlots(ctx, domain.SlotFilter{MasterID: &masterID})
	require.NoError(t, err)
	for _, s := range slots {
		require.Equal(t, masterID, s.MasterID)
	}
}

func TestSlotRepo_GetSlotByID_Found(t *testing.T) {
	db := NewTestDBWithMigrations()
	defer db.Close()

	masterRepo := NewMasterRepo(db)
	programRepo := NewProgramRepo(db)
	repo := NewSlotRepo(db, masterRepo, programRepo)
	ctx := context.Background()

	err := SeedData(db)
	require.NoError(t, err)

	slot, err := repo.GetSlotByID(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, slot.Program)
	require.NotNil(t, slot.Master)
	require.Greater(t, slot.TotalSpots, 0)
}

func TestSlotRepo_GetSlotByID_NotFound(t *testing.T) {
	db := NewTestDBWithMigrations()
	defer db.Close()

	masterRepo := NewMasterRepo(db)
	programRepo := NewProgramRepo(db)
	repo := NewSlotRepo(db, masterRepo, programRepo)
	ctx := context.Background()

	_, err := repo.GetSlotByID(ctx, 999)
	require.Error(t, err)
}

func TestSlotRepo_DecrementSpots(t *testing.T) {
	db := NewTestDBWithMigrations()
	defer db.Close()

	masterRepo := NewMasterRepo(db)
	programRepo := NewProgramRepo(db)
	repo := NewSlotRepo(db, masterRepo, programRepo)
	ctx := context.Background()

	now := time.Now().UTC().Format(time.RFC3339)
	end := time.Now().UTC().Add(2 * time.Hour).Format(time.RFC3339)

	_, err := db.Exec(`INSERT INTO masters (id, name, photo_url, level) VALUES (1, 'Test', 'http://test.com', 'опытный')`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO programs (id, name, max_capacity) VALUES (1, 'Test', 10)`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO slots (id, date_time, end_time, program_id, master_id, total_spots, available_spots) VALUES (1, $1, $2, 1, 1, 2, 2)`, now, end)
	require.NoError(t, err)

	err = repo.DecrementSpots(ctx, 1)
	require.NoError(t, err)

	var avail int
	err = db.Get(&avail, "SELECT available_spots FROM slots WHERE id = 1")
	require.NoError(t, err)
	require.Equal(t, 1, avail)
}

func TestSlotRepo_DecrementSpots_NoSpots(t *testing.T) {
	db := NewTestDBWithMigrations()
	defer db.Close()

	masterRepo := NewMasterRepo(db)
	programRepo := NewProgramRepo(db)
	repo := NewSlotRepo(db, masterRepo, programRepo)
	ctx := context.Background()

	now := time.Now().UTC().Format(time.RFC3339)
	end := time.Now().UTC().Add(2 * time.Hour).Format(time.RFC3339)

	_, err := db.Exec(`INSERT INTO masters (id, name, photo_url, level) VALUES (1, 'Test', 'http://test.com', 'опытный')`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO programs (id, name, max_capacity) VALUES (1, 'Test', 10)`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO slots (id, date_time, end_time, program_id, master_id, total_spots, available_spots) VALUES (1, $1, $2, 1, 1, 1, 0)`, now, end)
	require.NoError(t, err)

	err = repo.DecrementSpots(ctx, 1)
	require.Error(t, err)
}
