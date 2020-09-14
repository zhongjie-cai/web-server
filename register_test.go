package webserver

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestDoParameterReplacement_EmptyParameterType(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyParameterName = "some name"
	var dummyOriginalPath = "/some/original/path/with/{" + dummyParameterName + "}/in/it"
	var dummyParameterType ParameterType

	// mock
	createMock(t)

	// expect
	logAppRootFuncExpected = 1
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "register", category)
		assert.Equal(t, "doParameterReplacement", subcategory)
		assert.Equal(t, "Path parameter [%v] in path [%v] has no type specification; fallback to default.", messageFormat)
		assert.Equal(t, 2, len(parameters))
		assert.Equal(t, dummyParameterName, parameters[0])
		assert.Equal(t, dummyOriginalPath, parameters[1])
	}

	// SUT + act
	var result = doParameterReplacement(
		dummySession,
		dummyOriginalPath,
		dummyParameterName,
		dummyParameterType,
	)

	// assert
	assert.Equal(t, dummyOriginalPath, result)

	// verify
	verifyAll(t)
}

func TestDoParameterReplacement_ValidParameterType(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyParameterName = "some name"
	var dummyOriginalPath = "/some/original/path/with/{" + dummyParameterName + "}/in/it"
	var dummyParameterType = ParameterType("some type")
	var dummyResult = "/some/original/path/with/{" + dummyParameterName + ":" + string(dummyParameterType) + "}/in/it"

	// mock
	createMock(t)

	// expect
	fmtSprintfExpected = 2
	fmtSprintf = func(format string, a ...interface{}) string {
		fmtSprintfCalled++
		if fmtSprintfCalled == 1 {
			assert.Equal(t, "{%v}", format)
			assert.Equal(t, 1, len(a))
			assert.Equal(t, dummyParameterName, a[0])
		} else if fmtSprintfCalled == 2 {
			assert.Equal(t, "{%v:%v}", format)
			assert.Equal(t, 2, len(a))
			assert.Equal(t, dummyParameterName, a[0])
			assert.Equal(t, dummyParameterType, a[1])
		}
		return fmt.Sprintf(format, a...)
	}
	stringsReplaceExpected = 1
	stringsReplace = func(s, old, new string, n int) string {
		stringsReplaceCalled++
		assert.Equal(t, dummyOriginalPath, s)
		assert.Equal(t, "{"+dummyParameterName+"}", old)
		assert.Equal(t, "{"+dummyParameterName+":"+string(dummyParameterType)+"}", new)
		assert.Equal(t, -1, n)
		return strings.Replace(s, old, new, n)
	}

	// SUT + act
	var result = doParameterReplacement(
		dummySession,
		dummyOriginalPath,
		dummyParameterName,
		dummyParameterType,
	)

	// assert
	assert.Equal(t, dummyResult, result)

	// verify
	verifyAll(t)
}

func TestEvaluatePathWithParameters(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyOriginalPath = "some original path"
	var dummyParameterName1 = "some parameter name 1"
	var dummyParameterType1 = ParameterType("some parameter type 1")
	var dummyParameterName2 = "some parameter name 2"
	var dummyParameterType2 = ParameterType("some parameter type 2")
	var dummyParameterName3 = "some parameter name 3"
	var dummyParameterType3 = ParameterType("some parameter type 3")
	var dummyParameters = map[string]ParameterType{
		dummyParameterName1: dummyParameterType1,
		dummyParameterName2: dummyParameterType2,
		dummyParameterName3: dummyParameterType3,
	}
	var dummyUpdatedPath = "some updated path"

	// mock
	createMock(t)

	// expect
	doParameterReplacementFuncExpected = 3
	doParameterReplacementFunc = func(session *session, originalPath string, parameterName string, parameterType ParameterType) string {
		doParameterReplacementFuncCalled++
		assert.Equal(t, dummySession, session)
		if dummyParameterName1 == parameterName {
			assert.Equal(t, dummyParameterType1, parameterType)
			return dummyUpdatedPath
		} else if dummyParameterName2 == parameterName {
			assert.Equal(t, dummyParameterType2, parameterType)
			return dummyUpdatedPath
		} else if dummyParameterName3 == parameterName {
			assert.Equal(t, dummyParameterType3, parameterType)
			return dummyUpdatedPath
		}
		return ""
	}

	// SUT + act
	var result = evaluatePathWithParameters(
		dummySession,
		dummyOriginalPath,
		dummyParameters,
	)

	// assert
	assert.Equal(t, dummyUpdatedPath, result)

	// verify
	verifyAll(t)
}

