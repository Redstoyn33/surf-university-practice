package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/glini/backend/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestClientRepo_InsertClient(t *testing.T) {
	db := NewTestDBWithMigrations()
	defer db.Close()

	repo := NewClientRepo(db)
	ctx := context.Background()

	client, err := repo.InsertClient(ctx, "test@test.com", "hash123")
	require.NoError(t, err)
	require.Greater(t, client.ID, int64(0))
	require.Equal(t, "test@test.com", client.Login)
	require.Equal(t, "hash123", client.PasswordHash)
}

func TestClientRepo_GetClientByLogin(t *testing.T) {
	db := NewTestDBWithMigrations()
	defer db.Close()

	repo := NewClientRepo(db)
	ctx := context.Background()

	_, err := repo.InsertClient(ctx, "test@test.com", "hash123")
	require.NoError(t, err)

	client, err := repo.GetClientByLogin(ctx, "test@test.com")
	require.NoError(t, err)
	require.Equal(t, "test@test.com", client.Login)
}

func TestClientRepo_GetClientByLogin_NotFound(t *testing.T) {
	db := NewTestDBWithMigrations()
	defer db.Close()

	repo := NewClientRepo(db)
	ctx := context.Background()

	_, err := repo.GetClientByLogin(ctx, "nonexistent")
	require.Error(t, err)
	require.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestClientRepo_InsertClient_Duplicate(t *testing.T) {
	db := NewTestDBWithMigrations()
	defer db.Close()

	repo := NewClientRepo(db)
	ctx := context.Background()

	_, err := repo.InsertClient(ctx, "test@test.com", "hash1")
	require.NoError(t, err)

	_, err = repo.InsertClient(ctx, "test@test.com", "hash2")
	require.Error(t, err)
	require.True(t, errors.Is(err, domain.ErrDuplicate))
}

func TestClientRepo_GetClientByID(t *testing.T) {
	db := NewTestDBWithMigrations()
	defer db.Close()

	repo := NewClientRepo(db)
	ctx := context.Background()

	created, err := repo.InsertClient(ctx, "test@test.com", "hash123")
	require.NoError(t, err)

	client, err := repo.GetClientByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, client.ID)
	require.Equal(t, "test@test.com", client.Login)
}
