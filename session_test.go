package webserver

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"net/textproto"
	"net/url"
	"runtime"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/zhongjie-cai/gomocker/v2"
)

func TestSessionGetID_NilSessionObject(t *testing.T) {
	// SUT
	var dummySession *session

	// act
	var result = dummySession.GetID()

	// assert
	assert.Zero(t, result)
}

func TestSessionGetID_ValidSessionObject(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// act
	var result = dummySession.GetID()

	// assert
	assert.Equal(t, dummySessionID, result)
}

func TestSessionGetName_NilSessionObject(t *testing.T) {
	// SUT
	var dummySession *session

	// act
	var result = dummySession.GetName()

	// assert
	assert.Zero(t, result)
}

func TestSessionGetName_ValidSessionObject(t *testing.T) {
	// arrange
	var dummyName = "some name"

	// SUT
	var dummySession = &session{
		name: dummyName,
	}

	// act
	var result = dummySession.GetName()

	// assert
	assert.Equal(t, dummyName, result)
}

func TestSessionGetRequest_NilSessionObject(t *testing.T) {
	// SUT
	var dummySession *session

	// act
	var result = dummySession.GetRequest()

	// assert
	assert.Equal(t, defaultRequest, result)
}

func TestSessionGetRequest_NilRequest(t *testing.T) {
	// arrange
	var dummyHTTPRequest *http.Request

	// SUT
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// act
	var result = dummySession.GetRequest()

	// assert
	assert.Equal(t, defaultRequest, result)
}

func TestSessionGetRequest_ValidRequest(t *testing.T) {
	// arrange
	var dummyHTTPRequest = &http.Request{
		Method:     http.MethodGet,
		RequestURI: "http://localhost/",
		Header:     map[string][]string{},
	}

	// SUT
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// act
	var result = dummySession.GetRequest()

	// assert
	assert.Equal(t, dummyHTTPRequest, result)
}

func TestSessionGetResponseWriter_NilSessionObject(t *testing.T) {
	// SUT
	var dummySession *session

	// act
	var result = dummySession.GetResponseWriter()

	// assert
	assert.Equal(t, defaultResponseWriter, result)
}

func TestSessionGetResponseWriter_NilResponseWriter(t *testing.T) {
	// arrange
	var dummyResponseWriterObject *dummyResponseWriter

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(isInterfaceValueNil).Expects(dummyResponseWriterObject).Returns(true).Once()

	// SUT
	var dummySession = &session{
		responseWriter: dummyResponseWriterObject,
	}

	// act
	var result = dummySession.GetResponseWriter()

	// assert
	assert.Equal(t, defaultResponseWriter, result)
}