func TestEvaluateQueries(t *testing.T) {
	// arrange
	var dummyQueryName1 = "some query name 1"
	var dummyParameterType1 ParameterType
	var dummyQueryName2 = "some query name 2"
	var dummyParameterType2 = ParameterType("some parameter type 2")
	var dummyQueries = map[string]ParameterType{
		dummyQueryName1: dummyParameterType1,
		dummyQueryName2: dummyParameterType2,
	}
	var expectedResult = []string{
		dummyQueryName1, "{" + dummyQueryName1 + "}",
		dummyQueryName2, "{" + dummyQueryName2 + ":" + string(dummyParameterType2) + "}",
	}

	// mock
	createMock(t)

	// expect
	fmtSprintfExpected = 2
	fmtSprintf = func(format string, a ...interface{}) string {
		fmtSprintfCalled++
		return fmt.Sprintf(format, a...)
	}

	// SUT + act
	var result = evaluateQueries(
		dummyQueries,
	)

	// assert
	assert.Equal(t, 4, len(result))
	assert.ElementsMatch(t, expectedResult, result)

	// verify
	verifyAll(t)
}

type dummyCustomizationRegisterRoutes struct {
	dummyCustomization
	routes func() []Route
}

func (customization *dummyCustomizationRegisterRoutes) Routes() []Route {
	if customization.routes != nil {
		return customization.routes()
	}
	assert.Fail(customization.t, "Unexpected call to Routes")
	return nil
}

func TestRegisterRoutes_EmptyRoutes(t *testing.T) {
	// arrange
	var dummyPort = rand.Intn(65536)
	var dummyCustomizationRegisterRoutes = &dummyCustomizationRegisterRoutes{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		customization: dummyCustomizationRegisterRoutes,
	}
	var dummyRouter = &mux.Router{}
	var customizationRoutesExpected int
	var customizationRoutesCalled int
	var dummyRoutes []Route

	// mock
	createMock(t)

	// expect
	customizationRoutesExpected = 1
	dummyCustomizationRegisterRoutes.routes = func() []Route {
		customizationRoutesCalled++
		return dummyRoutes
	}
	logAppRootFuncExpected = 1
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "register", category)
		assert.Equal(t, "registerRoutes", subcategory)
		assert.Equal(t, "customization.Routes function empty: no routes returned!", messageFormat)
		assert.Equal(t, 0, len(parameters))
	}

	// SUT + act
	registerRoutes(
		dummyPort,
		dummySession,
		dummyRouter,
	)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationRoutesExpected, customizationRoutesCalled, "Unexpected number of calls to customization.Routes")
}

