package webserver

// import (
// 	"errors"
// 	"net/http"
// 	"testing"

// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/zhongjie-cai/gomocker/v2"
// )

// func TestGetName_Undefined(t *testing.T) {
// 	// arrange
// 	type router struct {
// 		chi.Router
// 	}
// 	var dummyRouter = &router{}

// 	// SUT
// 	var dummyRoute = dummyRouter.NewRoute()

// 	// act
// 	var result = getName(
// 		dummyRoute,
// 	)

// 	// assert
// 	assert.Zero(t, result)
// }

// func TestGetName_Defined(t *testing.T) {
// 	// arrange
// 	type router struct {
// 		chi.Router
// 	}
// 	var dummyRouter = &router{}

// 	// SUT
// 	var dummyRoute = dummyRouter.NewRoute().Name(
// 		"test",
// 	)

// 	// act
// 	var result = getName(
// 		dummyRoute,
// 	)

// 	// assert
// 	assert.Equal(t, "test", result)
// }

// func TestGetPathTemplate_Error(t *testing.T) {
// 	// arrange
// 	type router struct {
// 		chi.Router
// 	}
// 	var dummyRouter = &router{}

// 	// SUT
// 	var dummyRoute = dummyRouter.NewRoute()

// 	// act
// 	var result, err = getPathTemplate(
// 		dummyRoute,
// 	)

// 	// assert
// 	assert.Zero(t, result)
// 	assert.Equal(t, "mux: route doesn't have a path", err.Error())
// }

// func TestGetPathTemplate_Success(t *testing.T) {
// 	// arrange
// 	type router struct {
// 		chi.Router
// 	}
// 	var dummyRouter = &router{}

// 	// SUT
// 	var dummyRoute = dummyRouter.NewRoute().Path(
// 		"/foo/{bar}",
// 	)

// 	// act
// 	var result, err = getPathTemplate(
// 		dummyRoute,
// 	)

// 	// assert
// 	assert.Equal(t, "/foo/{bar}", result)
// 	assert.NoError(t, err)
// }

// func TestGetPathRegexp_Error(t *testing.T) {
// 	// arrange
// 	type router struct {
// 		chi.Router
// 	}
// 	var dummyRouter = &router{}

// 	// SUT
// 	var dummyRoute = dummyRouter.NewRoute()

// 	// act
// 	var result, err = getPathRegexp(
// 		dummyRoute,
// 	)

// 	// assert
// 	assert.Zero(t, result)
// 	assert.Equal(t, "mux: route does not have a path", err.Error())
// }

// func TestGetPathRegexp_Success(t *testing.T) {
// 	// arrange
// 	type router struct {
// 		chi.Router
// 	}
// 	var dummyRouter = &router{}

// 	// SUT
// 	var dummyRoute = dummyRouter.NewRoute().Path(
// 		"/foo/{bar}",
// 	)

// 	// act
// 	var result, err = getPathRegexp(
// 		dummyRoute,
// 	)

// 	// assert
// 	assert.Equal(t, "^/foo/(?P<v0>[^/]+)$", result)
// 	assert.NoError(t, err)
// }

// func TestGetQueriesTemplate_Undefined(t *testing.T) {
// 	// arrange
// 	type router struct {
// 		chi.Router
// 	}
// 	var dummyRouter = &router{}

// 	// SUT
// 	var dummyRoute = dummyRouter.NewRoute()

// 	// act
// 	var result = getQueriesTemplates(
// 		dummyRoute,
// 	)

// 	// assert
// 	assert.Zero(t, result)
// }

// func TestGetQueriesTemplate_Defined(t *testing.T) {
// 	// arrange
// 	type router struct {
// 		chi.Router
// 	}
// 	var dummyRouter = &router{}

// 	// SUT
// 	var dummyRoute = dummyRouter.NewRoute().Queries(
// 		"abc",
// 		"{def}",
// 		"xyz",
// 		"{zyx}",
// 	)

// 	// act
// 	var result = getQueriesTemplates(
// 		dummyRoute,
// 	)

// 	// assert
// 	assert.Equal(t, "abc={def}|xyz={zyx}", result)
// }

// func TestGetQueriesRegexp_Undefined(t *testing.T) {
// 	// arrange
// 	type router struct {
// 		chi.Router
// 	}
// 	var dummyRouter = &router{}

// 	// SUT
// 	var dummyRoute = dummyRouter.NewRoute()

// 	// act
// 	var result = getQueriesRegexp(
// 		dummyRoute,
// 	)

// 	// assert
// 	assert.Zero(t, result)
// }

