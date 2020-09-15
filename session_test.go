package webserver

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSessionGetID_NilSessionObject(t *testing.T) {
	// mock
	createMock(t)

	// SUT
	var dummySession *session

	// act
	var result = dummySession.GetID()

	// assert
	assert.Zero(t, result)

	// verify
	verifyAll(t)
}

func TestSessionGetID_ValidSessionObject(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()

	// mock
	createMock(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// act
	var result = dummySession.GetID()

	// assert
	assert.Equal(t, dummySessionID, result)

	// verify
	verifyAll(t)
}

func TestSessionGetName_NilSessionObject(t *testing.T) {
	// mock
	createMock(t)

	// SUT
	var dummySession *session

	// act
	var result = dummySession.GetName()

	// assert
	assert.Zero(t, result)

	// verify
	verifyAll(t)
}

func TestSessionGetName_ValidSessionObject(t *testing.T) {
	// arrange
	var dummyName = "some name"

	// mock
	createMock(t)

	// SUT
	var dummySession = &session{
		name: dummyName,
	}

	// act
	var result = dummySession.GetName()

	// assert
	assert.Equal(t, dummyName, result)

	// verify
	verifyAll(t)
}

func TestSessionGetRequest_NilSessionObject(t *testing.T) {
	// mock
	createMock(t)

	// SUT
	var dummySession *session

	// act
	var result = dummySession.GetRequest()

	// assert
	assert.Equal(t, defaultRequest, result)

	// verify
	verifyAll(t)
}

func TestSessionGetRequest_NilRequest(t *testing.T) {
	// arrange
	var dummyHTTPRequest *http.Request

	// mock
	createMock(t)

	// SUT
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// act
	var result = dummySession.GetRequest()

	// assert
	assert.Equal(t, defaultRequest, result)

	// verify
	verifyAll(t)
}

func TestSessionGetRequest_ValidRequest(t *testing.T) {
	// arrange
	var dummyHTTPRequest = &http.Request{
		Method:     http.MethodGet,
		RequestURI: "http://localhost/",
		Header:     map[string][]string{},
	}

	// mock
	createMock(t)

	// SUT
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// act
	var result = dummySession.GetRequest()

	// assert
	assert.Equal(t, dummyHTTPRequest, result)

	// verify
	verifyAll(t)
}

func TestSessionGetResponseWriter_NilSessionObject(t *testing.T) {
	// mock
	createMock(t)

	// SUT
	var dummySession *session

	// act
	var result = dummySession.GetResponseWriter()

	// assert
	assert.Equal(t, defaultResponseWriter, result)

	// verify
	verifyAll(t)
}

func TestSessionGetResponseWriter_NilResponseWriter(t *testing.T) {
	// arrange
	var dummyResponseWriterObject *dummyResponseWriter

	// mock
	createMock(t)

	// expect
	isInterfaceValueNilFuncExpected = 1
	isInterfaceValueNilFunc = func(i interface{}) bool {
		isInterfaceValueNilFuncCalled++
		assert.Equal(t, dummyResponseWriterObject, i)
		return true
	}

	// SUT
	var dummySession = &session{
		responseWriter: dummyResponseWriterObject,
	}

	// act
	var result = dummySession.GetResponseWriter()

	// assert
	assert.Equal(t, defaultResponseWriter, result)

	// verify
	verifyAll(t)
}

func TestSessionGetResponseWriter_ValidResponseWriter(t *testing.T) {
	// arrange
	var dummyResponseWriterObject = &dummyResponseWriter{}

	// mock
	createMock(t)

	// expect
	isInterfaceValueNilFuncExpected = 1
	isInterfaceValueNilFunc = func(i interface{}) bool {
		isInterfaceValueNilFuncCalled++
		assert.Equal(t, dummyResponseWriterObject, i)
		return false
	}

	// SUT
	var dummySession = &session{
		responseWriter: dummyResponseWriterObject,
	}

	// act
	var result = dummySession.GetResponseWriter()

	// assert
	assert.Equal(t, dummyResponseWriterObject, result)

	// verify
	verifyAll(t)
}

func TestSessionGetRequestBody_NilSession(t *testing.T) {
	// arrange
	var dummyDataTemplate int
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageSessionNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	}

	// SUT
	var dummySession *session

	// act
	var err = dummySession.GetRequestBody(
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Zero(t, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestSessionGetRequestBody_BodyEmpty(t *testing.T) {
	// arrange
	var dummyDataTemplate int
	var dummyHTTPRequest = &http.Request{}
	var dummyRequestBody string
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	createMock(t)

	// SUT
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// expect
	getRequestBodyFuncExpected = 1
	getRequestBodyFunc = func(httpRequest *http.Request) string {
		getRequestBodyFuncCalled++
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		return dummyRequestBody
	}
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeBadRequest, errorCode)
		assert.Equal(t, errorMessageRequestBodyEmpty, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	}

	// act
	var err = dummySession.GetRequestBody(
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Zero(t, dummyDataTemplate)

	// verify
	verifyAll(t)
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
	createMock(t)

	// SUT
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// expect
	getRequestBodyFuncExpected = 1
	getRequestBodyFunc = func(httpRequest *http.Request) string {
		getRequestBodyFuncCalled++
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		return dummyRequestBody
	}
	logEndpointRequestFuncExpected = 2
	logEndpointRequestFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logEndpointRequestFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Body", category)
		if logEndpointRequestFuncCalled == 1 {
			assert.Equal(t, "Content", subcategory)
			assert.Equal(t, dummyRequestBody, messageFormat)
			assert.Empty(t, parameters)
		} else if logEndpointRequestFuncCalled == 2 {
			assert.Equal(t, "UnmarshalError", subcategory)
			assert.Equal(t, "%+v", messageFormat)
			assert.Equal(t, 1, len(parameters))
			assert.Equal(t, dummyError, parameters[0])
		}
	}
	tryUnmarshalFuncExpected = 1
	tryUnmarshalFunc = func(value string, dataTemplate interface{}) error {
		tryUnmarshalFuncCalled++
		assert.Equal(t, dummyRequestBody, value)
		*(dataTemplate.(*int)) = dummyResult
		return dummyError
	}
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeBadRequest, errorCode)
		assert.Equal(t, errorMessageRequestBodyInvalid, errorMessage)
		assert.Equal(t, 1, len(innerErrors))
		assert.Equal(t, dummyError, innerErrors[0])
		return dummyAppError
	}

	// act
	var err = dummySession.GetRequestBody(
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestSessionGetRequestBody_BodyValid(t *testing.T) {
	// arrange
	var dummyDataTemplate int
	var dummyHTTPRequest = &http.Request{}
	var dummyRequestBody = "some request body"
	var dummyResult = rand.Int()

	// mock
	createMock(t)

	// SUT
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// expect
	getRequestBodyFuncExpected = 1
	getRequestBodyFunc = func(httpRequest *http.Request) string {
		getRequestBodyFuncCalled++
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		return dummyRequestBody
	}
	logEndpointRequestFuncExpected = 1
	logEndpointRequestFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logEndpointRequestFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Body", category)
		assert.Equal(t, "Content", subcategory)
		assert.Equal(t, dummyRequestBody, messageFormat)
		assert.Empty(t, parameters)
	}
	tryUnmarshalFuncExpected = 1
	tryUnmarshalFunc = func(value string, dataTemplate interface{}) error {
		tryUnmarshalFuncCalled++
		assert.Equal(t, dummyRequestBody, value)
		*(dataTemplate.(*int)) = dummyResult
		return nil
	}

	// act
	var err = dummySession.GetRequestBody(
		&dummyDataTemplate,
	)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestSessionGetRequestParameter_NilSession(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyDataTemplate int
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageSessionNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	}

	// SUT
	var dummySession *session

	// act
	var err = dummySession.GetRequestParameter(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)

	// verify
	verifyAll(t)
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
	createMock(t)

	// SUT
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// expect
	muxVarsExpected = 1
	muxVars = func(r *http.Request) map[string]string {
		muxVarsCalled++
		assert.Equal(t, dummyHTTPRequest, r)
		return dummyParameters
	}
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeBadRequest, errorCode)
		assert.Equal(t, errorMessageParameterNotFound, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	}

	// act
	var err = dummySession.GetRequestParameter(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)

	// verify
	verifyAll(t)
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
	createMock(t)

	// SUT
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// expect
	muxVarsExpected = 1
	muxVars = func(r *http.Request) map[string]string {
		muxVarsCalled++
		assert.Equal(t, dummyHTTPRequest, r)
		return dummyParameters
	}
	logEndpointRequestFuncExpected = 2
	logEndpointRequestFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logEndpointRequestFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Parameter", category)
		if logEndpointRequestFuncCalled == 1 {
			assert.Equal(t, dummyName, subcategory)
			assert.Equal(t, dummyValue, messageFormat)
			assert.Empty(t, parameters)
		} else if logEndpointRequestFuncCalled == 2 {
			assert.Equal(t, "UnmarshalError", subcategory)
			assert.Equal(t, "%+v", messageFormat)
			assert.Equal(t, 1, len(parameters))
			assert.Equal(t, dummyError, parameters[0])
		}
	}
	tryUnmarshalFuncExpected = 1
	tryUnmarshalFunc = func(value string, dataTemplate interface{}) error {
		tryUnmarshalFuncCalled++
		assert.Equal(t, dummyValue, value)
		*(dataTemplate.(*int)) = dummyResult
		return dummyError
	}
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeBadRequest, errorCode)
		assert.Equal(t, errorMessageParameterInvalid, errorMessage)
		assert.Equal(t, 1, len(innerErrors))
		assert.Equal(t, dummyError, innerErrors[0])
		return dummyAppError
	}

	// act
	var err = dummySession.GetRequestParameter(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)

	// verify
	verifyAll(t)
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
	createMock(t)

	// SUT
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// expect
	muxVarsExpected = 1
	muxVars = func(r *http.Request) map[string]string {
		muxVarsCalled++
		assert.Equal(t, dummyHTTPRequest, r)
		return dummyParameters
	}
	logEndpointRequestFuncExpected = 1
	logEndpointRequestFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logEndpointRequestFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Parameter", category)
		assert.Equal(t, dummyName, subcategory)
		assert.Equal(t, dummyValue, messageFormat)
		assert.Empty(t, parameters)
	}
	tryUnmarshalFuncExpected = 1
	tryUnmarshalFunc = func(value string, dataTemplate interface{}) error {
		tryUnmarshalFuncCalled++
		assert.Equal(t, dummyValue, value)
		*(dataTemplate.(*int)) = dummyResult
		return nil
	}

	// act
	var err = dummySession.GetRequestParameter(
		dummyName,
		&dummyDataTemplate,
	)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)

	// verify
	verifyAll(t)
}