func TestRegisterRoutes_ValidRoutes(t *testing.T) {
	// arrange
	var dummyPort = rand.Intn(65536)
	var dummyCustomizationRegisterRoutes = &dummyCustomizationRegisterRoutes{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		customization: dummyCustomizationRegisterRoutes,
	}
	var dummyRouter = &mux.Router{}
	var customizationRoutesExpected int
	var customizationRoutesCalled int
	var dummyEndpoint1 = "some endpoint 1"
	var dummyMethod1 = "some method 1"
	var dummyPath1 = "some path 1"
	var dummyParameters1 = map[string]ParameterType{
		"foo1": ParameterType("bar1"),
	}
	var dummyQueries1 = map[string]ParameterType{
		"test1": ParameterType("me1"),
	}
	var dummyActionFunc1 = func(Session) (interface{}, error) {
		return nil, nil
	}
	var dummyActionFunc1Pointer = fmt.Sprintf("%v", reflect.ValueOf(dummyActionFunc1))
	var dummyEndpoint2 = "some endpoint 2"
	var dummyMethod2 = "some method 2"
	var dummyPath2 = "some path 2"
	var dummyParameters2 = map[string]ParameterType{
		"foo2": ParameterType("bar2"),
	}
	var dummyQueries2 = map[string]ParameterType{
		"test2": ParameterType("me2"),
	}
	var dummyActionFunc2 = func(Session) (interface{}, error) {
		return nil, nil
	}
	var dummyActionFunc2Pointer = fmt.Sprintf("%v", reflect.ValueOf(dummyActionFunc2))
	var dummyRoutes = []Route{
		{
			Endpoint:   dummyEndpoint1,
			Method:     dummyMethod1,
			Path:       dummyPath1,
			Parameters: dummyParameters1,
			Queries:    dummyQueries1,
			ActionFunc: dummyActionFunc1,
		},
		{
			Endpoint:   dummyEndpoint2,
			Method:     dummyMethod2,
			Path:       dummyPath2,
			Parameters: dummyParameters2,
			Queries:    dummyQueries2,
			ActionFunc: dummyActionFunc2,
		},
	}
	var dummyEvaluatedPath1 = "some evaluated path 1"
	var dummyEvaluatedPath2 = "some evaluated path 2"
	var dummyEvaluatedQueries1 = []string{"some evaluated queries 1"}
	var dummyEvaluatedQueries2 = []string{"some evaluated queries 2"}

	// mock
	createMock(t)

	// expect
	customizationRoutesExpected = 1
	dummyCustomizationRegisterRoutes.routes = func() []Route {
		customizationRoutesCalled++
		return dummyRoutes
	}
	evaluatePathWithParametersFuncExpected = 2
	evaluatePathWithParametersFunc = func(session *session, path string, parameters map[string]ParameterType) string {
		evaluatePathWithParametersFuncCalled++
		assert.Equal(t, dummySession, session)
		if dummyPath1 == path {
			assert.Equal(t, dummyParameters1, parameters)
			return dummyEvaluatedPath1
		} else if dummyPath2 == path {
			assert.Equal(t, dummyParameters2, parameters)
			return dummyEvaluatedPath2
		}
		return ""
	}
	evaluateQueriesFuncExpected = 2
	evaluateQueriesFunc = func(queries map[string]ParameterType) []string {
		evaluateQueriesFuncCalled++
		if queries["test1"] == ParameterType("me1") {
			return dummyEvaluatedQueries1
		} else if queries["test2"] == ParameterType("me2") {
			return dummyEvaluatedQueries2
		}
		return nil
	}
	registerRouteFuncExpected = 2
	registerRouteFunc = func(router *mux.Router, endpoint string, method string, path string, queries []string, handlerFunc func(http.ResponseWriter, *http.Request), actionFunc ActionFunc, port int) *mux.Route {
		registerRouteFuncCalled++
		assert.Equal(t, dummyRouter, router)
		assert.Equal(t, fmt.Sprintf("%v", reflect.ValueOf(handleSession)), fmt.Sprintf("%v", reflect.ValueOf(handlerFunc)))
		assert.Equal(t, dummyPort, port)
		if registerRouteFuncCalled == 1 {
			assert.Equal(t, dummyEndpoint1, endpoint)
			assert.Equal(t, dummyMethod1, method)
			assert.Equal(t, dummyEvaluatedPath1, path)
			assert.Equal(t, dummyEvaluatedQueries1, queries)
			assert.Equal(t, dummyActionFunc1Pointer, fmt.Sprintf("%v", reflect.ValueOf(actionFunc)))
		} else if registerRouteFuncCalled == 2 {
			assert.Equal(t, dummyEndpoint2, endpoint)
			assert.Equal(t, dummyMethod2, method)
			assert.Equal(t, dummyEvaluatedPath2, path)
			assert.Equal(t, dummyEvaluatedQueries2, queries)
			assert.Equal(t, dummyActionFunc2Pointer, fmt.Sprintf("%v", reflect.ValueOf(actionFunc)))
		}
		return nil
	}

	// SUT + act
	registerRoutes(
		dummyPort,
		dummySession,
		dummyRouter,
	)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationRoutesExpected, customizationRoutesCalled, "Unexpected number of calls to customization.Routes")
}

type dummyCustomizationRegisterStatics struct {
	dummyCustomization
	statics func() []Static
}

func (customization *dummyCustomizationRegisterStatics) Statics() []Static {
	if customization.statics != nil {
		return customization.statics()
	}
	assert.Fail(customization.t, "Unexpected call to Statics")
	return nil
}

func TestRegisterStatics_EmptyStatics(t *testing.T) {
	// arrange
	var dummyCustomizationRegisterStatics = &dummyCustomizationRegisterStatics{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		customization: dummyCustomizationRegisterStatics,
	}
	var dummyRouter = &mux.Router{}
	var customizationStaticsExpected int
	var customizationStaticsCalled int
	var dummyStatics []Static

	// mock
	createMock(t)

	// expect
	customizationStaticsExpected = 1
	dummyCustomizationRegisterStatics.statics = func() []Static {
		customizationStaticsCalled++
		return dummyStatics
	}
	logAppRootFuncExpected = 1
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "register", category)
		assert.Equal(t, "registerStatics", subcategory)
		assert.Equal(t, "customization.Statics function empty: no static content returned!", messageFormat)
		assert.Equal(t, 0, len(parameters))
	}

	// SUT + act
	registerStatics(
		dummySession,
		dummyRouter,
	)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationStaticsExpected, customizationStaticsCalled, "Unexpected number of calls to customization.Statics")
}

