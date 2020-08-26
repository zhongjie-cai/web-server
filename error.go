package webserver

import "errors"

// internal errors
var (
	errSessionNil          error = errors.New("The session object is nil")
	errRouteRegistration   error = errors.New("The route registration failed")
	errRouteNotFound       error = errors.New("The route is not found")
	errRouteNotImplemented error = errors.New("The route is not implemented")
	errHostServer          error = errors.New("The server hosting failed")
)

// built-in errors
var (
	ErrRequestBodyEmpty   error = errors.New("The request body is empty")
	ErrRequestBodyInvalid error = errors.New("The request body is invalid")
	ErrParameterNotFound  error = errors.New("The request parameter is not found")
	ErrParameterInvalid   error = errors.New("The request parameter is invalid")
	ErrQueryNotFound      error = errors.New("The request query is not found")
	ErrQueryInvalid       error = errors.New("The request query is invalid")
	ErrHeaderNotFound     error = errors.New("The request header is not found")
	ErrHeaderInvalid      error = errors.New("The request header is invalid")
	ErrWebRequestNil      error = errors.New("The web request object is nil")
	ErrResponseInvalid    error = errors.New("The response body is invalid")
)

// other errors
var (
	ErrInvalidOperation error = errors.New("The underlying operation is invalid")
	ErrForbidden        error = errors.New("The underlying operation is forbidden")
	ErrNotImplemented   error = errors.New("The underlying operation is not implemented")
	ErrBadRequest       error = errors.New("The client request is invalid")
	ErrResourceNotFound error = errors.New("The requested resource is not found")
	ErrResourceLocked   error = errors.New("The requested resource is locked")
	ErrResourceConflict error = errors.New("The requested resource has data conflicts")
)
