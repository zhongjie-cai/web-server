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

func TestDoParameterReplacement_EmptyParameterType(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyParameterName = "some name"
	var dummyOriginalPath = "/some/original/path/with/{" + dummyParameterName + "}/in/it"
	var dummyParameterType ParameterType

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(logAppRoot).Expects(dummySession, "register", "doParameterReplacement",
		"Path parameter [%v] in path [%v] has no type specification; fallback to default.",
		dummyParameterName, dummyOriginalPath).Returns().Once()

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
	m.Mock(doParameterReplacement).Expects(dummySession, dummyOriginalPath, gomocker.Anything(), gomocker.Anything()).Returns(dummyUpdatedPath).Once()
	m.Mock(doParameterReplacement).Expects(dummySession, dummyUpdatedPath, gomocker.Anything(), gomocker.Anything()).Returns(dummyUpdatedPath).Twice()

	// SUT + act
	var result = evaluatePathWithParameters(
		dummySession,
		dummyOriginalPath,
		dummyParameters,
	)

	// assert
	assert.Equal(t, dummyUpdatedPath, result)
}

func TestRegisterRoutes_EmptyRoutes(t *testing.T) {
	// arrange
	var dummyApplication = &application{}
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	type router struct {
		chi.Router
	}
	var dummyRouter = &router{}
	var dummyRoutes []Route

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock((*DefaultCustomization).Routes).Expects(dummyCustomization).Returns(dummyRoutes).Once()
	m.Mock(logAppRoot).Expects(dummySession, "register", "registerRoutes",
		"customization.Routes function empty: no routes returned!").Returns().Once()

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
	type router struct {
		chi.Router
	}
	var dummyRouter = &router{}
	var dummyMethod1 = "some method 1"
	var dummyPath1 = "some path 1"
	var dummyParameters1 = map[string]ParameterType{
		"foo1": ParameterType("bar1"),
	}
	var dummyActionFunc1 = func(Session) (any, error) {
		return nil, nil
	}
	var dummyMethod2 = "some method 2"
	var dummyPath2 = "some path 2"
	var dummyParameters2 = map[string]ParameterType{
		"foo2": ParameterType("bar2"),
	}
	var dummyActionFunc2 = func(Session) (any, error) {
		return nil, nil
	}
	var dummyRoutes = []Route{
		{
			Method:     dummyMethod1,
			Path:       dummyPath1,
			Parameters: dummyParameters1,
			ActionFunc: dummyActionFunc1,
		},
		{
			Method:     dummyMethod2,
			Path:       dummyPath2,
			Parameters: dummyParameters2,
			ActionFunc: dummyActionFunc2,
		},
	}
	var dummyEvaluatedPath1 = "some evaluated path 1"
	var dummyEvaluatedPath2 = "some evaluated path 2"
	var dummyName1 = "some name 1"
	var dummyName2 = "some name 2"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock((*DefaultCustomization).Routes).Expects(dummyCustomization).Returns(dummyRoutes).Once()
	m.Mock(evaluatePathWithParameters).Expects(dummySession, dummyPath1, dummyParameters1).Returns(dummyEvaluatedPath1).Once()
	m.Mock(generateRouteName).Expects(dummyMethod1, dummyEvaluatedPath1).Returns(dummyName1).Once()
	m.Mock((*router).MethodFunc).Expects(dummyRouter, dummyMethod1, dummyEvaluatedPath1,
		gomocker.Matches(func(value any) bool { return functionPointerEquals(dummyApplication.handleSession, value) })).Returns().Once()
	m.Mock(evaluatePathWithParameters).Expects(dummySession, dummyPath2, dummyParameters2).Returns(dummyEvaluatedPath2).Once()
	m.Mock(generateRouteName).Expects(dummyMethod2, dummyEvaluatedPath2).Returns(dummyName2).Once()
	m.Mock((*router).MethodFunc).Expects(dummyRouter, dummyMethod2, dummyEvaluatedPath2,
		gomocker.Matches(func(value any) bool { return functionPointerEquals(dummyApplication.handleSession, value) })).Returns().Once()

	// SUT + act
	registerRoutes(
		dummyApplication,
		dummySession,
		dummyRouter,
	)

	// assert
	assert.Contains(t, dummyApplication.actionFuncMap, dummyName1)
	assert.Contains(t, dummyApplication.actionFuncMap, dummyName2)
}

