package webserver

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Customization holds all customization methods
type Customization interface {
	// BootstrapCustomization holds customization methods related to bootstrapping
	BootstrapCustomization
	// LoggingCustomization holds customization methods related to logging
	LoggingCustomization
	// HostingCustomization holds customization methods related to hosting
	HostingCustomization
	// HandlerCustomization holds customization methods related to handlers
	HandlerCustomization
	// WebRequestCustomization holds customization methods related to web requests
	WebRequestCustomization
}

// BootstrapCustomization holds customization methods related to bootstrapping
type BootstrapCustomization interface {
	// PreBootstrap is to customize the pre-processing logic before bootstrapping
	PreBootstrap() error

	// PostBootstrap is to customize the post-processing logic after bootstrapping
	PostBootstrap() error

	// AppClosing is to customize the application closing logic after server shutdown
	AppClosing() error
}

// LoggingCustomization holds customization methods related to logging
type LoggingCustomization interface {
	// Log is to customize the logging backend for the whole application
	Log(session Session, logType LogType, logLevel LogLevel, category, subcategory, description string)
}

// HostingCustomization holds customization methods related to hosting
type HostingCustomization interface {
	// ServerCert is to customize the server certificate for application; also determines the server hosting security option (HTTP v.s. HTTPS)
	ServerCert() *tls.Certificate

	// CaCertPool is to customize the CA cert pool for incoming client certificate validation; if not set or nil, no validation is conducted for incoming client certificates
	CaCertPool() *x509.CertPool

	// GraceShutdownWaitTime is to customize the graceful shutdown wait time for the application
	GraceShutdownWaitTime() time.Duration

	// Routes is to customize the routes registration
	Routes() []Route

	// Statics is to customize the static contents registration
	Statics() []Static

	// Middlewares is to customize the middlewares registration
	Middlewares() []MiddlewareFunc

	// InstrumentRouter is to customize the instrumentation on top of a fully configured router; usually useful for 3rd party monitoring tools such as new relic, etc.
	InstrumentRouter(router *mux.Router) *mux.Router

	// WrapHandler is to customize the overall wrapping of http.Handler before the server is configured
	WrapHandler(handler http.Handler) http.Handler
}

// HandlerCustomization holds customization methods related to handlers
type HandlerCustomization interface {
	// PreAction is to customize the pre-action used before each route action takes place, e.g. authorization, etc.
	PreAction(session Session) error

	// PostAction is to customize the post-action used after each route action takes place, e.g. finalization, etc.
	PostAction(session Session) error

	// InterpretSuccess is to customize how application interpret a response content into HTTP status code and corresponding response body
	InterpretSuccess(responseContent interface{}) (int, string)

	// InterpretError is to customize how application interpret an error into HTTP status code and corresponding response body
	InterpretError(err error) (int, string)

	// RecoverPanic is to customize the recovery of panic into a valid response and error in case it happens (for recoverable panic only)
	RecoverPanic(session Session, recoverResult interface{}) (interface{}, error)

	// NotFoundHandler is to customize the handler to be used when no route matches.
	NotFoundHandler() http.Handler

	// MethodNotAllowed is to customize the handler to be used when the request method does not match the route
	MethodNotAllowedHandler() http.Handler
}

// WebRequestCustomization holds customization methods related to web requests
type WebRequestCustomization interface {
	// ClientCert is to customize the client certificate for external requests; if not set or nil, no client certificate is sent to external web services
	ClientCert() *tls.Certificate

	// DefaultTimeout is to customize the default timeout for any webcall communications through HTTP/HTTPS by session
	DefaultTimeout() time.Duration

	// SkipServerCertVerification is to customize the skip of server certificate verification for any webcall communications through HTTP/HTTPS by session
	SkipServerCertVerification() bool

	// RoundTripper is to customize the creation of the HTTP transport for any webcall communications through HTTP/HTTPS by session
	RoundTripper(originalTransport http.RoundTripper) http.RoundTripper

	// WrapRequest is to customize the creation of the HTTP request for any webcall communications through HTTP/HTTPS by session; utilize this method if needed for new relic wrapping, etc.
	WrapRequest(session Session, httpRequest *http.Request) *http.Request
}

var (
	customizationDefault = &DefaultCustomization{}
)

// DefaultCustomization can be used for easier customization override
type DefaultCustomization struct{}

// PreBootstrap is to customize the pre-processing logic before bootstrapping
func (customization *DefaultCustomization) PreBootstrap() error {
	return nil
}

// PostBootstrap is to customize the post-processing logic after bootstrapping
func (customization *DefaultCustomization) PostBootstrap() error {
	return nil
}

// AppClosing is to customize the application closing logic after server shutdown
func (customization *DefaultCustomization) AppClosing() error {
	return nil
}

// Log is to customize the logging backend for the whole application
func (customization *DefaultCustomization) Log(session Session, logType LogType, logLevel LogLevel, category, subcategory, description string) {
	if isInterfaceValueNilFunc(session) {
		return
	}
	fmtPrintf(
		"<%v|%v> (%v|%v) [%v|%v] %v\n",
		session.GetID(),
		session.GetName(),
		logType,
		logLevel,
		category,
		subcategory,
		description,
	)
}

