package webserver

import (
	"net/http"
	"net/textproto"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	GetRequestBody(dataTemplate any) error

	// GetRequestParameter loads HTTP request parameter associated to session for given name and unmarshals the content to given data template
	GetRequestParameter(name string, dataTemplate any) error

	// GetRequestQueries loads HTTP request query strings associated to session for given name and unmarshals the content to given data template (must be a slice)
	GetRequestQueries(name string, dataTemplate any) error

	// GetRequestQuery loads HTTP request single query string associated to session for given name and unmarshals the content to given data template
	GetRequestQuery(name string, index int, dataTemplate any) error

	// GetRequestHeader loads HTTP request header strings associated to session for given name and unmarshals the content to given data template (must be a slice)
	GetRequestHeaders(name string, dataTemplate any) error

	// GetRequestHeader loads HTTP request single header string associated to session for given name and unmarshals the content to given data template
	GetRequestHeader(name string, index int, dataTemplate any) error
}

// SessionHTTPResponse is a subset of SessionHTTP interface, containing only HTTP response related methods
type SessionHTTPResponse interface {
	// GetResponseWriter returns the HTTP response writer object from session object for given session ID
	GetResponseWriter() http.ResponseWriter
}

// SessionAttachment is a subset of Session interface, containing only attachment related methods
type SessionAttachment interface {
	// Attach attaches any value object into the given session associated to the session ID
	Attach(name string, value any) bool

	// Detach detaches any value object from the given session associated to the session ID
	Detach(name string) bool

	// GetRawAttachment retrieves any value object from the given session associated to the session ID and returns the raw interface (consumer needs to manually cast, but works for struct with private fields)
	GetRawAttachment(name string) (any, bool)

	// GetAttachment retrieves any value object from the given session associated to the session ID and unmarshals the content to given data template (only works for structs with exported fields)
	GetAttachment(name string, dataTemplate any) bool
}

// SessionLogging is a subset of Session interface, containing only logging related methods
type SessionLogging interface {
	// LogMethodEnter sends a logging entry of MethodEnter log type for the given session associated to the session ID
	LogMethodEnter()

	// LogMethodParameter sends a logging entry of MethodParameter log type for the given session associated to the session ID
	LogMethodParameter(parameters ...any)

	// LogMethodLogic sends a logging entry of MethodLogic log type for the given session associated to the session ID
	LogMethodLogic(logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...any)

	// LogMethodReturn sends a logging entry of MethodReturn log type for the given session associated to the session ID
	LogMethodReturn(returns ...any)

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
	attachment     map[string]any
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
		isInterfaceValueNil(session.responseWriter) {
		return defaultResponseWriter
	}
	return session.responseWriter
}

// GetRequestBodyFromSession is a sugar-function to retrieve request body as an object via generics
func GetRequestBodyFromSession[T any](session Session) (*T, error) {
	var result T
	var err = session.GetRequestBody(&result)
	return &result, err
}

// GetRequestBody loads HTTP request body associated to session and unmarshals the content JSON to given data template
func (session *session) GetRequestBody(dataTemplate any) error {
	if session == nil {
		return newAppError(
			errorCodeGeneralFailure,
			errorMessageSessionNil,
			[]error{},
		)
	}
	var httpRequest = session.GetRequest()
	var requestBody = getRequestBody(
		httpRequest,
	)
	if requestBody == "" {
		return newAppError(
			errorCodeBadRequest,
			errorMessageRequestBodyEmpty,
			[]error{},
		)
	}
	logEndpointRequest(
		session,
		"Body",
		"Content",
		requestBody,
	)
	var unmarshalError = tryUnmarshal(
		requestBody,
		dataTemplate,
	)
	if unmarshalError != nil {
		logEndpointRequest(
			session,
			"Body",
			"UnmarshalError",
			"%+v",
			unmarshalError,
		)
		return newAppError(
			errorCodeBadRequest,
			errorMessageRequestBodyInvalid,
			[]error{unmarshalError},
		)
	}
	return nil
}