func TestSessionGetResponseWriter_ValidResponseWriter(t *testing.T) {
	// arrange
	var dummyResponseWriterObject = &dummyResponseWriter{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(isInterfaceValueNil).Expects(dummyResponseWriterObject).Returns(false).Once()

	// SUT
	var dummySession = &session{
		responseWriter: dummyResponseWriterObject,
	}

	// act
	var result = dummySession.GetResponseWriter()

	// assert
	assert.Equal(t, dummyResponseWriterObject, result)
}

func TestSessionGetRequestBody_NilSession(t *testing.T) {
	// arrange
	var dummyDataTemplate int
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageSessionNil, []error{}).Returns(dummyAppError).Once()

	// SUT
	var dummySession *session

	// act
	var err = dummySession.GetRequestBody(
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Zero(t, dummyDataTemplate)
}

func TestSessionGetRequestBody_BodyEmpty(t *testing.T) {
	// arrange
	var dummyDataTemplate int
	var dummyHTTPRequest = &http.Request{}
	var dummyRequestBody string
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// expect
	m.Mock(getRequestBody).Expects(dummyHTTPRequest).Returns(dummyRequestBody).Once()
	m.Mock(newAppError).Expects(errorCodeBadRequest, errorMessageRequestBodyEmpty, []error{}).Returns(dummyAppError).Once()

	// act
	var err = dummySession.GetRequestBody(
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Zero(t, dummyDataTemplate)
}

func TestSessionGetRequestBody_BodyInvalid(t *testing.T) {
	// arrange
	var dummyDataTemplate int
	var dummyHTTPRequest = &http.Request{}
	var dummyRequestBody = "some request body"
	var dummyError = errors.New("some error")
	var dummyResult = rand.Int()
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// expect
	m.Mock(getRequestBody).Expects(dummyHTTPRequest).Returns(dummyRequestBody).Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Body", "Content", dummyRequestBody).Returns().Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Body", "UnmarshalError", "%+v", dummyError).Returns().Once()
	m.Mock(tryUnmarshal).Expects(dummyRequestBody, gomocker.Anything()).Returns(dummyError).SideEffect(func(index int, params ...interface{}) {
		*(params[1].(*int)) = dummyResult
	}).Once()
	m.Mock(newAppError).Expects(errorCodeBadRequest, errorMessageRequestBodyInvalid, []error{dummyError}).Returns(dummyAppError).Once()

	// act
	var err = dummySession.GetRequestBody(
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)
}

func TestSessionGetRequestBody_BodyValid(t *testing.T) {
	// arrange
	var dummyDataTemplate int
	var dummyHTTPRequest = &http.Request{}
	var dummyRequestBody = "some request body"
	var dummyResult = rand.Int()

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// expect
	m.Mock(getRequestBody).Expects(dummyHTTPRequest).Returns(dummyRequestBody).Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Body", "Content", dummyRequestBody).Returns().Once()
	m.Mock(tryUnmarshal).Expects(dummyRequestBody, gomocker.Anything()).Returns(nil).SideEffect(func(index int, params ...interface{}) {
		*(params[1].(*int)) = dummyResult
	}).Once()

	// act
	var err = dummySession.GetRequestBody(
		&dummyDataTemplate,
	)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)
}

func TestSessionGetRequestParameter_NilSession(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyDataTemplate int
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageSessionNil, []error{}).Returns(dummyAppError).Once()

	// SUT
	var dummySession *session

	// act
	var err = dummySession.GetRequestParameter(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
}

func TestSessionGetRequestParameter_ParameterNotFound(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyDataTemplate int
	var dummyHTTPRequest = &http.Request{}
	var dummyParameters = map[string]string{
		"foo":  "bar",
		"test": "123",
	}
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// expect
	m.Mock(mux.Vars).Expects(dummyHTTPRequest).Returns(dummyParameters).Once()
	m.Mock(newAppError).Expects(errorCodeBadRequest, errorMessageParameterNotFound, []error{}).Returns(dummyAppError).Once()

	// act
	var err = dummySession.GetRequestParameter(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
}

func TestSessionGetRequestParameter_ParameterInvalid(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyValue = "some value"
	var dummyDataTemplate int
	var dummyHTTPRequest = &http.Request{}
	var dummyParameters = map[string]string{
		"foo":     "bar",
		"test":    "123",
		dummyName: dummyValue,
	}
	var dummyError = errors.New("some error")
	var dummyResult = rand.Int()
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// expect
	m.Mock(mux.Vars).Expects(dummyHTTPRequest).Returns(dummyParameters).Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Parameter", dummyName, dummyValue).Returns().Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Parameter", "UnmarshalError", "%+v", dummyError).Returns().Once()
	m.Mock(tryUnmarshal).Expects(dummyValue, gomocker.Anything()).Returns(dummyError).SideEffect(func(index int, params ...interface{}) {
		*(params[1].(*int)) = dummyResult
	}).Once()
	m.Mock(newAppError).Expects(errorCodeBadRequest, errorMessageParameterInvalid, []error{dummyError}).Returns(dummyAppError).Once()

	// act
	var err = dummySession.GetRequestParameter(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)
}

