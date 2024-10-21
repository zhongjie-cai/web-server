package webserver

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zhongjie-cai/gomocker"
)

func functionPointerEquals(t *testing.T, expectFunc interface{}, actualFunc interface{}) {
	var expectValue = fmt.Sprintf("%v", reflect.ValueOf(expectFunc))
	var actualValue = fmt.Sprintf("%v", reflect.ValueOf(actualFunc))
	assert.Equal(t, expectValue, actualValue)
}

type dummyResponseWriter struct {
	http.ResponseWriter
}

func TestInitiateSession(t *testing.T) {
	// arrange
	var dummyResponseWriter = &dummyResponseWriter{}
	var dummyHTTPRequest = &http.Request{Host: "some host"}
	var dummyCustomization = &DefaultCustomization{}
	var dummyAction = func(session Session) (interface{}, error) { return nil, nil }
	var dummyActionFuncMap = map[string]ActionFunc{
		"some key": dummyAction,
	}
	var dummyApplication = &application{
		customization: dummyCustomization,
		actionFuncMap: dummyActionFuncMap,
	}
	var dummyEndpoint = "some endpoint"
	var dummyRouteError = errors.New("some route error")
	var dummySessionID = uuid.New()

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(uuid.New, 1, func() uuid.UUID {
		return dummySessionID
	})
	m.ExpectFunc(getRouteInfo, 1, func(httpRequest *http.Request, actionFuncMap map[string]ActionFunc) (string, ActionFunc, error) {
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		assert.Equal(t, dummyActionFuncMap, actionFuncMap)
		return dummyEndpoint, dummyAction, dummyRouteError
	})

	// SUT + act
	var session, action, err = initiateSession(
		dummyApplication,
		dummyResponseWriter,
		dummyHTTPRequest,
	)

	// assert
	assert.NotNil(t, session)
	assert.Equal(t, dummySessionID, session.id)
	assert.Equal(t, dummyEndpoint, session.name)
	assert.Equal(t, dummyHTTPRequest, session.request)
	assert.Equal(t, dummyResponseWriter, session.responseWriter)
	assert.Empty(t, session.attachment)
	assert.Equal(t, dummyCustomization, session.customization)
	functionPointerEquals(t, dummyAction, action)
	assert.Equal(t, dummyRouteError, err)
}

func TestFinalizeSession(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyMethod = "some method"
	var dummyHTTPRequest = &http.Request{
		Method: dummyMethod,
	}
	var dummySession = &session{
		name:    dummyName,
		request: dummyHTTPRequest,
	}
	var dummyStartTime = getTimeNowUTC()
	var dummyRecoverResult = "some recover result"
	var dummyDuration = time.Duration(rand.Intn(100))

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(handlePanic, 1, func(session *session, recoverResult interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRecoverResult, recoverResult)
	})
	m.ExpectFunc(logEndpointExit, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyName, category)
		assert.Equal(t, dummyMethod, subcategory)
		assert.Equal(t, "%s", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyDuration, parameters[0])
	})
	m.ExpectFunc(time.Since, 1, func(tm time.Time) time.Duration {
		assert.Equal(t, dummyStartTime, tm)
		return dummyDuration
	})

	// SUT + act
	finalizeSession(
		dummySession,
		dummyStartTime,
		dummyRecoverResult,
	)
}

func dummyAction(session Session) (interface{}, error) {
	return nil, nil
}

func TestHandleAction_PreActionError(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	var dummyError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "PreAction", 1, func(self *DefaultCustomization, session Session) error {
		assert.Equal(t, dummyCustomization, self)
		assert.Equal(t, dummySession, session)
		return dummyError
	})
	m.ExpectFunc(writeResponse, 1, func(session *session, responseObject interface{}, responseError error) {
		assert.Equal(t, dummySession, session)
		assert.Nil(t, responseObject)
		assert.Equal(t, dummyError, responseError)
	})

	// SUT + act
	handleAction(
		dummySession,
		dummyAction,
	)
}

func TestHandleAction_ResponseError(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	var dummyResponseObject = rand.Int()
	var dummyError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "PreAction", 1, func(self *DefaultCustomization, session Session) error {
		assert.Equal(t, dummyCustomization, self)
		assert.Equal(t, dummySession, session)
		return nil
	})
	m.ExpectFunc(dummyAction, 1, func(session Session) (interface{}, error) {
		assert.Equal(t, dummySession, session)
		return dummyResponseObject, dummyError
	})
	m.ExpectFunc(writeResponse, 1, func(session *session, responseObject interface{}, responseError error) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyResponseObject, responseObject)
		assert.Equal(t, dummyError, responseError)
	})

	// SUT + act
	handleAction(
		dummySession,
		dummyAction,
	)
}