func TestSessionGetAllQueries_NoURL(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyHTTPRequest = &http.Request{}
	var dummySession = &session{
		request: dummyHTTPRequest,
	}

	// mock
	createMock(t)

	// SUT + act
	var result = getAllQueries(
		dummySession,
		dummyName,
	)

	// assert
	assert.Nil(t, result)

	// verify
	verifyAll(t)
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

	// mock
	createMock(t)

	// SUT + act
	var result = getAllQueries(
		dummySession,
		dummyName,
	)

	// assert
	assert.Nil(t, result)

	// verify
	verifyAll(t)
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

	// mock
	createMock(t)

	// SUT + act
	var result = getAllQueries(
		dummySession,
		dummyName,
	)

	// assert
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "me", result[0])
	assert.Equal(t, "you", result[1])

	// verify
	verifyAll(t)
}

func TestSessionGetRequestQuery_NilSession(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyIndex = rand.Intn(10)
	var dummyDataTemplate int
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageSessionNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	}

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

	// verify
	verifyAll(t)
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
	createMock(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	getAllQueriesFuncExpected = 1
	getAllQueriesFunc = func(session *session, name string) []string {
		getAllQueriesFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyName, name)
		return dummyQueries
	}
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeBadRequest, errorCode)
		assert.Equal(t, errorMessageQueryNotFound, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	}

	// act
	var err = dummySession.GetRequestQuery(
		dummyName,
		dummyIndex,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Zero(t, dummyDataTemplate)

	// verify
	verifyAll(t)
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
	createMock(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	getAllQueriesFuncExpected = 1
	getAllQueriesFunc = func(session *session, name string) []string {
		getAllQueriesFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyName, name)
		return dummyQueries
	}
	logEndpointRequestFuncExpected = 2
	logEndpointRequestFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logEndpointRequestFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Query", category)
		if logEndpointRequestFuncCalled == 1 {
			assert.Equal(t, dummyName, subcategory)
			assert.Equal(t, dummyQueries[dummyIndex], messageFormat)
			assert.Empty(t, parameters)
		} else if logEndpointRequestFuncCalled == 2 {
			assert.Equal(t, "UnmarshalError", subcategory)
			assert.Equal(t, "%+v", messageFormat)
			assert.Equal(t, 1, len(parameters))
			assert.Equal(t, dummyError, parameters[0])
		}
	}
	tryUnmarshalFuncExpected = 1
	tryUnmarshalFunc = func(value string, dataTemplate interface{}) error {
		tryUnmarshalFuncCalled++
		assert.Equal(t, dummyQueries[dummyIndex], value)
		*(dataTemplate.(*int)) = dummyResult
		return dummyError
	}
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeBadRequest, errorCode)
		assert.Equal(t, errorMessageQueryInvalid, errorMessage)
		assert.Equal(t, 1, len(innerErrors))
		assert.Equal(t, dummyError, innerErrors[0])
		return dummyAppError
	}

	// act
	var err = dummySession.GetRequestQuery(
		dummyName,
		dummyIndex,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)

	// verify
	verifyAll(t)
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
	createMock(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	getAllQueriesFuncExpected = 1
	getAllQueriesFunc = func(session *session, name string) []string {
		getAllQueriesFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyName, name)
		return dummyQueries
	}
	logEndpointRequestFuncExpected = 1
	logEndpointRequestFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logEndpointRequestFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Query", category)
		assert.Equal(t, dummyName, subcategory)
		assert.Equal(t, dummyQueries[dummyIndex], messageFormat)
		assert.Empty(t, parameters)
	}
	tryUnmarshalFuncExpected = 1
	tryUnmarshalFunc = func(value string, dataTemplate interface{}) error {
		tryUnmarshalFuncCalled++
		assert.Equal(t, dummyQueries[dummyIndex], value)
		*(dataTemplate.(*int)) = dummyResult
		return nil
	}

	// act
	var err = dummySession.GetRequestQuery(
		dummyName,
		dummyIndex,
		&dummyDataTemplate,
	)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	textprotoCanonicalMIMEHeaderKeyExpected = 1
	textprotoCanonicalMIMEHeaderKey = func(s string) string {
		textprotoCanonicalMIMEHeaderKeyCalled++
		assert.Equal(t, dummyName, s)
		return dummyCanonicalName
	}

	// SUT + act
	var result = getAllHeaders(
		dummySession,
		dummyName,
	)

	// assert
	assert.Nil(t, result)

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	textprotoCanonicalMIMEHeaderKeyExpected = 1
	textprotoCanonicalMIMEHeaderKey = func(s string) string {
		textprotoCanonicalMIMEHeaderKeyCalled++
		assert.Equal(t, dummyName, s)
		return dummyCanonicalName
	}

	// SUT + act
	var result = getAllHeaders(
		dummySession,
		dummyName,
	)

	// assert
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "me", result[0])
	assert.Equal(t, "you", result[1])

	// verify
	verifyAll(t)
}

