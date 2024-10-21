package webserver

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorCodeEnumHTTPStatusCode_GeneralFailure(t *testing.T) {
	// SUT
	var dummyErrorCode = errorCodeGeneralFailure

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusInternalServerError, result)
}

func TestErrorCodeEnumHTTPStatusCode_Unauthorized(t *testing.T) {
	// SUT
	var dummyErrorCode = errorCodeUnauthorized

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusUnauthorized, result)
}

func TestErrorCodeEnumHTTPStatusCode_InvalidOperation(t *testing.T) {
	// SUT
	var dummyErrorCode = errorCodeInvalidOperation

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusMethodNotAllowed, result)
}

func TestErrorCodeEnumHTTPStatusCode_BadRequest(t *testing.T) {
	// SUT
	var dummyErrorCode = errorCodeBadRequest

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusBadRequest, result)
}

func TestErrorCodeEnumHTTPStatusCode_NotFound(t *testing.T) {
	// SUT
	var dummyErrorCode = errorCodeNotFound

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusNotFound, result)
}

func TestErrorCodeEnumHTTPStatusCode_CircuitBreak(t *testing.T) {
	// SUT
	var dummyErrorCode = errorCodeCircuitBreak

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusForbidden, result)
}

func TestErrorCodeEnumHTTPStatusCode_OperationLock(t *testing.T) {
	// SUT
	var dummyErrorCode = errorCodeOperationLock

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusLocked, result)
}

func TestErrorCodeEnumHTTPStatusCode_AccessForbidden(t *testing.T) {
	// SUT
	var dummyErrorCode = errorCodeAccessForbidden

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusForbidden, result)
}

func TestErrorCodeEnumHTTPStatusCode_DataCorruption(t *testing.T) {
	// SUT
	var dummyErrorCode = errorCodeDataCorruption

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusConflict, result)
}

func TestErrorCodeEnumHTTPStatusCode_NotImplemented(t *testing.T) {
	// SUT
	var dummyErrorCode = errorCodeNotImplemented

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusNotImplemented, result)
}

func TestErrorCodeEnumHTTPStatusCode_OtherErrorCode(t *testing.T) {
	// SUT
	var dummyErrorCode = errorCode("some other error code")

	// act
	var result = dummyErrorCode.httpStatusCode()

	// assert
	assert.Equal(t, http.StatusInternalServerError, result)
}