// func TestGetQueriesRegexp_Defined(t *testing.T) {
// 	// arrange
// 	type router struct {
// 		chi.Router
// 	}
// 	var dummyRouter = &router{}

// 	// SUT
// 	var dummyRoute = dummyRouter.NewRoute().Queries(
// 		"abc",
// 		"{def}",
// 		"xyz",
// 		"{zyx}",
// 	)

// 	// act
// 	var result = getQueriesRegexp(
// 		dummyRoute,
// 	)

// 	// assert
// 	assert.Equal(t, "^abc=(?P<v0>.*)$|^xyz=(?P<v0>.*)$", result)
// }

// func TestGetMethods_Undefined(t *testing.T) {
// 	// arrange
// 	type router struct {
// 		chi.Router
// 	}
// 	var dummyRouter = &router{}

// 	// SUT
// 	var dummyRoute = dummyRouter.NewRoute()

// 	// act
// 	var result = getMethods(
// 		dummyRoute,
// 	)

// 	// assert
// 	assert.Zero(t, result)
// }

// func TestGetMethods_Defined(t *testing.T) {
// 	// arrange
// 	type router struct {
// 		chi.Router
// 	}
// 	var dummyRouter = &router{}

// 	// SUT
// 	var dummyRoute = dummyRouter.NewRoute().Methods(
// 		"GET",
// 		"PUT",
// 	)

// 	// act
// 	var result = getMethods(
// 		dummyRoute,
// 	)

// 	// assert
// 	assert.Equal(t, "GET|PUT", result)
// }

// func TestPrintRegisteredRouteDetails_TemplateError(t *testing.T) {
// 	// arrange
// 	var dummyRoute = &mux.Route{}
// 	var dummyRouter = &mux.Router{}
// 	var dummyAncestors = []*mux.Route{}
// 	var dummyPathTemplate string
// 	var dummyPathTemplateError = errors.New("some path template error")
// 	var dummyPathRegexp string
// 	var dummyPathRegexpError = errors.New("some path regexp error")

// 	// mock
// 	var m = gomocker.NewMocker(t)

// 	// expect
// 	m.Mock(getPathTemplate).Expects(dummyRoute).Returns(dummyPathTemplate, dummyPathTemplateError).Once()
// 	m.Mock(getPathRegexp).Expects(dummyRoute).Returns(dummyPathRegexp, dummyPathRegexpError).Once()

// 	// SUT + act
// 	var err = evaluateRoute(
// 		dummyRoute,
// 		dummyRouter,
// 		dummyAncestors,
// 	)

// 	// assert
// 	assert.Equal(t, dummyPathTemplateError, err)
// }

// func TestPrintRegisteredRouteDetails_RegexpError(t *testing.T) {
// 	// arrange
// 	var dummyRoute = &mux.Route{}
// 	var dummyRouter = &mux.Router{}
// 	var dummyAncestors = []*mux.Route{}
// 	var dummyPathTemplate string
// 	var dummyPathRegexp string
// 	var dummyPathRegexpError = errors.New("some path regexp error")

// 	// mock
// 	var m = gomocker.NewMocker(t)

// 	// expect
// 	m.Mock(getPathTemplate).Expects(dummyRoute).Returns(dummyPathTemplate, nil).Once()
// 	m.Mock(getPathRegexp).Expects(dummyRoute).Returns(dummyPathRegexp, dummyPathRegexpError).Once()

// 	// SUT + act
// 	var err = evaluateRoute(
// 		dummyRoute,
// 		dummyRouter,
// 		dummyAncestors,
// 	)

// 	// assert
// 	assert.Equal(t, dummyPathRegexpError, err)
// }

// func TestPrintRegisteredRouteDetails_Success(t *testing.T) {
// 	// arrange
// 	var dummyRoute = &mux.Route{}
// 	var dummyRouter = &mux.Router{}
// 	var dummyAncestors = []*mux.Route{}
// 	var dummyPathTemplate = "some path template"
// 	var dummyPathRegexp = "some path regexp"

// 	// mock
// 	var m = gomocker.NewMocker(t)

// 	// expect
// 	m.Mock(getPathTemplate).Expects(dummyRoute).Returns(dummyPathTemplate, nil).Once()
// 	m.Mock(getPathRegexp).Expects(dummyRoute).Returns(dummyPathRegexp, nil).Once()

// 	// SUT + act
// 	var err = evaluateRoute(
// 		dummyRoute,
// 		dummyRouter,
// 		dummyAncestors,
// 	)

// 	// assert
// 	assert.NoError(t, err)
// }

