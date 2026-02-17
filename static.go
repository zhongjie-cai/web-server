package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Static holds the registration information of a static content hosting
type Static struct {
	PathPrefix string
	Handler    http.Handler
}

// registerStatic wraps the go-chi static content handler
func registerStatic(
	router chi.Router,
	path string,
	handler http.Handler,
) {
	router.Handle(
		path,
		handler,
	)
}
