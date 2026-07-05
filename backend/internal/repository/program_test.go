package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProgramRepo_QueryPrograms(t *testing.T) {
	db := NewTestDBWithMigrations()
	defer db.Close()

	repo := NewProgramRepo(db)
	ctx := context.Background()

	err := SeedData(db)
	require.NoError(t, err)

	programs, err := repo.QueryPrograms(ctx)
	require.NoError(t, err)
	require.Len(t, programs, 3)
	require.Equal(t, "Лепка для новичков", programs[0].Name)
	require.Contains(t, programs[0].MasterIDs, int64(1))
}

func TestProgramRepo_GetProgramByID_Found(t *testing.T) {
	db := NewTestDBWithMigrations()
	defer db.Close()

	repo := NewProgramRepo(db)
	ctx := context.Background()

	err := SeedData(db)
	require.NoError(t, err)

	prog, err := repo.GetProgramByID(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, "Лепка для новичков", prog.Name)
	require.Len(t, prog.MasterIDs, 2)
}

func TestProgramRepo_GetProgramByID_NotFound(t *testing.T) {
	db := NewTestDBWithMigrations()
	defer db.Close()

	repo := NewProgramRepo(db)
	ctx := context.Background()

	_, err := repo.GetProgramByID(ctx, 999)
	require.Error(t, err)
}
