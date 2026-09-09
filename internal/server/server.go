package server

import (
	"net/http"

	"github.com/dmi3midd/simple-schedule/internal/config"
)

type Server struct {
	cfg *config.ServerConfig
}

func NewServer(
	cfg *config.ServerConfig,
) *http.Server {
	s := &Server{
		cfg: cfg,
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
