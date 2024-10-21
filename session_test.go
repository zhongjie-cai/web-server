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
	"github.com/zhongjie-cai/gomocker"
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
	m.ExpectFunc(isInterfaceValueNil, 1, func(i interface{}) bool {
		assert.Equal(t, dummyResponseWriterObject, i)
		return true
	})

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
	m.ExpectFunc(isInterfaceValueNil, 1, func(i interface{}) bool {
		assert.Equal(t, dummyResponseWriterObject, i)
		return false
	})

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
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageSessionNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	})

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
	m.ExpectFunc(getRequestBody, 1, func(httpRequest *http.Request) string {
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		return dummyRequestBody
	})
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeBadRequest, errorCode)
		assert.Equal(t, errorMessageRequestBodyEmpty, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	})

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
	m.ExpectFunc(getRequestBody, 1, func(httpRequest *http.Request) string {
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		return dummyRequestBody
	})
	m.ExpectFunc(logEndpointRequest, 2, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Body", category)
		if m.FuncCalledCount(logEndpointRequest) == 1 {
			assert.Equal(t, "Content", subcategory)
			assert.Equal(t, dummyRequestBody, messageFormat)
			assert.Empty(t, parameters)
		} else if m.FuncCalledCount(logEndpointRequest) == 2 {
			assert.Equal(t, "UnmarshalError", subcategory)
			assert.Equal(t, "%+v", messageFormat)
			assert.Equal(t, 1, len(parameters))
			assert.Equal(t, dummyError, parameters[0])
		}
	})
	m.ExpectFunc(tryUnmarshal, 1, func(value string, dataTemplate interface{}) error {
		assert.Equal(t, dummyRequestBody, value)
		*(dataTemplate.(*int)) = dummyResult
		return dummyError
	})
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeBadRequest, errorCode)
		assert.Equal(t, errorMessageRequestBodyInvalid, errorMessage)
		assert.Equal(t, 1, len(innerErrors))
		assert.Equal(t, dummyError, innerErrors[0])
		return dummyAppError
	})

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
	m.ExpectFunc(getRequestBody, 1, func(httpRequest *http.Request) string {
		assert.Equal(t, dummyHTTPRequest, httpRequest)
		return dummyRequestBody
	})
	m.ExpectFunc(logEndpointRequest, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Body", category)
		assert.Equal(t, "Content", subcategory)
		assert.Equal(t, dummyRequestBody, messageFormat)
		assert.Empty(t, parameters)
	})
	m.ExpectFunc(tryUnmarshal, 1, func(value string, dataTemplate interface{}) error {
		assert.Equal(t, dummyRequestBody, value)
		*(dataTemplate.(*int)) = dummyResult
		return nil
	})

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
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageSessionNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	})

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
	m.ExpectFunc(mux.Vars, 1, func(r *http.Request) map[string]string {
		assert.Equal(t, dummyHTTPRequest, r)
		return dummyParameters
	})
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeBadRequest, errorCode)
		assert.Equal(t, errorMessageParameterNotFound, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	})

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
	m.ExpectFunc(mux.Vars, 1, func(r *http.Request) map[string]string {
		assert.Equal(t, dummyHTTPRequest, r)
		return dummyParameters
	})
	m.ExpectFunc(logEndpointRequest, 2, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Parameter", category)
		if m.FuncCalledCount(logEndpointRequest) == 1 {
			assert.Equal(t, dummyName, subcategory)
			assert.Equal(t, dummyValue, messageFormat)
			assert.Empty(t, parameters)
		} else if m.FuncCalledCount(logEndpointRequest) == 2 {
			assert.Equal(t, "UnmarshalError", subcategory)
			assert.Equal(t, "%+v", messageFormat)
			assert.Equal(t, 1, len(parameters))
			assert.Equal(t, dummyError, parameters[0])
		}
	})
	m.ExpectFunc(tryUnmarshal, 1, func(value string, dataTemplate interface{}) error {
		assert.Equal(t, dummyValue, value)
		*(dataTemplate.(*int)) = dummyResult
		return dummyError
	})
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeBadRequest, errorCode)
		assert.Equal(t, errorMessageParameterInvalid, errorMessage)
		assert.Equal(t, 1, len(innerErrors))
		assert.Equal(t, dummyError, innerErrors[0])
		return dummyAppError
	})

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
	m.ExpectFunc(mux.Vars, 1, func(r *http.Request) map[string]string {
		assert.Equal(t, dummyHTTPRequest, r)
		return dummyParameters
	})
	m.ExpectFunc(logEndpointRequest, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Parameter", category)
		assert.Equal(t, dummyName, subcategory)
		assert.Equal(t, dummyValue, messageFormat)
		assert.Empty(t, parameters)
	})
	m.ExpectFunc(tryUnmarshal, 1, func(value string, dataTemplate interface{}) error {
		assert.Equal(t, dummyValue, value)
		*(dataTemplate.(*int)) = dummyResult
		return nil
	})

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
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageSessionNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	})

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
	m.ExpectFunc(getAllQueries, 1, func(session *session, name string) []string {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyName, name)
		return dummyQueries
	})
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeBadRequest, errorCode)
		assert.Equal(t, errorMessageQueryNotFound, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	})

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
	m.ExpectFunc(getAllQueries, 1, func(session *session, name string) []string {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyName, name)
		return dummyQueries
	})
	m.ExpectFunc(logEndpointRequest, 2, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Query", category)
		if m.FuncCalledCount(logEndpointRequest) == 1 {
			assert.Equal(t, dummyName, subcategory)
			assert.Equal(t, dummyQueries[dummyIndex], messageFormat)
			assert.Empty(t, parameters)
		} else if m.FuncCalledCount(logEndpointRequest) == 2 {
			assert.Equal(t, "UnmarshalError", subcategory)
			assert.Equal(t, "%+v", messageFormat)
			assert.Equal(t, 1, len(parameters))
			assert.Equal(t, dummyError, parameters[0])
		}
	})
	m.ExpectFunc(tryUnmarshal, 1, func(value string, dataTemplate interface{}) error {
		assert.Equal(t, dummyQueries[dummyIndex], value)
		*(dataTemplate.(*int)) = dummyResult
		return dummyError
	})
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeBadRequest, errorCode)
		assert.Equal(t, errorMessageQueryInvalid, errorMessage)
		assert.Equal(t, 1, len(innerErrors))
		assert.Equal(t, dummyError, innerErrors[0])
		return dummyAppError
	})

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
	m.ExpectFunc(getAllQueries, 1, func(session *session, name string) []string {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyName, name)
		return dummyQueries
	})
	m.ExpectFunc(logEndpointRequest, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Query", category)
		assert.Equal(t, dummyName, subcategory)
		assert.Equal(t, dummyQueries[dummyIndex], messageFormat)
		assert.Empty(t, parameters)
	})
	m.ExpectFunc(tryUnmarshal, 1, func(value string, dataTemplate interface{}) error {
		assert.Equal(t, dummyQueries[dummyIndex], value)
		*(dataTemplate.(*int)) = dummyResult
		return nil
	})

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
	m.ExpectFunc(textproto.CanonicalMIMEHeaderKey, 1, func(s string) string {
		assert.Equal(t, dummyName, s)
		return dummyCanonicalName
	})

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
	m.ExpectFunc(textproto.CanonicalMIMEHeaderKey, 1, func(s string) string {
		assert.Equal(t, dummyName, s)
		return dummyCanonicalName
	})

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
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageSessionNil, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	})

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
	m.ExpectFunc(getAllHeaders, 1, func(session *session, name string) []string {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyName, name)
		return dummyHeaders
	})
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeBadRequest, errorCode)
		assert.Equal(t, errorMessageHeaderNotFound, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	})

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
	m.ExpectFunc(getAllHeaders, 1, func(session *session, name string) []string {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyName, name)
		return dummyHeaders
	})
	m.ExpectFunc(logEndpointRequest, 2, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Header", category)
		if m.FuncCalledCount(logEndpointRequest) == 1 {
			assert.Equal(t, dummyName, subcategory)
			assert.Equal(t, dummyHeaders[dummyIndex], messageFormat)
			assert.Empty(t, parameters)
		} else if m.FuncCalledCount(logEndpointRequest) == 2 {
			assert.Equal(t, "UnmarshalError", subcategory)
			assert.Equal(t, "%+v", messageFormat)
			assert.Equal(t, 1, len(parameters))
			assert.Equal(t, dummyError, parameters[0])
		}
	})
	m.ExpectFunc(tryUnmarshal, 1, func(value string, dataTemplate interface{}) error {
		assert.Equal(t, dummyHeaders[dummyIndex], value)
		*(dataTemplate.(*int)) = dummyResult
		return dummyError
	})
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeBadRequest, errorCode)
		assert.Equal(t, errorMessageHeaderInvalid, errorMessage)
		assert.Equal(t, 1, len(innerErrors))
		assert.Equal(t, dummyError, innerErrors[0])
		return dummyAppError
	})

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
	m.ExpectFunc(getAllHeaders, 1, func(session *session, name string) []string {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyName, name)
		return dummyHeaders
	})
	m.ExpectFunc(logEndpointRequest, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "Header", category)
		assert.Equal(t, dummyName, subcategory)
		assert.Equal(t, dummyHeaders[dummyIndex], messageFormat)
		assert.Empty(t, parameters)
	})
	m.ExpectFunc(tryUnmarshal, 1, func(value string, dataTemplate interface{}) error {
		assert.Equal(t, dummyHeaders[dummyIndex], value)
		*(dataTemplate.(*int)) = dummyResult
		return nil
	})

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
	m.ExpectFunc(json.Marshal, 1, func(v interface{}) ([]byte, error) {
		assert.Equal(t, dummyValue, v)
		return nil, errors.New("some marshal error")
	})

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
	m.ExpectFunc(runtime.Caller, 1, func(skip int) (pc uintptr, file string, line int, ok bool) {
		assert.Equal(t, 3, skip)
		return dummyPC, dummyFile, dummyLine, dummyOK
	})

	// SUT + act
	var result = getMethodName()

	// assert
	assert.Equal(t, "?", result)
}

