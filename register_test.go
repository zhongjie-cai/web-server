package webserver

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/zhongjie-cai/gomocker"
)

func TestDoParameterReplacement_EmptyParameterType(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyParameterName = "some name"
	var dummyOriginalPath = "/some/original/path/with/{" + dummyParameterName + "}/in/it"
	var dummyParameterType ParameterType

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(logAppRoot, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "register", category)
		assert.Equal(t, "doParameterReplacement", subcategory)
		assert.Equal(t, "Path parameter [%v] in path [%v] has no type specification; fallback to default.", messageFormat)
		assert.Equal(t, 2, len(parameters))
		assert.Equal(t, dummyParameterName, parameters[0])
		assert.Equal(t, dummyOriginalPath, parameters[1])
	})

	// SUT + act
	var result = doParameterReplacement(
		dummySession,
		dummyOriginalPath,
		dummyParameterName,
		dummyParameterType,
	)

	// assert
	assert.Equal(t, dummyOriginalPath, result)
}

func TestDoParameterReplacement_ValidParameterType(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyParameterName = "some name"
	var dummyOriginalPath = "/some/original/path/with/{" + dummyParameterName + "}/in/it"
	var dummyParameterType = ParameterType("some type")
	var dummyResult = "/some/original/path/with/{" + dummyParameterName + ":" + string(dummyParameterType) + "}/in/it"

	// SUT + act
	var result = doParameterReplacement(
		dummySession,
		dummyOriginalPath,
		dummyParameterName,
		dummyParameterType,
	)

	// assert
	assert.Equal(t, dummyResult, result)
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
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(doParameterReplacement, 3, func(session *session, originalPath string, parameterName string, parameterType ParameterType) string {
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
	})

	// SUT + act
	var result = evaluatePathWithParameters(
		dummySession,
		dummyOriginalPath,
		dummyParameters,
	)

	// assert
	assert.Equal(t, dummyUpdatedPath, result)
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

	// SUT + act
	var result = evaluateQueries(
		dummyQueries,
	)

	// assert
	assert.Equal(t, 4, len(result))
	assert.ElementsMatch(t, expectedResult, result)
}

func TestRegisterRoutes_EmptyRoutes(t *testing.T) {
	// arrange
	var dummyApplication = &application{}
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	var dummyRouter = &mux.Router{}
	var dummyRoutes []Route

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "Routes", 1, func() []Route {
		return dummyRoutes
	})
	m.ExpectFunc(logAppRoot, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "register", category)
		assert.Equal(t, "registerRoutes", subcategory)
		assert.Equal(t, "customization.Routes function empty: no routes returned!", messageFormat)
		assert.Equal(t, 0, len(parameters))
	})

	// SUT + act
	registerRoutes(
		dummyApplication,
		dummySession,
		dummyRouter,
	)
}

func TestRegisterRoutes_ValidRoutes(t *testing.T) {
	// arrange
	var dummyApplication = &application{
		actionFuncMap: make(map[string]ActionFunc),
	}
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	var dummyRouter = &mux.Router{}
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
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "Routes", 1, func() []Route {
		return dummyRoutes
	})
	m.ExpectFunc(evaluatePathWithParameters, 2, func(session *session, path string, parameters map[string]ParameterType) string {
		assert.Equal(t, dummySession, session)
		if dummyPath1 == path {
			assert.Equal(t, dummyParameters1, parameters)
			return dummyEvaluatedPath1
		} else if dummyPath2 == path {
			assert.Equal(t, dummyParameters2, parameters)
			return dummyEvaluatedPath2
		}
		return ""
	})
	m.ExpectFunc(evaluateQueries, 2, func(queries map[string]ParameterType) []string {
		if queries["test1"] == ParameterType("me1") {
			return dummyEvaluatedQueries1
		} else if queries["test2"] == ParameterType("me2") {
			return dummyEvaluatedQueries2
		}
		return nil
	})
	m.ExpectFunc(registerRoute, 2, func(router *mux.Router, endpoint string, method string, path string, queries []string, handlerFunc func(http.ResponseWriter, *http.Request), actionFunc ActionFunc) (string, *mux.Route) {
		assert.Equal(t, dummyRouter, router)
		assert.Equal(t, fmt.Sprintf("%v", reflect.ValueOf(dummyApplication.handleSession)), fmt.Sprintf("%v", reflect.ValueOf(handlerFunc)))
		if m.FuncCalledCount(registerRoute) == 1 {
			assert.Equal(t, dummyEndpoint1, endpoint)
			assert.Equal(t, dummyMethod1, method)
			assert.Equal(t, dummyEvaluatedPath1, path)
			assert.Equal(t, dummyEvaluatedQueries1, queries)
			assert.Equal(t, dummyActionFunc1Pointer, fmt.Sprintf("%v", reflect.ValueOf(actionFunc)))
		} else if m.FuncCalledCount(registerRoute) == 2 {
			assert.Equal(t, dummyEndpoint2, endpoint)
			assert.Equal(t, dummyMethod2, method)
			assert.Equal(t, dummyEvaluatedPath2, path)
			assert.Equal(t, dummyEvaluatedQueries2, queries)
			assert.Equal(t, dummyActionFunc2Pointer, fmt.Sprintf("%v", reflect.ValueOf(actionFunc)))
		}
		return "", nil
	})

	// SUT + act
	registerRoutes(
		dummyApplication,
		dummySession,
		dummyRouter,
	)
}

