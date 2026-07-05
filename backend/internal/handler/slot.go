package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/glini/backend/internal/domain"
	"github.com/glini/backend/internal/repository"
)

type SlotHandler struct {
	slotRepo *repository.SlotRepo
}

func NewSlotHandler(slotRepo *repository.SlotRepo) *SlotHandler {
	return &SlotHandler{slotRepo: slotRepo}
}

func (h *SlotHandler) ListSlots(w http.ResponseWriter, r *http.Request) {
	filter := domain.SlotFilter{
		DateFrom: r.URL.Query().Get("dateFrom"),
		DateTo:   r.URL.Query().Get("dateTo"),
	}

	if masterIDStr := r.URL.Query().Get("masterId"); masterIDStr != "" {
		id, err := strconv.ParseInt(masterIDStr, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "некорректный masterId")
			return
		}
		filter.MasterID = &id
	}
	if programIDStr := r.URL.Query().Get("programId"); programIDStr != "" {
		id, err := strconv.ParseInt(programIDStr, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "некорректный programId")
			return
		}
		filter.ProgramID = &id
	}

	slots, err := h.slotRepo.QuerySlots(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}
	if slots == nil {
		slots = []domain.Slot{}
	}

	writeJSON(w, http.StatusOK, slots)
}

func (h *SlotHandler) GetSlot(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "некорректный id")
		return
	}

	slot, err := h.slotRepo.GetSlotByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "слот не найден")
			return
		}
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}

	writeJSON(w, http.StatusOK, slot)
}