func TestSessionGetRequestParameter_ParameterValid(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyValue = "some value"
	var dummyDataTemplate int
	var dummyHTTPRequest = &http.Request{}
	var dummyParameters = map[string]string{
		"foo":     "bar",
		"test":    "123",
		dummyName: dummyValue,
	}
	var dummyResult = rand.Int()

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// expect
	m.Mock(mux.Vars).Expects(dummyHTTPRequest).Returns(dummyParameters).Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Parameter", dummyName, dummyValue).Returns().Once()
	m.Mock(tryUnmarshal).Expects(dummyValue, gomocker.Anything()).Returns(nil).SideEffect(func(index int, params ...interface{}) {
		*(params[1].(*int)) = dummyResult
	}).Once()

	// act
	var err = dummySession.GetRequestParameter(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)
}

func TestSessionGetAllQueries_NoURL(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyHTTPRequest = &http.Request{}
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// SUT + act
	var result = getAllQueries(
		dummySession,
		dummyName,
	)

	// assert
	assert.Nil(t, result)
}

func TestSessionGetAllQueries_NotFound(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyHTTPRequest = &http.Request{
		URL: &url.URL{
			RawQuery: "test=me&test=you",
		},
	}
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// SUT + act
	var result = getAllQueries(
		dummySession,
		dummyName,
	)

	// assert
	assert.Nil(t, result)
}

func TestSessionGetAllQueries_HappyPath(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyHTTPRequest = &http.Request{
		URL: &url.URL{
			RawQuery: dummyName + "=me&" + dummyName + "=you",
		},
	}
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// SUT + act
	var result = getAllQueries(
		dummySession,
		dummyName,
	)

	// assert
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "me", result[0])
	assert.Equal(t, "you", result[1])
}

func TestSessionGetRequestQuery_NilSession(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyIndex = rand.Intn(10)
	var dummyDataTemplate int
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageSessionNil, []error{}).Returns(dummyAppError).Once()

	// SUT
	var dummySession *session

	// act
	var err = dummySession.GetRequestQuery(
		dummyName,
		dummyIndex,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Zero(t, dummyDataTemplate)
}

func TestSessionGetRequestQuery_QueryNotFound(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyName = "some name"
	var dummyDataTemplate int
	var dummyQueries = []string{
		"some query string 1",
		"some query string 2",
		"some query string 3",
	}
	var dummyIndex = rand.Intn(10) + len(dummyQueries)
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(getAllQueries).Expects(dummySession, dummyName).Returns(dummyQueries).Once()
	m.Mock(newAppError).Expects(errorCodeBadRequest, errorMessageQueryNotFound, []error{}).Returns(dummyAppError).Once()

	// act
	var err = dummySession.GetRequestQuery(
		dummyName,
		dummyIndex,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Zero(t, dummyDataTemplate)
}

func TestSessionGetRequestQuery_QueryInvalid(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyName = "some name"
	var dummyDataTemplate int
	var dummyQueries = []string{
		"some query string 1",
		"some query string 2",
		"some query string 3",
	}
	var dummyIndex = rand.Intn(len(dummyQueries))
	var dummyError = errors.New("some error")
	var dummyResult = rand.Int()
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(getAllQueries).Expects(dummySession, dummyName).Returns(dummyQueries).Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Query", dummyName, dummyQueries[dummyIndex]).Returns().Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Query", "UnmarshalError", "%+v", dummyError).Returns().Once()
	m.Mock(tryUnmarshal).Expects(dummyQueries[dummyIndex], gomocker.Anything()).Returns(dummyError).SideEffect(func(index int, params ...interface{}) {
		*(params[1].(*int)) = dummyResult
	}).Once()
	m.Mock(newAppError).Expects(errorCodeBadRequest, errorMessageQueryInvalid, []error{dummyError}).Returns(dummyAppError).Once()

	// act
	var err = dummySession.GetRequestQuery(
		dummyName,
		dummyIndex,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)
}

