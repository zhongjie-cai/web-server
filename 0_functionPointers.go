package webserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"net/url"
	"os/signal"
	"reflect"
	"regexp"
	"runtime"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// func pointers for injection / testing: appError.go
var (
	fmtSprint              = fmt.Sprint
	getErrorMessageFunc    = getErrorMessage
	printInnerErrorsFunc   = printInnerErrors
	errorsIs               = errors.Is
	equalsErrorFunc        = equalsError
	appErrorContainsFunc   = appErrorContains
	innerErrorContainsFunc = innerErrorContains
	cleanupInnerErrorsFunc = cleanupInnerErrors
	newAppErrorFunc        = newAppError
)

// func pointers for injection / testing: application.go
var (
	isInterfaceValueNilFunc   = isInterfaceValueNil
	uuidNew                   = uuid.New
	startApplicationFunc      = startApplication
	haltServerFunc            = haltServer
	preBootstrapingFunc       = preBootstraping
	bootstrapFunc             = bootstrap
	postBootstrapingFunc      = postBootstraping
	endApplicationFunc        = endApplication
	beginApplicationFunc      = beginApplication
	logAppRootFunc            = logAppRoot
	initializeHTTPClientsFunc = initializeHTTPClients
	hostServerFunc            = hostServer
)

// func pointers for injection / testing: customization.go
var (
	fmtPrintf              = fmt.Printf
	fmtSprintf             = fmt.Sprintf
	marshalIgnoreErrorFunc = marshalIgnoreError
	debugStack             = debug.Stack
	getRecoverErrorFunc    = getRecoverError
)

// func pointers for injection / testing: handler.go
var (
	stringsSplit         = strings.Split
	strconvAtoi          = strconv.Atoi
	getRequestedPortFunc = getRequestedPort
	getApplicationFunc   = getApplication
	getRouteInfoFunc     = getRouteInfo
	initiateSessionFunc  = initiateSession
	getTimeNowUTCFunc    = getTimeNowUTC
	finalizeSessionFunc  = finalizeSession
	logEndpointEnterFunc = logEndpointEnter
	logEndpointExitFunc  = logEndpointExit
	timeSince            = time.Since
	handlePanicFunc      = handlePanic
	writeResponseFunc    = writeResponse
	handleActionFunc     = handleAction
)

// func pointers for injection / testing: jsonutil.go
var (
	jsonNewEncoder                 = json.NewEncoder
	stringsTrimRight               = strings.TrimRight
	reflectTypeOf                  = reflect.TypeOf
	strconvParseBool               = strconv.ParseBool
	stringsToLower                 = strings.ToLower
	strconvParseInt                = strconv.ParseInt
	strconvParseFloat              = strconv.ParseFloat
	strconvParseUint               = strconv.ParseUint
	tryUnmarshalPrimitiveTypesFunc = tryUnmarshalPrimitiveTypes
	jsonUnmarshal                  = json.Unmarshal
	fmtErrorf                      = fmt.Errorf
)

// func pointers for injection / testing: logger.go
var (
	prepareLoggingFunc = prepareLogging
)

// func pointers for injection / testing: logType.go
var (
	sortStrings = sort.Strings
	stringsJoin = strings.Join
)

// func pointers for injection / testing: parameter.go
var (
	regexpMatchString = regexp.MatchString
)

// func pointers for injection / testing: pointerutil.go
var (
	reflectValueOf = reflect.ValueOf
)

// func pointers for injection / testing: register.go
var (
	stringsReplace                 = strings.Replace
	doParameterReplacementFunc     = doParameterReplacement
	evaluatePathWithParametersFunc = evaluatePathWithParameters
	evaluateQueriesFunc            = evaluateQueries
	registerRouteFunc              = registerRoute
	registerStaticFunc             = registerStatic
	addMiddlewareFunc              = addMiddleware
	muxNewRouter                   = mux.NewRouter
	registerRoutesFunc             = registerRoutes
	registerStaticsFunc            = registerStatics
	registerMiddlewaresFunc        = registerMiddlewares
	walkRegisteredRoutesFunc       = walkRegisteredRoutes
	registerErrorHandlersFunc      = registerErrorHandlers
)

// func pointers for injection / testing: request.go
var (
	ioutilReadAll   = ioutil.ReadAll
	ioutilNopCloser = ioutil.NopCloser
	bytesNewBuffer  = bytes.NewBuffer
)

// func pointers for injection / testing: response.go
var (
	shouldSkipHandlingFunc  = shouldSkipHandling
	constructResponseFunc   = constructResponse
	logEndpointResponseFunc = logEndpointResponse
	httpStatusText          = http.StatusText
	strconvItoa             = strconv.Itoa
)

// func pointers for injection / testing: route.go
var (
	getPathTemplateFunc   = getPathTemplate
	getPathRegexpFunc     = getPathRegexp
	evaluateRouteFunc     = evaluateRoute
	muxCurrentRoute       = mux.CurrentRoute
	getNameFunc           = getName
	getEndpointByNameFunc = getEndpointByName
)

// func pointers for injection / testing: server.go
var (
	instantiateRouterFunc    = instantiateRouter
	runServerFunc            = runServer
	createServerFunc         = createServer
	signalNotify             = signal.Notify
	listenAndServeFunc       = listenAndServe
	contextWithTimeout       = context.WithTimeout
	contextBackground        = context.Background
	shutdownServerFunc       = shutdownServer
	evaluateServerErrorsFunc = evaluateServerErrors
)

// func pointers for injection / testing: session.go
var (
	getRequestBodyFunc              = getRequestBody
	logEndpointRequestFunc          = logEndpointRequest
	tryUnmarshalFunc                = tryUnmarshal
	muxVars                         = mux.Vars
	getAllQueriesFunc               = getAllQueries
	textprotoCanonicalMIMEHeaderKey = textproto.CanonicalMIMEHeaderKey
	getAllHeadersFunc               = getAllHeaders
	jsonMarshal                     = json.Marshal
	runtimeCaller                   = runtime.Caller
	runtimeFuncForPC                = runtime.FuncForPC
	getMethodNameFunc               = getMethodName
	logMethodEnterFunc              = logMethodEnter
	logMethodParameterFunc          = logMethodParameter
	logMethodLogicFunc              = logMethodLogic
	logMethodReturnFunc             = logMethodReturn
	logMethodExitFunc               = logMethodExit
)

// func pointers for injection / testing: timeutil.go
var (
	timeNow = time.Now
)

// func pointers for injection / testing: webRequest.go
var (
	clientDoFunc            = clientDo
	timeSleep               = time.Sleep
	getHTTPTransportFunc    = getHTTPTransport
	urlQueryEscape          = url.QueryEscape
	createQueryStringFunc   = createQueryString
	generateRequestURLFunc  = generateRequestURL
	stringsNewReader        = strings.NewReader
	httpNewRequest          = http.NewRequest
	logWebcallStartFunc     = logWebcallStart
	logWebcallRequestFunc   = logWebcallRequest
	logWebcallResponseFunc  = logWebcallResponse
	logWebcallFinishFunc    = logWebcallFinish
	createHTTPRequestFunc   = createHTTPRequest
	getClientForRequestFunc = getClientForRequest
	clientDoWithRetryFunc   = clientDoWithRetry
	logErrorResponseFunc    = logErrorResponse
	logSuccessResponseFunc  = logSuccessResponse
	doRequestProcessingFunc = doRequestProcessing
	getDataTemplateFunc     = getDataTemplate
	parseResponseFunc       = parseResponse
)
