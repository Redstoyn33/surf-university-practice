package service

import (
	"context"
	"errors"
	"testing"

	"github.com/glini/backend/internal/domain"
	"github.com/glini/backend/internal/repository"
	"github.com/stretchr/testify/require"
)

func newTestAuthService(t *testing.T) *AuthService {
	t.Helper()
	db := repository.NewTestDBWithMigrations()
	t.Cleanup(func() { db.Close() })
	clientRepo := repository.NewClientRepo(db)
	return NewAuthService(clientRepo, "test-secret")
}

func TestAuthService_Register_Success(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	result, err := svc.Register(ctx, "test@test.com", "pass123")
	require.NoError(t, err)
	require.Greater(t, result.Client.ID, int64(0))
	require.Equal(t, "test@test.com", result.Client.Login)
	require.NotEmpty(t, result.Token)
}

func TestAuthService_Register_Duplicate(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	_, err := svc.Register(ctx, "test@test.com", "pass123")
	require.NoError(t, err)

	_, err = svc.Register(ctx, "test@test.com", "pass456")
	require.Error(t, err)
	require.True(t, errors.Is(err, domain.ErrConflict))
}

func TestAuthService_Register_ShortPassword(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	_, err := svc.Register(ctx, "test@test.com", "12345")
	require.Error(t, err)
	require.True(t, errors.Is(err, domain.ErrValidation))
}

func TestAuthService_Login_Success(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	_, err := svc.Register(ctx, "test@test.com", "pass123")
	require.NoError(t, err)

	result, err := svc.Login(ctx, "test@test.com", "pass123")
	require.NoError(t, err)
	require.Equal(t, "test@test.com", result.Client.Login)
	require.NotEmpty(t, result.Token)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	_, err := svc.Register(ctx, "test@test.com", "pass123")
	require.NoError(t, err)

	_, err = svc.Login(ctx, "test@test.com", "wrong")
	require.Error(t, err)
	require.True(t, errors.Is(err, domain.ErrForbidden))
}

func TestAuthService_Login_NotFound(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	_, err := svc.Login(ctx, "nonexistent", "pass123")
	require.Error(t, err)
	require.True(t, errors.Is(err, domain.ErrForbidden))
}
