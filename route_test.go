package webserver

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/zhongjie-cai/gomocker"
)

func TestGetName_Undefined(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// SUT
	var dummyRoute = dummyRouter.NewRoute()

	// act
	var result = getName(
		dummyRoute,
	)

	// assert
	assert.Zero(t, result)
}

func TestGetName_Defined(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// SUT
	var dummyRoute = dummyRouter.NewRoute().Name(
		"test",
	)

	// act
	var result = getName(
		dummyRoute,
	)

	// assert
	assert.Equal(t, "test", result)
}

func TestGetPathTemplate_Error(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// SUT
	var dummyRoute = dummyRouter.NewRoute()

	// act
	var result, err = getPathTemplate(
		dummyRoute,
	)

	// assert
	assert.Zero(t, result)
	assert.Equal(t, "mux: route doesn't have a path", err.Error())
}

func TestGetPathTemplate_Success(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// SUT
	var dummyRoute = dummyRouter.NewRoute().Path(
		"/foo/{bar}",
	)

	// act
	var result, err = getPathTemplate(
		dummyRoute,
	)

	// assert
	assert.Equal(t, "/foo/{bar}", result)
	assert.NoError(t, err)
}

func TestGetPathRegexp_Error(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// SUT
	var dummyRoute = dummyRouter.NewRoute()

	// act
	var result, err = getPathRegexp(
		dummyRoute,
	)

	// assert
	assert.Zero(t, result)
	assert.Equal(t, "mux: route does not have a path", err.Error())
}

func TestGetPathRegexp_Success(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// SUT
	var dummyRoute = dummyRouter.NewRoute().Path(
		"/foo/{bar}",
	)

	// act
	var result, err = getPathRegexp(
		dummyRoute,
	)

	// assert
	assert.Equal(t, "^/foo/(?P<v0>[^/]+)$", result)
	assert.NoError(t, err)
}

func TestGetQueriesTemplate_Undefined(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// SUT
	var dummyRoute = dummyRouter.NewRoute()

	// act
	var result = getQueriesTemplates(
		dummyRoute,
	)

	// assert
	assert.Zero(t, result)
}

func TestGetQueriesTemplate_Defined(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// SUT
	var dummyRoute = dummyRouter.NewRoute().Queries(
		"abc",
		"{def}",
		"xyz",
		"{zyx}",
	)

	// act
	var result = getQueriesTemplates(
		dummyRoute,
	)

	// assert
	assert.Equal(t, "abc={def}|xyz={zyx}", result)
}

func TestGetQueriesRegexp_Undefined(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// SUT
	var dummyRoute = dummyRouter.NewRoute()

	// act
	var result = getQueriesRegexp(
		dummyRoute,
	)

	// assert
	assert.Zero(t, result)
}

func TestGetQueriesRegexp_Defined(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// SUT
	var dummyRoute = dummyRouter.NewRoute().Queries(
		"abc",
		"{def}",
		"xyz",
		"{zyx}",
	)

	// act
	var result = getQueriesRegexp(
		dummyRoute,
	)

	// assert
	assert.Equal(t, "^abc=(?P<v0>.*)$|^xyz=(?P<v0>.*)$", result)
}

func TestGetMethods_Undefined(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// SUT
	var dummyRoute = dummyRouter.NewRoute()

	// act
	var result = getMethods(
		dummyRoute,
	)

	// assert
	assert.Zero(t, result)
}

func TestGetMethods_Defined(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// SUT
	var dummyRoute = dummyRouter.NewRoute().Methods(
		"GET",
		"PUT",
	)

	// act
	var result = getMethods(
		dummyRoute,
	)

	// assert
	assert.Equal(t, "GET|PUT", result)
}

func TestPrintRegisteredRouteDetails_TemplateError(t *testing.T) {
	// arrange
	var dummyRoute = &mux.Route{}
	var dummyRouter = &mux.Router{}
	var dummyAncestors = []*mux.Route{}
	var dummyPathTemplate string
	var dummyPathTemplateError = errors.New("some path template error")
	var dummyPathRegexp string
	var dummyPathRegexpError = errors.New("some path regexp error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(getPathTemplate, 1, func(route *mux.Route) (string, error) {
		assert.Equal(t, dummyRoute, route)
		return dummyPathTemplate, dummyPathTemplateError
	})
	m.ExpectFunc(getPathRegexp, 1, func(route *mux.Route) (string, error) {
		assert.Equal(t, dummyRoute, route)
		return dummyPathRegexp, dummyPathRegexpError
	})

	// SUT + act
	var err = evaluateRoute(
		dummyRoute,
		dummyRouter,
		dummyAncestors,
	)

	// assert
	assert.Equal(t, dummyPathTemplateError, err)
}

