package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	s.activityHandler.RegisterRoutes(mux)
	s.slotHandler.RegisterRoutes(mux)
	return mux
}