// GetRequestParameterFromSession is a sugar-function to retrieve request parameter as an object via generics
func GetRequestParameterFromSession[T any](session Session, name string) (*T, error) {
	var result T
	var err = session.GetRequestParameter(name, &result)
	return &result, err
}

// GetRequestParameter loads HTTP request parameter associated to session for given name and unmarshals the content to given data template
func (session *session) GetRequestParameter(name string, dataTemplate any) error {
	if session == nil {
		return newAppError(
			errorCodeGeneralFailure,
			errorMessageSessionNil,
			[]error{},
		)
	}
	var httpRequest = session.GetRequest()
	var parameters = mux.Vars(
		httpRequest,
	)
	var value, found = parameters[name]
	if !found {
		return newAppError(
			errorCodeBadRequest,
			errorMessageParameterNotFound,
			[]error{},
		)
	}
	logEndpointRequest(
		session,
		"Parameter",
		name,
		value,
	)
	var unmarshalError = tryUnmarshal(
		value,
		dataTemplate,
	)
	if unmarshalError != nil {
		logEndpointRequest(
			session,
			"Parameter",
			"UnmarshalError",
			"%+v",
			unmarshalError,
		)
		return newAppError(
			errorCodeBadRequest,
			errorMessageParameterInvalid,
			[]error{unmarshalError},
		)
	}
	return nil
}

func getAllQueries(session *session, name string) []string {
	var results = make([]string, 0)
	var httpRequest = session.GetRequest()
	if httpRequest.URL == nil {
		return results
	}
	var queries, found = httpRequest.URL.Query()[name]
	if !found {
		return results
	}
	for _, query := range queries {
		var items = strings.Split(query, ",")
		results = append(results, items...)
	}
	return results
}

// GetRequestQueriesFromSession is a sugar-function to retrieve request queries as a slice via generics
func GetRequestQueriesFromSession[T any](session Session, name string) ([]T, error) {
	var results []T
	var err = session.GetRequestQueries(name, &results)
	return results, err
}

// GetRequestQueries loads HTTP request query strings associated to session for given name and unmarshals the content to given data template (must be a slice)
func (session *session) GetRequestQueries(name string, dataTemplate any) error {
	if session == nil {
		return newAppError(
			errorCodeGeneralFailure,
			errorMessageSessionNil,
			[]error{},
		)
	}
	var vTemplate = reflect.ValueOf(dataTemplate)
	var tTemplate = vTemplate.Type()
	if vTemplate.Type().Kind() != reflect.Pointer {
		return newAppError(
			errorCodeGeneralFailure,
			errorMessageDataTemplateInvalid,
			[]error{},
		)
	}
	vTemplate = reflect.Indirect(vTemplate)
	tTemplate = vTemplate.Type()
	if tTemplate.Kind() != reflect.Slice {
		return newAppError(
			errorCodeGeneralFailure,
			errorMessageDataTemplateInvalid,
			[]error{},
		)
	}
	var tItem = tTemplate.Elem()
	var queries = getAllQueries(
		session,
		name,
	)
	for _, query := range queries {
		var vItem = reflect.New(tItem)
		logEndpointRequest(
			session,
			"Query",
			name,
			query,
		)
		var unmarshalError = tryUnmarshal(
			query,
			vItem.Interface(),
		)
		if unmarshalError != nil {
			logEndpointRequest(
				session,
				"Query",
				"UnmarshalError",
				"%+v",
				unmarshalError,
			)
			return newAppError(
				errorCodeBadRequest,
				errorMessageQueryInvalid,
				[]error{unmarshalError},
			)
		}
		vTemplate.Set(reflect.Append(vTemplate, vItem.Elem()))
	}
	return nil
}

