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
	"github.com/zhongjie-cai/gomocker/v2"
)

func functionPointerEquals(expectFunc any, actualFunc any) bool {
	var expectValue = fmt.Sprintf("%v", reflect.ValueOf(expectFunc))
	var actualValue = fmt.Sprintf("%v", reflect.ValueOf(actualFunc))
	return expectValue == actualValue
}

func assertFunctionEquals(t *testing.T, expectFunc any, actualFunc any) {
	assert.True(t, functionPointerEquals(expectFunc, actualFunc))
}

type dummyResponseWriter struct {
	http.ResponseWriter
}

func TestInitiateSession(t *testing.T) {
	// arrange
	var dummyResponseWriter = &dummyResponseWriter{}
	var dummyHTTPRequest = &http.Request{Host: "some host"}
	var dummyCustomization = &DefaultCustomization{}
	var dummyAction = func(session Session) (any, error) { return nil, nil }
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
	m.Mock(uuid.New).Expects().Returns(dummySessionID).Once()
	m.Mock(getRouteInfo).Expects(dummyHTTPRequest, dummyActionFuncMap).
		Returns(dummyEndpoint, dummyAction, dummyRouteError).Once()

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
	assertFunctionEquals(t, dummyAction, action)
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
	m.Mock(handlePanic).Expects(dummySession, dummyRecoverResult).Returns().Once()
	m.Mock(logEndpointExit).Expects(dummySession, dummyName, dummyMethod,
		"%s", dummyDuration).Returns().Once()
	m.Mock(time.Since).Expects(dummyStartTime).Returns(dummyDuration).Once()

	// SUT + act
	finalizeSession(
		dummySession,
		dummyStartTime,
		dummyRecoverResult,
	)
}

func dummyAction(session Session) (any, error) {
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
	m.Mock((*DefaultCustomization).PreAction).Expects(dummyCustomization, dummySession).Returns(dummyError).Once()
	m.Mock(writeResponse).Expects(dummySession, nil, dummyError).Returns().Once()

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
	m.Mock((*DefaultCustomization).PreAction).Expects(dummyCustomization, dummySession).Returns(nil).Once()
	m.Mock(dummyAction).Expects(dummySession).Returns(dummyResponseObject, dummyError).Once()
	m.Mock(writeResponse).Expects(dummySession, dummyResponseObject, dummyError).Returns().Once()

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
	m.Mock((*DefaultCustomization).PreAction).Expects(dummyCustomization, dummySession).Returns(nil).Once()
	m.Mock(dummyAction).Expects(dummySession).Returns(dummyResponseObject, nil).Once()
	m.Mock((*DefaultCustomization).PostAction).Expects(dummyCustomization, dummySession).Returns(dummyError).Once()
	m.Mock(writeResponse).Expects(dummySession, dummyResponseObject, dummyError).Returns().Once()

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
	var dummyAction = func(session Session) (any, error) { return nil, nil }
	var dummyRouteError = errors.New("some route error")
	var dummyStartTime = time.Now()

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(initiateSession).Expects(dummyApplication, dummyResponseWriter, dummyHTTPRequest).Returns(dummySession, dummyAction, dummyRouteError).Once()
	m.Mock(logEndpointEnter).Expects(dummySession, dummyName, dummyMethod, "").Returns().Once()
	m.Mock(getTimeNowUTC).Expects().Returns(dummyStartTime).Once()
	m.Mock(finalizeSession).Expects(dummySession, dummyStartTime, recover()).Returns().Once()
	m.Mock(writeResponse).Expects(dummySession, nil, dummyRouteError).Returns().Once()

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
	var dummyAction = func(session Session) (any, error) { return nil, nil }
	var dummyStartTime = time.Now()

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(initiateSession).Expects(dummyApplication, dummyResponseWriter, dummyHTTPRequest).Returns(dummySession, dummyAction, nil).Once()
	m.Mock(logEndpointEnter).Expects(dummySession, dummyName, dummyMethod, "").Returns().Once()
	m.Mock(getTimeNowUTC).Expects().Returns(dummyStartTime).Once()
	m.Mock(finalizeSession).Expects(dummySession, dummyStartTime, recover()).Returns().Once()
	m.Mock(handleAction).Expects(dummySession, gomocker.Matches(func(value any) bool {
		return functionPointerEquals(dummyAction, value)
	})).Returns().Once()

	// SUT + act
	dummyApplication.handleSession(
		dummyResponseWriter,
		dummyHTTPRequest,
	)
}