func TestPrintRegisteredRouteDetails_RegexpError(t *testing.T) {
	// arrange
	var dummyRoute = &mux.Route{}
	var dummyRouter = &mux.Router{}
	var dummyAncestors = []*mux.Route{}
	var dummyPathTemplate string
	var dummyPathRegexp string
	var dummyPathRegexpError = errors.New("some path regexp error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(getPathTemplate, 1, func(route *mux.Route) (string, error) {
		assert.Equal(t, dummyRoute, route)
		return dummyPathTemplate, nil
	})
	m.ExpectFunc(getPathRegexp, 1, func(route *mux.Route) (string, error) {
		assert.Equal(t, dummyRoute, route)
		return dummyPathRegexp, dummyPathRegexpError
	})

	// SUT + act
	var err = evaluateRoute(
		dummyRoute,
		dummyRouter,
		dummyAncestors,
	)

	// assert
	assert.Equal(t, dummyPathRegexpError, err)
}

func TestPrintRegisteredRouteDetails_Success(t *testing.T) {
	// arrange
	var dummyRoute = &mux.Route{}
	var dummyRouter = &mux.Router{}
	var dummyAncestors = []*mux.Route{}
	var dummyPathTemplate = "some path template"
	var dummyPathRegexp = "some path regexp"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(getPathTemplate, 1, func(route *mux.Route) (string, error) {
		assert.Equal(t, dummyRoute, route)
		return dummyPathTemplate, nil
	})
	m.ExpectFunc(getPathRegexp, 1, func(route *mux.Route) (string, error) {
		assert.Equal(t, dummyRoute, route)
		return dummyPathRegexp, nil
	})

	// SUT + act
	var err = evaluateRoute(
		dummyRoute,
		dummyRouter,
		dummyAncestors,
	)

	// assert
	assert.NoError(t, err)
}

func TestWalkRegisteredRoutes_Error(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyRouter = &mux.Router{}
	var dummyError = errors.New("some error")
	var dummyAppError = &appError{Message: "some error message"}

	// stub
	dummyRouter.HandleFunc("/", func(http.ResponseWriter, *http.Request) {})

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(evaluateRoute, 1, func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		return dummyError
	})
	m.ExpectFunc(logAppRoot, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "route", category)
		assert.Equal(t, "walkRegisteredRoutes", subcategory)
		assert.Equal(t, "Failure: %+v", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyError, parameters[0])
	})
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageRouteRegistration, errorMessage)
		assert.Equal(t, 1, len(innerErrors))
		assert.Equal(t, dummyError, innerErrors[0])
		return dummyAppError
	})

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
	var dummyRouter = &mux.Router{}

	// stub
	dummyRouter.HandleFunc("/", func(http.ResponseWriter, *http.Request) {})

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(evaluateRoute, 1, func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		return nil
	})

	// SUT + act
	var err = walkRegisteredRoutes(
		dummySession,
		dummyRouter,
	)

	// assert
	assert.NoError(t, err)
}

func TestRegisterRoute(t *testing.T) {
	// arrange
	var dummyEndpoint = "some endpoint"
	var dummyMethod = "SOME METHOD"
	var dummyName = "some endpoint:SOME METHOD"
	var dummyPath = "/foo/{bar}"
	var dummyQueries = []string{"test", "{test}"}
	var dummyRouter = &mux.Router{}
	var dummyRoute = &mux.Route{}

	// stub
	var dummyHandlerFuncExpected = 0
	var dummyHandlerFuncCalled = 0
	var dummyHandlerFunc = func(http.ResponseWriter, *http.Request) {
		dummyHandlerFuncCalled++
	}
	var dummyActionFuncExpected = 0
	var dummyActionFuncCalled = 0
	var dummyActionFunc = func(Session) (interface{}, error) {
		dummyActionFuncCalled++
		return nil, nil
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyRouter, "HandleFunc", 1, func(self *mux.Router, path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
		assert.Equal(t, dummyRouter, self)
		assert.Equal(t, dummyPath, path)
		functionPointerEquals(t, dummyHandlerFunc, f)
		return dummyRoute
	})
	m.ExpectMethod(dummyRoute, "Methods", 1, func(self *mux.Route, methods ...string) *mux.Route {
		assert.Equal(t, dummyRoute, self)
		assert.Equal(t, 1, len(methods))
		assert.Equal(t, dummyMethod, methods[0])
		return dummyRoute
	})
	m.ExpectMethod(dummyRoute, "Queries", 1, func(self *mux.Route, pairs ...string) *mux.Route {
		assert.Equal(t, dummyRoute, self)
		assert.Equal(t, dummyQueries, pairs)
		return dummyRoute
	})
	m.ExpectMethod(dummyRoute, "Name", 1, func(self *mux.Route, name string) *mux.Route {
		assert.Equal(t, dummyRoute, self)
		assert.Equal(t, dummyName, name)
		return dummyRoute
	})

	// SUT + act
	var name, route = registerRoute(
		dummyRouter,
		dummyEndpoint,
		dummyMethod,
		dummyPath,
		dummyQueries,
		dummyHandlerFunc,
		dummyActionFunc,
	)

	// assert
	assert.Equal(t, dummyName, name)
	assert.Equal(t, dummyRoute, route)
	assert.Equal(t, dummyHandlerFuncExpected, dummyHandlerFuncCalled)
	assert.Equal(t, dummyActionFuncExpected, dummyActionFuncCalled)
}

