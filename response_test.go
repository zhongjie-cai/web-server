package webserver

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyCustomizationConstructResponse struct {
	dummyCustomization
	interpretError   func(err error) (int, string)
	interpretSuccess func(responseContent interface{}) (int, string)
}

func (customization *dummyCustomizationConstructResponse) InterpretError(err error) (int, string) {
	if customization.interpretError != nil {
		return customization.interpretError(err)
	}
	assert.Fail(customization.t, "Unexpected call to InterpretError")
	return 0, ""
}

func (customization *dummyCustomizationConstructResponse) InterpretSuccess(responseContent interface{}) (int, string) {
	if customization.interpretSuccess != nil {
		return customization.interpretSuccess(responseContent)
	}
	assert.Fail(customization.t, "Unexpected call to InterpretSuccess")
	return 0, ""
}

func TestConstructResponse_ResponseError(t *testing.T) {
	// arrange
	var dummyCustomizationConstructResponse = &dummyCustomizationConstructResponse{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		customization: dummyCustomizationConstructResponse,
	}
	var dummyResponseObject = rand.Int()
	var dummyResponseError = errors.New("some response error")
	var customizationInterpretErrorExpected int
	var customizationInterpretErrorCalled int
	var customizationInterpretSuccessExpected int
	var customizationInterpretSuccessCalled int
	var dummyCode = rand.Int()
	var dummyMessage = "some message"

	// mock
	createMock(t)

	// expect
	customizationInterpretErrorExpected = 1
	dummyCustomizationConstructResponse.interpretError = func(err error) (int, string) {
		customizationInterpretErrorCalled++
		assert.Equal(t, dummyResponseError, err)
		return dummyCode, dummyMessage
	}

	// SUT + act
	var code, message = constructResponse(
		dummySession,
		dummyResponseObject,
		dummyResponseError,
	)

	// assert
	assert.Equal(t, dummyCode, code)
	assert.Equal(t, dummyMessage, message)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationInterpretErrorExpected, customizationInterpretErrorCalled, "Unexpected number of calls to method customization.InterpretError")
	assert.Equal(t, customizationInterpretSuccessExpected, customizationInterpretSuccessCalled, "Unexpected number of calls to method customization.InterpretSuccess")
}

func TestConstructResponse_Success(t *testing.T) {
	// arrange
	var dummyCustomizationConstructResponse = &dummyCustomizationConstructResponse{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		customization: dummyCustomizationConstructResponse,
	}
	var dummyResponseObject = rand.Int()
	var customizationInterpretErrorExpected int
	var customizationInterpretErrorCalled int
	var customizationInterpretSuccessExpected int
	var customizationInterpretSuccessCalled int
	var dummyCode = rand.Int()
	var dummyMessage = "some message"

	// mock
	createMock(t)

	// expect
	customizationInterpretSuccessExpected = 1
	dummyCustomizationConstructResponse.interpretSuccess = func(responseObject interface{}) (int, string) {
		customizationInterpretSuccessCalled++
		assert.Equal(t, dummyResponseObject, responseObject)
		return dummyCode, dummyMessage
	}

	// SUT + act
	var code, message = constructResponse(
		dummySession,
		dummyResponseObject,
		nil,
	)

	// assert
	assert.Equal(t, dummyCode, code)
	assert.Equal(t, dummyMessage, message)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationInterpretErrorExpected, customizationInterpretErrorCalled, "Unexpected number of calls to method customization.InterpretError")
	assert.Equal(t, customizationInterpretSuccessExpected, customizationInterpretSuccessCalled, "Unexpected number of calls to method customization.InterpretSuccess")
}

type dummyResponseWriterWriteResponse struct {
	dummyResponseWriter
	header      func() http.Header
	write       func([]byte) (int, error)
	writeHeader func(int)
}

func (drw *dummyResponseWriterWriteResponse) Header() http.Header {
	if drw.header != nil {
		return drw.header()
	}
	assert.Fail(drw.t, "Unexpected number of calls to Header")
	return nil
}

