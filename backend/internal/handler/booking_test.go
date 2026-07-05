package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/glini/backend/internal/middleware"
	"github.com/glini/backend/internal/repository"
	"github.com/glini/backend/internal/service"
	"github.com/stretchr/testify/require"
)

type bookingTestEnv struct {
	h        *BookingHandler
	token    string
	clientID int64
	authH    *AuthHandler
}

func newBookingHandlerWithAuth(t *testing.T) *bookingTestEnv {
	t.Helper()
	db := repository.NewTestDBWithMigrations()
	t.Cleanup(func() { db.Close() })

	masterRepo := repository.NewMasterRepo(db)
	programRepo := repository.NewProgramRepo(db)
	slotRepo := repository.NewSlotRepo(db, masterRepo, programRepo)
	bookingRepo := repository.NewBookingRepo(db, slotRepo, masterRepo, programRepo)
	clientRepo := repository.NewClientRepo(db)

	authService := service.NewAuthService(clientRepo, "test-secret")
	bookingService := service.NewBookingService(bookingRepo, slotRepo)

	authH := NewAuthHandler(authService)
	bookingH := NewBookingHandler(bookingService)

	regBody := map[string]string{"login": "test@test.com", "password": "pass123"}
	b, _ := json.Marshal(regBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	authH.Register(w, req)

	var loginResp loginResponse
	json.Unmarshal(w.Body.Bytes(), &loginResp)

	future := time.Now().UTC().Add(24 * time.Hour)
	_, err := db.Exec(`INSERT INTO masters (id, name, photo_url, level) VALUES (1, 'Test Master', 'http://test.com', 'опытный')`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO programs (id, name, max_capacity) VALUES (1, 'Test Program', 10)`)
	require.NoError(t, err)
	_, err = db.Exec(
		`INSERT INTO slots (id, date_time, end_time, program_id, master_id, total_spots, available_spots)
		 VALUES (1, $1, $2, 1, 1, 5, 5)`,
		future.Format(time.RFC3339), future.Add(2*time.Hour).Format(time.RFC3339))
	require.NoError(t, err)

	return &bookingTestEnv{
		h:        bookingH,
		token:    loginResp.Token,
		clientID: loginResp.Client.ID,
		authH:    authH,
	}
}

func bookingRequest(method, target string, body []byte, clientID int64) *http.Request {
	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if clientID > 0 {
		ctx := context.WithValue(req.Context(), middleware.ClientIDKey, clientID)
		req = req.WithContext(ctx)
	}
	return req
}

func bookingRequestWithID(method, target, id string, body []byte, clientID int64) *http.Request {
	req := bookingRequest(method, target, body, clientID)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)
	return req
}

func TestBookingHandler_Create_201(t *testing.T) {
	env := newBookingHandlerWithAuth(t)

	body := map[string]any{"slotId": 1, "rentalSelected": false}
	b, _ := json.Marshal(body)
	req := bookingRequest(http.MethodPost, "/bookings", b, env.clientID)
	w := httptest.NewRecorder()
	env.h.Create(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
}

func TestBookingHandler_Create_401(t *testing.T) {
	env := newBookingHandlerWithAuth(t)

	body := map[string]any{"slotId": 1, "rentalSelected": false}
	b, _ := json.Marshal(body)
	req := bookingRequest(http.MethodPost, "/bookings", b, 0)
	w := httptest.NewRecorder()
	env.h.Create(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestBookingHandler_Cancel_200(t *testing.T) {
	env := newBookingHandlerWithAuth(t)

	body := map[string]any{"slotId": 1, "rentalSelected": false}
	b, _ := json.Marshal(body)
	req := bookingRequest(http.MethodPost, "/bookings", b, env.clientID)
	w := httptest.NewRecorder()
	env.h.Create(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var created struct {
		ID int64 `json:"id"`
	}
	json.Unmarshal(w.Body.Bytes(), &created)

	bookingIDStr := strconv.FormatInt(created.ID, 10)
	req2 := bookingRequestWithID(http.MethodPatch, "/bookings/"+bookingIDStr+"/cancel", bookingIDStr, nil, env.clientID)
	w2 := httptest.NewRecorder()
	env.h.Cancel(w2, req2)

	require.Equal(t, http.StatusOK, w2.Code)
}

func TestBookingHandler_Cancel_401(t *testing.T) {
	env := newBookingHandlerWithAuth(t)

	body := map[string]any{"slotId": 1, "rentalSelected": false}
	b, _ := json.Marshal(body)
	req := bookingRequest(http.MethodPost, "/bookings", b, env.clientID)
	w := httptest.NewRecorder()
	env.h.Create(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	req2 := bookingRequestWithID(http.MethodPatch, "/bookings/1/cancel", "1", nil, 0)
	w2 := httptest.NewRecorder()
	env.h.Cancel(w2, req2)

	require.Equal(t, http.StatusUnauthorized, w2.Code)
}

var _ = middleware.ClientIDFromContext
var _ = service.NewAuthService
