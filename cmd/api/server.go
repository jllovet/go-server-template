package api

import (
	"log"
	"net/http"

	"github.com/jllovet/go-server-template/config"
	"github.com/jllovet/go-server-template/internal/todo"
)

type Server struct {
	service todo.Service
	config  *config.Config
	logger  *log.Logger
}

func NewServer(service todo.Service, config *config.Config, logger *log.Logger) *Server {
	return &Server{
		service: service,
		config:  config,
		logger:  logger,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.routes().ServeHTTP(w, r)
}