func TestRegisterStatics_EmptyStatics(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	type router struct {
		chi.Router
	}
	var dummyRouter = &router{}
	var dummyStatics []Static

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock((*DefaultCustomization).Statics).Expects(dummyCustomization).Returns(dummyStatics).Once()
	m.Mock(logAppRoot).Expects(dummySession, "register", "registerStatics",
		"customization.Statics function empty: no static content returned!").Returns().Once()

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
	type router struct {
		chi.Router
	}
	var dummyRouter = &router{}
	var dummyPathPrefix1 = "some path prefix 1"
	var dummyHandler1 = &dummyHandler{}
	var dummyPathPrefix2 = "some path prefix 2"
	var dummyHandler2 = &dummyHandler{}
	var dummyStatics = []Static{
		{
			PathPrefix: dummyPathPrefix1,
			Handler:    dummyHandler1,
		},
		{
			PathPrefix: dummyPathPrefix2,
			Handler:    dummyHandler2,
		},
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock((*DefaultCustomization).Statics).Expects(dummyCustomization).Returns(dummyStatics).Once()
	m.Mock((*router).Handle).Expects(dummyRouter, dummyPathPrefix1, gomocker.Matches(func(value any) bool {
		return functionPointerEquals(dummyHandler1, value)
	})).Returns().Once()
	m.Mock((*router).Handle).Expects(dummyRouter, dummyPathPrefix2, gomocker.Matches(func(value any) bool {
		return functionPointerEquals(dummyHandler2, value)
	})).Returns().Once()

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
	type router struct {
		chi.Router
	}
	var dummyRouter = &router{}
	var dummyMiddlewares []func(http.Handler) http.Handler

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock((*DefaultCustomization).Middlewares).Expects(dummyCustomization).Returns(dummyMiddlewares).Once()
	m.Mock(logAppRoot).Expects(dummySession, "register", "registerMiddlewares",
		"customization.Middlewares function empty: no middleware returned!").Returns().Once()

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
	type router struct {
		chi.Router
	}
	var dummyRouter = &router{}
	var dummyMiddleware1 = func(http.Handler) http.Handler { return nil }
	var dummyMiddleware2 = func(http.Handler) http.Handler { return nil }
	var dummyMiddlewares = []func(http.Handler) http.Handler{
		dummyMiddleware1,
		dummyMiddleware2,
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock((*DefaultCustomization).Middlewares).Expects(dummyCustomization).Returns(dummyMiddlewares).Once()
	m.Mock((*router).Use).Expects(dummyRouter, gomocker.Matches(func(value any) bool {
		return functionPointerEquals(dummyMiddleware1, value)
	})).Returns().Once()
	m.Mock((*router).Use).Expects(dummyRouter, gomocker.Matches(func(value any) bool {
		return functionPointerEquals(dummyMiddleware2, value)
	})).Returns().Once()

	// SUT + act
	registerMiddlewares(
		dummySession,
		dummyRouter,
	)
}

func TestRegisterErrorHandlers_NotNils(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	type router struct {
		chi.Router
	}
	var dummyRouter = &router{}
	var dummyMethodNotAllowedHandler = func(http.ResponseWriter, *http.Request) {}
	var dummyNotFoundHandler = func(http.ResponseWriter, *http.Request) {}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock((*DefaultCustomization).MethodNotAllowedHandler).Expects(dummyCustomization).Returns(dummyMethodNotAllowedHandler).Once()
	m.Mock((*router).MethodNotAllowed).Expects(dummyRouter, gomocker.Matches(func(value any) bool {
		return functionPointerEquals(dummyMethodNotAllowedHandler, value)
	})).Returns().Once()
	m.Mock((*DefaultCustomization).NotFoundHandler).Expects(dummyCustomization).Returns(dummyNotFoundHandler).Once()
	m.Mock((*router).NotFound).Expects(dummyRouter, gomocker.Matches(func(value any) bool {
		return functionPointerEquals(dummyNotFoundHandler, value)
	})).Returns().Once()

	// SUT + act
	registerErrorHandlers(
		dummyCustomization,
		dummyRouter,
	)
}

func TestInstantiateRouter_RouterError(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummyApplication = &application{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	var dummyRouter = &chi.Mux{}
	var dummyError = errors.New("some error")
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(chi.NewRouter).Expects().Returns(dummyRouter).Once()
	m.Mock(registerRoutes).Expects(dummyApplication, dummySession, dummyRouter).Returns().Once()
	m.Mock(registerStatics).Expects(dummySession, dummyRouter).Returns().Once()
	m.Mock(registerMiddlewares).Expects(dummySession, dummyRouter).Returns().Once()
	m.Mock(walkRegisteredRoutes).Expects(dummySession, dummyRouter).Returns(dummyError).Once()
	m.Mock(logAppRoot).Expects(dummySession, "register", "instantiateRouter", "%+v", dummyError).Returns().Once()
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageRouteRegistration, dummyError).Returns(dummyAppError).Once()

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
	var dummyRouter = &chi.Mux{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(chi.NewRouter).Expects().Returns(dummyRouter).Once()
	m.Mock(registerRoutes).Expects(dummyApplication, dummySession, dummyRouter).Returns().Once()
	m.Mock(registerStatics).Expects(dummySession, dummyRouter).Returns().Once()
	m.Mock(registerMiddlewares).Expects(dummySession, dummyRouter).Returns().Once()
	m.Mock(walkRegisteredRoutes).Expects(dummySession, dummyRouter).Returns(nil).Once()
	m.Mock(registerErrorHandlers).Expects(dummyCustomization, dummyRouter).Returns().Once()
	m.Mock((*DefaultCustomization).InstrumentRouter).Expects(dummyCustomization, dummyRouter).Returns(dummyRouter).Once()

	// SUT + act
	var result, err = instantiateRouter(
		dummyApplication,
		dummySession,
	)

	// assert
	assert.Equal(t, dummyRouter, result)
	assert.NoError(t, err)
}
