package webserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_SessionNil(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = errSessionNil.Error()

	// assert
	assert.Equal(t, "The session object is nil", message)

	// verify
	verifyAll(t)
}

func TestError_RouteRegistration(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = errRouteRegistration.Error()

	// assert
	assert.Equal(t, "The route registration failed", message)

	// verify
	verifyAll(t)
}

func TestError_RouteNotFound(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = errRouteNotFound.Error()

	// assert
	assert.Equal(t, "The route is not found", message)

	// verify
	verifyAll(t)
}

func TestError_RouteNotImplemented(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = errRouteNotImplemented.Error()

	// assert
	assert.Equal(t, "The route is not implemented", message)

	// verify
	verifyAll(t)
}

func TestError_HostServer(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = errHostServer.Error()

	// assert
	assert.Equal(t, "The server hosting failed", message)

	// verify
	verifyAll(t)
}

func TestError_RequestBodyEmpty(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = ErrRequestBodyEmpty.Error()

	// assert
	assert.Equal(t, "The request body is empty", message)

	// verify
	verifyAll(t)
}

func TestError_RequestBodyInvalid(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = ErrRequestBodyInvalid.Error()

	// assert
	assert.Equal(t, "The request body is invalid", message)

	// verify
	verifyAll(t)
}

func TestError_ParameterNotFound(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = ErrParameterNotFound.Error()

	// assert
	assert.Equal(t, "The request parameter is not found", message)

	// verify
	verifyAll(t)
}

func TestError_ParameterInvalid(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = ErrParameterInvalid.Error()

	// assert
	assert.Equal(t, "The request parameter is invalid", message)

	// verify
	verifyAll(t)
}

func TestError_QueryNotFound(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = ErrQueryNotFound.Error()

	// assert
	assert.Equal(t, "The request query is not found", message)

	// verify
	verifyAll(t)
}

func TestError_QueryInvalid(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = ErrQueryInvalid.Error()

	// assert
	assert.Equal(t, "The request query is invalid", message)

	// verify
	verifyAll(t)
}

func TestError_HeaderNotFound(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = ErrHeaderNotFound.Error()

	// assert
	assert.Equal(t, "The request header is not found", message)

	// verify
	verifyAll(t)
}

func TestError_HeaderInvalid(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = ErrHeaderInvalid.Error()

	// assert
	assert.Equal(t, "The request header is invalid", message)

	// verify
	verifyAll(t)
}

func TestError_WebRequestNil(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = ErrWebRequestNil.Error()

	// assert
	assert.Equal(t, "The web request object is nil", message)

	// verify
	verifyAll(t)
}

func TestError_ResponseInvalid(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = ErrResponseInvalid.Error()

	// assert
	assert.Equal(t, "The response body is invalid", message)

	// verify
	verifyAll(t)
}

func TestError_InvalidOperation(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = ErrInvalidOperation.Error()

	// assert
	assert.Equal(t, "The underlying operation is invalid", message)

	// verify
	verifyAll(t)
}

func TestError_Forbidden(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = ErrForbidden.Error()

	// assert
	assert.Equal(t, "The underlying operation is forbidden", message)

	// verify
	verifyAll(t)
}

func TestError_NotImplemented(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = ErrNotImplemented.Error()

	// assert
	assert.Equal(t, "The underlying operation is not implemented", message)

	// verify
	verifyAll(t)
}

func TestError_BadRequest(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = ErrBadRequest.Error()

	// assert
	assert.Equal(t, "The client request is invalid", message)

	// verify
	verifyAll(t)
}

func TestError_ResourceNotFound(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = ErrResourceNotFound.Error()

	// assert
	assert.Equal(t, "The requested resource is not found", message)

	// verify
	verifyAll(t)
}

func TestError_ResourceLocked(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = ErrResourceLocked.Error()

	// assert
	assert.Equal(t, "The requested resource is locked", message)

	// verify
	verifyAll(t)
}

func TestError_ResourceConflict(t *testing.T) {
	// mock
	createMock(t)

	// SUT + act
	var message = ErrResourceConflict.Error()

	// assert
	assert.Equal(t, "The requested resource has data conflicts", message)

	// verify
	verifyAll(t)
}
