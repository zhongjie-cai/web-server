package webserver

import (
	"net/http"

	"github.com/google/uuid"
)

var (
	defaultRequest        *http.Request       = &http.Request{}
	defaultResponseWriter http.ResponseWriter = &nilResponseWriter{}
)

// Session is the storage for the current HTTP request session, containing information needed for logging, monitoring, etc.
type Session interface {
	SessionMeta
	SessionHTTP
	SessionAttachment
	SessionLogging
	SessionWebcall
}

// SessionMeta is a subset of Session interface, containing only meta data related methods
type SessionMeta interface {
	// GetID returns the ID of this registered session object
	GetID() uuid.UUID

	// GetName returns the name registered to session object for given session ID
	GetName() string
}

// SessionHTTP is a subset of Session interface, containing only HTTP request & response related methods
type SessionHTTP interface {
	SessionHTTPRequest
	SessionHTTPResponse
}

// SessionHTTPRequest is a subset of Session interface, containing only HTTP request related methods
type SessionHTTPRequest interface {
	// GetRequest returns the HTTP request object from session object for given session ID
	GetRequest() *http.Request

	// GetRequestBody loads HTTP request body associated to session and unmarshals the content JSON to given data template
	GetRequestBody(dataTemplate interface{}) error

	// GetRequestParameter loads HTTP request parameter associated to session for given name and unmarshals the content to given data template
	GetRequestParameter(name string, dataTemplate interface{}) error

	// GetRequestQuery loads HTTP request single query string associated to session for given name and unmarshals the content to given data template
	GetRequestQuery(name string, index int, dataTemplate interface{}) error

	// GetRequestHeader loads HTTP request single header string associated to session for given name and unmarshals the content to given data template
	GetRequestHeader(name string, index int, dataTemplate interface{}) error
}

// SessionHTTPResponse is a subset of SessionHTTP interface, containing only HTTP response related methods
type SessionHTTPResponse interface {
	// GetResponseWriter returns the HTTP response writer object from session object for given session ID
	GetResponseWriter() http.ResponseWriter
}

// SessionAttachment is a subset of Session interface, containing only attachment related methods
type SessionAttachment interface {
	// Attach attaches any value object into the given session associated to the session ID
	Attach(name string, value interface{}) bool

	// Detach detaches any value object from the given session associated to the session ID
	Detach(name string) bool

	// GetRawAttachment retrieves any value object from the given session associated to the session ID and returns the raw interface (consumer needs to manually cast, but works for struct with private fields)
	GetRawAttachment(name string) (interface{}, bool)

	// GetAttachment retrieves any value object from the given session associated to the session ID and unmarshals the content to given data template (only works for structs with exported fields)
	GetAttachment(name string, dataTemplate interface{}) bool
}

// SessionLogging is a subset of Session interface, containing only logging related methods
type SessionLogging interface {
	// LogMethodEnter sends a logging entry of MethodEnter log type for the given session associated to the session ID
	LogMethodEnter()

	// LogMethodParameter sends a logging entry of MethodParameter log type for the given session associated to the session ID
	LogMethodParameter(parameters ...interface{})

	// LogMethodLogic sends a logging entry of MethodLogic log type for the given session associated to the session ID
	LogMethodLogic(logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{})

	// LogMethodReturn sends a logging entry of MethodReturn log type for the given session associated to the session ID
	LogMethodReturn(returns ...interface{})

	// LogMethodExit sends a logging entry of MethodExit log type for the given session associated to the session ID
	LogMethodExit()
}

// SessionWebcall is a subset of Session interface, containing only webcall related methods
type SessionWebcall interface {
	// CreateWebcallRequest generates a webcall request object to the targeted external web service for the given session associated to the session ID
	CreateWebcallRequest(method string, url string, payload string, sendClientCert bool) WebRequest
}

type session struct {
	id             uuid.UUID
	name           string
	request        *http.Request
	responseWriter http.ResponseWriter
	attachment     map[string]interface{}
	customization  Customization
}

// GetID returns the ID of this registered session object
func (session *session) GetID() uuid.UUID {
	if session == nil {
		return uuid.Nil
	}
	return session.id
}

// GetName returns the name registered to session object for given session ID
func (session *session) GetName() string {
	if session == nil {
		return ""
	}
	return session.name
}

// GetRequest returns the HTTP request object from session object for given session ID
func (session *session) GetRequest() *http.Request {
	if session == nil ||
		session.request == nil {
		return defaultRequest
	}
	return session.request
}

