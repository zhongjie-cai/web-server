package webserver

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestGetRequestedPort_NilRequest(t *testing.T) {
	// arrange
	var dummyHTTPRequest *http.Request

	// mock
	createMock(t)

	// SUT + act
	var result = getRequestedPort(
		dummyHTTPRequest,
	)

	// assert
	assert.Zero(t, result)

	// verify
	verifyAll(t)
}

func TestGetRequestedPort_NoPort(t *testing.T) {
	// arrange
	var dummyHost = "some host"
	var dummyHTTPRequest = &http.Request{
		Host: dummyHost,
	}

	// mock
	createMock(t)

	// expect
	stringsSplitExpected = 1
	stringsSplit = func(s, sep string) []string {
		stringsSplitCalled++
		assert.Equal(t, dummyHost, s)
		assert.Equal(t, ":", sep)
		return strings.Split(s, sep)
	}

	// SUT + act
	var result = getRequestedPort(
		dummyHTTPRequest,
	)

	// assert
	assert.Zero(t, result)

	// verify
	verifyAll(t)
}

func TestGetRequestedPort_InvalidPort(t *testing.T) {
	// arrange
	var dummyPort = "some port"
	var dummyHost = "some host:" + dummyPort
	var dummyHTTPRequest = &http.Request{
		Host: dummyHost,
	}

	// mock
	createMock(t)

	// expect
	stringsSplitExpected = 1
	stringsSplit = func(s, sep string) []string {
		stringsSplitCalled++
		assert.Equal(t, dummyHost, s)
		assert.Equal(t, ":", sep)
		return strings.Split(s, sep)
	}
	strconvAtoiExpected = 1
	strconvAtoi = func(s string) (int, error) {
		strconvAtoiCalled++
		assert.Equal(t, dummyPort, s)
		return strconv.Atoi(s)
	}

	// SUT + act
	var result = getRequestedPort(
		dummyHTTPRequest,
	)

	// assert
	assert.Zero(t, result)

	// verify
	verifyAll(t)
}

func TestGetRequestedPort_Success(t *testing.T) {
	// arrange
	var dummyPortValue = rand.Intn(65536)
	var dummyPortString = strconv.Itoa(dummyPortValue)
	var dummyHost = "some host:" + dummyPortString
	var dummyHTTPRequest = &http.Request{
		Host: dummyHost,
	}

	// mock
	createMock(t)

	// expect
	stringsSplitExpected = 1
	stringsSplit = func(s, sep string) []string {
		stringsSplitCalled++
		assert.Equal(t, dummyHost, s)
		assert.Equal(t, ":", sep)
		return strings.Split(s, sep)
	}
	strconvAtoiExpected = 1
	strconvAtoi = func(s string) (int, error) {
		strconvAtoiCalled++
		assert.Equal(t, dummyPortString, s)
		return strconv.Atoi(s)
	}

	// SUT + act
	var result = getRequestedPort(
		dummyHTTPRequest,
	)

	// assert
	assert.Equal(t, dummyPortValue, result)

	// verify
	verifyAll(t)
}

func TestInitiateSession(t *testing.T) {
	// arrange
	var dummyResponseWriter = &dummyResponseWriter{t: t}
	var dummyHTTPRequest = &http.Request{Host: "some host"}
	var dummyPort = rand.Intn(65536)
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
	getRequestedPortFuncExpected = 1
	getRequestedPortFunc = func(httpRequest *http.Request) int {
		getRequestedPortFuncCalled++
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		return dummyPort
	}
	getApplicationFuncExpected = 1
	getApplicationFunc = func(port int) *application {
		getApplicationFuncCalled++
		assert.Equal(t, dummyPort, port)
		return dummyApplication
	}
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
	var dummyDuration = time.Duration(rand.Intn(100))

	// mock
	createMock(t)

	// expect
	handlePanicFuncExpected = 1
	handlePanicFunc = func(session *session, recoverResult interface{}) {
		handlePanicFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, recover(), recoverResult)
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
	initiateSessionFunc = func(responseWriter http.ResponseWriter, httpRequest *http.Request) (*session, ActionFunc, error) {
		initiateSessionFuncCalled++
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
	finalizeSessionFunc = func(session *session, startTime time.Time) {
		finalizeSessionFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyStartTime, startTime)
	}
	writeResponseFuncExpected = 1
	writeResponseFunc = func(session *session, responseObject interface{}, responseError error) {
		writeResponseFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Nil(t, responseObject)
		assert.Equal(t, dummyRouteError, responseError)
	}

	// SUT + act
	handleSession(
		dummyResponseWriter,
		dummyHTTPRequest,
	)

	// verify
	verifyAll(t)
}

func TestHandleSession_Success(t *testing.T) {
	// arrange
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
	initiateSessionFunc = func(responseWriter http.ResponseWriter, httpRequest *http.Request) (*session, ActionFunc, error) {
		initiateSessionFuncCalled++
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
	finalizeSessionFunc = func(session *session, startTime time.Time) {
		finalizeSessionFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyStartTime, startTime)
	}
	handleActionFuncExpected = 1
	handleActionFunc = func(session *session, action ActionFunc) {
		handleActionFuncCalled++
		assert.Equal(t, dummySession, session)
		functionPointerEquals(t, dummyAction, action)
	}

	// SUT + act
	handleSession(
		dummyResponseWriter,
		dummyHTTPRequest,
	)

	// verify
	verifyAll(t)
}
