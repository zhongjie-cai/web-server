package webserver

import "errors"

var (
	errSessionNil         error = errors.New("The session object is nil")
	errRouteRegistion     error = errors.New("The route registration failed")
	errRouteNotFound      error = errors.New("The route is not found")
	errHostServer         error = errors.New("The server hosting failed")
	errRequestBodyEmpty   error = errors.New("The request body is empty")
	errRequestBodyInvalid error = errors.New("The request body is invalid")
	errParameterNotFound  error = errors.New("The request parameter is not found")
	errParameterInvalid   error = errors.New("The request parameter is invalid")
	errQueryNotFound      error = errors.New("The request query is not found")
	errQueryInvalid       error = errors.New("The request query is invalid")
	errHeaderNotFound     error = errors.New("The request header is not found")
	errHeaderInvalid      error = errors.New("The request header is invalid")
	errWebRequestNil      error = errors.New("The web request object is nil")
	errResponseInvalid    error = errors.New("The response body is invalid")
)
