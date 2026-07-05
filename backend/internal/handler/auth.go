package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/glini/backend/internal/domain"
	"github.com/glini/backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type loginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	Token  string        `json:"token"`
	Client domain.Client `json:"client"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}

	result, err := h.authService.Register(r.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, domain.ErrConflict) {
			writeError(w, http.StatusConflict, "логин уже занят")
			return
		}
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}

	writeJSON(w, http.StatusCreated, loginResponse{
		Token:  result.Token,
		Client: result.Client,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}

	result, err := h.authService.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			writeError(w, http.StatusUnauthorized, "неверный логин или пароль")
			return
		}
		if errors.Is(err, domain.ErrValidation) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{
		Token:  result.Token,
		Client: result.Client,
	})
}
