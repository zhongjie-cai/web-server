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

	// mock
	createMock(t)

	// SUT
	var router = mux.NewRouter()

	// act
	addMiddleware(
		router,
		dummyMiddleware,
	)

	// verify
	verifyAll(t)
}