func (drw *dummyResponseWriterWriteResponse) Write(bytes []byte) (int, error) {
	if drw.write != nil {
		return drw.write(bytes)
	}
	assert.Fail(drw.t, "Unexpected number of calls to Write")
	return 0, nil
}

func (drw *dummyResponseWriterWriteResponse) WriteHeader(statusCode int) {
	if drw.writeHeader != nil {
		drw.writeHeader(statusCode)
		return
	}
	assert.Fail(drw.t, "Unexpected number of calls to WriteHeader")
}

func TestWriteResponse(t *testing.T) {
	// arrange
	var dummyResponseWriterWriteResponse = &dummyResponseWriterWriteResponse{
		dummyResponseWriter: dummyResponseWriter{t: t},
	}
	var dummySession = &session{
		responseWriter: dummyResponseWriterWriteResponse,
	}
	var dummyResponseObject = rand.Int()
	var dummyResponseError = errors.New("some response error")
	var responseWriterHeaderExpected int
	var responseWriterHeaderCalled int
	var responseWriterWriteHeaderExpected int
	var responseWriterWriteHeaderCalled int
	var responseWriterWriteExpected int
	var responseWriterWriteCalled int
	var dummyCode = rand.Int()
	var dummyMessage = "some message"
	var dummyStatusText = "some status text"
	var dummyCodeString = strconv.Itoa(dummyCode)
	var dummyHeader = make(http.Header)

	// mock
	createMock(t)

	// expect
	constructResponseFuncExpected = 1
	constructResponseFunc = func(session *session, responseObject interface{}, responseError error) (int, string) {
		constructResponseFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyResponseObject, responseObject)
		assert.Equal(t, dummyResponseError, responseError)
		return dummyCode, dummyMessage
	}
	httpStatusTextExpected = 1
	httpStatusText = func(code int) string {
		httpStatusTextCalled++
		assert.Equal(t, dummyCode, code)
		return dummyStatusText
	}
	strconvItoaExpected = 1
	strconvItoa = func(i int) string {
		strconvItoaCalled++
		assert.Equal(t, dummyCode, i)
		return dummyCodeString
	}
	logEndpointResponseFuncExpected = 1
	logEndpointResponseFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logEndpointResponseFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyStatusText, category)
		assert.Equal(t, dummyCodeString, subcategory)
		assert.Equal(t, dummyMessage, messageFormat)
		assert.Empty(t, parameters)
	}
	isInterfaceValueNilFuncExpected = 1
	isInterfaceValueNilFunc = func(i interface{}) bool {
		isInterfaceValueNilFuncCalled++
		assert.Equal(t, dummyResponseWriterWriteResponse, i)
		return false
	}
	responseWriterHeaderExpected = 1
	dummyResponseWriterWriteResponse.header = func() http.Header {
		responseWriterHeaderCalled++
		return dummyHeader
	}
	responseWriterWriteHeaderExpected = 1
	dummyResponseWriterWriteResponse.writeHeader = func(statusCode int) {
		responseWriterWriteHeaderCalled++
		assert.Equal(t, dummyCode, statusCode)
	}
	responseWriterWriteExpected = 1
	dummyResponseWriterWriteResponse.write = func(bytes []byte) (int, error) {
		responseWriterWriteCalled++
		assert.Equal(t, []byte(dummyMessage), bytes)
		return rand.Int(), errors.New("some error")
	}

	// SUT + act
	writeResponse(
		dummySession,
		dummyResponseObject,
		dummyResponseError,
	)

	// assert
	assert.Equal(t, 1, len(dummyHeader))
	assert.Equal(t, ContentTypeJSON, dummyHeader.Get("Content-Type"))

	// verify
	verifyAll(t)
	assert.Equal(t, responseWriterHeaderExpected, responseWriterHeaderCalled, "Unexpected number of calls to method responseWriter.Header")
	assert.Equal(t, responseWriterWriteHeaderExpected, responseWriterWriteHeaderCalled, "Unexpected number of calls to method responseWriter.WriteHeader")
	assert.Equal(t, responseWriterWriteExpected, responseWriterWriteCalled, "Unexpected number of calls to method responseWriter.Write")
}
