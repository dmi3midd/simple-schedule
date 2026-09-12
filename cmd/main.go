package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dmi3midd/shkvcache"
	"github.com/dmi3midd/simple-schedule/internal/config"
	"github.com/dmi3midd/simple-schedule/internal/domain"
	"github.com/dmi3midd/simple-schedule/internal/logger"
	"github.com/dmi3midd/simple-schedule/internal/postgres"
	"github.com/dmi3midd/simple-schedule/internal/repository"
	"github.com/dmi3midd/simple-schedule/internal/server"
	"github.com/dmi3midd/simple-schedule/internal/server/handlers"
	"github.com/dmi3midd/simple-schedule/internal/service"
	"github.com/go-playground/validator/v10"
)

// @title           SimpleSchedule API
// @version         1.0
// @description     SimpleSchedule service API.
// @host            localhost:2811
// @BasePath        /
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

	// Initialize logger with configured level
	logger.Setup(cfg.Log.Level)

	// Postgres
	pg, err := postgres.New(&cfg.Postgres)
	if err != nil {
		slog.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer pg.Close()

	db := pg.GetDB()

	// Repositories
	tagRepo := repository.NewTagRepository(db)
	weekRepo := repository.NewWeekRepository(db)
	slotRepo := repository.NewSlotRepository(db)

	// Cache
	scheduleCache, err := shkvcache.NewCache[domain.WeekSchedule](ctx, shkvcache.DefaultOptions())
	if err != nil {
		slog.Error("failed to create cache", slog.Any("error", err))
		os.Exit(1)
	}
	defer scheduleCache.Close()

	// Services
	tagService := service.NewTagService(tagRepo)
	weekService := service.NewWeekService(weekRepo, slotRepo, *scheduleCache)
	slotService := service.NewSlotService(slotRepo, weekRepo, tagRepo)

	// Validator
	val := validator.New()

	// Handlers
	tagHandler := handlers.NewTagHandler(tagService, val)
	weekHandler := handlers.NewWeekHandler(weekService, val)
	slotHandler := handlers.NewSlotHandler(slotService, val)

	// Create and start server
	srv := server.NewServer(
		&cfg.Server,
		tagHandler,
		weekHandler,
		slotHandler,
	)
	slog.Info(
		"server is running",
		slog.String("address", cfg.Server.Address),
	)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to run server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	<-ctx.Done()
	slog.Info("received shutdown signal, stopping application...")
}