func TestSessionGetRequestHeader_NilSession(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyDataTemplate int
	var dummyIndex = rand.Intn(10)
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	createMock(t)

	// expect
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageSessionNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	}

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

	// verify
	verifyAll(t)
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
	createMock(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	getAllHeadersFuncExpected = 1
	getAllHeadersFunc = func(session *session, name string) []string {
		getAllHeadersFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyName, name)
		return dummyHeaders
	}
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeBadRequest, errorCode)
		assert.Equal(t, errorMessageHeaderNotFound, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	}

	// act
	var err = dummySession.GetRequestHeader(
		dummyName,
		dummyIndex,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Zero(t, dummyDataTemplate)

	// verify
	verifyAll(t)
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
	createMock(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	getAllHeadersFuncExpected = 1
	getAllHeadersFunc = func(session *session, name string) []string {
		getAllHeadersFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyName, name)
		return dummyHeaders
	}
	logEndpointRequestFuncExpected = 2
	logEndpointRequestFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logEndpointRequestFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Header", category)
		if logEndpointRequestFuncCalled == 1 {
			assert.Equal(t, dummyName, subcategory)
			assert.Equal(t, dummyHeaders[dummyIndex], messageFormat)
			assert.Empty(t, parameters)
		} else if logEndpointRequestFuncCalled == 2 {
			assert.Equal(t, "UnmarshalError", subcategory)
			assert.Equal(t, "%+v", messageFormat)
			assert.Equal(t, 1, len(parameters))
			assert.Equal(t, dummyError, parameters[0])
		}
	}
	tryUnmarshalFuncExpected = 1
	tryUnmarshalFunc = func(value string, dataTemplate interface{}) error {
		tryUnmarshalFuncCalled++
		assert.Equal(t, dummyHeaders[dummyIndex], value)
		*(dataTemplate.(*int)) = dummyResult
		return dummyError
	}
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeBadRequest, errorCode)
		assert.Equal(t, errorMessageHeaderInvalid, errorMessage)
		assert.Equal(t, 1, len(innerErrors))
		assert.Equal(t, dummyError, innerErrors[0])
		return dummyAppError
	}

	// act
	var err = dummySession.GetRequestHeader(
		dummyName,
		dummyIndex,
		&dummyDataTemplate,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)

	// verify
	verifyAll(t)
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
	createMock(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	getAllHeadersFuncExpected = 1
	getAllHeadersFunc = func(session *session, name string) []string {
		getAllHeadersFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyName, name)
		return dummyHeaders
	}
	logEndpointRequestFuncExpected = 1
	logEndpointRequestFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logEndpointRequestFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Header", category)
		assert.Equal(t, dummyName, subcategory)
		assert.Equal(t, dummyHeaders[dummyIndex], messageFormat)
		assert.Empty(t, parameters)
	}
	tryUnmarshalFuncExpected = 1
	tryUnmarshalFunc = func(value string, dataTemplate interface{}) error {
		tryUnmarshalFuncCalled++
		assert.Equal(t, dummyHeaders[dummyIndex], value)
		*(dataTemplate.(*int)) = dummyResult
		return nil
	}

	// act
	var err = dummySession.GetRequestHeader(
		dummyName,
		dummyIndex,
		&dummyDataTemplate,
	)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, dummyResult, dummyDataTemplate)

	// verify
	verifyAll(t)
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

	// mock
	createMock(t)

	// SUT
	var dummySession *session

	// act
	var result = dummySession.Attach(
		dummyName,
		dummyValue,
	)

	// assert
	assert.False(t, result)

	// verify
	verifyAll(t)
}

