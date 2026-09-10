package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dmi3midd/simple-schedule/internal/config"
	"github.com/dmi3midd/simple-schedule/internal/postgres"
	"github.com/dmi3midd/simple-schedule/internal/server"
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

	// Postgres
	pg, err := postgres.New(&cfg.Postgres)
	if err != nil {
		slog.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer pg.Close()

	// Create and start server
	server := server.NewServer(
		&cfg.Server,
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
