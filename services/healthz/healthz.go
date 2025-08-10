package healthz

import (
	"log"
	"net/http"

	"github.com/jllovet/go-server-template/encoding"
)

// Returns status of server
func HandleHealthCheck(logger *log.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			encoding.Encode(w, http.StatusOK, map[string]any{"status": "ok"})
		},
	)
}
