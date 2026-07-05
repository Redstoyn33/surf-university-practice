package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/glini/backend/internal"
	"github.com/glini/backend/internal/config"
	"github.com/glini/backend/internal/handler"
	"github.com/glini/backend/internal/repository"
	"github.com/glini/backend/internal/service"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
	slog.Info("starting glini backend")

	cfg := config.Load()

	db, err := repository.NewDB(cfg.DBPath)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	clientRepo := repository.NewClientRepo(db)
	masterRepo := repository.NewMasterRepo(db)
	programRepo := repository.NewProgramRepo(db)
	slotRepo := repository.NewSlotRepo(db, masterRepo, programRepo)

	authService := service.NewAuthService(clientRepo, cfg.AuthSecret)
	bookingRepo := repository.NewBookingRepo(db, slotRepo, masterRepo, programRepo)
	bookingService := service.NewBookingService(bookingRepo, slotRepo)
	ratingRepo := repository.NewRatingRepo(db)
	ratingService := service.NewRatingService(ratingRepo, bookingRepo, slotRepo)

	authHandler := handler.NewAuthHandler(authService)
	slotHandler := handler.NewSlotHandler(slotRepo)
	masterHandler := handler.NewMasterHandler(masterRepo)
	programHandler := handler.NewProgramHandler(programRepo)
	bookingHandler := handler.NewBookingHandler(bookingService)
	ratingHandler := handler.NewRatingHandler(ratingService)

	router := internal.NewRouter(internal.RouterDeps{
		AuthHandler:    authHandler,
		SlotHandler:    slotHandler,
		MasterHandler:  masterHandler,
		ProgramHandler: programHandler,
		BookingHandler: bookingHandler,
		RatingHandler:  ratingHandler,
		AuthSecret:     cfg.AuthSecret,
	})

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("server listening", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
	slog.Info("server stopped")
}
