package server

import (
	"io/fs"
	"net/http"
	"strings"

	_ "github.com/dmi3midd/simple-schedule/docs"
	"github.com/dmi3midd/simple-schedule/internal/shared/httputils/apierror"
	webui "github.com/dmi3midd/simple-schedule/web"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Swagger UI
	mux.HandleFunc("GET /swagger/", httpSwagger.Handler())

	// Tag routes
	mux.HandleFunc("GET /api/tags", apierror.ErrorHandler(s.tagHandler.GetAll))
	mux.HandleFunc("GET /api/tags/{id}", apierror.ErrorHandler(s.tagHandler.GetById))
	mux.HandleFunc("POST /api/tags", apierror.ErrorHandler(s.tagHandler.Create))
	mux.HandleFunc("PUT /api/tags/{id}", apierror.ErrorHandler(s.tagHandler.Update))
	mux.HandleFunc("DELETE /api/tags/{id}", apierror.ErrorHandler(s.tagHandler.Delete))

	// Week routes
	mux.HandleFunc("GET /api/weeks", apierror.ErrorHandler(s.weekHandler.GetAll))
	mux.HandleFunc("GET /api/weeks/{id}", apierror.ErrorHandler(s.weekHandler.GetById))
	mux.HandleFunc("GET /api/weeks/{id}/schedule", apierror.ErrorHandler(s.weekHandler.GetSchedule))
	mux.HandleFunc("POST /api/weeks", apierror.ErrorHandler(s.weekHandler.Create))
	mux.HandleFunc("PUT /api/weeks/{id}", apierror.ErrorHandler(s.weekHandler.Update))
	mux.HandleFunc("DELETE /api/weeks/{id}", apierror.ErrorHandler(s.weekHandler.Delete))

	// Slot routes
	mux.HandleFunc("GET /api/slots/{id}", apierror.ErrorHandler(s.slotHandler.GetById))
	mux.HandleFunc("GET /api/weeks/{id}/slots", apierror.ErrorHandler(s.slotHandler.GetByWeekId))
	mux.HandleFunc("POST /api/slots", apierror.ErrorHandler(s.slotHandler.Create))
	mux.HandleFunc("PUT /api/slots/{id}", apierror.ErrorHandler(s.slotHandler.Update))
	mux.HandleFunc("DELETE /api/slots/{id}", apierror.ErrorHandler(s.slotHandler.Delete))

	// Web UI — embedded SPA
	distFS, _ := fs.Sub(webui.StaticFS, "dist")
	staticHandler := http.FileServer(http.FS(distFS))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Try to serve the static file directly (JS, CSS, images, etc.)
		if path != "/" && !strings.HasPrefix(path, "/api") && !strings.HasPrefix(path, "/swagger") {
			if f, err := fs.Stat(distFS, strings.TrimPrefix(path, "/")); err == nil && !f.IsDir() {
				staticHandler.ServeHTTP(w, r)
				return
			}
		}

		// Fallback: serve index.html for SPA client-side routing
		r.URL.Path = "/"
		staticHandler.ServeHTTP(w, r)
	})

	return loggingMiddleware(corsMiddleware(mux))
}