func TestRegisterStatics_EmptyStatics(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	var dummyRouter = &mux.Router{}
	var dummyStatics []Static

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "Statics", 1, func() []Static {
		return dummyStatics
	})
	m.ExpectFunc(logAppRoot, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "register", category)
		assert.Equal(t, "registerStatics", subcategory)
		assert.Equal(t, "customization.Statics function empty: no static content returned!", messageFormat)
		assert.Equal(t, 0, len(parameters))
	})

	// SUT + act
	registerStatics(
		dummySession,
		dummyRouter,
	)
}

type dummyHandler struct {
	http.Handler
}

func TestRegisterStatics_ValidStatics(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	var dummyRouter = &mux.Router{}
	var dummyName1 = "some name 1"
	var dummyPathPrefix1 = "some path prefix 1"
	var dummyHandler1 = &dummyHandler{}
	var dummyName2 = "some name 2"
	var dummyPathPrefix2 = "some path prefix 2"
	var dummyHandler2 = &dummyHandler{}
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
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "Statics", 1, func() []Static {
		return dummyStatics
	})
	m.ExpectFunc(registerStatic, 2, func(router *mux.Router, name string, path string, handler http.Handler) *mux.Route {
		assert.Equal(t, dummyRouter, router)
		if m.FuncCalledCount(registerStatic) == 1 {
			assert.Equal(t, dummyName1, name)
			assert.Equal(t, dummyPathPrefix1, path)
			assert.Equal(t, dummyHandler1, handler)
		} else if m.FuncCalledCount(registerStatic) == 2 {
			assert.Equal(t, dummyName2, name)
			assert.Equal(t, dummyPathPrefix2, path)
			assert.Equal(t, dummyHandler2, handler)
		}
		return nil
	})

	// SUT + act
	registerStatics(
		dummySession,
		dummyRouter,
	)
}

func TestRegisterMiddlewares_EmptyMiddlewares(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	var dummyRouter = &mux.Router{}
	var dummyMiddlewares []MiddlewareFunc

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "Middlewares", 1, func() []MiddlewareFunc {
		return dummyMiddlewares
	})
	m.ExpectFunc(logAppRoot, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "register", category)
		assert.Equal(t, "registerMiddlewares", subcategory)
		assert.Equal(t, "customization.Middlewares function empty: no middleware returned!", messageFormat)
		assert.Equal(t, 0, len(parameters))
	})

	// SUT + act
	registerMiddlewares(
		dummySession,
		dummyRouter,
	)
}

