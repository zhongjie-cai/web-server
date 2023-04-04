package webserver

import (
	"errors"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInitiateSession(t *testing.T) {
	// arrange
	var dummyResponseWriter = &dummyResponseWriter{t: t}
	var dummyHTTPRequest = &http.Request{Host: "some host"}
	var dummyCustomization = &dummyCustomization{t: t}
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
	createMock(t)

	// expect
	uuidNewExpected = 1
	uuidNew = func() uuid.UUID {
		uuidNewCalled++
		return dummySessionID
	}
	getRouteInfoFuncExpected = 1
	getRouteInfoFunc = func(httpRequest *http.Request, actionFuncMap map[string]ActionFunc) (string, ActionFunc, error) {
		getRouteInfoFuncCalled++
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		assert.Equal(t, dummyActionFuncMap, actionFuncMap)
		return dummyEndpoint, dummyAction, dummyRouteError
	}

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

	// verify
	verifyAll(t)
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
	var dummyStartTime = time.Now().UTC()
	var dummyRecoverResult = "some recover result"
	var dummyDuration = time.Duration(rand.Intn(100))

	// mock
	createMock(t)

	// expect
	handlePanicFuncExpected = 1
	handlePanicFunc = func(session *session, recoverResult interface{}) {
		handlePanicFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRecoverResult, recoverResult)
	}
	logEndpointExitFuncExpected = 1
	logEndpointExitFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logEndpointExitFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyName, category)
		assert.Equal(t, dummyMethod, subcategory)
		assert.Equal(t, "%s", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyDuration, parameters[0])
	}
	timeSinceExpected = 1
	timeSince = func(tm time.Time) time.Duration {
		timeSinceCalled++
		assert.Equal(t, dummyStartTime, tm)
		return dummyDuration
	}

	// SUT + act
	finalizeSession(
		dummySession,
		dummyStartTime,
		dummyRecoverResult,
	)

	// verify
	verifyAll(t)
}

type dummyCustomizationHandleAction struct {
	dummyCustomization
	preAction  func(session Session) error
	postAction func(session Session) error
}

func (dummySession *dummyCustomizationHandleAction) PreAction(session Session) error {
	if dummySession.preAction != nil {
		return dummySession.preAction(session)
	}
	assert.Fail(dummySession.t, "Unexpected call to PreAction")
	return nil
}

func (dummySession *dummyCustomizationHandleAction) PostAction(session Session) error {
	if dummySession.postAction != nil {
		return dummySession.postAction(session)
	}
	assert.Fail(dummySession.t, "Unexpected call to PostAction")
	return nil
}

func TestHandleAction_PreActionError(t *testing.T) {
	// arrange
	var dummyCustomizationHandleAction = &dummyCustomizationHandleAction{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		customization: dummyCustomizationHandleAction,
	}
	var customizationPreActionExpected int
	var customizationPreActionCalled int
	var dummyActionExpected int
	var dummyActionCalled int
	var dummyAction ActionFunc
	var customizationPostActionExpected int
	var customizationPostActionCalled int
	var dummyError = errors.New("some error")

	// mock
	createMock(t)

	// expect
	customizationPreActionExpected = 1
	dummyCustomizationHandleAction.preAction = func(session Session) error {
		customizationPreActionCalled++
		assert.Equal(t, dummySession, session)
		return dummyError
	}
	writeResponseFuncExpected = 1
	writeResponseFunc = func(session *session, responseObject interface{}, responseError error) {
		writeResponseFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Nil(t, responseObject)
		assert.Equal(t, dummyError, responseError)
	}

	// SUT + act
	handleAction(
		dummySession,
		dummyAction,
	)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationPreActionExpected, customizationPreActionCalled, "Unexpected number of calls to method customization.PreAction")
	assert.Equal(t, dummyActionExpected, dummyActionCalled, "Unexpected number of calls to method dummyAction")
	assert.Equal(t, customizationPostActionExpected, customizationPostActionCalled, "Unexpected number of calls to method customization.PostAction")
}

func TestHandleAction_ResponseError(t *testing.T) {
	// arrange
	var dummyCustomizationHandleAction = &dummyCustomizationHandleAction{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		customization: dummyCustomizationHandleAction,
	}
	var customizationPreActionExpected int
	var customizationPreActionCalled int
	var dummyActionExpected int
	var dummyActionCalled int
	var dummyAction ActionFunc
	var customizationPostActionExpected int
	var customizationPostActionCalled int
	var dummyResponseObject = rand.Int()
	var dummyError = errors.New("some error")

	// mock
	createMock(t)

	// expect
	customizationPreActionExpected = 1
	dummyCustomizationHandleAction.preAction = func(session Session) error {
		customizationPreActionCalled++
		assert.Equal(t, dummySession, session)
		return nil
	}
	dummyActionExpected = 1
	dummyAction = func(session Session) (interface{}, error) {
		dummyActionCalled++
		assert.Equal(t, dummySession, session)
		return dummyResponseObject, dummyError
	}
	writeResponseFuncExpected = 1
	writeResponseFunc = func(session *session, responseObject interface{}, responseError error) {
		writeResponseFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyResponseObject, responseObject)
		assert.Equal(t, dummyError, responseError)
	}

	// SUT + act
	handleAction(
		dummySession,
		dummyAction,
	)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationPreActionExpected, customizationPreActionCalled, "Unexpected number of calls to method customization.PreAction")
	assert.Equal(t, dummyActionExpected, dummyActionCalled, "Unexpected number of calls to method dummyAction")
	assert.Equal(t, customizationPostActionExpected, customizationPostActionCalled, "Unexpected number of calls to method customization.PostAction")
}