// GetResponseWriter returns the HTTP response writer object from session object for given session ID
func (session *session) GetResponseWriter() http.ResponseWriter {
	if session == nil ||
		isInterfaceValueNilFunc(session.responseWriter) {
		return defaultResponseWriter
	}
	return session.responseWriter
}

// GetRequestBody loads HTTP request body associated to session and unmarshals the content JSON to given data template
func (session *session) GetRequestBody(dataTemplate interface{}) error {
	if session == nil {
		return newAppErrorFunc(
			errorCodeGeneralFailure,
			errorMessageSessionNil,
			[]error{},
		)
	}
	var httpRequest = session.GetRequest()
	var requestBody = getRequestBodyFunc(
		httpRequest,
	)
	if requestBody == "" {
		return newAppErrorFunc(
			errorCodeBadRequest,
			errorMessageRequestBodyEmpty,
			[]error{},
		)
	}
	logEndpointRequestFunc(
		session,
		"Body",
		"Content",
		requestBody,
	)
	var unmarshalError = tryUnmarshalFunc(
		requestBody,
		dataTemplate,
	)
	if unmarshalError != nil {
		logEndpointRequestFunc(
			session,
			"Body",
			"UnmarshalError",
			"%+v",
			unmarshalError,
		)
		return newAppErrorFunc(
			errorCodeBadRequest,
			errorMessageRequestBodyInvalid,
			[]error{unmarshalError},
		)
	}
	return nil
}

// GetRequestParameter loads HTTP request parameter associated to session for given name and unmarshals the content to given data template
func (session *session) GetRequestParameter(name string, dataTemplate interface{}) error {
	if session == nil {
		return newAppErrorFunc(
			errorCodeGeneralFailure,
			errorMessageSessionNil,
			[]error{},
		)
	}
	var httpRequest = session.GetRequest()
	var parameters = muxVars(
		httpRequest,
	)
	var value, found = parameters[name]
	if !found {
		return newAppErrorFunc(
			errorCodeBadRequest,
			errorMessageParameterNotFound,
			[]error{},
		)
	}
	logEndpointRequestFunc(
		session,
		"Parameter",
		name,
		value,
	)
	var unmarshalError = tryUnmarshalFunc(
		value,
		dataTemplate,
	)
	if unmarshalError != nil {
		logEndpointRequestFunc(
			session,
			"Parameter",
			"UnmarshalError",
			"%+v",
			unmarshalError,
		)
		return newAppErrorFunc(
			errorCodeBadRequest,
			errorMessageParameterInvalid,
			[]error{unmarshalError},
		)
	}
	return nil
}

func getAllQueries(session *session, name string) []string {
	var httpRequest = session.GetRequest()
	if httpRequest.URL == nil {
		return nil
	}
	var queries, found = httpRequest.URL.Query()[name]
	if !found {
		return nil
	}
	return queries
}

// GetRequestQuery loads HTTP request single query string associated to session for given name and unmarshals the content to given data template
func (session *session) GetRequestQuery(name string, index int, dataTemplate interface{}) error {
	if session == nil {
		return newAppErrorFunc(
			errorCodeGeneralFailure,
			errorMessageSessionNil,
			[]error{},
		)
	}
	var queries = getAllQueriesFunc(
		session,
		name,
	)
	if len(queries) <= index {
		return newAppErrorFunc(
			errorCodeBadRequest,
			errorMessageQueryNotFound,
			[]error{},
		)
	}
	var value = queries[index]
	logEndpointRequestFunc(
		session,
		"Query",
		name,
		value,
	)
	var unmarshalError = tryUnmarshalFunc(
		value,
		dataTemplate,
	)
	if unmarshalError != nil {
		logEndpointRequestFunc(
			session,
			"Query",
			"UnmarshalError",
			"%+v",
			unmarshalError,
		)
		return newAppErrorFunc(
			errorCodeBadRequest,
			errorMessageQueryInvalid,
			[]error{unmarshalError},
		)
	}
	return nil
}

func getAllHeaders(session *session, name string) []string {
	var httpRequest = session.GetRequest()
	var canonicalName = textprotoCanonicalMIMEHeaderKey(name)
	var headers, found = httpRequest.Header[canonicalName]
	if !found {
		return nil
	}
	return headers
}

