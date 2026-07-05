package internal

import (
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/glini/backend/internal/handler"
	"github.com/glini/backend/internal/middleware"
)

type RouterDeps struct {
	AuthHandler    *handler.AuthHandler
	SlotHandler    *handler.SlotHandler
	MasterHandler  *handler.MasterHandler
	ProgramHandler *handler.ProgramHandler
	BookingHandler *handler.BookingHandler
	AuthSecret     string
}

func NewRouter(deps RouterDeps) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)

	authMw := middleware.Auth(deps.AuthSecret)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", deps.AuthHandler.Register)
		r.Post("/login", deps.AuthHandler.Login)
	})

	r.Route("/slots", func(r chi.Router) {
		r.Get("/", deps.SlotHandler.ListSlots)
		r.Get("/{id}", deps.SlotHandler.GetSlot)
	})

	r.Route("/bookings", func(r chi.Router) {
		r.Use(authMw)
		r.Post("/", deps.BookingHandler.Create)
		r.Get("/", deps.BookingHandler.ListMy)
		r.Get("/{id}", deps.BookingHandler.Get)
		r.Patch("/{id}/cancel", deps.BookingHandler.Cancel)
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
