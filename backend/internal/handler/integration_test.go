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
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/glini/backend/internal/middleware"
	"github.com/glini/backend/internal/repository"
	"github.com/glini/backend/internal/service"
	"github.com/stretchr/testify/require"
)

func newIntegrationEnv(t *testing.T) (*chi.Mux, context.Context, int64, string) {
	t.Helper()
	db := repository.NewTestDBWithMigrations()
	t.Cleanup(func() { db.Close() })

	masterRepo := repository.NewMasterRepo(db)
	programRepo := repository.NewProgramRepo(db)
	slotRepo := repository.NewSlotRepo(db, masterRepo, programRepo)
	bookingRepo := repository.NewBookingRepo(db, slotRepo, masterRepo, programRepo)
	clientRepo := repository.NewClientRepo(db)
	ratingRepo := repository.NewRatingRepo(db)

	authService := service.NewAuthService(clientRepo, "test-secret")
	bookingService := service.NewBookingService(bookingRepo, slotRepo)
	ratingService := service.NewRatingService(ratingRepo, bookingRepo, slotRepo)

	authHandler := NewAuthHandler(authService)
	slotHandler := NewSlotHandler(slotRepo)
	masterHandler := NewMasterHandler(masterRepo)
	programHandler := NewProgramHandler(programRepo)
	bookingHandler := NewBookingHandler(bookingService)
	ratingHandler := NewRatingHandler(ratingService)

	r := chi.NewRouter()
	authMw := middleware.Auth("test-secret")

	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)
	r.Get("/slots", slotHandler.ListSlots)
	r.Get("/slots/{id}", slotHandler.GetSlot)
	r.Get("/masters", masterHandler.List)
	r.Get("/masters/{id}", masterHandler.Get)
	r.Get("/programs", programHandler.List)
	r.Get("/programs/{id}", programHandler.Get)
	r.With(authMw).Post("/bookings", bookingHandler.Create)
	r.With(authMw).Get("/bookings", bookingHandler.ListMy)
	r.With(authMw).Get("/bookings/{id}", bookingHandler.Get)
	r.With(authMw).Patch("/bookings/{id}/cancel", bookingHandler.Cancel)
	r.With(authMw).Post("/ratings", ratingHandler.Create)

	_, err := db.Exec(`INSERT INTO masters (id, name, photo_url, level) VALUES (1, 'Тестовый Мастер', 'http://test.com/photo.jpg', 'опытный')`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO programs (id, name, max_capacity) VALUES (1, 'Тестовая Программа', 6)`)
	require.NoError(t, err)

	future := time.Now().UTC().Add(24 * time.Hour)
	_, err = db.Exec(
		`INSERT INTO slots (id, date_time, end_time, program_id, master_id, total_spots, available_spots, rental_available, rental_price)
		 VALUES (1, $1, $2, 1, 1, 5, 5, 1, 500)`,
		future.Format(time.RFC3339), future.Add(2*time.Hour).Format(time.RFC3339))
	require.NoError(t, err)

	ctx := context.Background()
	return r, ctx, int64(0), ""
}

func integrationRequest(method, target, id string, body []byte, token string) *http.Request {
	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if id != "" {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	}
	return req
}

func TestIntegration_FullFlow(t *testing.T) {
	r, _, _, _ := newIntegrationEnv(t)

	register := func(login, password string) (string, int64) {
		body := map[string]string{"login": login, "password": password}
		b, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, integrationRequest(http.MethodPost, "/auth/register", "", b, ""))
		require.Equal(t, http.StatusCreated, w.Code)
		var resp loginResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		return resp.Token, resp.Client.ID
	}

	token, clientID := register("integration@test.com", "pass123")

	bookSlot := func(slotID int64) int64 {
		body := map[string]any{"slotId": slotID, "rentalSelected": false}
		b, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, integrationRequest(http.MethodPost, "/bookings", "", b, token))
		require.Equal(t, http.StatusCreated, w.Code)
		var created struct {
			ID     int64 `json:"id"`
			Status string `json:"status"`
		}
		json.Unmarshal(w.Body.Bytes(), &created)
		require.Equal(t, "активна", created.Status)
		return created.ID
	}

	getBookings := func() int {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, integrationRequest(http.MethodGet, "/bookings", "", nil, token))
		require.Equal(t, http.StatusOK, w.Code)
		var bookings []struct {
			ID     int64 `json:"id"`
			Status string `json:"status"`
		}
		json.Unmarshal(w.Body.Bytes(), &bookings)
		return len(bookings)
	}

	cancelBooking := func(bookingID int64) {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, integrationRequest(http.MethodPatch, "/bookings/"+strconv.FormatInt(bookingID, 10)+"/cancel", strconv.FormatInt(bookingID, 10), nil, token))
		require.Equal(t, http.StatusOK, w.Code)
	}

	slots := func() {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, integrationRequest(http.MethodGet, "/slots", "", nil, ""))
		require.Equal(t, http.StatusOK, w.Code)
	}

	masters := func() {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, integrationRequest(http.MethodGet, "/masters", "", nil, ""))
		require.Equal(t, http.StatusOK, w.Code)
	}

	programs := func() {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, integrationRequest(http.MethodGet, "/programs", "", nil, ""))
		require.Equal(t, http.StatusOK, w.Code)
	}

	_ = clientID

	slots()
	masters()
	programs()

	bookingID1 := bookSlot(1)
	require.Equal(t, 1, getBookings())

	cancelBooking(bookingID1)

	bookingID2 := bookSlot(1)

	_ = bookingID2
	_ = chimw.RealIP
	_ = middleware.ClientIDFromContext
}