func TestSessionAttach_NoAttachment(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyValue = dummyAttachment{
		ID:   uuid.New(),
		Foo:  "bar",
		Test: rand.Intn(100),
	}

	// mock
	createMock(t)

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

	// verify
	verifyAll(t)
}

func TestSessionAttach_WithAttachment(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyValue = dummyAttachment{
		ID:   uuid.New(),
		Foo:  "bar",
		Test: rand.Intn(100),
	}

	// mock
	createMock(t)

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

	// verify
	verifyAll(t)
}

func TestSessionDetach_NilSessionObject(t *testing.T) {
	// arrange
	var dummyName = "some name"

	// mock
	createMock(t)

	// SUT
	var dummySession *session

	// act
	var result = dummySession.Detach(
		dummyName,
	)

	// assert
	assert.False(t, result)

	// verify
	verifyAll(t)
}

func TestSessionDetach_NoAttachment(t *testing.T) {
	// arrange
	var dummyName = "some name"

	// mock
	createMock(t)

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

	// verify
	verifyAll(t)
}

func TestSessionDetach_WithAttachment(t *testing.T) {
	// arrange
	var dummyName = "some name"

	// mock
	createMock(t)

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

	// verify
	verifyAll(t)
}

func TestSessionGetRawAttachment_NoSession(t *testing.T) {
	// arrange
	var dummyName = "some name"

	// mock
	createMock(t)

	// SUT
	var dummySession *session

	// act
	var result, found = dummySession.GetRawAttachment(
		dummyName,
	)

	// assert
	assert.Nil(t, result)
	assert.False(t, found)

	// verify
	verifyAll(t)
}