func TestHandleAction_HappyPath(t *testing.T) {
	// arrange
	var dummyCustomizationHandleAction = &dummyCustomizationHandleAction{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		customization: dummyCustomizationHandleAction,
	}
	var customizationPreActionExpected int
	var customizationPreActionCalled int
	var dummyActionExpected int
	var dummyActionCalled int
	var dummyAction ActionFunc
	var customizationPostActionExpected int
	var customizationPostActionCalled int
	var dummyResponseObject = rand.Int()
	var dummyError = errors.New("some error")

	// mock
	createMock(t)

	// expect
	customizationPreActionExpected = 1
	dummyCustomizationHandleAction.preAction = func(session Session) error {
		customizationPreActionCalled++
		assert.Equal(t, dummySession, session)
		return nil
	}
	dummyActionExpected = 1
	dummyAction = func(session Session) (interface{}, error) {
		dummyActionCalled++
		assert.Equal(t, dummySession, session)
		return dummyResponseObject, nil
	}
	customizationPostActionExpected = 1
	dummyCustomizationHandleAction.postAction = func(session Session) error {
		customizationPostActionCalled++
		assert.Equal(t, dummySession, session)
		return dummyError
	}
	writeResponseFuncExpected = 1
	writeResponseFunc = func(session *session, responseObject interface{}, responseError error) {
		writeResponseFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyResponseObject, responseObject)
		assert.Equal(t, dummyError, responseError)
	}

	// SUT + act
	handleAction(
		dummySession,
		dummyAction,
	)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationPreActionExpected, customizationPreActionCalled, "Unexpected number of calls to method customization.PreAction")
	assert.Equal(t, dummyActionExpected, dummyActionCalled, "Unexpected number of calls to method dummyAction")
	assert.Equal(t, customizationPostActionExpected, customizationPostActionCalled, "Unexpected number of calls to method customization.PostAction")
}

func TestHandleSession_RouteError(t *testing.T) {
	// arrange
	var dummyApplication = &application{}
	var dummyResponseWriter = &dummyResponseWriter{t: t}
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
	createMock(t)

	// expect
	initiateSessionFuncExpected = 1
	initiateSessionFunc = func(app *application, responseWriter http.ResponseWriter, httpRequest *http.Request) (*session, ActionFunc, error) {
		initiateSessionFuncCalled++
		assert.Equal(t, dummyApplication, app)
		assert.Equal(t, dummyResponseWriter, responseWriter)
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		return dummySession, dummyAction, dummyRouteError
	}
	logEndpointEnterFuncExpected = 1
	logEndpointEnterFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logEndpointEnterFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyName, category)
		assert.Equal(t, dummyMethod, subcategory)
		assert.Zero(t, messageFormat)
		assert.Empty(t, parameters)
	}
	getTimeNowUTCFuncExpected = 1
	getTimeNowUTCFunc = func() time.Time {
		getTimeNowUTCFuncCalled++
		return dummyStartTime
	}
	finalizeSessionFuncExpected = 1
	finalizeSessionFunc = func(session *session, startTime time.Time, recoverResult interface{}) {
		finalizeSessionFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyStartTime, startTime)
		assert.Equal(t, recover(), recoverResult)
	}
	writeResponseFuncExpected = 1
	writeResponseFunc = func(session *session, responseObject interface{}, responseError error) {
		writeResponseFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Nil(t, responseObject)
		assert.Equal(t, dummyRouteError, responseError)
	}

	// SUT + act
	dummyApplication.handleSession(
		dummyResponseWriter,
		dummyHTTPRequest,
	)

	// verify
	verifyAll(t)
}

func TestHandleSession_Success(t *testing.T) {
	// arrange
	var dummyApplication = &application{}
	var dummyResponseWriter = &dummyResponseWriter{t: t}
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
	createMock(t)

	// expect
	initiateSessionFuncExpected = 1
	initiateSessionFunc = func(app *application, responseWriter http.ResponseWriter, httpRequest *http.Request) (*session, ActionFunc, error) {
		initiateSessionFuncCalled++
		assert.Equal(t, dummyApplication, app)
		assert.Equal(t, dummyResponseWriter, responseWriter)
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		return dummySession, dummyAction, nil
	}
	logEndpointEnterFuncExpected = 1
	logEndpointEnterFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logEndpointEnterFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyName, category)
		assert.Equal(t, dummyMethod, subcategory)
		assert.Zero(t, messageFormat)
		assert.Empty(t, parameters)
	}
	getTimeNowUTCFuncExpected = 1
	getTimeNowUTCFunc = func() time.Time {
		getTimeNowUTCFuncCalled++
		return dummyStartTime
	}
	finalizeSessionFuncExpected = 1
	finalizeSessionFunc = func(session *session, startTime time.Time, recoverResult interface{}) {
		finalizeSessionFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyStartTime, startTime)
		assert.Equal(t, recover(), recoverResult)
	}
	handleActionFuncExpected = 1
	handleActionFunc = func(session *session, action ActionFunc) {
		handleActionFuncCalled++
		assert.Equal(t, dummySession, session)
		functionPointerEquals(t, dummyAction, action)
	}

	// SUT + act
	dummyApplication.handleSession(
		dummyResponseWriter,
		dummyHTTPRequest,
	)

	// verify
	verifyAll(t)
}
