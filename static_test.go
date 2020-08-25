package webserver

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestRegisterStatic(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyPath = "/foo/"
	var dummyHandler = &dummyHandler{t}

	// mock
	createMock(t)

	// SUT
	var router = mux.NewRouter()

	// act
	var route = registerStatic(
		router,
		dummyName,
		dummyPath,
		dummyHandler,
	)
	var name = route.GetName()
	var pathTemplate, _ = route.GetPathTemplate()
	var handler = route.GetHandler()

	// assert
	assert.Equal(t, dummyName, name)
	assert.Equal(t, dummyPath, pathTemplate)
	assert.Equal(t, dummyHandler, handler)

	// verify
	verifyAll(t)
}