func TestSessionGetRawAttachment_NoAttachment(t *testing.T) {
	// arrange
	var dummyName = "some name"

	// mock
	createMock(t)

	// SUT
	var dummySession = &session{}

	// act
	var result, found = dummySession.GetRawAttachment(
		dummyName,
	)

	// assert
	assert.Nil(t, result)
	assert.False(t, found)

	// verify
	verifyAll(t)
}

func TestSessionGetRawAttachment_Success(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyValue = dummyAttachment{
		Foo:  "bar",
		Test: rand.Intn(100),
		ID:   uuid.New(),
	}

	// mock
	createMock(t)

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

	// verify
	verifyAll(t)
}

func TestSessionGetAttachment_NoSession(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyDataTemplate dummyAttachment

	// mock
	createMock(t)

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

	// verify
	verifyAll(t)
}

func TestSessionGetAttachment_NoAttachment(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyDataTemplate dummyAttachment

	// mock
	createMock(t)

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

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	jsonMarshalExpected = 1
	jsonMarshal = func(v interface{}) ([]byte, error) {
		jsonMarshalCalled++
		assert.Equal(t, dummyValue, v)
		return nil, errors.New("some marshal error")
	}

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

	// verify
	verifyAll(t)
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

	// mock
	createMock(t)

	// expect
	jsonMarshalExpected = 1
	jsonMarshal = func(v interface{}) ([]byte, error) {
		jsonMarshalCalled++
		assert.Equal(t, dummyValue, v)
		return json.Marshal(v)
	}
	jsonUnmarshalExpected = 1
	jsonUnmarshal = func(data []byte, v interface{}) error {
		jsonUnmarshalCalled++
		return json.Unmarshal(data, v)
	}

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

	// verify
	verifyAll(t)
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

	// mock
	createMock(t)

	// expect
	jsonMarshalExpected = 1
	jsonMarshal = func(v interface{}) ([]byte, error) {
		jsonMarshalCalled++
		assert.Equal(t, dummyValue, v)
		return json.Marshal(v)
	}
	jsonUnmarshalExpected = 1
	jsonUnmarshal = func(data []byte, v interface{}) error {
		jsonUnmarshalCalled++
		return json.Unmarshal(data, v)
	}

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

	// verify
	verifyAll(t)
}

