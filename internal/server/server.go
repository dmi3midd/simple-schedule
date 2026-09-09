package server

import (
	"net/http"

	"github.com/dmi3midd/simple-schedule/internal/config"
	"github.com/dmi3midd/simple-schedule/internal/server/handlers"
)

type Server struct {
	cfg             *config.ServerConfig
	activityHandler *handlers.ActivityHandler
	slotHandler     *handlers.SlotHandler
}

func NewServer(
	cfg *config.ServerConfig,
	activityHandler *handlers.ActivityHandler,
	slotHandler *handlers.SlotHandler,
) *http.Server {
	s := &Server{
		cfg:             cfg,
		activityHandler: activityHandler,
		slotHandler:     slotHandler,
	}
	router := s.RegisterRoutes()
	return &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}
