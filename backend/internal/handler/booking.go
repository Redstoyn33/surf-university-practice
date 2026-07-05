package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/glini/backend/internal/domain"
	"github.com/glini/backend/internal/middleware"
	"github.com/glini/backend/internal/service"
)

type BookingHandler struct {
	bookingService *service.BookingService
}

func NewBookingHandler(bookingService *service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

type createBookingRequest struct {
	SlotID         int64 `json:"slotId"`
	RentalSelected bool  `json:"rentalSelected"`
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	clientID, ok := middleware.ClientIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	var req createBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}

	if req.SlotID <= 0 {
		writeError(w, http.StatusBadRequest, "некорректный slotId")
		return
	}

	booking, err := h.bookingService.CreateBooking(r.Context(), clientID, req.SlotID, req.RentalSelected)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "слот не найден")
			return
		}
		if errors.Is(err, domain.ErrConflict) {
			writeError(w, http.StatusConflict, "нет свободных мест или двойная бронь")
			return
		}
		if errors.Is(err, domain.ErrValidation) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}

	writeJSON(w, http.StatusCreated, booking)
}

func (h *BookingHandler) ListMy(w http.ResponseWriter, r *http.Request) {
	clientID, ok := middleware.ClientIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	statusFilter := r.URL.Query().Get("status")
	bookings, err := h.bookingService.GetMyBookings(r.Context(), clientID, statusFilter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}
	if bookings == nil {
		bookings = []domain.Booking{}
	}

	writeJSON(w, http.StatusOK, bookings)
}

func (h *BookingHandler) Get(w http.ResponseWriter, r *http.Request) {
	clientID, ok := middleware.ClientIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "некорректный id")
		return
	}

	booking, err := h.bookingService.GetBookingByID(r.Context(), id, clientID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "бронь не найдена")
			return
		}
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}

	writeJSON(w, http.StatusOK, booking)
}

func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	clientID, ok := middleware.ClientIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "некорректный id")
		return
	}

	booking, err := h.bookingService.CancelBooking(r.Context(), id, clientID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "бронь не найдена")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			writeError(w, http.StatusNotFound, "бронь не найдена")
			return
		}
		if errors.Is(err, domain.ErrNotActive) {
			writeError(w, http.StatusUnprocessableEntity, "бронь уже отменена")
			return
		}
		if errors.Is(err, domain.ErrValidation) {
			writeError(w, http.StatusUnprocessableEntity, "отмена невозможна — менее 4 часов до начала")
			return
		}
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}

	writeJSON(w, http.StatusOK, booking)
}