func TestRegisterStatics_ValidStatics(t *testing.T) {
	// arrange
	var dummyCustomizationRegisterStatics = &dummyCustomizationRegisterStatics{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		customization: dummyCustomizationRegisterStatics,
	}
	var dummyRouter = &mux.Router{}
	var customizationStaticsExpected int
	var customizationStaticsCalled int
	var dummyName1 = "some name 1"
	var dummyPathPrefix1 = "some path prefix 1"
	var dummyHandler1 = &dummyHandler{t}
	var dummyName2 = "some name 2"
	var dummyPathPrefix2 = "some path prefix 2"
	var dummyHandler2 = &dummyHandler{t}
	var dummyStatics = []Static{
		{
			Name:       dummyName1,
			PathPrefix: dummyPathPrefix1,
			Handler:    dummyHandler1,
		},
		{
			Name:       dummyName2,
			PathPrefix: dummyPathPrefix2,
			Handler:    dummyHandler2,
		},
	}

	// mock
	createMock(t)

	// expect
	customizationStaticsExpected = 1
	dummyCustomizationRegisterStatics.statics = func() []Static {
		customizationStaticsCalled++
		return dummyStatics
	}
	registerStaticFuncExpected = 2
	registerStaticFunc = func(router *mux.Router, name string, path string, handler http.Handler) *mux.Route {
		registerStaticFuncCalled++
		assert.Equal(t, dummyRouter, router)
		if registerStaticFuncCalled == 1 {
			assert.Equal(t, dummyName1, name)
			assert.Equal(t, dummyPathPrefix1, path)
			assert.Equal(t, dummyHandler1, handler)
		} else if registerStaticFuncCalled == 2 {
			assert.Equal(t, dummyName2, name)
			assert.Equal(t, dummyPathPrefix2, path)
			assert.Equal(t, dummyHandler2, handler)
		}
		return nil
	}

	// SUT + act
	registerStatics(
		dummySession,
		dummyRouter,
	)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationStaticsExpected, customizationStaticsCalled, "Unexpected number of calls to customization.Statics")
}

type dummyCustomizationRegisterMiddlewares struct {
	dummyCustomization
	middlewares func() []MiddlewareFunc
}

func (customization *dummyCustomizationRegisterMiddlewares) Middlewares() []MiddlewareFunc {
	if customization.middlewares != nil {
		return customization.middlewares()
	}
	assert.Fail(customization.t, "Unexpected call to Middlewares")
	return nil
}

func TestRegisterMiddlewares_EmptyMiddlewares(t *testing.T) {
	// arrange
	var dummyCustomizationRegisterMiddlewares = &dummyCustomizationRegisterMiddlewares{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		customization: dummyCustomizationRegisterMiddlewares,
	}
	var dummyRouter = &mux.Router{}
	var customizationMiddlewaresExpected int
	var customizationMiddlewaresCalled int
	var dummyMiddlewares []MiddlewareFunc

	// mock
	createMock(t)

	// expect
	customizationMiddlewaresExpected = 1
	dummyCustomizationRegisterMiddlewares.middlewares = func() []MiddlewareFunc {
		customizationMiddlewaresCalled++
		return dummyMiddlewares
	}
	logAppRootFuncExpected = 1
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "register", category)
		assert.Equal(t, "registerMiddlewares", subcategory)
		assert.Equal(t, "customization.Middlewares function empty: no middleware returned!", messageFormat)
		assert.Equal(t, 0, len(parameters))
	}

	// SUT + act
	registerMiddlewares(
		dummySession,
		dummyRouter,
	)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationMiddlewaresExpected, customizationMiddlewaresCalled, "Unexpected number of calls to customization.Middlewares")
}