func TestSessionGetMethodName_HappyPath(t *testing.T) {
	// arrange
	var dummyPC = uintptr(rand.Int())
	var dummyFile = "some file"
	var dummyLine = rand.Int()
	var dummyOK = true
	var dummyFn = &runtime.Func{}
	var dummyName = "some name"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(runtime.Caller, 1, func(skip int) (pc uintptr, file string, line int, ok bool) {
		assert.Equal(t, 3, skip)
		return dummyPC, dummyFile, dummyLine, dummyOK
	})
	m.ExpectFunc(runtime.FuncForPC, 1, func(pc uintptr) *runtime.Func {
		assert.Equal(t, dummyPC, pc)
		return dummyFn
	})
	m.ExpectMethod(dummyFn, "Name", 1, func(self *runtime.Func) string {
		assert.Equal(t, dummyFn, self)
		return dummyName
	})

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
	m.ExpectFunc(getMethodName, 1, func() string {
		return dummyMethodName
	})
	m.ExpectFunc(logMethodEnter, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyMethodName, category)
		assert.Zero(t, subcategory)
		assert.Zero(t, messageFormat)
		assert.Empty(t, parameters)
	})

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
	m.ExpectFunc(getMethodName, 1, func() string {
		return dummyMethodName
	})
	m.ExpectFunc(logMethodParameter, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyMethodName, category)
		assert.Equal(t, "0", subcategory)
		assert.Equal(t, "%v", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyParameters[0], parameters[0])
	}).ExpectFunc(logMethodParameter, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyMethodName, category)
		assert.Equal(t, "1", subcategory)
		assert.Equal(t, "%v", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyParameters[1], parameters[0])
	}).ExpectFunc(logMethodParameter, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyMethodName, category)
		assert.Equal(t, "2", subcategory)
		assert.Equal(t, "%v", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyParameters[2], parameters[0])
	})

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
	m.ExpectFunc(logMethodLogic, 1, func(session *session, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyLogLevel, logLevel)
		assert.Equal(t, dummyCategory, category)
		assert.Equal(t, dummySubcategory, subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 3, len(parameters))
		assert.Equal(t, dummyParameter1, parameters[0])
		assert.Equal(t, dummyParameter2, parameters[1])
		assert.Equal(t, dummyParameter3, parameters[2])
	})

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
	m.ExpectFunc(getMethodName, 1, func() string {
		return dummyMethodName
	})
	m.ExpectFunc(logMethodReturn, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyMethodName, category)
		assert.Equal(t, "0", subcategory)
		assert.Equal(t, "%v", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyReturns[0], parameters[0])
	}).ExpectFunc(logMethodReturn, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyMethodName, category)
		assert.Equal(t, "1", subcategory)
		assert.Equal(t, "%v", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyReturns[1], parameters[0])
	}).ExpectFunc(logMethodReturn, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyMethodName, category)
		assert.Equal(t, "2", subcategory)
		assert.Equal(t, "%v", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyReturns[2], parameters[0])
	})

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
	m.ExpectFunc(getMethodName, 1, func() string {
		return dummyMethodName
	})
	m.ExpectFunc(logMethodExit, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyMethodName, category)
		assert.Zero(t, subcategory)
		assert.Zero(t, messageFormat)
		assert.Empty(t, parameters)
	})

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
