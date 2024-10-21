package webserver

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
)

func TestAddMiddleware(t *testing.T) {
	// arrange
	var dummyMiddleware = func(next http.Handler) http.Handler {
		return next
	}

	// SUT
	var router = mux.NewRouter()

	// act
	addMiddleware(
		router,
		dummyMiddleware,
	)
}
