package webserver

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Static holds the registration information of a static content hosting
type Static struct {
	Name       string
	PathPrefix string
	Handler    http.Handler
}

// registerStatic wraps the mux static content handler
func registerStatic(
	router *mux.Router,
	name string,
	path string,
	handler http.Handler,
) *mux.Route {
	return router.PathPrefix(
		path,
	).Handler(
		handler,
	).Name(
		name,
	)
}
