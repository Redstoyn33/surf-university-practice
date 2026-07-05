package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/glini/backend/internal/domain"
	"github.com/glini/backend/internal/repository"
)

type ProgramHandler struct {
	programRepo *repository.ProgramRepo
}

func NewProgramHandler(programRepo *repository.ProgramRepo) *ProgramHandler {
	return &ProgramHandler{programRepo: programRepo}
}

func (h *ProgramHandler) List(w http.ResponseWriter, r *http.Request) {
	programs, err := h.programRepo.QueryPrograms(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}
	if programs == nil {
		programs = []domain.Program{}
	}

	writeJSON(w, http.StatusOK, programs)
}

func (h *ProgramHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "некорректный id")
		return
	}

	program, err := h.programRepo.GetProgramByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "программа не найдена")
			return
		}
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}

	writeJSON(w, http.StatusOK, program)
}