func TestSessionGetRequestQuery_QueryValid(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyName = "some name"
	var dummyDataTemplate int
	var dummyQueries = []string{
		"some query string 1",
		"some query string 2",
		"some query string 3",
	}
	var dummyIndex = rand.Intn(len(dummyQueries))
	var dummyResult = rand.Int()

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(getAllQueries).Expects(dummySession, dummyName).Returns(dummyQueries).Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Query", dummyName, dummyQueries[dummyIndex]).Returns().Once()
	m.Mock(tryUnmarshal).Expects(dummyQueries[dummyIndex], gomocker.Anything()).Returns(nil).SideEffect(func(index int, params ...interface{}) {
		*(params[1].(*int)) = dummyResult
	}).Once()

	// act
	var err = dummySession.GetRequestQuery(
		dummyName,
		dummyIndex,
		&dummyDataTemplate,
	)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)
}

func TestSessionGetAllHeaders_NotFound(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyCanonicalName = "some canonical name"
	var dummyHTTPRequest = &http.Request{
		Header: http.Header{},
	}
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// stub
	dummyHTTPRequest.Header.Add("test", "me")
	dummyHTTPRequest.Header.Add("test", "you")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(textproto.CanonicalMIMEHeaderKey).Expects(dummyName).Returns(dummyCanonicalName).Once()

	// SUT + act
	var result = getAllHeaders(
		dummySession,
		dummyName,
	)

	// assert
	assert.Nil(t, result)
}

func TestSessionGetAllHeaders_HappyPath(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyCanonicalName = "some canonical name"
	var dummyHTTPRequest = &http.Request{
		Header: http.Header{},
	}
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// stub
	dummyHTTPRequest.Header.Add(dummyCanonicalName, "me")
	dummyHTTPRequest.Header.Add(dummyCanonicalName, "you")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(textproto.CanonicalMIMEHeaderKey).Expects(dummyName).Returns(dummyCanonicalName).Once()

	// SUT + act
	var result = getAllHeaders(
		dummySession,
		dummyName,
	)

	// assert
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "me", result[0])
	assert.Equal(t, "you", result[1])
}

func TestSessionGetRequestHeader_NilSession(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyDataTemplate int
	var dummyIndex = rand.Intn(10)
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageSessionNil, []error{}).Returns(dummyAppError).Once()

	// SUT
	var dummySession *session

	// act
	var err = dummySession.GetRequestHeader(
		dummyName,
		dummyIndex,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Zero(t, dummyDataTemplate)
}

func TestSessionGetRequestHeader_HeaderNotFound(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyName = "some name"
	var dummyDataTemplate int
	var dummyHeaders = []string{
		"some header string 1",
		"some header string 2",
		"some header string 3",
	}
	var dummyIndex = rand.Intn(10) + len(dummyHeaders)
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(getAllHeaders).Expects(dummySession, dummyName).Returns(dummyHeaders).Once()
	m.Mock(newAppError).Expects(errorCodeBadRequest, errorMessageHeaderNotFound, []error{}).Returns(dummyAppError).Once()

	// act
	var err = dummySession.GetRequestHeader(
		dummyName,
		dummyIndex,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Zero(t, dummyDataTemplate)
}

func TestSessionGetRequestHeader_HeaderInvalid(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyName = "some name"
	var dummyDataTemplate int
	var dummyHeaders = []string{
		"some header string 1",
		"some header string 2",
		"some header string 3",
	}
	var dummyIndex = rand.Intn(len(dummyHeaders))
	var dummyError = errors.New("some error")
	var dummyResult = rand.Int()
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(getAllHeaders).Expects(dummySession, dummyName).Returns(dummyHeaders).Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Header", dummyName, dummyHeaders[dummyIndex]).Returns().Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Header", "UnmarshalError", "%+v", dummyError).Returns().Once()
	m.Mock(tryUnmarshal).Expects(dummyHeaders[dummyIndex], gomocker.Anything()).Returns(dummyError).SideEffect(func(index int, params ...interface{}) {
		*(params[1].(*int)) = dummyResult
	}).Once()
	m.Mock(newAppError).Expects(errorCodeBadRequest, errorMessageHeaderInvalid, []error{dummyError}).Returns(dummyAppError).Once()

	// act
	var err = dummySession.GetRequestHeader(
		dummyName,
		dummyIndex,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)
}