func TestRegisterMiddlewares_ValidMiddlewares(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	var dummyRouter = &mux.Router{}
	var dummyMiddleware1 = func(http.Handler) http.Handler { return nil }
	var dummyMiddleware1Pointer = fmt.Sprintf("%v", reflect.ValueOf(dummyMiddleware1))
	var dummyMiddleware2 = func(http.Handler) http.Handler { return nil }
	var dummyMiddleware2Pointer = fmt.Sprintf("%v", reflect.ValueOf(dummyMiddleware2))
	var dummyMiddlewares = []MiddlewareFunc{
		dummyMiddleware1,
		dummyMiddleware2,
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "Middlewares", 1, func() []MiddlewareFunc {
		return dummyMiddlewares
	})
	m.ExpectFunc(addMiddleware, 2, func(router *mux.Router, middleware MiddlewareFunc) {
		assert.Equal(t, dummyRouter, router)
		var middlewarePointer = fmt.Sprintf("%v", reflect.ValueOf(middleware))
		if m.FuncCalledCount(addMiddleware) == 1 {
			assert.Equal(t, dummyMiddleware1Pointer, middlewarePointer)
		} else if m.FuncCalledCount(addMiddleware) == 2 {
			assert.Equal(t, dummyMiddleware2Pointer, middlewarePointer)
		}
	})

	// SUT + act
	registerMiddlewares(
		dummySession,
		dummyRouter,
	)
}

func TestRegisterErrorHandlers(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummyRouter = &mux.Router{}
	var dummyMethodNotAllowedHandler = &dummyHandler{}
	var dummyNotFoundHandler = &dummyHandler{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "MethodNotAllowedHandler", 1, func() http.Handler {
		return dummyMethodNotAllowedHandler
	})
	m.ExpectMethod(dummyCustomization, "NotFoundHandler", 1, func() http.Handler {
		return dummyNotFoundHandler
	})

	// SUT + act
	registerErrorHandlers(
		dummyCustomization,
		dummyRouter,
	)

	// assert
	assert.Equal(t, dummyMethodNotAllowedHandler, dummyRouter.MethodNotAllowedHandler)
	assert.Equal(t, dummyNotFoundHandler, dummyRouter.NotFoundHandler)
}

func TestInstantiateRouter_RouterError(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummyApplication = &application{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	var dummyRouter = &mux.Router{KeepContext: rand.Intn(100) > 50}
	var dummyError = errors.New("some error")
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(mux.NewRouter, 1, func() *mux.Router {
		return dummyRouter
	})
	m.ExpectFunc(registerRoutes, 1, func(app *application, session *session, router *mux.Router) {
		assert.Equal(t, dummyApplication, app)
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
	})
	m.ExpectFunc(registerStatics, 1, func(session *session, router *mux.Router) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
	})
	m.ExpectFunc(registerMiddlewares, 1, func(session *session, router *mux.Router) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
	})
	m.ExpectFunc(walkRegisteredRoutes, 1, func(session *session, router *mux.Router) error {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
		return dummyError
	})
	m.ExpectFunc(logAppRoot, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "register", category)
		assert.Equal(t, "instantiateRouter", subcategory)
		assert.Equal(t, "%+v", messageFormat)
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
	var result, err = instantiateRouter(
		dummyApplication,
		dummySession,
	)

	// assert
	assert.Equal(t, dummyRouter, result)
	assert.Equal(t, dummyAppError, err)
}

func TestInstantiateRouter_HappyPath(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummyApplication = &application{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	var dummyRouter = &mux.Router{KeepContext: rand.Intn(100) > 50}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(mux.NewRouter, 1, func() *mux.Router {
		return dummyRouter
	})
	m.ExpectFunc(registerRoutes, 1, func(app *application, session *session, router *mux.Router) {
		assert.Equal(t, dummyApplication, app)
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
	})
	m.ExpectFunc(registerStatics, 1, func(session *session, router *mux.Router) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
	})
	m.ExpectFunc(registerMiddlewares, 1, func(session *session, router *mux.Router) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
	})
	m.ExpectFunc(walkRegisteredRoutes, 1, func(session *session, router *mux.Router) error {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
		return nil
	})
	m.ExpectFunc(registerErrorHandlers, 1, func(customization Customization, router *mux.Router) {
		assert.Equal(t, dummyCustomization, customization)
		assert.Equal(t, dummyRouter, router)
	})
	m.ExpectMethod(dummyCustomization, "InstrumentRouter", 1, func(self *DefaultCustomization, router *mux.Router) *mux.Router {
		assert.Equal(t, dummyCustomization, self)
		assert.Equal(t, dummyRouter, router)
		return dummyRouter
	})

	// SUT + act
	var result, err = instantiateRouter(
		dummyApplication,
		dummySession,
	)

	// assert
	assert.Equal(t, dummyRouter, result)
	assert.NoError(t, err)
}