func TestSessionGetMethodName_UnknownCaller(t *testing.T) {
	// arrange
	var dummyPC = uintptr(rand.Int())
	var dummyFile = "some file"
	var dummyLine = rand.Int()
	var dummyOK = false

	// mock
	createMock(t)

	// expect
	runtimeCallerExpected = 1
	runtimeCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		runtimeCallerCalled++
		assert.Equal(t, 3, skip)
		return dummyPC, dummyFile, dummyLine, dummyOK
	}

	// SUT + act
	var result = getMethodName()

	// assert
	assert.Equal(t, "?", result)

	// verify
	verifyAll(t)
}

func TestSessionGetMethodName_HappyPath(t *testing.T) {
	// mock
	createMock(t)

	// expect
	runtimeCallerExpected = 1
	runtimeCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		runtimeCallerCalled++
		assert.Equal(t, 3, skip)
		return runtime.Caller(2)
	}
	runtimeFuncForPCExpected = 1
	runtimeFuncForPC = func(pc uintptr) *runtime.Func {
		runtimeFuncForPCCalled++
		assert.NotZero(t, pc)
		return runtime.FuncForPC(pc)
	}

	// SUT + act
	var result = getMethodName()

	// assert
	assert.Contains(t, result, "TestSessionGetMethodName_HappyPath")

	// verify
	verifyAll(t)
}

