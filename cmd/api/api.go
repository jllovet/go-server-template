package api

import (
	"log"
	"net/http"

	"github.com/jllovet/go-server-template/config"
)

func NewServer(
	config *config.Config,
	logger *log.Logger,
	// myStore *myStore,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		config,
		logger,
		// myStore,
	)
	var handler http.Handler = mux
	// handler = someMiddleware(handler)
	// handler = someMiddleware2(handler)
	// handler = someMiddleware3(handler)
	return handler
}
