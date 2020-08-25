package webserver

import (
	"errors"
	"math/rand"
	"net/http"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetName_Undefined(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// mock
	createMock(t)

	// SUT
	var dummyRoute = dummyRouter.NewRoute()

	// act
	var result = getName(
		dummyRoute,
	)

	// assert
	assert.Zero(t, result)

	// verify
	verifyAll(t)
}

func TestGetName_Defined(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// mock
	createMock(t)

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

	// verify
	verifyAll(t)
}

func TestGetPathTemplate_Error(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// mock
	createMock(t)

	// SUT
	var dummyRoute = dummyRouter.NewRoute()

	// act
	var result, err = getPathTemplate(
		dummyRoute,
	)

	// assert
	assert.Zero(t, result)
	assert.Equal(t, "mux: route doesn't have a path", err.Error())

	// verify
	verifyAll(t)
}

func TestGetPathTemplate_Success(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// mock
	createMock(t)

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

	// verify
	verifyAll(t)
}

func TestGetPathRegexp_Error(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// mock
	createMock(t)

	// SUT
	var dummyRoute = dummyRouter.NewRoute()

	// act
	var result, err = getPathRegexp(
		dummyRoute,
	)

	// assert
	assert.Zero(t, result)
	assert.Equal(t, "mux: route does not have a path", err.Error())

	// verify
	verifyAll(t)
}

func TestGetPathRegexp_Success(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// mock
	createMock(t)

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

	// verify
	verifyAll(t)
}

func TestGetQueriesTemplate_Undefined(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// mock
	createMock(t)

	// expect
	stringsJoinExpected = 1
	stringsJoin = func(a []string, sep string) string {
		stringsJoinCalled++
		return strings.Join(a, sep)
	}

	// SUT
	var dummyRoute = dummyRouter.NewRoute()

	// act
	var result = getQueriesTemplates(
		dummyRoute,
	)

	// assert
	assert.Zero(t, result)

	// verify
	verifyAll(t)
}

func TestGetQueriesTemplate_Defined(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// mock
	createMock(t)

	// expect
	stringsJoinExpected = 1
	stringsJoin = func(a []string, sep string) string {
		stringsJoinCalled++
		return strings.Join(a, sep)
	}

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

	// verify
	verifyAll(t)
}

func TestGetQueriesRegexp_Undefined(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// mock
	createMock(t)

	// expect
	stringsJoinExpected = 1
	stringsJoin = func(a []string, sep string) string {
		stringsJoinCalled++
		return strings.Join(a, sep)
	}

	// SUT
	var dummyRoute = dummyRouter.NewRoute()

	// act
	var result = getQueriesRegexp(
		dummyRoute,
	)

	// assert
	assert.Zero(t, result)

	// verify
	verifyAll(t)
}

func TestGetQueriesRegexp_Defined(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// mock
	createMock(t)

	// expect
	stringsJoinExpected = 1
	stringsJoin = func(a []string, sep string) string {
		stringsJoinCalled++
		return strings.Join(a, sep)
	}

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

	// verify
	verifyAll(t)
}

func TestGetMethods_Undefined(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// mock
	createMock(t)

	// expect
	stringsJoinExpected = 1
	stringsJoin = func(a []string, sep string) string {
		stringsJoinCalled++
		return strings.Join(a, sep)
	}

	// SUT
	var dummyRoute = dummyRouter.NewRoute()

	// act
	var result = getMethods(
		dummyRoute,
	)

	// assert
	assert.Zero(t, result)

	// verify
	verifyAll(t)
}

func TestGetMethods_Defined(t *testing.T) {
	// arrange
	var dummyRouter = mux.NewRouter()

	// mock
	createMock(t)

	// expect
	stringsJoinExpected = 1
	stringsJoin = func(a []string, sep string) string {
		stringsJoinCalled++
		return strings.Join(a, sep)
	}

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

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	getPathTemplateFuncExpected = 1
	getPathTemplateFunc = func(route *mux.Route) (string, error) {
		getPathTemplateFuncCalled++
		assert.Equal(t, dummyRoute, route)
		return dummyPathTemplate, dummyPathTemplateError
	}
	getPathRegexpFuncExpected = 1
	getPathRegexpFunc = func(route *mux.Route) (string, error) {
		getPathRegexpFuncCalled++
		assert.Equal(t, dummyRoute, route)
		return dummyPathRegexp, dummyPathRegexpError
	}

	// SUT + act
	var err = evaluateRoute(
		dummyRoute,
		dummyRouter,
		dummyAncestors,
	)

	// assert
	assert.Equal(t, dummyPathTemplateError, err)

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	getPathTemplateFuncExpected = 1
	getPathTemplateFunc = func(route *mux.Route) (string, error) {
		getPathTemplateFuncCalled++
		assert.Equal(t, dummyRoute, route)
		return dummyPathTemplate, nil
	}
	getPathRegexpFuncExpected = 1
	getPathRegexpFunc = func(route *mux.Route) (string, error) {
		getPathRegexpFuncCalled++
		assert.Equal(t, dummyRoute, route)
		return dummyPathRegexp, dummyPathRegexpError
	}

	// SUT + act
	var err = evaluateRoute(
		dummyRoute,
		dummyRouter,
		dummyAncestors,
	)

	// assert
	assert.Equal(t, dummyPathRegexpError, err)

	// verify
	verifyAll(t)
}

func TestPrintRegisteredRouteDetails_Success(t *testing.T) {
	// arrange
	var dummyRoute = &mux.Route{}
	var dummyRouter = &mux.Router{}
	var dummyAncestors = []*mux.Route{}
	var dummyPathTemplate = "some path template"
	var dummyPathRegexp = "some path regexp"

	// mock
	createMock(t)

	// expect
	getPathTemplateFuncExpected = 1
	getPathTemplateFunc = func(route *mux.Route) (string, error) {
		getPathTemplateFuncCalled++
		assert.Equal(t, dummyRoute, route)
		return dummyPathTemplate, nil
	}
	getPathRegexpFuncExpected = 1
	getPathRegexpFunc = func(route *mux.Route) (string, error) {
		getPathRegexpFuncCalled++
		assert.Equal(t, dummyRoute, route)
		return dummyPathRegexp, nil
	}

	// SUT + act
	var err = evaluateRoute(
		dummyRoute,
		dummyRouter,
		dummyAncestors,
	)

	// assert
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestWalkRegisteredRoutes_Error(t *testing.T) {
	// arrange
	var dummySession = &session{name: "some name"}
	var dummyRouter = &mux.Router{}
	var dummyError = errors.New("some error")

	// stub
	dummyRouter.HandleFunc("/", func(http.ResponseWriter, *http.Request) {})

	// mock
	createMock(t)

	// expect
	evaluateRouteFuncExpected = 1
	evaluateRouteFunc = func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		evaluateRouteFuncCalled++
		return dummyError
	}
	logAppRootFuncExpected = 1
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "route", category)
		assert.Equal(t, "walkRegisteredRoutes", subcategory)
		assert.Equal(t, "Failure: %+v", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyError, parameters[0])
	}

	// SUT + act
	var err = walkRegisteredRoutes(
		dummySession,
		dummyRouter,
	)

	// assert
	assert.Equal(t, errRouteRegistration, err)

	// verify
	verifyAll(t)
}