// GetRequestQueryFromSession is a sugar-function to retrieve request query as an object via generics
func GetRequestQueryFromSession[T any](session Session, name string, index int) (*T, error) {
	var result T
	var err = session.GetRequestQuery(name, index, &result)
	return &result, err
}

// GetRequestQuery loads HTTP request single query string associated to session for given name and unmarshals the content to given data template
func (session *session) GetRequestQuery(name string, index int, dataTemplate any) error {
	if session == nil {
		return newAppError(
			errorCodeGeneralFailure,
			errorMessageSessionNil,
			[]error{},
		)
	}
	var queries = getAllQueries(
		session,
		name,
	)
	if len(queries) <= index {
		return newAppError(
			errorCodeBadRequest,
			errorMessageQueryNotFound,
			[]error{},
		)
	}
	var value = queries[index]
	logEndpointRequest(
		session,
		"Query",
		name,
		value,
	)
	var unmarshalError = tryUnmarshal(
		value,
		dataTemplate,
	)
	if unmarshalError != nil {
		logEndpointRequest(
			session,
			"Query",
			"UnmarshalError",
			"%+v",
			unmarshalError,
		)
		return newAppError(
			errorCodeBadRequest,
			errorMessageQueryInvalid,
			[]error{unmarshalError},
		)
	}
	return nil
}

func getAllHeaders(session *session, name string) []string {
	var httpRequest = session.GetRequest()
	var canonicalName = textproto.CanonicalMIMEHeaderKey(name)
	var headers, found = httpRequest.Header[canonicalName]
	if !found {
		return nil
	}
	return headers
}

// GetRequestHeadersFromSession is a sugar-function to retrieve request headers as a slice via generics
func GetRequestHeadersFromSession[T any](session Session, name string) ([]T, error) {
	var results []T
	var err = session.GetRequestHeaders(name, &results)
	return results, err
}

// GetRequestHeader loads HTTP request header strings associated to session for given name and unmarshals the content to given data template (must be a slice)
func (session *session) GetRequestHeaders(name string, dataTemplate any) error {
	if session == nil {
		return newAppError(
			errorCodeGeneralFailure,
			errorMessageSessionNil,
			[]error{},
		)
	}
	var vTemplate = reflect.ValueOf(dataTemplate)
	var tTemplate = vTemplate.Type()
	if vTemplate.Type().Kind() != reflect.Pointer {
		return newAppError(
			errorCodeGeneralFailure,
			errorMessageDataTemplateInvalid,
			[]error{},
		)
	}
	vTemplate = reflect.Indirect(vTemplate)
	tTemplate = vTemplate.Type()
	if tTemplate.Kind() != reflect.Slice {
		return newAppError(
			errorCodeGeneralFailure,
			errorMessageDataTemplateInvalid,
			[]error{},
		)
	}
	var tItem = tTemplate.Elem()
	var headers = getAllHeaders(
		session,
		name,
	)
	for _, header := range headers {
		var vItem = reflect.New(tItem)
		logEndpointRequest(
			session,
			"Header",
			name,
			header,
		)
		var unmarshalError = tryUnmarshal(
			header,
			vItem.Interface(),
		)
		if unmarshalError != nil {
			logEndpointRequest(
				session,
				"Header",
				"UnmarshalError",
				"%+v",
				unmarshalError,
			)
			return newAppError(
				errorCodeBadRequest,
				errorMessageHeaderInvalid,
				[]error{unmarshalError},
			)
		}
		vTemplate.Set(reflect.Append(vTemplate, vItem.Elem()))
	}
	return nil
}

// GetRequestHeaderFromSession is a sugar-function to retrieve request header as an object via generics
func GetRequestHeaderFromSession[T any](session Session, name string, index int) (*T, error) {
	var result T
	var err = session.GetRequestHeader(name, index, &result)
	return &result, err
}

