package server

import (
	"fmt"
	"net/http"

	"github.com/jllovet/go-server-template/internal/healthz"
)

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()

	// Health endpoints
	mux.Handle("GET /healthz", healthz.HandleHealthCheck(s.logger))
	mux.Handle("GET /ready", s.handleReady())

	// Versioned API example
	mux.Handle("GET /api/v0/hello", s.handleHello())

	// Todo endpoints
	mux.HandleFunc("POST /api/v1/todos", s.handleCreateTodo())
	mux.HandleFunc("GET /api/v1/todos", s.handleListTodos())
	mux.HandleFunc("GET /api/v1/todos/{id}", s.handleGetTodo())
	mux.HandleFunc("PATCH /api/v1/todos/{id}", s.handleUpdateTodoTitle())
	mux.HandleFunc("POST /api/v1/todos/{id}/complete", s.handleMarkTodoComplete())
	mux.HandleFunc("POST /api/v1/todos/{id}/incomplete", s.handleMarkTodoIncomplete())
	mux.HandleFunc("DELETE /api/v1/todos/{id}", s.handleDeleteTodo())

	// Default 404
	mux.Handle("/", http.NotFoundHandler())

	return s.loggingMiddleware(mux)
}

func (s *Server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.encode(w, http.StatusOK, map[string]string{
			"message": fmt.Sprintf("Hello from %s:%s", s.config.Host, s.config.Port),
		})
	}
}

// Returns whether the server is ready
func (s *Server) handleReady() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check dependencies here. If all ready return {"ready": true}
		s.encode(w, http.StatusOK, map[string]any{"ready": true})
	}
}