func TestRegisterMiddlewares_ValidMiddlewares(t *testing.T) {
	// arrange
	var dummyCustomizationRegisterMiddlewares = &dummyCustomizationRegisterMiddlewares{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		customization: dummyCustomizationRegisterMiddlewares,
	}
	var dummyRouter = &mux.Router{}
	var customizationMiddlewaresExpected int
	var customizationMiddlewaresCalled int
	var dummyMiddleware1 = func(http.Handler) http.Handler { return nil }
	var dummyMiddleware1Pointer = fmt.Sprintf("%v", reflect.ValueOf(dummyMiddleware1))
	var dummyMiddleware2 = func(http.Handler) http.Handler { return nil }
	var dummyMiddleware2Pointer = fmt.Sprintf("%v", reflect.ValueOf(dummyMiddleware2))
	var dummyMiddlewares = []MiddlewareFunc{
		dummyMiddleware1,
		dummyMiddleware2,
	}

	// mock
	createMock(t)

	// expect
	customizationMiddlewaresExpected = 1
	dummyCustomizationRegisterMiddlewares.middlewares = func() []MiddlewareFunc {
		customizationMiddlewaresCalled++
		return dummyMiddlewares
	}
	addMiddlewareFuncExpected = 2
	addMiddlewareFunc = func(router *mux.Router, middleware MiddlewareFunc) {
		addMiddlewareFuncCalled++
		assert.Equal(t, dummyRouter, router)
		var middlewarePointer = fmt.Sprintf("%v", reflect.ValueOf(middleware))
		if addMiddlewareFuncCalled == 1 {
			assert.Equal(t, dummyMiddleware1Pointer, middlewarePointer)
		} else if addMiddlewareFuncCalled == 2 {
			assert.Equal(t, dummyMiddleware2Pointer, middlewarePointer)
		}
	}

	// SUT + act
	registerMiddlewares(
		dummySession,
		dummyRouter,
	)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationMiddlewaresExpected, customizationMiddlewaresCalled, "Unexpected number of calls to customization.Middlewares")
}

type dummyCustomizationErrorHandlers struct {
	dummyCustomization
	methodNotAllowedHandler func() http.Handler
	notFoundHandler         func() http.Handler
}

func (customization *dummyCustomizationErrorHandlers) MethodNotAllowedHandler() http.Handler {
	if customization.methodNotAllowedHandler != nil {
		return customization.methodNotAllowedHandler()
	}
	assert.Fail(customization.t, "Unexpected call to MethodNotAllowedHandler")
	return nil
}

func (customization *dummyCustomizationErrorHandlers) NotFoundHandler() http.Handler {
	if customization.notFoundHandler != nil {
		return customization.notFoundHandler()
	}
	assert.Fail(customization.t, "Unexpected call to NotFoundHandler")
	return nil
}

func TestRegisterErrorHandlers(t *testing.T) {
	// arrange
	var dummyCustomizationErrorHandlers = &dummyCustomizationErrorHandlers{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummyRouter = &mux.Router{}
	var customizationMethodNotAllowedHandlerExpected int
	var customizationMethodNotAllowedHandlerCalled int
	var customizationNotFoundHandlerExpected int
	var customizationNotFoundHandlerCalled int
	var dummyMethodNotAllowedHandler = &dummyHandler{}
	var dummyNotFoundHandler = &dummyHandler{}

	// mock
	createMock(t)

	// expect
	customizationMethodNotAllowedHandlerExpected = 1
	dummyCustomizationErrorHandlers.methodNotAllowedHandler = func() http.Handler {
		customizationMethodNotAllowedHandlerCalled++
		return dummyMethodNotAllowedHandler
	}
	customizationNotFoundHandlerExpected = 1
	dummyCustomizationErrorHandlers.notFoundHandler = func() http.Handler {
		customizationNotFoundHandlerCalled++
		return dummyNotFoundHandler
	}

	// SUT + act
	registerErrorHandlers(
		dummyCustomizationErrorHandlers,
		dummyRouter,
	)

	// assert
	assert.Equal(t, dummyMethodNotAllowedHandler, dummyRouter.MethodNotAllowedHandler)
	assert.Equal(t, dummyNotFoundHandler, dummyRouter.NotFoundHandler)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationMethodNotAllowedHandlerExpected, customizationMethodNotAllowedHandlerCalled, "Unexpected number of calls to customization.MethodNotFoundHandler")
	assert.Equal(t, customizationNotFoundHandlerExpected, customizationNotFoundHandlerCalled, "Unexpected number of calls to customization.NotFoundHandler")
}

type dummyCustomizationInitRouter struct {
	dummyCustomization
	instrumentRouter func(router *mux.Router) *mux.Router
}

