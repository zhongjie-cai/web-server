package webserver

import (
	"errors"
	"math/rand"
	"net/http"
	"net/textproto"
	"net/url"
	"runtime"
	"testing"

	"github.com/google/uuid"
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

func TestGetRequestBodyFromSession_HappyPath(t *testing.T) {
	// arrange
	var dummyResult = rand.Int()
	var dummyError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var sut = &session{}

	// expect
	m.Mock((*session).GetRequestBody).Expects(sut, gomocker.Anything()).Returns(dummyError).SideEffects(
		gomocker.ParamSideEffect(1, 2, func(value *int) { *value = dummyResult })).Once()

	//  act
	var result, err = GetRequestBodyFromSession[int](sut)

	// assert
	assert.Equal(t, dummyResult, *result)
	assert.Equal(t, dummyError, err)
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
	m.Mock(tryUnmarshal).Expects(dummyRequestBody, gomocker.Anything()).Returns(dummyError).SideEffects(
		gomocker.ParamSideEffect(1, 2, func(value *int) { *value = dummyResult })).Once()
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
	m.Mock(tryUnmarshal).Expects(dummyRequestBody, gomocker.Anything()).Returns(nil).SideEffects(
		gomocker.ParamSideEffect(1, 2, func(value *int) { *value = dummyResult })).Once()

	// act
	var err = dummySession.GetRequestBody(
		&dummyDataTemplate,
	)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)
}

func TestGetRequestParameterFromSession_HappyPath(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyResult = rand.Int()
	var dummyError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var sut = &session{}

	// expect
	m.Mock((*session).GetRequestParameter).Expects(sut, dummyName, gomocker.Anything()).Returns(dummyError).SideEffects(
		gomocker.ParamSideEffect(1, 3, func(value *int) { *value = dummyResult })).Once()

	//  act
	var result, err = GetRequestParameterFromSession[int](sut, dummyName)

	// assert
	assert.Equal(t, dummyResult, *result)
	assert.Equal(t, dummyError, err)
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
	m.Mock((*http.Request).PathValue).Expects(dummyHTTPRequest, dummyName).Returns(dummyParameters).Once()
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
	m.Mock((*http.Request).PathValue).Expects(dummyHTTPRequest, dummyName).Returns(dummyParameters).Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Parameter", dummyName, dummyValue).Returns().Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Parameter", "UnmarshalError", "%+v", dummyError).Returns().Once()
	m.Mock(tryUnmarshal).Expects(dummyValue, gomocker.Anything()).Returns(dummyError).SideEffects(
		gomocker.ParamSideEffect(1, 2, func(value *int) { *value = dummyResult })).Once()
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
	m.Mock((*http.Request).PathValue).Expects(dummyHTTPRequest, dummyName).Returns(dummyParameters).Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Parameter", dummyName, dummyValue).Returns().Once()
	m.Mock(tryUnmarshal).Expects(dummyValue, gomocker.Anything()).Returns(nil).SideEffects(
		gomocker.ParamSideEffect(1, 2, func(value *int) { *value = dummyResult })).Once()

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
	assert.Empty(t, result)
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
	assert.Empty(t, result)
}