func TestSessionGetRequestHeader_HeaderValid(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyName = "some name"
	var dummyDataTemplate int
	var dummyHeaders = []string{
		"some header string 1",
		"some header string 2",
		"some header string 3",
	}
	var dummyIndex = rand.Intn(len(dummyHeaders))
	var dummyResult = rand.Int()

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(getAllHeaders).Expects(dummySession, dummyName).Returns(dummyHeaders).Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Header", dummyName, dummyHeaders[dummyIndex]).Returns().Once()
	m.Mock(tryUnmarshal).Expects(dummyHeaders[dummyIndex], gomocker.Anything()).Returns(nil).SideEffect(func(index int, params ...interface{}) {
		*(params[1].(*int)) = dummyResult
	}).Once()

	// act
	var err = dummySession.GetRequestHeader(
		dummyName,
		dummyIndex,
		&dummyDataTemplate,
	)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)
}

type dummyAttachment struct {
	ID   uuid.UUID
	Foo  string
	Test int
}

func TestSessionAttach_NilSessionObject(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyValue = dummyAttachment{
		ID:   uuid.New(),
		Foo:  "bar",
		Test: rand.Intn(100),
	}

	// SUT
	var dummySession *session

	// act
	var result = dummySession.Attach(
		dummyName,
		dummyValue,
	)

	// assert
	assert.False(t, result)
}

