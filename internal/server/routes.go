package server

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/dmi3midd/simple-schedule/docs"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	s.activityHandler.RegisterRoutes(mux)
	s.slotHandler.RegisterRoutes(mux)

	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	return corsMiddleware(mux)
}
