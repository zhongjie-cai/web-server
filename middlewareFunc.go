package webserver

import (
	"net/http"

	"github.com/gorilla/mux"
)

// MiddlewareFunc warps around mux.MiddlewareFunc, which receives an http.Handler and returns another http.Handler.
// Typically, the returned handler is a closure which does something with the http.ResponseWriter and http.Request passed
// to it, and then calls the handler passed as parameter to the MiddlewareFunc.
type MiddlewareFunc mux.MiddlewareFunc

// addMiddleware wraps the mux middleware addition function
func addMiddleware(
	router *mux.Router,
	middleware MiddlewareFunc,
) {
	router.Use((func(http.Handler) http.Handler)(middleware))
}