func TestWalkRegisteredRoutes_Success(t *testing.T) {
	// arrange
	var dummySession = &session{name: "some name"}
	var dummyRouter = &mux.Router{}

	// stub
	dummyRouter.HandleFunc("/", func(http.ResponseWriter, *http.Request) {})

	// mock
	createMock(t)

	// expect
	evaluateRouteFuncExpected = 1
	evaluateRouteFunc = func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		evaluateRouteFuncCalled++
		return nil
	}

	// SUT + act
	var err = walkRegisteredRoutes(
		dummySession,
		dummyRouter,
	)

	// assert
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestRegisterRoute(t *testing.T) {
	// arrange
	var dummyEndpoint = "some endpoint"
	var dummyMethod = "SOME METHOD"
	var dummyName = "some name"
	var dummyPath = "/foo/{bar}"
	var dummyQueries = []string{"test", "{test}"}
	var dummyQueriesTemplates = []string{"test={test}"}
	var dummyPort = rand.Intn(65536)
	var dummyApplication = &application{
		actionFuncMap: map[string]ActionFunc{},
	}

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
	createMock(t)

	// expect
	fmtSprintfExpected = 1
	fmtSprintf = func(format string, a ...interface{}) string {
		fmtSprintfCalled++
		assert.Equal(t, "%v:%v", format)
		assert.Equal(t, 2, len(a))
		assert.Equal(t, dummyEndpoint, a[0])
		assert.Equal(t, dummyMethod, a[1])
		return dummyName
	}
	getApplicationFuncExpected = 1
	getApplicationFunc = func(port int) *application {
		getApplicationFuncCalled++
		assert.Equal(t, dummyPort, port)
		return dummyApplication
	}

	// SUT
	var router = mux.NewRouter()

	// act
	var route = registerRoute(
		router,
		dummyEndpoint,
		dummyMethod,
		dummyPath,
		dummyQueries,
		dummyHandlerFunc,
		dummyActionFunc,
		dummyPort,
	)
	var name = route.GetName()
	var methods, _ = route.GetMethods()
	var pathTemplate, _ = route.GetPathTemplate()
	var queriesTemplate, _ = route.GetQueriesTemplates()

	// assert
	assert.Equal(t, dummyName, name)
	assert.Equal(t, 1, len(methods))
	assert.Equal(t, dummyMethod, methods[0])
	assert.Equal(t, dummyPath, pathTemplate)
	assert.Equal(t, dummyQueriesTemplates, queriesTemplate)
	assert.Equal(t, dummyHandlerFuncExpected, dummyHandlerFuncCalled)
	assert.Equal(t, dummyActionFuncExpected, dummyActionFuncCalled)
	functionPointerEquals(t, dummyActionFunc, dummyApplication.actionFuncMap[dummyName])

	// verify
	verifyAll(t)
}