// GetRequestHeader loads HTTP request single header string associated to session for given name and unmarshals the content to given data template
func (session *session) GetRequestHeader(name string, index int, dataTemplate any) error {
	if session == nil {
		return newAppError(
			errorCodeGeneralFailure,
			errorMessageSessionNil,
			[]error{},
		)
	}
	var headers = getAllHeaders(
		session,
		name,
	)
	if len(headers) <= index {
		return newAppError(
			errorCodeBadRequest,
			errorMessageHeaderNotFound,
			[]error{},
		)
	}
	var value = headers[index]
	logEndpointRequest(
		session,
		"Header",
		name,
		value,
	)
	var unmarshalError = tryUnmarshal(
		value,
		dataTemplate,
	)
	if unmarshalError != nil {
		logEndpointRequest(
			session,
			"Header",
			"UnmarshalError",
			"%+v",
			unmarshalError,
		)
		return newAppError(
			errorCodeBadRequest,
			errorMessageHeaderInvalid,
			[]error{unmarshalError},
		)
	}
	return nil
}

// Attach attaches any value object into the given session associated to the session ID
func (session *session) Attach(name string, value any) bool {
	if session == nil {
		return false
	}
	if session.attachment == nil {
		session.attachment = map[string]any{}
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
func (session *session) GetRawAttachment(name string) (any, bool) {
	if session == nil {
		return nil, false
	}
	var attachment, found = session.attachment[name]
	if !found {
		return nil, false
	}
	return attachment, true
}

// GetAttachmentFromSession is a sugar-function to retrieve attachment as an object via generics
func GetAttachmentFromSession[T any](session Session, name string) (*T, bool) {
	var attachment, found = session.GetRawAttachment(name)
	if !found {
		return new(T), false
	}
	var result, ok = attachment.(T)
	if !ok {
		return new(T), false
	}
	return &result, true
}

// GetAttachment retrieves any value object from the given session associated to the session ID and unmarshals the content to given data template
func (session *session) GetAttachment(name string, dataTemplate any) bool {
	if session == nil {
		return false
	}
	var attachment, found = session.GetRawAttachment(name)
	if !found {
		return false
	}
	var vTemplate = reflect.ValueOf(dataTemplate)
	var tTemplate = vTemplate.Type()
	if tTemplate.Kind() != reflect.Pointer {
		return false
	}
	vTemplate = reflect.Indirect(vTemplate)
	vTemplate.Set(reflect.ValueOf(attachment))
	return true
}

func getMethodName() string {
	var pc, _, _, ok = runtime.Caller(3)
	if !ok {
		return "?"
	}
	var fn = runtime.FuncForPC(pc)
	return fn.Name()
}

// LogMethodEnter sends a logging entry of MethodEnter log type for the given session associated to the session ID
func (session *session) LogMethodEnter() {
	var methodName = getMethodName()
	logMethodEnter(
		session,
		methodName,
		"",
		"",
	)
}

// LogMethodParameter sends a logging entry of MethodParameter log type for the given session associated to the session ID
func (session *session) LogMethodParameter(parameters ...any) {
	var methodName = getMethodName()
	for index, parameter := range parameters {
		logMethodParameter(
			session,
			methodName,
			strconv.Itoa(index),
			"%v",
			parameter,
		)
	}
}

// LogMethodLogic sends a logging entry of MethodLogic log type for the given session associated to the session ID
func (session *session) LogMethodLogic(logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...any) {
	logMethodLogic(
		session,
		logLevel,
		category,
		subcategory,
		messageFormat,
		parameters...,
	)
}

// LogMethodReturn sends a logging entry of MethodReturn log type for the given session associated to the session ID
func (session *session) LogMethodReturn(returns ...any) {
	var methodName = getMethodName()
	for index, returnValue := range returns {
		logMethodReturn(
			session,
			methodName,
			strconv.Itoa(index),
			"%v",
			returnValue,
		)
	}
}

// LogMethodExit sends a logging entry of MethodExit log type for the given session associated to the session ID
func (session *session) LogMethodExit() {
	var methodName = getMethodName()
	logMethodExit(
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
