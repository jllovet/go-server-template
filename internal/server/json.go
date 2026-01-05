package server

import (
	"encoding/json"
	"net/http"

	"github.com/jllovet/go-server-template/logger"
)

func (s *Server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		logger.FromContext(r.Context()).Error("json decoding failed", "error", err)
		return err
	}
	return nil
}

func (s *Server) encode(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		s.logger.Error("json encoding failed", "error", err)
	}
}