func TestDefaultActionFunc(t *testing.T) {
	// arrange
	var dummySessionObject = &dummySession{t}

	// mock
	createMock(t)

	// SUT + act
	var result, err = defaultActionFunc(
		dummySessionObject,
	)

	// assert
	assert.Nil(t, result)
	assert.Equal(t, errRouteNotImplemented, err)

	// verify
	verifyAll(t)
}

func TestGetEndpointByName_NoSeparator(t *testing.T) {
	// arrange
	var dummyName = "some name"

	// mock
	createMock(t)

	// expect
	stringsSplitExpected = 1
	stringsSplit = func(s string, sep string) []string {
		stringsSplitCalled++
		assert.Equal(t, dummyName, s)
		assert.Equal(t, ":", sep)
		return strings.Split(s, sep)
	}

	// SUT + act
	result := getEndpointByName(
		dummyName,
	)

	// assert
	assert.Equal(t, dummyName, result)

	// verify
	verifyAll(t)
}

func TestGetEndpointByName_WithSeparator(t *testing.T) {
	// arrange
	var dummyEndpoint = "some endpoint"
	var dummyName = dummyEndpoint + ":some name"

	// mock
	createMock(t)

	// expect
	stringsSplitExpected = 1
	stringsSplit = func(s string, sep string) []string {
		stringsSplitCalled++
		assert.Equal(t, dummyName, s)
		assert.Equal(t, ":", sep)
		return strings.Split(s, sep)
	}

	// SUT + act
	result := getEndpointByName(
		dummyName,
	)

	// assert
	assert.Equal(t, dummyEndpoint, result)

	// verify
	verifyAll(t)
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

	// mock
	createMock(t)

	// expect
	muxCurrentRouteExpected = 1
	muxCurrentRoute = func(httpRequest *http.Request) *mux.Route {
		muxCurrentRouteCalled++
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		return dummyRoute
	}

	// SUT + act
	var name, action, err = getRouteInfo(
		dummyHTTPRequest,
		dummyActionFuncMap,
	)

	// assert
	assert.Zero(t, name)
	assert.Nil(t, action)
	assert.Equal(t, errRouteNotFound, err)

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	muxCurrentRouteExpected = 1
	muxCurrentRoute = func(httpRequest *http.Request) *mux.Route {
		muxCurrentRouteCalled++
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		return dummyRoute
	}
	getNameFuncExpected = 1
	getNameFunc = func(route *mux.Route) string {
		getNameFuncCalled++
		assert.Equal(t, dummyRoute, route)
		return dummyName
	}
	getEndpointByNameFuncExpected = 1
	getEndpointByNameFunc = func(name string) string {
		getEndpointByNameFuncCalled++
		assert.Equal(t, dummyName, name)
		return dummyEndpoint
	}

	// SUT + act
	var endpoint, action, err = getRouteInfo(
		dummyHTTPRequest,
		dummyActionFuncMap,
	)

	// assert
	assert.Equal(t, dummyEndpoint, endpoint)
	functionPointerEquals(t, defaultActionFunc, action)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	muxCurrentRouteExpected = 1
	muxCurrentRoute = func(httpRequest *http.Request) *mux.Route {
		muxCurrentRouteCalled++
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		return dummyRoute
	}
	getNameFuncExpected = 1
	getNameFunc = func(route *mux.Route) string {
		getNameFuncCalled++
		assert.Equal(t, dummyRoute, route)
		return dummyName
	}
	getEndpointByNameFuncExpected = 1
	getEndpointByNameFunc = func(name string) string {
		getEndpointByNameFuncCalled++
		assert.Equal(t, dummyName, name)
		return dummyEndpoint
	}

	// SUT + act
	var endpoint, action, err = getRouteInfo(
		dummyHTTPRequest,
		dummyActionFuncMap,
	)

	// assert
	assert.Equal(t, dummyEndpoint, endpoint)
	functionPointerEquals(t, dummyAction, action)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
	assert.Equal(t, dummyActionExpected, dummyActionCalled, "Unexpected number of calls to dummyAction")
}
