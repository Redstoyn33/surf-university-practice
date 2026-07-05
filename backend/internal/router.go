package internal

import (
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/glini/backend/internal/handler"
)

type RouterDeps struct {
	AuthHandler    *handler.AuthHandler
	SlotHandler    *handler.SlotHandler
	MasterHandler  *handler.MasterHandler
	ProgramHandler *handler.ProgramHandler
}

func NewRouter(deps RouterDeps) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", deps.AuthHandler.Register)
		r.Post("/login", deps.AuthHandler.Login)
	})

	r.Route("/slots", func(r chi.Router) {
		r.Get("/", deps.SlotHandler.ListSlots)
		r.Get("/{id}", deps.SlotHandler.GetSlot)
	})

	r.Route("/bookings", func(r chi.Router) {
	})

	r.Route("/ratings", func(r chi.Router) {
	})

	r.Route("/masters", func(r chi.Router) {
		r.Get("/", deps.MasterHandler.List)
		r.Get("/{id}", deps.MasterHandler.Get)
	})

	r.Route("/programs", func(r chi.Router) {
		r.Get("/", deps.ProgramHandler.List)
		r.Get("/{id}", deps.ProgramHandler.Get)
	})

	return r
}