func TestSessionAttach_NoAttachment(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyValue = dummyAttachment{
		ID:   uuid.New(),
		Foo:  "bar",
		Test: rand.Intn(100),
	}

	// SUT
	var dummySession = &session{
		attachment: nil,
	}

	// act
	var result = dummySession.Attach(
		dummyName,
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.Equal(t, dummyValue, dummySession.attachment[dummyName])
}

func TestSessionAttach_WithAttachment(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyValue = dummyAttachment{
		ID:   uuid.New(),
		Foo:  "bar",
		Test: rand.Intn(100),
	}

	// SUT
	var dummySession = &session{
		attachment: map[string]interface{}{
			dummyName: "some value",
		},
	}

	// act
	var result = dummySession.Attach(
		dummyName,
		dummyValue,
	)

	// assert
	assert.True(t, result)
	assert.Equal(t, dummyValue, dummySession.attachment[dummyName])
}

func TestSessionDetach_NilSessionObject(t *testing.T) {
	// arrange
	var dummyName = "some name"

	// SUT
	var dummySession *session

	// act
	var result = dummySession.Detach(
		dummyName,
	)

	// assert
	assert.False(t, result)
}

func TestSessionDetach_NoAttachment(t *testing.T) {
	// arrange
	var dummyName = "some name"

	// SUT
	var dummySession = &session{
		attachment: nil,
	}

	// act
	var result = dummySession.Detach(
		dummyName,
	)

	// assert
	assert.True(t, result)
	var _, found = dummySession.attachment[dummyName]
	assert.False(t, found)
}

func TestSessionDetach_WithAttachment(t *testing.T) {
	// arrange
	var dummyName = "some name"

	// SUT
	var dummySession = &session{
		attachment: map[string]interface{}{
			dummyName: "some value",
		},
	}

	// act
	var result = dummySession.Detach(
		dummyName,
	)

	// assert
	assert.True(t, result)
	var _, found = dummySession.attachment[dummyName]
	assert.False(t, found)
}

func TestSessionGetRawAttachment_NoSession(t *testing.T) {
	// arrange
	var dummyName = "some name"

	// SUT
	var dummySession *session

	// act
	var result, found = dummySession.GetRawAttachment(
		dummyName,
	)

	// assert
	assert.Nil(t, result)
	assert.False(t, found)
}

func TestSessionGetRawAttachment_NoAttachment(t *testing.T) {
	// arrange
	var dummyName = "some name"

	// SUT
	var dummySession = &session{}

	// act
	var result, found = dummySession.GetRawAttachment(
		dummyName,
	)

	// assert
	assert.Nil(t, result)
	assert.False(t, found)
}

func TestSessionGetRawAttachment_Success(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyValue = dummyAttachment{
		Foo:  "bar",
		Test: rand.Intn(100),
		ID:   uuid.New(),
	}

	// SUT
	var dummySession = &session{
		attachment: map[string]interface{}{
			dummyName: dummyValue,
		},
	}

	// act
	var result, found = dummySession.GetRawAttachment(
		dummyName,
	)

	// assert
	assert.Equal(t, dummyValue, result)
	assert.True(t, found)
}

func TestSessionGetAttachment_NoSession(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyDataTemplate dummyAttachment

	// SUT
	var dummySession *session

	// act
	var result = dummySession.GetAttachment(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.False(t, result)
	assert.Zero(t, dummyDataTemplate)
}

func TestSessionGetAttachment_NoAttachment(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyDataTemplate dummyAttachment

	// SUT
	var dummySession = &session{}

	// act
	var result = dummySession.GetAttachment(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.False(t, result)
	assert.Zero(t, dummyDataTemplate)
}

func TestSessionGetAttachment_MarshalError(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyValue = dummyAttachment{
		Foo:  "bar",
		Test: rand.Intn(100),
		ID:   uuid.New(),
	}
	var dummyDataTemplate dummyAttachment

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(json.Marshal).Expects(dummyValue).Returns(nil, errors.New("some marshal error")).Once()

	// SUT
	var dummySession = &session{
		attachment: map[string]interface{}{
			dummyName: dummyValue,
		},
	}

	// act
	var result = dummySession.GetAttachment(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.False(t, result)
	assert.Zero(t, dummyDataTemplate)
}

func TestSessionGetAttachment_UnmarshalError(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyValue = dummyAttachment{
		Foo:  "bar",
		Test: rand.Intn(100),
		ID:   uuid.New(),
	}
	var dummyDataTemplate int

	// SUT
	var dummySession = &session{
		attachment: map[string]interface{}{
			dummyName: dummyValue,
		},
	}

	// act
	var result = dummySession.GetAttachment(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.False(t, result)
	assert.Zero(t, dummyDataTemplate)
}

func TestSessionGetAttachment_Success(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyValue = dummyAttachment{
		Foo:  "bar",
		Test: rand.Intn(100),
		ID:   uuid.New(),
	}
	var dummyDataTemplate dummyAttachment

	// SUT
	var dummySession = &session{
		attachment: map[string]interface{}{
			dummyName: dummyValue,
		},
	}

	// act
	var result = dummySession.GetAttachment(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.True(t, result)
	assert.Equal(t, dummyValue, dummyDataTemplate)
}

func TestSessionGetMethodName_UnknownCaller(t *testing.T) {
	// arrange
	var dummyPC = uintptr(rand.Int())
	var dummyFile = "some file"
	var dummyLine = rand.Int()
	var dummyOK = false

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(runtime.Caller).Expects(3).Returns(dummyPC, dummyFile, dummyLine, dummyOK).Once()

	// SUT + act
	var result = getMethodName()

	// assert
	assert.Equal(t, "?", result)
}

func TestSessionGetMethodName_HappyPath(t *testing.T) {
	// arrange
	var dummyName = "runtime.goexit"

	// SUT + act
	var result = getMethodName()

	// assert
	assert.Contains(t, dummyName, result)
}

func TestSessionLogMethodEnter(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyMethodName = "some method name"

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(getMethodName).Expects().Returns(dummyMethodName).Once()
	m.Mock(logMethodEnter).Expects(dummySession, dummyMethodName, "", "").Returns().Once()

	// act
	dummySession.LogMethodEnter()
}

func TestSessionLogMethodParameter(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyParameter1 = "foo"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("test")
	var dummyParameters = []interface{}{
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	}
	var dummyMethodName = "some method name"

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(getMethodName).Expects().Returns(dummyMethodName).Once()
	m.Mock(logMethodParameter).Expects(dummySession, dummyMethodName, "0", "%v", dummyParameters[0]).Returns().Once()
	m.Mock(logMethodParameter).Expects(dummySession, dummyMethodName, "1", "%v", dummyParameters[1]).Returns().Once()
	m.Mock(logMethodParameter).Expects(dummySession, dummyMethodName, "2", "%v", dummyParameters[2]).Returns().Once()

	// act
	dummySession.LogMethodParameter(
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)
}

func TestSessionLogMethodLogic(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyLogLevel = LogLevel(rand.Int())
	var dummyCategory = "some category"
	var dummySubcategory = "some subcategory"
	var dummyMessageFormat = "some message format"
	var dummyParameter1 = "foo"
	var dummyParameter2 = rand.Int()
	var dummyParameter3 = errors.New("test")

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(logMethodLogic).Expects(dummySession, dummyLogLevel,
		dummyCategory, dummySubcategory, dummyMessageFormat,
		dummyParameter1, dummyParameter2, dummyParameter3).Returns().Once()

	// act
	dummySession.LogMethodLogic(
		dummyLogLevel,
		dummyCategory,
		dummySubcategory,
		dummyMessageFormat,
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)
}

func TestSessionLogMethodReturn(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyReturn1 = "foo"
	var dummyReturn2 = rand.Int()
	var dummyReturn3 = errors.New("test")
	var dummyReturns = []interface{}{
		dummyReturn1,
		dummyReturn2,
		dummyReturn3,
	}
	var dummyMethodName = "some method name"

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(getMethodName).Expects().Returns(dummyMethodName).Once()
	m.Mock(logMethodReturn).Expects(dummySession, dummyMethodName, "0", "%v", dummyReturns[0]).Returns().Once()
	m.Mock(logMethodReturn).Expects(dummySession, dummyMethodName, "1", "%v", dummyReturns[1]).Returns().Once()
	m.Mock(logMethodReturn).Expects(dummySession, dummyMethodName, "2", "%v", dummyReturns[2]).Returns().Once()

	// act
	dummySession.LogMethodReturn(
		dummyReturn1,
		dummyReturn2,
		dummyReturn3,
	)
}

func TestSessionLogMethodExit(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyMethodName = "some method name"

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(getMethodName).Expects().Returns(dummyMethodName).Once()
	m.Mock(logMethodExit).Expects(dummySession, dummyMethodName, "", "").Returns().Once()

	// act
	dummySession.LogMethodExit()
}

func TestSessionCreateWebcallRequest(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyMethod = "some method"
	var dummyURL = "some URL"
	var dummyPayload = "some payload"
	var dummySendClientCert = rand.Intn(100) < 50

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// act
	var result = dummySession.CreateWebcallRequest(
		dummyMethod,
		dummyURL,
		dummyPayload,
		dummySendClientCert,
	)
	var webrequest, ok = result.(*webRequest)

	// assert
	assert.True(t, ok)
	assert.Equal(t, dummySession, webrequest.session)
	assert.Equal(t, dummyMethod, webrequest.method)
	assert.Equal(t, dummyURL, webrequest.url)
	assert.Equal(t, dummyPayload, webrequest.payload)
	assert.NotNil(t, webrequest.query)
	assert.Empty(t, webrequest.query)
	assert.NotNil(t, webrequest.header)
	assert.Empty(t, webrequest.header)
	assert.Zero(t, webrequest.connRetry)
	assert.Nil(t, webrequest.httpRetry)
	assert.Equal(t, dummySendClientCert, webrequest.sendClientCert)
	assert.Zero(t, webrequest.retryDelay)
	assert.Empty(t, webrequest.dataReceivers)
}