func (customization *dummyCustomizationInitRouter) InstrumentRouter(router *mux.Router) *mux.Router {
	if customization.instrumentRouter != nil {
		return customization.instrumentRouter(router)
	}
	assert.Fail(customization.t, "Unexpected call to InstrumentRouter")
	return nil
}

func TestInstantiateRouter_RouterError(t *testing.T) {
	// arrange
	var dummyCustomizationInitRouter = &dummyCustomizationInitRouter{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummyPort = rand.Intn(65536)
	var dummySession = &session{
		customization: dummyCustomizationInitRouter,
	}
	var dummyRouter = &mux.Router{KeepContext: rand.Intn(100) > 50}
	var dummyError = errors.New("some error")
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	createMock(t)

	// expect
	muxNewRouterExpected = 1
	muxNewRouter = func() *mux.Router {
		muxNewRouterCalled++
		return dummyRouter
	}
	registerRoutesFuncExpected = 1
	registerRoutesFunc = func(port int, session *session, router *mux.Router) {
		registerRoutesFuncCalled++
		assert.Equal(t, dummyPort, port)
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
	}
	registerStaticsFuncExpected = 1
	registerStaticsFunc = func(session *session, router *mux.Router) {
		registerStaticsFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
	}
	registerMiddlewaresFuncExpected = 1
	registerMiddlewaresFunc = func(session *session, router *mux.Router) {
		registerMiddlewaresFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
	}
	walkRegisteredRoutesFuncExpected = 1
	walkRegisteredRoutesFunc = func(session *session, router *mux.Router) error {
		walkRegisteredRoutesFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
		return dummyError
	}
	logAppRootFuncExpected = 1
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "register", category)
		assert.Equal(t, "instantiateRouter", subcategory)
		assert.Equal(t, "%+v", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyError, parameters[0])
	}
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessageor string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageRouteRegistration, errorMessageor)
		assert.Equal(t, 1, len(innerErrors))
		assert.Equal(t, dummyError, innerErrors[0])
		return dummyAppError
	}

	// SUT + act
	var result, err = instantiateRouter(
		dummyPort,
		dummySession,
	)

	// assert
	assert.Equal(t, dummyRouter, result)
	assert.Equal(t, dummyAppError, err)

	// verify
	verifyAll(t)
}

func TestInstantiateRouter_HappyPath(t *testing.T) {
	// arrange
	var dummyCustomizationInitRouter = &dummyCustomizationInitRouter{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummyPort = rand.Intn(65536)
	var dummySession = &session{
		customization: dummyCustomizationInitRouter,
	}
	var dummyRouter = &mux.Router{KeepContext: rand.Intn(100) > 50}
	var customizationInstrumentRouterExpected int
	var customizationInstrumentRouterCalled int

	// mock
	createMock(t)

	// expect
	muxNewRouterExpected = 1
	muxNewRouter = func() *mux.Router {
		muxNewRouterCalled++
		return dummyRouter
	}
	registerRoutesFuncExpected = 1
	registerRoutesFunc = func(port int, session *session, router *mux.Router) {
		registerRoutesFuncCalled++
		assert.Equal(t, dummyPort, port)
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
	}
	registerStaticsFuncExpected = 1
	registerStaticsFunc = func(session *session, router *mux.Router) {
		registerStaticsFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
	}
	registerMiddlewaresFuncExpected = 1
	registerMiddlewaresFunc = func(session *session, router *mux.Router) {
		registerMiddlewaresFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
	}
	walkRegisteredRoutesFuncExpected = 1
	walkRegisteredRoutesFunc = func(session *session, router *mux.Router) error {
		walkRegisteredRoutesFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
		return nil
	}
	registerErrorHandlersFuncExpected = 1
	registerErrorHandlersFunc = func(customization Customization, router *mux.Router) {
		registerErrorHandlersFuncCalled++
		assert.Equal(t, dummyCustomizationInitRouter, customization)
		assert.Equal(t, dummyRouter, router)
	}
	customizationInstrumentRouterExpected = 1
	dummyCustomizationInitRouter.instrumentRouter = func(router *mux.Router) *mux.Router {
		customizationInstrumentRouterCalled++
		assert.Equal(t, dummyRouter, router)
		return dummyRouter
	}

	// SUT + act
	var result, err = instantiateRouter(
		dummyPort,
		dummySession,
	)

	// assert
	assert.Equal(t, dummyRouter, result)
	assert.NoError(t, err)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationInstrumentRouterExpected, customizationInstrumentRouterCalled, "Unexpected number of calls to customization.InstrumentRouter")
}