func TestSessionLogMethodEnter(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyMethodName = "some method name"

	// mock
	createMock(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	getMethodNameFuncExpected = 1
	getMethodNameFunc = func() string {
		getMethodNameFuncCalled++
		return dummyMethodName
	}
	logMethodEnterFuncExpected = 1
	logMethodEnterFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logMethodEnterFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyMethodName, category)
		assert.Zero(t, subcategory)
		assert.Zero(t, messageFormat)
		assert.Empty(t, parameters)
	}

	// act
	dummySession.LogMethodEnter()

	// verify
	verifyAll(t)
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
	createMock(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	getMethodNameFuncExpected = 1
	getMethodNameFunc = func() string {
		getMethodNameFuncCalled++
		return dummyMethodName
	}
	strconvItoaExpected = 3
	strconvItoa = func(i int) string {
		strconvItoaCalled++
		return strconv.Itoa(i)
	}
	logMethodParameterFuncExpected = 3
	logMethodParameterFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logMethodParameterFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyMethodName, category)
		assert.Equal(t, strconv.Itoa(logMethodParameterFuncCalled-1), subcategory)
		assert.Equal(t, "%v", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyParameters[logMethodParameterFuncCalled-1], parameters[0])
	}

	// act
	dummySession.LogMethodParameter(
		dummyParameter1,
		dummyParameter2,
		dummyParameter3,
	)

	// verify
	verifyAll(t)
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
	createMock(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	logMethodLogicFuncExpected = 1
	logMethodLogicFunc = func(session *session, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logMethodLogicFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyLogLevel, logLevel)
		assert.Equal(t, dummyCategory, category)
		assert.Equal(t, dummySubcategory, subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 3, len(parameters))
		assert.Equal(t, dummyParameter1, parameters[0])
		assert.Equal(t, dummyParameter2, parameters[1])
		assert.Equal(t, dummyParameter3, parameters[2])
	}

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

	// verify
	verifyAll(t)
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
	createMock(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	getMethodNameFuncExpected = 1
	getMethodNameFunc = func() string {
		getMethodNameFuncCalled++
		return dummyMethodName
	}
	strconvItoaExpected = 3
	strconvItoa = func(i int) string {
		strconvItoaCalled++
		return strconv.Itoa(i)
	}
	logMethodReturnFuncExpected = 3
	logMethodReturnFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logMethodReturnFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyMethodName, category)
		assert.Equal(t, strconv.Itoa(logMethodReturnFuncCalled-1), subcategory)
		assert.Equal(t, "%v", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyReturns[logMethodReturnFuncCalled-1], parameters[0])
	}

	// act
	dummySession.LogMethodReturn(
		dummyReturn1,
		dummyReturn2,
		dummyReturn3,
	)

	// verify
	verifyAll(t)
}

func TestSessionLogMethodExit(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyMethodName = "some method name"

	// mock
	createMock(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// expect
	getMethodNameFuncExpected = 1
	getMethodNameFunc = func() string {
		getMethodNameFuncCalled++
		return dummyMethodName
	}
	logMethodExitFuncExpected = 1
	logMethodExitFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logMethodExitFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyMethodName, category)
		assert.Zero(t, subcategory)
		assert.Zero(t, messageFormat)
		assert.Empty(t, parameters)
	}

	// act
	dummySession.LogMethodExit()

	// verify
	verifyAll(t)
}

func TestSessionCreateWebcallRequest(t *testing.T) {
	// arrange
	var dummySessionID = uuid.New()
	var dummyMethod = "some method"
	var dummyURL = "some URL"
	var dummyPayload = "some payload"
	var dummyHeader = map[string]string{
		"foo":  "bar",
		"test": "123",
	}
	var dummySendClientCert = rand.Intn(100) < 50

	// mock
	createMock(t)

	// SUT
	var dummySession = &session{
		id: dummySessionID,
	}

	// act
	var result = dummySession.CreateWebcallRequest(
		dummyMethod,
		dummyURL,
		dummyPayload,
		dummyHeader,
		dummySendClientCert,
	)
	var webrequest, ok = result.(*webRequest)

	// assert
	assert.True(t, ok)
	assert.Equal(t, dummySession, webrequest.session)
	assert.Equal(t, dummyMethod, webrequest.method)
	assert.Equal(t, dummyURL, webrequest.url)
	assert.Equal(t, dummyPayload, webrequest.payload)
	assert.Equal(t, dummyHeader, webrequest.header)
	assert.Zero(t, webrequest.connRetry)
	assert.Nil(t, webrequest.httpRetry)
	assert.Equal(t, dummySendClientCert, webrequest.sendClientCert)
	assert.Zero(t, webrequest.retryDelay)

	// verify
	verifyAll(t)
}