// ServerCert is to customize the server certificate for application; also determines the server hosting security option (HTTP v.s. HTTPS)
func (customization *DefaultCustomization) ServerCert() *tls.Certificate {
	return nil
}

// CaCertPool is to customize the CA cert pool for incoming client certificate validation; if not set or nil, no validation is conducted for incoming client certificates
func (customization *DefaultCustomization) CaCertPool() *x509.CertPool {
	return nil
}

// GraceShutdownWaitTime is to customize the graceful shutdown wait time for the application
func (customization *DefaultCustomization) GraceShutdownWaitTime() time.Duration {
	return 3 * time.Minute
}

// Routes is to customize the routes registration
func (customization *DefaultCustomization) Routes() []Route {
	return []Route{}
}

// Statics is to customize the static contents registration
func (customization *DefaultCustomization) Statics() []Static {
	return []Static{}
}

// Middlewares is to customize the middlewares registration
func (customization *DefaultCustomization) Middlewares() []MiddlewareFunc {
	return []MiddlewareFunc{}
}

// InstrumentRouter is to customize the instrumentation on top of a fully configured router; usually useful for 3rd party monitoring tools such as new relic, etc.
func (customization *DefaultCustomization) InstrumentRouter(router *mux.Router) *mux.Router {
	return router
}

// WrapHandler is to customize the overall wrapping of http.Handler before the server is configured
func (customization *DefaultCustomization) WrapHandler(handler http.Handler) http.Handler {
	return handler
}

// PreAction is to customize the pre-action used before each route action takes place, e.g. authorization, etc.
func (customization *DefaultCustomization) PreAction(session Session) error {
	return nil
}

// PostAction is to customize the post-action used after each route action takes place, e.g. finalization, etc.
func (customization *DefaultCustomization) PostAction(session Session) error {
	return nil
}

// InterpretSuccess is to customize how application interpret a response content into HTTP status code and corresponding response body
func (customization *DefaultCustomization) InterpretSuccess(responseContent interface{}) (int, string) {
	if isInterfaceValueNilFunc(responseContent) {
		return http.StatusNoContent, ""
	}
	var responseMessage = marshalIgnoreErrorFunc(responseContent)
	if responseMessage == "" {
		return http.StatusNoContent, ""
	}
	return http.StatusOK, responseMessage
}

// InterpretError is to customize how application interpret an error into HTTP status code and corresponding status message
func (customization *DefaultCustomization) InterpretError(err error) (int, string) {
	var typedError, isTyped = err.(AppHTTPError)
	if !isTyped {
		return http.StatusInternalServerError,
			err.Error()
	}
	return typedError.HTTPStatusCode(),
		typedError.HTTPResponseMessage()
}

func getRecoverError(recoverResult interface{}) error {
	var err, ok = recoverResult.(error)
	if !ok {
		err = fmtErrorf("Endpoint panic: %v", recoverResult)
	}
	return err
}

// RecoverPanic is to customize the recovery of panic into a valid response and error in case it happens (for recoverable panic only)
func (customization *DefaultCustomization) RecoverPanic(session Session, recoverResult interface{}) (interface{}, error) {
	var err = getRecoverErrorFunc(
		recoverResult,
	)
	session.LogMethodLogic(
		LogLevelError,
		"RecoverPanic",
		session.GetName(),
		"Error: %+v\nCallstack: %v",
		err,
		string(
			debugStack(),
		),
	)
	return nil, err
}

// NotFoundHandler is to customize the handler to be used when no route matches.
func (customization *DefaultCustomization) NotFoundHandler() http.Handler {
	return nil
}

// MethodNotAllowedHandler is to customize the handler to be used when the request method does not match the route
func (customization *DefaultCustomization) MethodNotAllowedHandler() http.Handler {
	return nil
}

// ClientCert is to customize the client certificate for external requests; if not set or nil, no client certificate is sent to external web services
func (customization *DefaultCustomization) ClientCert() *tls.Certificate {
	return nil
}

// DefaultTimeout is to customize the default timeout for any webcall communications through HTTP/HTTPS by session
func (customization *DefaultCustomization) DefaultTimeout() time.Duration {
	return 3 * time.Minute
}

// SkipServerCertVerification is to customize the skip of server certificate verification for any webcall communications through HTTP/HTTPS by session
func (customization *DefaultCustomization) SkipServerCertVerification() bool {
	return false
}

// RoundTripper is to customize the creation of the HTTP transport for any webcall communications through HTTP/HTTPS by session
func (customization *DefaultCustomization) RoundTripper(originalTransport http.RoundTripper) http.RoundTripper {
	return originalTransport
}

// WrapRequest is to customize the creation of the HTTP request for any webcall communications through HTTP/HTTPS by session; utilize this method if needed for new relic wrapping, etc.
func (customization *DefaultCustomization) WrapRequest(session Session, httpRequest *http.Request) *http.Request {
	return httpRequest
}