func TestDefaultActionFunc(t *testing.T) {
	// arrange
	var dummySession = &session{}
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeNotImplemented, errorCode)
		assert.Equal(t, "No corresponding action function configured; falling back to default", errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	})

	// SUT + act
	var result, err = defaultActionFunc(
		dummySession,
	)

	// assert
	assert.Nil(t, result)
	assert.Equal(t, dummyAppError, err)
}

func TestGetEndpointByName_NoSeparator(t *testing.T) {
	// arrange
	var dummyName = "some name"

	// SUT + act
	result := getEndpointByName(
		dummyName,
	)

	// assert
	assert.Equal(t, dummyName, result)
}

func TestGetEndpointByName_WithSeparator(t *testing.T) {
	// arrange
	var dummyEndpoint = "some endpoint"
	var dummyName = dummyEndpoint + ":some name"

	// SUT + act
	result := getEndpointByName(
		dummyName,
	)

	// assert
	assert.Equal(t, dummyEndpoint, result)
}

func TestGetRouteInfo_NilRoute(t *testing.T) {
	// arrange
	var dummyHTTPRequest = &http.Request{
		Method:     http.MethodGet,
		RequestURI: "http://localhost/",
		Header:     map[string][]string{},
	}
	var dummyActionFuncMap = map[string]ActionFunc{}
	var dummyRoute *mux.Route
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(mux.CurrentRoute, 1, func(httpRequest *http.Request) *mux.Route {
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		return dummyRoute
	})
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeNotFound, errorCode)
		assert.Equal(t, "No corresponding route configured for path", errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	})

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
	var dummyActionFuncMap = map[string]ActionFunc{}
	var dummyRoute = &mux.Route{}
	var dummyName = "some name"
	var dummyEndpoint = "some endpoint"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(mux.CurrentRoute, 1, func(httpRequest *http.Request) *mux.Route {
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		return dummyRoute
	})
	m.ExpectFunc(getName, 1, func(route *mux.Route) string {
		assert.Equal(t, dummyRoute, route)
		return dummyName
	})
	m.ExpectFunc(getEndpointByName, 1, func(name string) string {
		assert.Equal(t, dummyName, name)
		return dummyEndpoint
	})

	// SUT + act
	var endpoint, action, err = getRouteInfo(
		dummyHTTPRequest,
		dummyActionFuncMap,
	)

	// assert
	assert.Equal(t, dummyEndpoint, endpoint)
	functionPointerEquals(t, defaultActionFunc, action)
	assert.NoError(t, err)
}

func TestGetRouteInfo_ValidRoute(t *testing.T) {
	// arrange
	var dummyHTTPRequest = &http.Request{
		Method:     http.MethodGet,
		RequestURI: "http://localhost/",
		Header:     map[string][]string{},
	}
	var dummyActionFuncMap = map[string]ActionFunc{}
	var dummyRoute = &mux.Route{}
	var dummyName = "some name"
	var dummyActionExpected = 0
	var dummyActionCalled = 0
	var dummyAction = func(Session) (interface{}, error) {
		dummyActionCalled++
		return nil, nil
	}
	var dummyEndpoint = "some endpoint"

	// stub
	dummyActionFuncMap[dummyName] = dummyAction

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(mux.CurrentRoute, 1, func(httpRequest *http.Request) *mux.Route {
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		return dummyRoute
	})
	m.ExpectFunc(getName, 1, func(route *mux.Route) string {
		assert.Equal(t, dummyRoute, route)
		return dummyName
	})
	m.ExpectFunc(getEndpointByName, 1, func(name string) string {
		assert.Equal(t, dummyName, name)
		return dummyEndpoint
	})

	// SUT + act
	var endpoint, action, err = getRouteInfo(
		dummyHTTPRequest,
		dummyActionFuncMap,
	)

	// assert
	assert.Equal(t, dummyEndpoint, endpoint)
	functionPointerEquals(t, dummyAction, action)
	assert.NoError(t, err)
	assert.Equal(t, dummyActionExpected, dummyActionCalled, "Unexpected number of calls to dummyAction")
}
