package webserver

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Customization holds all customization methods
type Customization interface {
	BootstrapCustomization
	LoggingCustomization
	HostingCustomization
	HandlerCustomization
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
}

// HandlerCustomization holds customization methods related to handlers
type HandlerCustomization interface {
	// PreAction is to customize the pre-actiontion used before each route action takes place, e.g. authorization, etc.
	PreAction(session Session) error

	// PostAction is to customize the post-actiontion used after each route action takes place, e.g. finalization, etc.
	PostAction(session Session) error

	// InterpretError is to customize how application interpret an error into HTTP status code and corresponding status message
	InterpretError(err error) (int, string)

	// NotFoundHandler is to customize the handler to be used when no route matches.
	NotFoundHandler() http.Handler

	// MethodNotAllowed is to customize the handler to be used when the request method does not match the route
	MethodNotAllowedHandler() http.Handler
}

// WebRequestCustomization holds customization methods related to web requests
type WebRequestCustomization interface {
	// ClientCert is to customize the client certificate for external requests; if not set or nil, no client certificate is sent to external web services
	ClientCert() *tls.Certificate

	// DefaultTimeout is to customize the default timeout for any network communications through HTTP/HTTPS by session
	DefaultTimeout() time.Duration

	// SkipServerCertVerification is to customize the skip of server certificate verification for any network communications through HTTP/HTTPS by session
	SkipServerCertVerification() bool

	// RoundTripper is to customize the creation of the HTTP transport for any network communications through HTTP/HTTPS by session
	RoundTripper(originalTransport http.RoundTripper) http.RoundTripper

	// WrapRequest is to customize the creation of the HTTP request for any network communications through HTTP/HTTPS by session; utilize this method if needed for new relic wrapping, etc.
	WrapRequest(session Session, httpRequest *http.Request) *http.Request
}

var (
	customizationDefault = &defaultCustomization{}
)

type defaultCustomization struct{}

// PreBootstrap is to customize the pre-processing logic before bootstrapping
func (customization *defaultCustomization) PreBootstrap() error {
	return nil
}

// PostBootstrap is to customize the post-processing logic after bootstrapping
func (customization *defaultCustomization) PostBootstrap() error {
	return nil
}

// AppClosing is to customize the application closing logic after server shutdown
func (customization *defaultCustomization) AppClosing() error {
	return nil
}

// Log is to customize the logging backend for the whole application
func (customization *defaultCustomization) Log(session Session, logType LogType, logLevel LogLevel, category, subcategory, description string) {
	if isInterfaceValueNil(session) {
		return
	}
	fmt.Println(
		fmt.Sprintf(
			"<%v|%v> (%v|%v) [%v|%v] %v",
			session.GetID(),
			session.GetName(),
			logType,
			logLevel,
			category,
			subcategory,
			description,
		),
	)
}

// ServerCert is to customize the server certificate for application; also determines the server hosting security option (HTTP v.s. HTTPS)
func (customization *defaultCustomization) ServerCert() *tls.Certificate {
	return nil
}

// CaCertPool is to customize the CA cert pool for incoming client certificate validation; if not set or nil, no validation is conducted for incoming client certificates
func (customization *defaultCustomization) CaCertPool() *x509.CertPool {
	return nil
}

// GraceShutdownWaitTime is to customize the graceful shutdown wait time for the application
func (customization *defaultCustomization) GraceShutdownWaitTime() time.Duration {
	return 3 * time.Minute
}

// Routes is to customize the routes registration
func (customization *defaultCustomization) Routes() []Route {
	return []Route{}
}

// Statics is to customize the static contents registration
func (customization *defaultCustomization) Statics() []Static {
	return []Static{}
}

// Middlewares is to customize the middlewares registration
func (customization *defaultCustomization) Middlewares() []MiddlewareFunc {
	return []MiddlewareFunc{}
}

// InstrumentRouter is to customize the instrumentation on top of a fully configured router; usually useful for 3rd party monitoring tools such as new relic, etc.
func (customization *defaultCustomization) InstrumentRouter(router *mux.Router) *mux.Router {
	return router
}

// PreAction is to customize the pre-actiontion used before each route action takes place, e.g. authorization, etc.
func (customization *defaultCustomization) PreAction(session Session) error {
	return nil
}

// PostAction is to customize the post-actiontion used after each route action takes place, e.g. finalization, etc.
func (customization *defaultCustomization) PostAction(session Session) error {
	return nil
}

// InterpretError is to customize how application interpret an error into HTTP status code and corresponding status message
func (customization *defaultCustomization) InterpretError(err error) (int, string) {
	var statusCode int
	switch err {
	case errHeaderInvalid:
		statusCode = http.StatusBadRequest
	default:
		statusCode = http.StatusInternalServerError
	}
	var statusMessage = fmt.Sprintf(
		"%+v",
		err,
	)
	return statusCode, statusMessage
}

// NotFoundHandler is to customize the handler to be used when no route matches.
func (customization *defaultCustomization) NotFoundHandler() http.Handler {
	return nil
}

// MethodNotAllowed is to customize the handler to be used when the request method does not match the route
func (customization *defaultCustomization) MethodNotAllowedHandler() http.Handler {
	return nil
}

// ClientCert is to customize the client certificate for external requests; if not set or nil, no client certificate is sent to external web services
func (customization *defaultCustomization) ClientCert() *tls.Certificate {
	return nil
}

// DefaultTimeout is to customize the default timeout for any network communications through HTTP/HTTPS by session
func (customization *defaultCustomization) DefaultTimeout() time.Duration {
	return 3 * time.Minute
}

// SkipServerCertVerification is to customize the skip of server certificate verification for any network communications through HTTP/HTTPS by session
func (customization *defaultCustomization) SkipServerCertVerification() bool {
	return false
}

// RoundTripper is to customize the creation of the HTTP transport for any network communications through HTTP/HTTPS by session
func (customization *defaultCustomization) RoundTripper(originalTransport http.RoundTripper) http.RoundTripper {
	return originalTransport
}

// WrapRequest is to customize the creation of the HTTP request for any network communications through HTTP/HTTPS by session; utilize this method if needed for new relic wrapping, etc.
func (customization *defaultCustomization) WrapRequest(session Session, httpRequest *http.Request) *http.Request {
	return httpRequest
}
