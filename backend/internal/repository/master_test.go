package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMasterRepo_QueryMasters(t *testing.T) {
	db := NewTestDBWithMigrations()
	defer db.Close()

	repo := NewMasterRepo(db)
	ctx := context.Background()

	err := SeedData(db)
	require.NoError(t, err)

	masters, err := repo.QueryMasters(ctx)
	require.NoError(t, err)
	require.Len(t, masters, 3)
	require.Equal(t, "Анна Кузнецова", masters[0].Name)
	require.Contains(t, masters[0].ProgramIDs, int64(1))
}

func TestMasterRepo_GetMasterByID_Found(t *testing.T) {
	db := NewTestDBWithMigrations()
	defer db.Close()

	repo := NewMasterRepo(db)
	ctx := context.Background()

	err := SeedData(db)
	require.NoError(t, err)

	master, err := repo.GetMasterByID(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, "Анна Кузнецова", master.Name)
	require.Len(t, master.ProgramIDs, 3)
}

func TestMasterRepo_GetMasterByID_NotFound(t *testing.T) {
	db := NewTestDBWithMigrations()
	defer db.Close()

	repo := NewMasterRepo(db)
	ctx := context.Background()

	_, err := repo.GetMasterByID(ctx, 999)
	require.Error(t, err)
}