// func TestWalkRegisteredRoutes_Error(t *testing.T) {
// 	// arrange
// 	var dummySession = &session{id: uuid.New()}
// 	var dummyRouter = &mux.Router{}
// 	var dummyError = errors.New("some error")
// 	var dummyAppError = &appError{Message: "some error message"}

// 	// stub
// 	dummyRouter.HandleFunc("/", func(http.ResponseWriter, *http.Request) {})

// 	// mock
// 	var m = gomocker.NewMocker(t)

// 	// expect
// 	m.Mock((*mux.Router).Walk).Expects(dummyRouter, gomocker.Matches(func(value any) bool {
// 		return functionPointerEquals(evaluateRoute, value)
// 	})).Returns(dummyError).Once()
// 	m.Mock(logAppRoot).Expects(dummySession, "route", "walkRegisteredRoutes", "Failure: %+v", dummyError).Returns().Once()
// 	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageRouteRegistration, []error{dummyError}).Returns(dummyAppError).Once()

// 	// SUT + act
// 	var err = walkRegisteredRoutes(
// 		dummySession,
// 		dummyRouter,
// 	)

// 	// assert
// 	assert.Equal(t, dummyAppError, err)
// }

// func TestWalkRegisteredRoutes_Success(t *testing.T) {
// 	// arrange
// 	var dummySession = &session{id: uuid.New()}
// 	var dummyRouter = &mux.Router{}

// 	// stub
// 	dummyRouter.HandleFunc("/", func(http.ResponseWriter, *http.Request) {})

// 	// mock
// 	var m = gomocker.NewMocker(t)

// 	// expect
// 	m.Mock((*mux.Router).Walk).Expects(dummyRouter, gomocker.Matches(func(value any) bool {
// 		return functionPointerEquals(evaluateRoute, value)
// 	})).Returns(nil).Once()

// 	// SUT + act
// 	var err = walkRegisteredRoutes(
// 		dummySession,
// 		dummyRouter,
// 	)

// 	// assert
// 	assert.NoError(t, err)
// }

// func TestRegisterRoute(t *testing.T) {
// 	// arrange
// 	var dummyEndpoint = "some endpoint"
// 	var dummyMethod = "SOME METHOD"
// 	var dummyName = "some endpoint:SOME METHOD"
// 	var dummyPath = "/foo/{bar}"
// 	var dummyQueries = []string{"test", "{test}"}
// 	var dummyRouter = &mux.Router{}
// 	var dummyRoute = &mux.Route{}

// 	// stub
// 	var dummyHandlerFuncExpected = 0
// 	var dummyHandlerFuncCalled = 0
// 	var dummyHandlerFunc = func(http.ResponseWriter, *http.Request) {
// 		dummyHandlerFuncCalled++
// 	}

// 	// mock
// 	var m = gomocker.NewMocker(t)

// 	// expect
// 	m.Mock((*mux.Router).HandleFunc).Expects(dummyRouter, dummyPath, gomocker.Matches(func(value any) bool {
// 		return functionPointerEquals(dummyHandlerFunc, value)
// 	})).Returns(dummyRoute).Once()
// 	m.Mock((*mux.Route).Methods).Expects(dummyRoute, dummyMethod).Returns(dummyRoute).Once()
// 	m.Mock((*mux.Route).Queries).Expects(dummyRoute, "test", "{test}").Returns(dummyRoute).Once()
// 	m.Mock((*mux.Route).Name).Expects(dummyRoute, dummyName).Returns(dummyRoute).Once()

// 	// SUT + act
// 	var name, route = registerRoute(
// 		dummyRouter,
// 		dummyEndpoint,
// 		dummyMethod,
// 		dummyPath,
// 		dummyQueries,
// 		dummyHandlerFunc,
// 	)

// 	// assert
// 	assert.Equal(t, dummyName, name)
// 	assert.Equal(t, dummyRoute, route)
// 	assert.Equal(t, dummyHandlerFuncExpected, dummyHandlerFuncCalled)
// }

// func TestDefaultActionFunc(t *testing.T) {
// 	// arrange
// 	var dummySession = &session{}
// 	var dummyAppError = &appError{Message: "some error message"}

// 	// mock
// 	var m = gomocker.NewMocker(t)

// 	// expect
// 	m.Mock(newAppError).Expects(errorCodeNotImplemented, "No corresponding action function configured; falling back to default", []error{}).Returns(dummyAppError).Once()

// 	// SUT + act
// 	var result, err = defaultActionFunc(
// 		dummySession,
// 	)

// 	// assert
// 	assert.Nil(t, result)
// 	assert.Equal(t, dummyAppError, err)
// }

// func TestGetEndpointByName_NoSeparator(t *testing.T) {
// 	// arrange
// 	var dummyName = "some name"

