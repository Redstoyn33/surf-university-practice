package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glini/backend/internal/repository"
	"github.com/glini/backend/internal/service"
	"github.com/stretchr/testify/require"
)

func newTestAuthHandler(t *testing.T) *AuthHandler {
	t.Helper()
	db := repository.NewTestDBWithMigrations()
	t.Cleanup(func() { db.Close() })
	clientRepo := repository.NewClientRepo(db)
	authService := service.NewAuthService(clientRepo, "test-secret")
	return NewAuthHandler(authService)
}

func TestAuthHandler_Register_201(t *testing.T) {
	h := newTestAuthHandler(t)

	body := map[string]string{"login": "test@test.com", "password": "pass123"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var resp loginResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.NotEmpty(t, resp.Token)
	require.Equal(t, "test@test.com", resp.Client.Login)
}

func TestAuthHandler_Register_400(t *testing.T) {
	h := newTestAuthHandler(t)

	body := map[string]string{"login": "test@test.com", "password": "12345"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Register_409(t *testing.T) {
	h := newTestAuthHandler(t)

	body := map[string]string{"login": "test@test.com", "password": "pass123"}
	b, _ := json.Marshal(body)

	req1 := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(b))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	h.Register(w1, req1)
	require.Equal(t, http.StatusCreated, w1.Code)

	req2 := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(b))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	h.Register(w2, req2)
	require.Equal(t, http.StatusConflict, w2.Code)
}

func TestAuthHandler_Login_200(t *testing.T) {
	h := newTestAuthHandler(t)

	regBody := map[string]string{"login": "test@test.com", "password": "pass123"}
	b, _ := json.Marshal(regBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Register(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	loginBody := map[string]string{"login": "test@test.com", "password": "pass123"}
	b2, _ := json.Marshal(loginBody)
	req2 := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(b2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	h.Login(w2, req2)

	require.Equal(t, http.StatusOK, w2.Code)
	var resp loginResponse
	err := json.Unmarshal(w2.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.NotEmpty(t, resp.Token)
}

func TestAuthHandler_Login_401(t *testing.T) {
	h := newTestAuthHandler(t)

	body := map[string]string{"login": "test@test.com", "password": "pass123"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Register(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	loginBody := map[string]string{"login": "test@test.com", "password": "wrong"}
	b2, _ := json.Marshal(loginBody)
	req2 := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(b2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	h.Login(w2, req2)

	require.Equal(t, http.StatusUnauthorized, w2.Code)
}
