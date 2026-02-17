package webserver

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/zhongjie-cai/gomocker/v2"
)

func TestRegisterStatic(t *testing.T) {
	// arrange
	var dummyPath = "/foo/"
	var dummyHandler = &dummyHandler{}
	type router struct {
		chi.Router
	}
	var dummyRouter = &router{}

	// mock
	var mock = gomocker.NewMocker(t)

	// expect
	mock.Mock((*router).Route).Expects(dummyRouter, dummyPath, dummyHandler).Returns().Once()

	// SUT + act
	registerStatic(
		dummyRouter,
		dummyPath,
		dummyHandler,
	)
}