// 	// SUT + act
// 	result := getEndpointByName(
// 		dummyName,
// 	)

// 	// assert
// 	assert.Equal(t, dummyName, result)
// }

// func TestGetEndpointByName_WithSeparator(t *testing.T) {
// 	// arrange
// 	var dummyEndpoint = "some endpoint"
// 	var dummyName = dummyEndpoint + ":some name"

// 	// SUT + act
// 	result := getEndpointByName(
// 		dummyName,
// 	)

// 	// assert
// 	assert.Equal(t, dummyEndpoint, result)
// }

// func TestGetRouteInfo_NilRoute(t *testing.T) {
// 	// arrange
// 	var dummyHTTPRequest = &http.Request{
// 		Method:     http.MethodGet,
// 		RequestURI: "http://localhost/",
// 		Header:     map[string][]string{},
// 	}
// 	var dummyActionFuncMap = map[string]ActionFunc{}
// 	var dummyRoute *mux.Route
// 	var dummyAppError = &appError{Message: "some error message"}

// 	// mock
// 	var m = gomocker.NewMocker(t)

// 	// expect
// 	m.Mock(mux.CurrentRoute).Expects(dummyHTTPRequest).Returns(dummyRoute).Once()
// 	m.Mock(newAppError).Expects(errorCodeNotFound, "No corresponding route configured for path", []error{}).Returns(dummyAppError).Once()

// 	// SUT + act
// 	var name, action, err = getRouteInfo(
// 		dummyHTTPRequest,
// 		dummyActionFuncMap,
// 	)

// 	// assert
// 	assert.Zero(t, name)
// 	assert.Nil(t, action)
// 	assert.Equal(t, dummyAppError, err)
// }

// func TestGetRouteInfo_RouteNotFound(t *testing.T) {
// 	// arrange
// 	var dummyHTTPRequest = &http.Request{
// 		Method:     http.MethodGet,
// 		RequestURI: "http://localhost/",
// 		Header:     map[string][]string{},
// 	}
// 	var dummyActionFuncMap = map[string]ActionFunc{}
// 	var dummyRoute = &mux.Route{}
// 	var dummyName = "some name"
// 	var dummyEndpoint = "some endpoint"

// 	// mock
// 	var m = gomocker.NewMocker(t)

// 	// expect
// 	m.Mock(mux.CurrentRoute).Expects(dummyHTTPRequest).Returns(dummyRoute).Once()
// 	m.Mock(getName).Expects(dummyRoute).Returns(dummyName).Once()
// 	m.Mock(getEndpointByName).Expects(dummyName).Returns(dummyEndpoint).Once()

// 	// SUT + act
// 	var endpoint, action, err = getRouteInfo(
// 		dummyHTTPRequest,
// 		dummyActionFuncMap,
// 	)

// 	// assert
// 	assert.Equal(t, dummyEndpoint, endpoint)
// 	assertFunctionEquals(t, defaultActionFunc, action)
// 	assert.NoError(t, err)
// }

// func TestGetRouteInfo_ValidRoute(t *testing.T) {
// 	// arrange
// 	var dummyHTTPRequest = &http.Request{
// 		Method:     http.MethodGet,
// 		RequestURI: "http://localhost/",
// 		Header:     map[string][]string{},
// 	}
// 	var dummyActionFuncMap = map[string]ActionFunc{}
// 	var dummyRoute = &mux.Route{}
// 	var dummyName = "some name"
// 	var dummyActionExpected = 0
// 	var dummyActionCalled = 0
// 	var dummyAction = func(Session) (any, error) {
// 		dummyActionCalled++
// 		return nil, nil
// 	}
// 	var dummyEndpoint = "some endpoint"

// 	// stub
// 	dummyActionFuncMap[dummyName] = dummyAction

// 	// mock
// 	var m = gomocker.NewMocker(t)

// 	// expect
// 	m.Mock(mux.CurrentRoute).Expects(dummyHTTPRequest).Returns(dummyRoute).Once()
// 	m.Mock(getName).Expects(dummyRoute).Returns(dummyName).Once()
// 	m.Mock(getEndpointByName).Expects(dummyName).Returns(dummyEndpoint).Once()

// 	// SUT + act
// 	var endpoint, action, err = getRouteInfo(
// 		dummyHTTPRequest,
// 		dummyActionFuncMap,
// 	)

// 	// assert
// 	assert.Equal(t, dummyEndpoint, endpoint)
// 	assertFunctionEquals(t, dummyAction, action)
// 	assert.NoError(t, err)
// 	assert.Equal(t, dummyActionExpected, dummyActionCalled, "Unexpected number of calls to dummyAction")
// }