func TestHandleAction_HappyPath(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	var dummyResponseObject = rand.Int()
	var dummyError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "PreAction", 1, func(self *DefaultCustomization, session Session) error {
		assert.Equal(t, dummyCustomization, self)
		assert.Equal(t, dummySession, session)
		return nil
	})
	m.ExpectFunc(dummyAction, 1, func(session Session) (interface{}, error) {
		assert.Equal(t, dummySession, session)
		return dummyResponseObject, nil
	})
	m.ExpectMethod(dummyCustomization, "PostAction", 1, func(self *DefaultCustomization, session Session) error {
		assert.Equal(t, dummyCustomization, self)
		assert.Equal(t, dummySession, session)
		return dummyError
	})
	m.ExpectFunc(writeResponse, 1, func(session *session, responseObject interface{}, responseError error) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyResponseObject, responseObject)
		assert.Equal(t, dummyError, responseError)
	})

	// SUT + act
	handleAction(
		dummySession,
		dummyAction,
	)
}

func TestHandleSession_RouteError(t *testing.T) {
	// arrange
	var dummyApplication = &application{}
	var dummyResponseWriter = &dummyResponseWriter{}
	var dummyMethod = "some method"
	var dummyHTTPRequest = &http.Request{
		Method: dummyMethod,
	}
	var dummyName = "some name"
	var dummySession = &session{
		name: dummyName,
	}
	var dummyAction = func(session Session) (interface{}, error) { return nil, nil }
	var dummyRouteError = errors.New("some route error")
	var dummyStartTime = time.Now()

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(initiateSession, 1, func(app *application, responseWriter http.ResponseWriter, httpRequest *http.Request) (*session, ActionFunc, error) {
		assert.Equal(t, dummyApplication, app)
		assert.Equal(t, dummyResponseWriter, responseWriter)
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		return dummySession, dummyAction, dummyRouteError
	})
	m.ExpectFunc(logEndpointEnter, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyName, category)
		assert.Equal(t, dummyMethod, subcategory)
		assert.Zero(t, messageFormat)
		assert.Empty(t, parameters)
	})
	m.ExpectFunc(getTimeNowUTC, 1, func() time.Time {
		return dummyStartTime
	})
	m.ExpectFunc(finalizeSession, 1, func(session *session, startTime time.Time, recoverResult interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyStartTime, startTime)
		assert.Equal(t, recover(), recoverResult)
	})
	m.ExpectFunc(writeResponse, 1, func(session *session, responseObject interface{}, responseError error) {
		assert.Equal(t, dummySession, session)
		assert.Nil(t, responseObject)
		assert.Equal(t, dummyRouteError, responseError)
	})

	// SUT + act
	dummyApplication.handleSession(
		dummyResponseWriter,
		dummyHTTPRequest,
	)
}

func TestHandleSession_Success(t *testing.T) {
	// arrange
	var dummyApplication = &application{}
	var dummyResponseWriter = &dummyResponseWriter{}
	var dummyMethod = "some method"
	var dummyHTTPRequest = &http.Request{
		Method: dummyMethod,
	}
	var dummyName = "some name"
	var dummySession = &session{
		name: dummyName,
	}
	var dummyAction = func(session Session) (interface{}, error) { return nil, nil }
	var dummyStartTime = time.Now()

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(initiateSession, 1, func(app *application, responseWriter http.ResponseWriter, httpRequest *http.Request) (*session, ActionFunc, error) {
		assert.Equal(t, dummyApplication, app)
		assert.Equal(t, dummyResponseWriter, responseWriter)
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		return dummySession, dummyAction, nil
	})
	m.ExpectFunc(logEndpointEnter, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyName, category)
		assert.Equal(t, dummyMethod, subcategory)
		assert.Zero(t, messageFormat)
		assert.Empty(t, parameters)
	})
	m.ExpectFunc(getTimeNowUTC, 1, func() time.Time {
		return dummyStartTime
	})
	m.ExpectFunc(finalizeSession, 1, func(session *session, startTime time.Time, recoverResult interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyStartTime, startTime)
		assert.Equal(t, recover(), recoverResult)
	})
	m.ExpectFunc(handleAction, 1, func(session *session, action ActionFunc) {
		assert.Equal(t, dummySession, session)
		functionPointerEquals(t, dummyAction, action)
	})

	// SUT + act
	dummyApplication.handleSession(
		dummyResponseWriter,
		dummyHTTPRequest,
	)
}
