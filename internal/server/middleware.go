package server

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"github.com/jllovet/go-server-template/logger"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Ensure Request ID exists
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			b := make([]byte, 16)
			_, _ = rand.Read(b)
			reqID = hex.EncodeToString(b)
		}
		w.Header().Set("X-Request-ID", reqID)

		// Create a logger with request-specific attributes
		log := s.logger.With(
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("request_id", reqID),
		)

		// Inject logger into context
		ctx := logger.WithContext(r.Context(), log)

		// Wrap ResponseWriter to capture status code
		ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(ww, r.WithContext(ctx))

		log.Info("req completed",
			slog.Int("status", ww.status),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}
