package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/glini/backend/internal/domain"
	"github.com/glini/backend/internal/middleware"
	"github.com/glini/backend/internal/service"
)

type RatingHandler struct {
	ratingService *service.RatingService
}

func NewRatingHandler(ratingService *service.RatingService) *RatingHandler {
	return &RatingHandler{ratingService: ratingService}
}

type createRatingRequest struct {
	MasterID int64 `json:"masterId"`
	SlotID   int64 `json:"slotId"`
	Score    int   `json:"score"`
}

func (h *RatingHandler) Create(w http.ResponseWriter, r *http.Request) {
	clientID, ok := middleware.ClientIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	var req createRatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}

	if req.MasterID <= 0 || req.SlotID <= 0 || req.Score < 1 || req.Score > 5 {
		writeError(w, http.StatusBadRequest, "некорректные данные")
		return
	}

	rating, err := h.ratingService.CreateRating(r.Context(), clientID, req.MasterID, req.SlotID, req.Score)
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		if errors.Is(err, domain.ErrConflict) {
			writeError(w, http.StatusConflict, "оценка уже оставлена")
			return
		}
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}

	writeJSON(w, http.StatusCreated, rating)
}
