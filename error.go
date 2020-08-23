package webserver

import "errors"

// internal errors
var (
	errSessionNil     error = errors.New("The session object is nil")
	errRouteRegistion error = errors.New("The route registration failed")
	errRouteNotFound  error = errors.New("The route is not found")
	errHostServer     error = errors.New("The server hosting failed")
)

// external errors
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
