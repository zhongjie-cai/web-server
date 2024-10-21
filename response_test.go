package webserver

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhongjie-cai/gomocker"
)

func TestSkipResponseHandling(t *testing.T) {
	// SUT + act
	var result, err = SkipResponseHandling()

	// assert
	assert.IsType(t, skipResponseHandlingDummy{}, result)
	assert.NoError(t, err)
}

func TestShouldSkipHandling_HasError(t *testing.T) {
	// arrange
	var dummyResponseObject skipResponseHandlingDummy
	var dummyResponseError = errors.New("some error")

	// SUT + act
	var result = shouldSkipHandling(
		dummyResponseObject,
		dummyResponseError,
	)

	// assert
	assert.False(t, result)
}

func TestShouldSkipHandling_Yes(t *testing.T) {
	// arrange
	var dummyResponseObject skipResponseHandlingDummy
	var dummyResponseError error

	// SUT + act
	var result = shouldSkipHandling(
		dummyResponseObject,
		dummyResponseError,
	)

	// assert
	assert.True(t, result)
}

func TestShouldSkipHandling_No(t *testing.T) {
	// arrange
	var dummyResponseObject = rand.Int()
	var dummyResponseError error

	// SUT + act
	var result = shouldSkipHandling(
		dummyResponseObject,
		dummyResponseError,
	)

	// assert
	assert.False(t, result)
}

func TestConstructResponse_ResponseError(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	var dummyResponseObject = rand.Int()
	var dummyResponseError = errors.New("some response error")
	var dummyCode = rand.Int()
	var dummyMessage = "some message"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "InterpretError", 1, func(self *DefaultCustomization, err error) (int, string) {
		assert.Equal(t, dummyCustomization, self)
		assert.Equal(t, dummyResponseError, err)
		return dummyCode, dummyMessage
	})

	// SUT + act
	var code, message = constructResponse(
		dummySession,
		dummyResponseObject,
		dummyResponseError,
	)

	// assert
	assert.Equal(t, dummyCode, code)
	assert.Equal(t, dummyMessage, message)
}

func TestConstructResponse_Success(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		customization: dummyCustomization,
	}
	var dummyResponseObject = rand.Int()
	var dummyCode = rand.Int()
	var dummyMessage = "some message"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "InterpretSuccess", 1, func(self *DefaultCustomization, responseObject interface{}) (int, string) {
		assert.Equal(t, dummyCustomization, self)
		assert.Equal(t, dummyResponseObject, responseObject)
		return dummyCode, dummyMessage
	})

	// SUT + act
	var code, message = constructResponse(
		dummySession,
		dummyResponseObject,
		nil,
	)

	// assert
	assert.Equal(t, dummyCode, code)
	assert.Equal(t, dummyMessage, message)
}

func TestWriteResponse_SkipHandling(t *testing.T) {
	// arrange
	var dummyResponseWriter = &dummyResponseWriter{}
	var dummySession = &session{
		responseWriter: dummyResponseWriter,
	}
	var dummyResponseObject = rand.Int()
	var dummyResponseError = errors.New("some response error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(shouldSkipHandling, 1, func(responseObject interface{}, responseError error) bool {
		assert.Equal(t, dummyResponseObject, responseObject)
		assert.Equal(t, dummyResponseError, responseError)
		return true
	})
	m.ExpectFunc(logEndpointResponse, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "None", category)
		assert.Equal(t, "-1", subcategory)
		assert.Equal(t, "Skipped response handling", messageFormat)
		assert.Empty(t, parameters)
	})

	// SUT + act
	writeResponse(
		dummySession,
		dummyResponseObject,
		dummyResponseError,
	)
}

func TestWriteResponse_HappyPath(t *testing.T) {
	// arrange
	var dummyResponseWriterInstance = &dummyResponseWriter{}
	var dummySession = &session{
		responseWriter: dummyResponseWriterInstance,
	}
	var dummyResponseObject = rand.Int()
	var dummyResponseError = errors.New("some response error")
	var dummyCode = rand.Int()
	var dummyMessage = "some message"
	var dummyStatusText = "some status text"
	var dummyCodeString = strconv.Itoa(dummyCode)
	var dummyHeader = make(http.Header)

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(shouldSkipHandling, 1, func(responseObject interface{}, responseError error) bool {
		assert.Equal(t, dummyResponseObject, responseObject)
		assert.Equal(t, dummyResponseError, responseError)
		return false
	})
	m.ExpectFunc(constructResponse, 1, func(session *session, responseObject interface{}, responseError error) (int, string) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyResponseObject, responseObject)
		assert.Equal(t, dummyResponseError, responseError)
		return dummyCode, dummyMessage
	})
	m.ExpectFunc(http.StatusText, 1, func(code int) string {
		assert.Equal(t, dummyCode, code)
		return dummyStatusText
	})
	m.ExpectFunc(strconv.Itoa, 1, func(i int) string {
		assert.Equal(t, dummyCode, i)
		return dummyCodeString
	})
	m.ExpectFunc(logEndpointResponse, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyStatusText, category)
		assert.Equal(t, dummyCodeString, subcategory)
		assert.Equal(t, dummyMessage, messageFormat)
		assert.Empty(t, parameters)
	})
	m.ExpectFunc(isInterfaceValueNil, 1, func(i interface{}) bool {
		assert.Equal(t, dummyResponseWriterInstance, i)
		return false
	})
	m.ExpectMethod(dummyResponseWriterInstance, "Header", 1, func(self *dummyResponseWriter) http.Header {
		assert.Equal(t, dummyResponseWriterInstance, self)
		return dummyHeader
	})
	m.ExpectMethod(dummyResponseWriterInstance, "WriteHeader", 1, func(self *dummyResponseWriter, statusCode int) {
		assert.Equal(t, dummyResponseWriterInstance, self)
		assert.Equal(t, dummyCode, statusCode)
	})
	m.ExpectMethod(dummyResponseWriterInstance, "Write", 1, func(self *dummyResponseWriter, bytes []byte) (int, error) {
		assert.Equal(t, dummyResponseWriterInstance, self)
		assert.Equal(t, []byte(dummyMessage), bytes)
		return rand.Int(), errors.New("some error")
	})

	// SUT + act
	writeResponse(
		dummySession,
		dummyResponseObject,
		dummyResponseError,
	)

	// assert
	assert.Equal(t, 1, len(dummyHeader))
	assert.Equal(t, ContentTypeJSON, dummyHeader.Get("Content-Type"))
}