func TestSessionGetAllQueries_HappyPath(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyHTTPRequest = &http.Request{
		URL: &url.URL{
			RawQuery: dummyName + "=foo,bar&" + dummyName + "=test",
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
	assert.Equal(t, 3, len(result))
	assert.Equal(t, "foo", result[0])
	assert.Equal(t, "bar", result[1])
	assert.Equal(t, "test", result[2])
}

func TestGetRequestQueriesFromSession_HappyPath(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyResults = []int{rand.Int(), rand.Int(), rand.Int()}
	var dummyError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var sut = &session{}

	// expect
	m.Mock((*session).GetRequestQueries).Expects(sut, dummyName, gomocker.Anything()).Returns(dummyError).SideEffects(
		gomocker.ParamSideEffect(1, 3, func(value *[]int) { *value = dummyResults })).Once()

	//  act
	var results, err = GetRequestQueriesFromSession[int](sut, dummyName)

	// assert
	assert.Equal(t, dummyResults, results)
	assert.Equal(t, dummyError, err)
}

func TestSessionGetRequestQueries_NilSession(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyDataTemplate []int
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession *session

	// expect
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageSessionNil, []error{}).Returns(dummyAppError).Once()

	// act
	var err = dummySession.GetRequestQueries(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Zero(t, dummyDataTemplate)
}

func TestSessionGetRequestQueries_DataTemplateNotAPointer(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyName = "some name"
	var dummyDataTemplate []int
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageDataTemplateInvalid, []error{}).Returns(dummyAppError).Once()

	// act
	var err = dummySession.GetRequestQueries(
		dummyName,
		dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Zero(t, dummyDataTemplate)
}

func TestSessionGetRequestQueries_DataTemplateNotASlice(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyName = "some name"
	var dummyDataTemplate int
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageDataTemplateInvalid, []error{}).Returns(dummyAppError).Once()

	// act
	var err = dummySession.GetRequestQueries(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Zero(t, dummyDataTemplate)
}

func TestSessionGetRequestQueries_UnmarshalError(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyName = "some name"
	var dummyDataTemplate []int
	var dummyQueries = []string{
		"foo",
		"456",
		"789",
	}
	var dummyError = errors.New("some error")
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(getAllQueries).Expects(dummySession, dummyName).Returns(dummyQueries).Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Query", dummyName, dummyQueries[0]).Returns().Once()
	m.Mock(tryUnmarshal).Expects(dummyQueries[0], gomocker.Anything()).Returns(dummyError).Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Query", "UnmarshalError", "%+v", dummyError).Returns().Once()
	m.Mock(newAppError).Expects(errorCodeBadRequest, errorMessageQueryInvalid, []error{dummyError}).Returns(dummyAppError).Once()

	// act
	var err = dummySession.GetRequestQueries(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Empty(t, dummyDataTemplate)
}

func TestSessionGetRequestQueries_Success(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyName = "some name"
	var dummyDataTemplate []int
	var dummyQueries = []string{
		"123",
		"456",
		"789",
	}
	var dummyResult = []int{123, 456, 789}

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(getAllQueries).Expects(dummySession, dummyName).Returns(dummyQueries).Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Query", dummyName, dummyQueries[0]).Returns().Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Query", dummyName, dummyQueries[1]).Returns().Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Query", dummyName, dummyQueries[2]).Returns().Once()

	// act
	var err = dummySession.GetRequestQueries(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)
}

func TestGetRequestQueryFromSession_HappyPath(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyIndex = rand.Int()
	var dummyResult = rand.Int()
	var dummyError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var sut = &session{}

	// expect
	m.Mock((*session).GetRequestQuery).Expects(sut, dummyName, dummyIndex, gomocker.Anything()).Returns(dummyError).SideEffects(
		gomocker.ParamSideEffect(1, 4, func(value *int) { *value = dummyResult })).Once()

	//  act
	var result, err = GetRequestQueryFromSession[int](sut, dummyName, dummyIndex)

	// assert
	assert.Equal(t, dummyResult, *result)
	assert.Equal(t, dummyError, err)
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
	m.Mock(tryUnmarshal).Expects(dummyQueries[dummyIndex], gomocker.Anything()).Returns(dummyError).SideEffects(
		gomocker.ParamSideEffect(1, 2, func(value *int) { *value = dummyResult })).Once()
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
	m.Mock(tryUnmarshal).Expects(dummyQueries[dummyIndex], gomocker.Anything()).Returns(nil).SideEffects(
		gomocker.ParamSideEffect(1, 2, func(value *int) { *value = dummyResult })).Once()

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

func TestGetRequestHeadersFromSession_HappyPath(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyResults = []int{rand.Int(), rand.Int(), rand.Int()}
	var dummyError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var sut = &session{}

	// expect
	m.Mock((*session).GetRequestHeaders).Expects(sut, dummyName, gomocker.Anything()).Returns(dummyError).SideEffects(
		gomocker.ParamSideEffect(1, 3, func(value *[]int) { *value = dummyResults })).Once()

	//  act
	var results, err = GetRequestHeadersFromSession[int](sut, dummyName)

	// assert
	assert.Equal(t, dummyResults, results)
	assert.Equal(t, dummyError, err)
}

func TestSessionGetRequestHeaders_NilSession(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyDataTemplate []int
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession *session

	// expect
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageSessionNil, []error{}).Returns(dummyAppError).Once()

	// act
	var err = dummySession.GetRequestHeaders(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Zero(t, dummyDataTemplate)
}

func TestSessionGetRequestHeaders_DataTemplateNotAPointer(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyName = "some name"
	var dummyDataTemplate []int
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageDataTemplateInvalid, []error{}).Returns(dummyAppError).Once()

	// act
	var err = dummySession.GetRequestHeaders(
		dummyName,
		dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Zero(t, dummyDataTemplate)
}

func TestSessionGetRequestHeaders_DataTemplateNotASlice(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyName = "some name"
	var dummyDataTemplate int
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageDataTemplateInvalid, []error{}).Returns(dummyAppError).Once()

	// act
	var err = dummySession.GetRequestHeaders(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Zero(t, dummyDataTemplate)
}

func TestSessionGetRequestHeaders_UnmarshalError(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyName = "some name"
	var dummyDataTemplate []int
	var dummyHeaders = []string{
		"foo",
		"456",
		"789",
	}
	var dummyError = errors.New("some error")
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(getAllHeaders).Expects(dummySession, dummyName).Returns(dummyHeaders).Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Header", dummyName, dummyHeaders[0]).Returns().Once()
	m.Mock(tryUnmarshal).Expects(dummyHeaders[0], gomocker.Anything()).Returns(dummyError).Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Header", "UnmarshalError", "%+v", dummyError).Returns().Once()
	m.Mock(newAppError).Expects(errorCodeBadRequest, errorMessageHeaderInvalid, []error{dummyError}).Returns(dummyAppError).Once()

	// act
	var err = dummySession.GetRequestHeaders(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Empty(t, dummyDataTemplate)
}

func TestSessionGetRequestHeaders_Success(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyName = "some name"
	var dummyDataTemplate []int
	var dummyHeaders = []string{
		"123",
		"456",
		"789",
	}
	var dummyResult = []int{123, 456, 789}

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	m.Mock(getAllHeaders).Expects(dummySession, dummyName).Returns(dummyHeaders).Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Header", dummyName, dummyHeaders[0]).Returns().Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Header", dummyName, dummyHeaders[1]).Returns().Once()
	m.Mock(logEndpointRequest).Expects(dummySession, "Header", dummyName, dummyHeaders[2]).Returns().Once()

	// act
	var err = dummySession.GetRequestHeaders(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)
}

func TestGetRequestHeaderFromSession_HappyPath(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyIndex = rand.Int()
	var dummyResult = rand.Int()
	var dummyError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var sut = &session{}

	// expect
	m.Mock((*session).GetRequestHeader).Expects(sut, dummyName, dummyIndex, gomocker.Anything()).Returns(dummyError).SideEffects(
		gomocker.ParamSideEffect(1, 4, func(value *int) { *value = dummyResult })).Once()

	//  act
	var result, err = GetRequestHeaderFromSession[int](sut, dummyName, dummyIndex)

	// assert
	assert.Equal(t, dummyResult, *result)
	assert.Equal(t, dummyError, err)
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
	m.Mock(tryUnmarshal).Expects(dummyHeaders[dummyIndex], gomocker.Anything()).Returns(dummyError).SideEffects(
		gomocker.ParamSideEffect(1, 2, func(value *int) { *value = dummyResult })).Once()
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
	m.Mock(tryUnmarshal).Expects(dummyHeaders[dummyIndex], gomocker.Anything()).Returns(nil).SideEffects(
		gomocker.ParamSideEffect(1, 2, func(value *int) { *value = dummyResult })).Once()

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
	var dummyValue = &dummyAttachment{
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
	var dummyValue = &dummyAttachment{
		ID:   uuid.New(),
		Foo:  "bar",
		Test: rand.Intn(100),
	}

	// SUT
	var dummySession = &session{
		attachment: map[string]any{
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
		attachment: map[string]any{
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
		attachment: map[string]any{
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

func TestGetAttachmentFromSession_NotFound(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyResult = rand.Int()

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var sut = &session{}

	// expect
	m.Mock((*session).GetRawAttachment).Expects(sut, dummyName).Returns(dummyResult, false).Once()

	//  act
	var result, found = GetAttachmentFromSession[int](sut, dummyName)

	// assert
	assert.Zero(t, *result)
	assert.False(t, found)
}

func TestGetAttachmentFromSession_TypeMismatch(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyResult = "some value"

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var sut = &session{}

	// expect
	m.Mock((*session).GetRawAttachment).Expects(sut, dummyName).Returns(dummyResult, true).Once()

	//  act
	var result, found = GetAttachmentFromSession[int](sut, dummyName)

	// assert
	assert.Zero(t, *result)
	assert.False(t, found)
}

func TestGetAttachmentFromSession_Success(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyResult = rand.Int()

	// mock
	var m = gomocker.NewMocker(t)

	// SUT
	var sut = &session{}

	// expect
	m.Mock((*session).GetRawAttachment).Expects(sut, dummyName).Returns(dummyResult, true).Once()

	//  act
	var result, found = GetAttachmentFromSession[int](sut, dummyName)

	// assert
	assert.Equal(t, dummyResult, *result)
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

func TestSessionGetAttachment_DataTemplateNotAPointer(t *testing.T) {
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
		attachment: map[string]any{
			dummyName: dummyValue,
		},
	}

	// act
	var result = dummySession.GetAttachment(
		dummyName,
		dummyDataTemplate,
	)

	// assert
	assert.False(t, result)
	assert.NotEqual(t, dummyValue, dummyDataTemplate)
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
		attachment: map[string]any{
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
	var dummyParameters = []any{
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
	var dummyReturns = []any{
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
