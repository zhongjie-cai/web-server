package webserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

// func pointers for injection / testing: panic.go
var (
	debugStack          = debug.Stack
	getRecoverErrorFunc = getRecoverError
	getDebugStackFunc   = getDebugStack
)

// func pointers for injection / testing: parameter.go
var (
	regexpMatchString = regexp.MatchString
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
