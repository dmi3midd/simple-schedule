package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"

	"github.com/dmi3midd/simple-schedule/internal/config"
	"github.com/dmi3midd/simple-schedule/internal/server"
	"github.com/dmi3midd/simple-schedule/internal/server/handlers"
	"github.com/dmi3midd/simple-schedule/internal/service"
)

func main() {
	// Root context with signal cancellation for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	// Initialize services
	activityService := service.NewActivityService()
	slotService := service.NewSlotService(activityService)
	activityService.SetSlotCascadeDeleter(slotService)

	// Initialize validator and handlers
	validate := validator.New()
	activityHandler := handlers.NewActivityHandler(activityService, validate)
	slotHandler := handlers.NewSlotHandler(slotService, validate)

	// Create and start server
	server := server.NewServer(
		&cfg.Server,
		activityHandler,
		slotHandler,
	)
	slog.Info(
		"server is running",
		slog.String("address", cfg.Server.Address),
	)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to run server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	<-ctx.Done()
	slog.Info("received shutdown signal, stopping application...")
}
