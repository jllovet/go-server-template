package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jllovet/go-server-template/config"
	"github.com/jllovet/go-server-template/encoding"
	"github.com/jllovet/go-server-template/services/healthz"
)

func addRoutes(
	mux *http.ServeMux,
	cfg *config.Config,
	logger *log.Logger,
) {
	// Health endpoints
	mux.Handle("GET /healthz", healthz.HandleHealthCheck(logger))
	mux.Handle("GET /ready", HandleReady(logger))

	// Versioned API example
	mux.Handle("GET /api/v0/hello", HandleHello(cfg))

	// Default 404
	mux.Handle("/", http.NotFoundHandler())
}

func HandleHello(cfg *config.Config) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			encoding.Encode(w, http.StatusOK, map[string]string{
				"message": fmt.Sprintf("Hello from %s:%s", cfg.Host, cfg.Port),
			})
		},
	)
}

// Returns whether the server is ready
func HandleReady(logger *log.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Check dependencies here. If all ready return {"ready": true}
			encoding.Encode(w, http.StatusOK, map[string]any{"ready": true})
		},
	)
}
