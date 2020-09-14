package webserver

import (
	"net/http"
)

// web server built-in error messages
const (
	errorMessageSessionNil         = "The session object is nil"
	errorMessageRouteRegistration  = "The route registration failed"
	errorMessageHostServer         = "The server hosting failed"
	errorMessageRequestBodyEmpty   = "The request body is empty"
	errorMessageRequestBodyInvalid = "The request body is invalid"
	errorMessageParameterNotFound  = "The request parameter is not found"
	errorMessageParameterInvalid   = "The request parameter is invalid"
	errorMessageQueryNotFound      = "The request query is not found"
	errorMessageQueryInvalid       = "The request query is invalid"
	errorMessageHeaderNotFound     = "The request header is not found"
	errorMessageHeaderInvalid      = "The request header is invalid"
	errorMessageWebRequestNil      = "The web request object is nil"
	errorMessageResponseInvalid    = "The response body is invalid"
)

type errorCode string

// These are the integer values of the enum that corresponds to a given error
const (
	errorCodeGeneralFailure   errorCode = "GeneralFailure"
	errorCodeUnauthorized     errorCode = "Unauthorized"
	errorCodeInvalidOperation errorCode = "InvalidOperation"
	errorCodeBadRequest       errorCode = "BadRequest"
	errorCodeNotFound         errorCode = "NotFound"
	errorCodeCircuitBreak     errorCode = "CircuitBreak"
	errorCodeOperationLock    errorCode = "OperationLock"
	errorCodeAccessForbidden  errorCode = "AccessForbidden"
	errorCodeDataCorruption   errorCode = "DataCorruption"
	errorCodeNotImplemented   errorCode = "NotImplemented"
)

func (errorCode errorCode) httpStatusCode() int {
	var statusCode int
	switch errorCode {
	case errorCodeGeneralFailure:
		statusCode = http.StatusInternalServerError
	case errorCodeUnauthorized:
		statusCode = http.StatusUnauthorized
	case errorCodeInvalidOperation:
		statusCode = http.StatusMethodNotAllowed
	case errorCodeBadRequest:
		statusCode = http.StatusBadRequest
	case errorCodeNotFound:
		statusCode = http.StatusNotFound
	case errorCodeCircuitBreak:
		statusCode = http.StatusForbidden
	case errorCodeOperationLock:
		statusCode = http.StatusLocked
	case errorCodeAccessForbidden:
		statusCode = http.StatusForbidden
	case errorCodeDataCorruption:
		statusCode = http.StatusConflict
	case errorCodeNotImplemented:
		statusCode = http.StatusNotImplemented
	default:
		statusCode = http.StatusInternalServerError
	}
	return statusCode
}
