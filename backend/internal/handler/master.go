package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/glini/backend/internal/domain"
	"github.com/glini/backend/internal/repository"
)

type MasterHandler struct {
	masterRepo *repository.MasterRepo
}

func NewMasterHandler(masterRepo *repository.MasterRepo) *MasterHandler {
	return &MasterHandler{masterRepo: masterRepo}
}

func (h *MasterHandler) List(w http.ResponseWriter, r *http.Request) {
	masters, err := h.masterRepo.QueryMasters(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}
	if masters == nil {
		masters = []domain.Master{}
	}

	writeJSON(w, http.StatusOK, masters)
}

func (h *MasterHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "некорректный id")
		return
	}

	master, err := h.masterRepo.GetMasterByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "мастер не найден")
			return
		}
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}

	writeJSON(w, http.StatusOK, master)
}
