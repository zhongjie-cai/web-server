package webserver

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zhongjie-cai/gomocker/v2"
)

func TestEvaluateRoute_NilHandler(t *testing.T) {
	// arrange
	var dummyMethod = "some_method"
	var dummyRoute = "some_route"
	var dummyHandler http.Handler
	var dummyMiddlewares []func(http.Handler) http.Handler
	var dummyErrorMessage = "Invalid handler for some_method:some_route"

	// SUT + act
	var err = evaluateRoute(
		dummyMethod,
		dummyRoute,
		dummyHandler,
		dummyMiddlewares...,
	)

	// assert
	assert.Error(t, err, dummyErrorMessage)
}

func TestEvaluateRoute_NilMiddleware(t *testing.T) {
	// arrange
	var dummyMethod = "some_method"
	var dummyRoute = "some_route"
	type handler struct {
		http.Handler
	}
	var dummyHandler = &handler{}
	var dummyMiddlewares = []func(http.Handler) http.Handler{
		func(h http.Handler) http.Handler { return h },
		nil,
		func(h http.Handler) http.Handler { return h },
	}
	var dummyErrorMessage = "Invalid middleware for some_method:some_route @ #2"

	// SUT + act
	var err = evaluateRoute(
		dummyMethod,
		dummyRoute,
		dummyHandler,
		dummyMiddlewares...,
	)

	// assert
	assert.Error(t, err, dummyErrorMessage)
}

func TestEvaluateRoute_Success(t *testing.T) {
	// arrange
	var dummyMethod = "some_method"
	var dummyRoute = "some_route"
	type handler struct {
		http.Handler
	}
	var dummyHandler = &handler{}
	var dummyMiddlewares = []func(http.Handler) http.Handler{
		func(h http.Handler) http.Handler { return h },
		func(h http.Handler) http.Handler { return h },
		func(h http.Handler) http.Handler { return h },
	}

	// SUT + act
	var err = evaluateRoute(
		dummyMethod,
		dummyRoute,
		dummyHandler,
		dummyMiddlewares...,
	)

	// assert
	assert.NoError(t, err)
}

func TestWalkRegisteredRoutes_Error(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	type router struct {
		chi.Router
	}
	var dummyRouter = &router{}
	var dummyError = errors.New("some error")
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(chi.Walk).Expects(dummyRouter, gomocker.Matches(func(value any) bool {
		return functionPointerEquals(evaluateRoute, value)
	})).Returns(dummyError).Once()
	m.Mock(logAppRoot).Expects(dummySession, "route", "walkRegisteredRoutes", "Failure: %+v", dummyError).Returns().Once()
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageRouteRegistration, dummyError).Returns(dummyAppError).Once()

	// SUT + act
	var err = walkRegisteredRoutes(
		dummySession,
		dummyRouter,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
}

func TestWalkRegisteredRoutes_Success(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	type router struct {
		chi.Router
	}
	var dummyRouter = &router{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(chi.Walk).Expects(dummyRouter, gomocker.Matches(func(value any) bool {
		return functionPointerEquals(evaluateRoute, value)
	})).Returns(nil).Once()

	// SUT + act
	var err = walkRegisteredRoutes(
		dummySession,
		dummyRouter,
	)

	// assert
	assert.NoError(t, err)
}

func TestGenerateRouteName(t *testing.T) {
	// arrange
	var dummyMethod = "some_method"
	var dummyPattern = "some_pattern"

	// SUT + act
	var result = generateRouteName(dummyMethod, dummyPattern)

	// assert
	assert.Equal(t, "some_method:some_pattern", result)
}

func TestGetRouteInfo_NilRoute(t *testing.T) {
	// arrange
	var dummyHTTPRequest = &http.Request{
		Method:     http.MethodGet,
		RequestURI: "http://localhost/",
		Header:     map[string][]string{},
	}
	var dummyActionFuncMap = map[string]ActionFunc{}
	var dummyCtx *chi.Context
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(chi.RouteContext).Expects(dummyHTTPRequest.Context()).Returns(dummyCtx).Once()
	m.Mock(newAppError).Expects(errorCodeDataCorruption, "No go-chi context found in HTTP request").Returns(dummyAppError).Once()

	// SUT + act
	var name, action, err = getRouteInfo(
		dummyHTTPRequest,
		dummyActionFuncMap,
	)

	// assert
	assert.Zero(t, name)
	assert.Nil(t, action)
	assert.Equal(t, dummyAppError, err)
}

func TestGetRouteInfo_RouteNotFound(t *testing.T) {
	// arrange
	var dummyHTTPRequest = &http.Request{
		Method:     http.MethodGet,
		RequestURI: "http://localhost/",
		Header:     map[string][]string{},
	}
	var dummyPattern1 = "some pattern 1"
	var dummyPattern2 = "some pattern 2"
	var dummyCtx = &chi.Context{
		RouteMethod:   dummyHTTPRequest.Method,
		RoutePatterns: []string{dummyPattern1, dummyPattern2},
	}
	var dummyName1 = "some name 1"
	var dummyName2 = "some name 2"
	var dummyActionFuncMap = map[string]ActionFunc{}
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(chi.RouteContext).Expects(dummyHTTPRequest.Context()).Returns(dummyCtx).Once()
	m.Mock(generateRouteName).Expects(dummyHTTPRequest.Method, dummyPattern1).Returns(dummyName1).Once()
	m.Mock(generateRouteName).Expects(dummyHTTPRequest.Method, dummyPattern2).Returns(dummyName2).Once()
	m.Mock(newAppError).Expects(errorCodeNotFound, "No corresponding route configured for path: http://localhost/").Returns(dummyAppError).Once()

	// SUT + act
	var endpoint, action, err = getRouteInfo(
		dummyHTTPRequest,
		dummyActionFuncMap,
	)

	// assert
	assert.Equal(t, "", endpoint)
	assert.Nil(t, action)
	assert.Equal(t, dummyAppError, err)
}

func TestGetRouteInfo_ValidRoute(t *testing.T) {
	// arrange
	var dummyHTTPRequest = &http.Request{
		Method:     http.MethodGet,
		RequestURI: "http://localhost/",
		Header:     map[string][]string{},
	}
	var dummyPattern1 = "some pattern 1"
	var dummyPattern2 = "some pattern 2"
	var dummyCtx = &chi.Context{
		RouteMethod:   dummyHTTPRequest.Method,
		RoutePatterns: []string{dummyPattern1, dummyPattern2},
	}
	var dummyName1 = "some name 1"
	var dummyName2 = "some name 2"
	var dummyAction = func(Session) (any, error) { return nil, nil }
	var dummyActionFuncMap = map[string]ActionFunc{
		dummyName2: dummyAction,
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(chi.RouteContext).Expects(dummyHTTPRequest.Context()).Returns(dummyCtx).Once()
	m.Mock(generateRouteName).Expects(dummyHTTPRequest.Method, dummyPattern1).Returns(dummyName1).Once()
	m.Mock(generateRouteName).Expects(dummyHTTPRequest.Method, dummyPattern2).Returns(dummyName2).Once()

	// SUT + act
	var endpoint, action, err = getRouteInfo(
		dummyHTTPRequest,
		dummyActionFuncMap,
	)

	// assert
	assert.Equal(t, dummyName2, endpoint)
	assertFunctionEquals(t, dummyAction, action)
	assert.NoError(t, err)
}
