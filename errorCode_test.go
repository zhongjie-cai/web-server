package webserver

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorCodeEnumHTTPStatusCode_GeneralFailure(t *testing.T) {
	// mock
	createMock(t)

	// SUT
	var dummyErrorCode = errorCodeGeneralFailure

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusInternalServerError, result)

	// verify
	verifyAll(t)
}

func TestErrorCodeEnumHTTPStatusCode_Unauthorized(t *testing.T) {
	// mock
	createMock(t)

	// SUT
	var dummyErrorCode = errorCodeUnauthorized

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusUnauthorized, result)

	// verify
	verifyAll(t)
}

func TestErrorCodeEnumHTTPStatusCode_InvalidOperation(t *testing.T) {
	// mock
	createMock(t)

	// SUT
	var dummyErrorCode = errorCodeInvalidOperation

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusMethodNotAllowed, result)

	// verify
	verifyAll(t)
}

func TestErrorCodeEnumHTTPStatusCode_BadRequest(t *testing.T) {
	// mock
	createMock(t)

	// SUT
	var dummyErrorCode = errorCodeBadRequest

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusBadRequest, result)

	// verify
	verifyAll(t)
}

func TestErrorCodeEnumHTTPStatusCode_NotFound(t *testing.T) {
	// mock
	createMock(t)

	// SUT
	var dummyErrorCode = errorCodeNotFound

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusNotFound, result)

	// verify
	verifyAll(t)
}

func TestErrorCodeEnumHTTPStatusCode_CircuitBreak(t *testing.T) {
	// mock
	createMock(t)

	// SUT
	var dummyErrorCode = errorCodeCircuitBreak

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusForbidden, result)

	// verify
	verifyAll(t)
}

func TestErrorCodeEnumHTTPStatusCode_OperationLock(t *testing.T) {
	// mock
	createMock(t)

	// SUT
	var dummyErrorCode = errorCodeOperationLock

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusLocked, result)

	// verify
	verifyAll(t)
}

func TestErrorCodeEnumHTTPStatusCode_AccessForbidden(t *testing.T) {
	// mock
	createMock(t)

	// SUT
	var dummyErrorCode = errorCodeAccessForbidden

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusForbidden, result)

	// verify
	verifyAll(t)
}

func TestErrorCodeEnumHTTPStatusCode_DataCorruption(t *testing.T) {
	// mock
	createMock(t)

	// SUT
	var dummyErrorCode = errorCodeDataCorruption

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusConflict, result)

	// verify
	verifyAll(t)
}

func TestErrorCodeEnumHTTPStatusCode_NotImplemented(t *testing.T) {
	// mock
	createMock(t)

	// SUT
	var dummyErrorCode = errorCodeNotImplemented

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusNotImplemented, result)

	// verify
	verifyAll(t)
}

func TestErrorCodeEnumHTTPStatusCode_OtherErrorCode(t *testing.T) {
	// mock
	createMock(t)

	// SUT
	var dummyErrorCode = errorCode("some other error code")

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusInternalServerError, result)

	// verify
	verifyAll(t)
}