// GetRequestHeader loads HTTP request single header string associated to session for given name and unmarshals the content to given data template
func (session *session) GetRequestHeader(name string, index int, dataTemplate interface{}) error {
	if session == nil {
		return newAppErrorFunc(
			errorCodeGeneralFailure,
			errorMessageSessionNil,
			[]error{},
		)
	}
	var headers = getAllHeadersFunc(
		session,
		name,
	)
	if len(headers) <= index {
		return newAppErrorFunc(
			errorCodeBadRequest,
			errorMessageHeaderNotFound,
			[]error{},
		)
	}
	var value = headers[index]
	logEndpointRequestFunc(
		session,
		"Header",
		name,
		value,
	)
	var unmarshalError = tryUnmarshalFunc(
		value,
		dataTemplate,
	)
	if unmarshalError != nil {
		logEndpointRequestFunc(
			session,
			"Header",
			"UnmarshalError",
			"%+v",
			unmarshalError,
		)
		return newAppErrorFunc(
			errorCodeBadRequest,
			errorMessageHeaderInvalid,
			[]error{unmarshalError},
		)
	}
	return nil
}

// Attach attaches any value object into the given session associated to the session ID
func (session *session) Attach(name string, value interface{}) bool {
	if session == nil {
		return false
	}
	if session.attachment == nil {
		session.attachment = map[string]interface{}{}
	}
	session.attachment[name] = value
	return true
}

// Detach detaches any value object from the given session associated to the session ID
func (session *session) Detach(name string) bool {
	if session == nil {
		return false
	}
	if session.attachment != nil {
		delete(session.attachment, name)
	}
	return true
}

// GetRawAttachment retrieves any value object from the given session associated to the session ID and returns the raw interface (consumer needs to manually cast, but works for struct with private fields)
func (session *session) GetRawAttachment(name string) (interface{}, bool) {
	if session == nil {
		return nil, false
	}
	var attachment, found = session.attachment[name]
	if !found {
		return nil, false
	}
	return attachment, true
}

// GetAttachment retrieves any value object from the given session associated to the session ID and unmarshals the content to given data template
func (session *session) GetAttachment(name string, dataTemplate interface{}) bool {
	if session == nil {
		return false
	}
	var attachment, found = session.GetRawAttachment(name)
	if !found {
		return false
	}
	var bytes, marshalError = jsonMarshal(attachment)
	if marshalError != nil {
		return false
	}
	var unmarshalError = jsonUnmarshal(
		bytes,
		dataTemplate,
	)
	return unmarshalError == nil
}

func getMethodName() string {
	var pc, _, _, ok = runtimeCaller(3)
	if !ok {
		return "?"
	}
	var fn = runtimeFuncForPC(pc)
	return fn.Name()
}

// LogMethodEnter sends a logging entry of MethodEnter log type for the given session associated to the session ID
func (session *session) LogMethodEnter() {
	var methodName = getMethodNameFunc()
	logMethodEnterFunc(
		session,
		methodName,
		"",
		"",
	)
}

// LogMethodParameter sends a logging entry of MethodParameter log type for the given session associated to the session ID
func (session *session) LogMethodParameter(parameters ...interface{}) {
	var methodName = getMethodNameFunc()
	for index, parameter := range parameters {
		logMethodParameterFunc(
			session,
			methodName,
			strconvItoa(index),
			"%v",
			parameter,
		)
	}
}

// LogMethodLogic sends a logging entry of MethodLogic log type for the given session associated to the session ID
func (session *session) LogMethodLogic(logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	logMethodLogicFunc(
		session,
		logLevel,
		category,
		subcategory,
		messageFormat,
		parameters...,
	)
}

// LogMethodReturn sends a logging entry of MethodReturn log type for the given session associated to the session ID
func (session *session) LogMethodReturn(returns ...interface{}) {
	var methodName = getMethodNameFunc()
	for index, returnValue := range returns {
		logMethodReturnFunc(
			session,
			methodName,
			strconvItoa(index),
			"%v",
			returnValue,
		)
	}
}

// LogMethodExit sends a logging entry of MethodExit log type for the given session associated to the session ID
func (session *session) LogMethodExit() {
	var methodName = getMethodNameFunc()
	logMethodExitFunc(
		session,
		methodName,
		"",
		"",
	)
}

// CreateWebcallRequest generates a webcall request object to the targeted external web service for the given session associated to the session ID
func (session *session) CreateWebcallRequest(
	method string,
	url string,
	payload string,
	sendClientCert bool,
) WebRequest {
	return &webRequest{
		session,
		method,
		url,
		payload,
		map[string][]string{},
		map[string][]string{},
		0,
		nil,
		sendClientCert,
		0,
		[]dataReceiver{},
	}
}
