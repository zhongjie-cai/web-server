package webserver

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"os/signal"
	"reflect"
	"regexp"
	"runtime"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var (
	fmtSprintExpected                       int
	fmtSprintCalled                         int
	getErrorMessageFuncExpected             int
	getErrorMessageFuncCalled               int
	printInnerErrorsFuncExpected            int
	printInnerErrorsFuncCalled              int
	errorsIsExpected                        int
	errorsIsCalled                          int
	equalsErrorFuncExpected                 int
	equalsErrorFuncCalled                   int
	appErrorContainsFuncExpected            int
	appErrorContainsFuncCalled              int
	innerErrorContainsFuncExpected          int
	innerErrorContainsFuncCalled            int
	cleanupInnerErrorsFuncExpected          int
	cleanupInnerErrorsFuncCalled            int
	newAppErrorFuncExpected                 int
	newAppErrorFuncCalled                   int
	isInterfaceValueNilFuncExpected         int
	isInterfaceValueNilFuncCalled           int
	uuidNewExpected                         int
	uuidNewCalled                           int
	startApplicationFuncExpected            int
	startApplicationFuncCalled              int
	haltServerFuncExpected                  int
	haltServerFuncCalled                    int
	preBootstrapingFuncExpected             int
	preBootstrapingFuncCalled               int
	bootstrapFuncExpected                   int
	bootstrapFuncCalled                     int
	postBootstrapingFuncExpected            int
	postBootstrapingFuncCalled              int
	endApplicationFuncExpected              int
	endApplicationFuncCalled                int
	beginApplicationFuncExpected            int
	beginApplicationFuncCalled              int
	logAppRootFuncExpected                  int
	logAppRootFuncCalled                    int
	initializeHTTPClientsFuncExpected       int
	initializeHTTPClientsFuncCalled         int
	hostServerFuncExpected                  int
	hostServerFuncCalled                    int
	fmtPrintfExpected                       int
	fmtPrintfCalled                         int
	fmtSprintfExpected                      int
	fmtSprintfCalled                        int
	marshalIgnoreErrorFuncExpected          int
	marshalIgnoreErrorFuncCalled            int
	debugStackExpected                      int
	debugStackCalled                        int
	getRecoverErrorFuncExpected             int
	getRecoverErrorFuncCalled               int
	stringsSplitExpected                    int
	stringsSplitCalled                      int
	strconvAtoiExpected                     int
	strconvAtoiCalled                       int
	getRouteInfoFuncExpected                int
	getRouteInfoFuncCalled                  int
	initiateSessionFuncExpected             int
	initiateSessionFuncCalled               int
	getTimeNowUTCFuncExpected               int
	getTimeNowUTCFuncCalled                 int
	finalizeSessionFuncExpected             int
	finalizeSessionFuncCalled               int
	logEndpointEnterFuncExpected            int
	logEndpointEnterFuncCalled              int
	logEndpointExitFuncExpected             int
	logEndpointExitFuncCalled               int
	timeSinceExpected                       int
	timeSinceCalled                         int
	handlePanicFuncExpected                 int
	handlePanicFuncCalled                   int
	writeResponseFuncExpected               int
	writeResponseFuncCalled                 int
	handleActionFuncExpected                int
	handleActionFuncCalled                  int
	jsonNewEncoderExpected                  int
	jsonNewEncoderCalled                    int
	stringsTrimRightExpected                int
	stringsTrimRightCalled                  int
	jsonUnmarshalExpected                   int
	jsonUnmarshalCalled                     int
	fmtErrorfExpected                       int
	fmtErrorfCalled                         int
	reflectTypeOfExpected                   int
	reflectTypeOfCalled                     int
	stringsToLowerExpected                  int
	stringsToLowerCalled                    int
	strconvParseBoolExpected                int
	strconvParseBoolCalled                  int
	strconvParseIntExpected                 int
	strconvParseIntCalled                   int
	strconvParseFloatExpected               int
	strconvParseFloatCalled                 int
	strconvParseUintExpected                int
	strconvParseUintCalled                  int
	tryUnmarshalPrimitiveTypesFuncExpected  int
	tryUnmarshalPrimitiveTypesFuncCalled    int
	prepareLoggingFuncExpected              int
	prepareLoggingFuncCalled                int
	sortStringsExpected                     int
	sortStringsCalled                       int
	stringsJoinExpected                     int
	stringsJoinCalled                       int
	regexpMatchStringExpected               int
	regexpMatchStringCalled                 int
	reflectValueOfExpected                  int
	reflectValueOfCalled                    int
	stringsReplaceExpected                  int
	stringsReplaceCalled                    int
	doParameterReplacementFuncExpected      int
	doParameterReplacementFuncCalled        int
	evaluatePathWithParametersFuncExpected  int
	evaluatePathWithParametersFuncCalled    int
	evaluateQueriesFuncExpected             int
	evaluateQueriesFuncCalled               int
	registerRouteFuncExpected               int
	registerRouteFuncCalled                 int
	registerStaticFuncExpected              int
	registerStaticFuncCalled                int
	addMiddlewareFuncExpected               int
	addMiddlewareFuncCalled                 int
	muxNewRouterExpected                    int
	muxNewRouterCalled                      int
	registerRoutesFuncExpected              int
	registerRoutesFuncCalled                int
	registerStaticsFuncExpected             int
	registerStaticsFuncCalled               int
	registerMiddlewaresFuncExpected         int
	registerMiddlewaresFuncCalled           int
	walkRegisteredRoutesFuncExpected        int
	walkRegisteredRoutesFuncCalled          int
	registerErrorHandlersFuncExpected       int
	registerErrorHandlersFuncCalled         int
	ioutilReadAllExpected                   int
	ioutilReadAllCalled                     int
	ioutilNopCloserExpected                 int
	ioutilNopCloserCalled                   int
	bytesNewBufferExpected                  int
	bytesNewBufferCalled                    int
	skipResponseHandlingFuncExpected        int
	skipResponseHandlingFuncCalled          int
	constructResponseFuncExpected           int
	constructResponseFuncCalled             int
	logEndpointResponseFuncExpected         int
	logEndpointResponseFuncCalled           int
	httpStatusTextExpected                  int
	httpStatusTextCalled                    int
	strconvItoaExpected                     int
	strconvItoaCalled                       int
	getPathTemplateFuncExpected             int
	getPathTemplateFuncCalled               int
	getPathRegexpFuncExpected               int
	getPathRegexpFuncCalled                 int
	evaluateRouteFuncExpected               int
	evaluateRouteFuncCalled                 int
	muxCurrentRouteExpected                 int
	muxCurrentRouteCalled                   int
	getNameFuncExpected                     int
	getNameFuncCalled                       int
	getEndpointByNameFuncExpected           int
	getEndpointByNameFuncCalled             int
	instantiateRouterFuncExpected           int
	instantiateRouterFuncCalled             int
	runServerFuncExpected                   int
	runServerFuncCalled                     int
	createServerFuncExpected                int
	createServerFuncCalled                  int
	signalNotifyExpected                    int
	signalNotifyCalled                      int
	listenAndServeFuncExpected              int
	listenAndServeFuncCalled                int
	contextWithTimeoutExpected              int
	contextWithTimeoutCalled                int
	contextBackgroundExpected               int
	contextBackgroundCalled                 int
	shutdownServerFuncExpected              int
	shutdownServerFuncCalled                int
	evaluateServerErrorsFuncExpected        int
	evaluateServerErrorsFuncCalled          int
	getRequestBodyFuncExpected              int
	getRequestBodyFuncCalled                int
	logEndpointRequestFuncExpected          int
	logEndpointRequestFuncCalled            int
	tryUnmarshalFuncExpected                int
	tryUnmarshalFuncCalled                  int
	muxVarsExpected                         int
	muxVarsCalled                           int
	getAllQueriesFuncExpected               int
	getAllQueriesFuncCalled                 int
	textprotoCanonicalMIMEHeaderKeyExpected int
	textprotoCanonicalMIMEHeaderKeyCalled   int
	getAllHeadersFuncExpected               int
	getAllHeadersFuncCalled                 int
	jsonMarshalExpected                     int
	jsonMarshalCalled                       int
	runtimeCallerExpected                   int
	runtimeCallerCalled                     int
	runtimeFuncForPCExpected                int
	runtimeFuncForPCCalled                  int
	getMethodNameFuncExpected               int
	getMethodNameFuncCalled                 int
	logMethodEnterFuncExpected              int
	logMethodEnterFuncCalled                int
	logMethodParameterFuncExpected          int
	logMethodParameterFuncCalled            int
	logMethodLogicFuncExpected              int
	logMethodLogicFuncCalled                int
	logMethodReturnFuncExpected             int
	logMethodReturnFuncCalled               int
	logMethodExitFuncExpected               int
	logMethodExitFuncCalled                 int
	timeNowExpected                         int
	timeNowCalled                           int
	clientDoFuncExpected                    int
	clientDoFuncCalled                      int
	timeSleepExpected                       int
	timeSleepCalled                         int
	getHTTPTransportFuncExpected            int
	getHTTPTransportFuncCalled              int
	urlQueryEscapeExpected                  int
	urlQueryEscapeCalled                    int
	createQueryStringFuncExpected           int
	createQueryStringFuncCalled             int
	generateRequestURLFuncExpected          int
	generateRequestURLFuncCalled            int
	stringsNewReaderExpected                int
	stringsNewReaderCalled                  int
	httpNewRequestExpected                  int
	httpNewRequestCalled                    int
	logWebcallStartFuncExpected             int
	logWebcallStartFuncCalled               int
	logWebcallRequestFuncExpected           int
	logWebcallRequestFuncCalled             int
	logWebcallResponseFuncExpected          int
	logWebcallResponseFuncCalled            int
	logWebcallFinishFuncExpected            int
	logWebcallFinishFuncCalled              int
	createHTTPRequestFuncExpected           int
	createHTTPRequestFuncCalled             int
	getClientForRequestFuncExpected         int
	getClientForRequestFuncCalled           int
	clientDoWithRetryFuncExpected           int
	clientDoWithRetryFuncCalled             int
	logErrorResponseFuncExpected            int
	logErrorResponseFuncCalled              int
	logSuccessResponseFuncExpected          int
	logSuccessResponseFuncCalled            int
	doRequestProcessingFuncExpected         int
	doRequestProcessingFuncCalled           int
	getDataTemplateFuncExpected             int
	getDataTemplateFuncCalled               int
	parseResponseFuncExpected               int
	parseResponseFuncCalled                 int
)

func createMock(t *testing.T) {
	fmtSprintExpected = 0
	fmtSprintCalled = 0
	fmtSprint = func(a ...interface{}) string {
		fmtSprintCalled++
		return ""
	}
	getErrorMessageFuncExpected = 0
	getErrorMessageFuncCalled = 0
	getErrorMessageFunc = func(err error) string {
		getErrorMessageFuncCalled++
		return ""
	}
	printInnerErrorsFuncExpected = 0
	printInnerErrorsFuncCalled = 0
	printInnerErrorsFunc = func(innerErrors []*appError) string {
		printInnerErrorsFuncCalled++
		return ""
	}
	errorsIsExpected = 0
	errorsIsCalled = 0
	errorsIs = func(err, target error) bool {
		errorsIsCalled++
		return false
	}
	equalsErrorFuncExpected = 0
	equalsErrorFuncCalled = 0
	equalsErrorFunc = func(err, target error) bool {
		equalsErrorFuncCalled++
		return false
	}
	appErrorContainsFuncExpected = 0
	appErrorContainsFuncCalled = 0
	appErrorContainsFunc = func(appError AppError, err error) bool {
		appErrorContainsFuncCalled++
		return false
	}
	innerErrorContainsFuncExpected = 0
	innerErrorContainsFuncCalled = 0
	innerErrorContainsFunc = func(innerErrors []*appError, err error) bool {
		innerErrorContainsFuncCalled++
		return false
	}
	cleanupInnerErrorsFuncExpected = 0
	cleanupInnerErrorsFuncCalled = 0
	cleanupInnerErrorsFunc = func(innerErrors []error) []*appError {
		cleanupInnerErrorsFuncCalled++
		return nil
	}
	newAppErrorFuncExpected = 0
	newAppErrorFuncCalled = 0
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		return nil
	}
	isInterfaceValueNilFuncExpected = 0
	isInterfaceValueNilFuncCalled = 0
	isInterfaceValueNilFunc = func(i interface{}) bool {
		isInterfaceValueNilFuncCalled++
		return false
	}
	uuidNewExpected = 0
	uuidNewCalled = 0
	uuidNew = func() uuid.UUID {
		uuidNewCalled++
		return uuid.Nil
	}
	startApplicationFuncExpected = 0
	startApplicationFuncCalled = 0
	startApplicationFunc = func(app *application) {
		startApplicationFuncCalled++
	}
	haltServerFuncExpected = 0
	haltServerFuncCalled = 0
	haltServerFunc = func(shutdownSignal chan os.Signal) {
		haltServerFuncCalled++
	}
	preBootstrapingFuncExpected = 0
	preBootstrapingFuncCalled = 0
	preBootstrapingFunc = func(app *application) bool {
		preBootstrapingFuncCalled++
		return false
	}
	bootstrapFuncExpected = 0
	bootstrapFuncCalled = 0
	bootstrapFunc = func(app *application) {
		bootstrapFuncCalled++
	}
	postBootstrapingFuncExpected = 0
	postBootstrapingFuncCalled = 0
	postBootstrapingFunc = func(app *application) bool {
		postBootstrapingFuncCalled++
		return false
	}
	endApplicationFuncExpected = 0
	endApplicationFuncCalled = 0
	endApplicationFunc = func(app *application) {
		endApplicationFuncCalled++
	}
	beginApplicationFuncExpected = 0
	beginApplicationFuncCalled = 0
	beginApplicationFunc = func(app *application) {
		beginApplicationFuncCalled++
	}
	logAppRootFuncExpected = 0
	logAppRootFuncCalled = 0
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
	}
	initializeHTTPClientsFuncExpected = 0
	initializeHTTPClientsFuncCalled = 0
	initializeHTTPClientsFunc = func(webcallTimeout time.Duration, skipServerCertVerification bool, clientCertificate *tls.Certificate, roundTripperWrapper func(originalTransport http.RoundTripper) http.RoundTripper) {
		initializeHTTPClientsFuncCalled++
	}
	hostServerFuncExpected = 0
	hostServerFuncCalled = 0
	hostServerFunc = func(app *application, session *session, shutdownSignal chan os.Signal, started *bool) error {
		hostServerFuncCalled++
		return nil
	}
	fmtPrintfExpected = 0
	fmtPrintfCalled = 0
	fmtPrintf = func(format string, a ...interface{}) (n int, err error) {
		fmtPrintfCalled++
		return 0, nil
	}
	fmtSprintfExpected = 0
	fmtSprintfCalled = 0
	fmtSprintf = func(format string, a ...interface{}) string {
		fmtSprintfCalled++
		return ""
	}
	marshalIgnoreErrorFuncExpected = 0
	marshalIgnoreErrorFuncCalled = 0
	marshalIgnoreErrorFunc = func(v interface{}) string {
		marshalIgnoreErrorFuncCalled++
		return ""
	}
	debugStackExpected = 0
	debugStackCalled = 0
	debugStack = func() []byte {
		debugStackCalled++
		return nil
	}
	getRecoverErrorFuncExpected = 0
	getRecoverErrorFuncCalled = 0
	getRecoverErrorFunc = func(recoverResult interface{}) error {
		getRecoverErrorFuncCalled++
		return nil
	}
	stringsSplitExpected = 0
	stringsSplitCalled = 0
	stringsSplit = func(s, sep string) []string {
		stringsSplitCalled++
		return nil
	}
	strconvAtoiExpected = 0
	strconvAtoiCalled = 0
	strconvAtoi = func(s string) (int, error) {
		strconvAtoiCalled++
		return 0, nil
	}
	getRouteInfoFuncExpected = 0
	getRouteInfoFuncCalled = 0
	getRouteInfoFunc = func(httpRequest *http.Request, actionFuncMap map[string]ActionFunc) (string, ActionFunc, error) {
		getRouteInfoFuncCalled++
		return "", nil, nil
	}
	initiateSessionFuncExpected = 0
	initiateSessionFuncCalled = 0
	initiateSessionFunc = func(app *application, responseWriter http.ResponseWriter, httpRequest *http.Request) (*session, ActionFunc, error) {
		initiateSessionFuncCalled++
		return nil, nil, nil
	}
	getTimeNowUTCFuncExpected = 0
	getTimeNowUTCFuncCalled = 0
	getTimeNowUTCFunc = func() time.Time {
		getTimeNowUTCFuncCalled++
		return time.Time{}
	}
	finalizeSessionFuncExpected = 0
	finalizeSessionFuncCalled = 0
	finalizeSessionFunc = func(session *session, startTime time.Time, recoverResult interface{}) {
		finalizeSessionFuncCalled++
	}
	logEndpointEnterFuncExpected = 0
	logEndpointEnterFuncCalled = 0
	logEndpointEnterFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logEndpointEnterFuncCalled++
	}
	logEndpointExitFuncExpected = 0
	logEndpointExitFuncCalled = 0
	logEndpointExitFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logEndpointExitFuncCalled++
	}
	timeSinceExpected = 0
	timeSinceCalled = 0
	timeSince = func(t time.Time) time.Duration {
		timeSinceCalled++
		return 0
	}
	handlePanicFuncExpected = 0
	handlePanicFuncCalled = 0
	handlePanicFunc = func(session *session, recoverResult interface{}) {
		handlePanicFuncCalled++
	}
	writeResponseFuncExpected = 0
	writeResponseFuncCalled = 0
	writeResponseFunc = func(session *session, responseObject interface{}, responseError error) {
		writeResponseFuncCalled++
	}
	handleActionFuncExpected = 0
	handleActionFuncCalled = 0
	handleActionFunc = func(session *session, action ActionFunc) {
		handleActionFuncCalled++
	}
	jsonNewEncoderExpected = 0
	jsonNewEncoderCalled = 0
	jsonNewEncoder = func(w io.Writer) *json.Encoder {
		jsonNewEncoderCalled++
		return nil
	}
	stringsTrimRightExpected = 0
	stringsTrimRightCalled = 0
	stringsTrimRight = func(s string, cutset string) string {
		stringsTrimRightCalled++
		return ""
	}
	jsonUnmarshalExpected = 0
	jsonUnmarshalCalled = 0
	jsonUnmarshal = func(data []byte, v interface{}) error {
		jsonUnmarshalCalled++
		return nil
	}
	fmtErrorfExpected = 0
	fmtErrorfCalled = 0
	fmtErrorf = func(format string, a ...interface{}) error {
		fmtErrorfCalled++
		return nil
	}
	reflectTypeOfExpected = 0
	reflectTypeOfCalled = 0
	reflectTypeOf = func(i interface{}) reflect.Type {
		reflectTypeOfCalled++
		return nil
	}
	stringsToLowerExpected = 0
	stringsToLowerCalled = 0
	stringsToLower = func(s string) string {
		stringsToLowerCalled++
		return ""
	}
	strconvParseBoolExpected = 0
	strconvParseBoolCalled = 0
	strconvParseBool = func(str string) (bool, error) {
		strconvParseBoolCalled++
		return false, nil
	}
	strconvParseIntExpected = 0
	strconvParseIntCalled = 0
	strconvParseInt = func(s string, base int, bitSize int) (int64, error) {
		strconvParseIntCalled++
		return 0, nil
	}
	strconvParseFloatExpected = 0
	strconvParseFloatCalled = 0
	strconvParseFloat = func(s string, bitSize int) (float64, error) {
		strconvParseFloatCalled++
		return 0, nil
	}
	strconvParseUintExpected = 0
	strconvParseUintCalled = 0
	strconvParseUint = func(s string, base int, bitSize int) (uint64, error) {
		strconvParseUintCalled++
		return 0, nil
	}
	tryUnmarshalPrimitiveTypesFuncExpected = 0
	tryUnmarshalPrimitiveTypesFuncCalled = 0
	tryUnmarshalPrimitiveTypesFunc = func(value string, dataTemplate interface{}) bool {
		tryUnmarshalPrimitiveTypesFuncCalled++
		return false
	}
	prepareLoggingFuncExpected = 0
	prepareLoggingFuncCalled = 0
	prepareLoggingFunc = func(session *session, logType LogType, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		prepareLoggingFuncCalled++
	}
	sortStringsExpected = 0
	sortStringsCalled = 0
	sortStrings = func(a []string) {
		sortStringsCalled++
	}
	stringsJoinExpected = 0
	stringsJoinCalled = 0
	stringsJoin = func(a []string, sep string) string {
		stringsJoinCalled++
		return ""
	}
	regexpMatchStringExpected = 0
	regexpMatchStringCalled = 0
	regexpMatchString = func(pattern string, s string) (bool, error) {
		regexpMatchStringCalled++
		return false, nil
	}
	reflectValueOfExpected = 0
	reflectValueOfCalled = 0
	reflectValueOf = func(i interface{}) reflect.Value {
		reflectValueOfCalled++
		return reflect.Value{}
	}
	stringsReplaceExpected = 0
	stringsReplaceCalled = 0
	stringsReplace = func(s, old, new string, n int) string {
		stringsReplaceCalled++
		return ""
	}
	doParameterReplacementFuncExpected = 0
	doParameterReplacementFuncCalled = 0
	doParameterReplacementFunc = func(session *session, originalPath string, parameterName string, parameterType ParameterType) string {
		doParameterReplacementFuncCalled++
		return ""
	}
	evaluatePathWithParametersFuncExpected = 0
	evaluatePathWithParametersFuncCalled = 0
	evaluatePathWithParametersFunc = func(session *session, path string, parameters map[string]ParameterType) string {
		evaluatePathWithParametersFuncCalled++
		return ""
	}
	evaluateQueriesFuncExpected = 0
	evaluateQueriesFuncCalled = 0
	evaluateQueriesFunc = func(queries map[string]ParameterType) []string {
		evaluateQueriesFuncCalled++
		return nil
	}
	registerRouteFuncExpected = 0
	registerRouteFuncCalled = 0
	registerRouteFunc = func(router *mux.Router, endpoint string, method string, path string, queries []string, handleFunc func(http.ResponseWriter, *http.Request), actionFunc ActionFunc) (string, *mux.Route) {
		registerRouteFuncCalled++
		return "", nil
	}
	registerStaticFuncExpected = 0
	registerStaticFuncCalled = 0
	registerStaticFunc = func(router *mux.Router, name string, path string, handler http.Handler) *mux.Route {
		registerStaticFuncCalled++
		return nil
	}
	addMiddlewareFuncExpected = 0
	addMiddlewareFuncCalled = 0
	addMiddlewareFunc = func(router *mux.Router, middleware MiddlewareFunc) {
		addMiddlewareFuncCalled++
	}
	muxNewRouterExpected = 0
	muxNewRouterCalled = 0
	muxNewRouter = func() *mux.Router {
		muxNewRouterCalled++
		return nil
	}
	registerRoutesFuncExpected = 0
	registerRoutesFuncCalled = 0
	registerRoutesFunc = func(app *application, session *session, router *mux.Router) {
		registerRoutesFuncCalled++
	}
	registerStaticsFuncExpected = 0
	registerStaticsFuncCalled = 0
	registerStaticsFunc = func(session *session, router *mux.Router) {
		registerStaticsFuncCalled++
	}
	registerMiddlewaresFuncExpected = 0
	registerMiddlewaresFuncCalled = 0
	registerMiddlewaresFunc = func(session *session, router *mux.Router) {
		registerMiddlewaresFuncCalled++
	}
	walkRegisteredRoutesFuncExpected = 0
	walkRegisteredRoutesFuncCalled = 0
	walkRegisteredRoutesFunc = func(session *session, router *mux.Router) error {
		walkRegisteredRoutesFuncCalled++
		return nil
	}
	registerErrorHandlersFuncExpected = 0
	registerErrorHandlersFuncCalled = 0
	registerErrorHandlersFunc = func(customization Customization, router *mux.Router) {
		registerErrorHandlersFuncCalled++
	}
	ioutilReadAllExpected = 0
	ioutilReadAllCalled = 0
	ioutilReadAll = func(r io.Reader) ([]byte, error) {
		ioutilReadAllCalled++
		return nil, nil
	}
	ioutilNopCloserExpected = 0
	ioutilNopCloserCalled = 0
	ioutilNopCloser = func(r io.Reader) io.ReadCloser {
		ioutilNopCloserCalled++
		return nil
	}
	bytesNewBufferExpected = 0
	bytesNewBufferCalled = 0
	bytesNewBuffer = func(buf []byte) *bytes.Buffer {
		bytesNewBufferCalled++
		return nil
	}
	skipResponseHandlingFuncExpected = 0
	skipResponseHandlingFuncCalled = 0
	shouldSkipHandlingFunc = func(responseObject interface{}, responseError error) bool {
		skipResponseHandlingFuncCalled++
		return false
	}
	constructResponseFuncExpected = 0
	constructResponseFuncCalled = 0
	constructResponseFunc = func(session *session, responseObject interface{}, responseError error) (int, string) {
		constructResponseFuncCalled++
		return 0, ""
	}
	logEndpointResponseFuncExpected = 0
	logEndpointResponseFuncCalled = 0
	logEndpointResponseFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logEndpointResponseFuncCalled++
	}
	httpStatusTextExpected = 0
	httpStatusTextCalled = 0
	httpStatusText = func(code int) string {
		httpStatusTextCalled++
		return ""
	}
	strconvItoaExpected = 0
	strconvItoaCalled = 0
	strconvItoa = func(i int) string {
		strconvItoaCalled++
		return ""
	}
	getPathTemplateFuncExpected = 0
	getPathTemplateFuncCalled = 0
	getPathTemplateFunc = func(route *mux.Route) (string, error) {
		getPathTemplateFuncCalled++
		return "", nil
	}
	getPathRegexpFuncExpected = 0
	getPathRegexpFuncCalled = 0
	getPathRegexpFunc = func(route *mux.Route) (string, error) {
		getPathRegexpFuncCalled++
		return "", nil
	}
	evaluateRouteFuncExpected = 0
	evaluateRouteFuncCalled = 0
	evaluateRouteFunc = func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		evaluateRouteFuncCalled++
		return nil
	}
	muxCurrentRouteExpected = 0
	muxCurrentRouteCalled = 0
	muxCurrentRoute = func(r *http.Request) *mux.Route {
		muxCurrentRouteCalled++
		return nil
	}
	getNameFuncExpected = 0
	getNameFuncCalled = 0
	getNameFunc = func(route *mux.Route) string {
		getNameFuncCalled++
		return ""
	}
	getEndpointByNameFuncExpected = 0
	getEndpointByNameFuncCalled = 0
	getEndpointByNameFunc = func(name string) string {
		getEndpointByNameFuncCalled++
		return ""
	}
	instantiateRouterFuncExpected = 0
	instantiateRouterFuncCalled = 0
	instantiateRouterFunc = func(app *application, session *session) (*mux.Router, error) {
		instantiateRouterFuncCalled++
		return nil, nil
	}
	runServerFuncExpected = 0
	runServerFuncCalled = 0
	runServerFunc = func(address string, session *session, router *mux.Router, shutdownSignal chan os.Signal, started *bool) bool {
		runServerFuncCalled++
		return false
	}
	createServerFuncExpected = 0
	createServerFuncCalled = 0
	createServerFunc = func(address string, session *session, router *mux.Router) (*http.Server, bool) {
		createServerFuncCalled++
		return nil, false
	}
	signalNotifyExpected = 0
	signalNotifyCalled = 0
	signalNotify = func(c chan<- os.Signal, sig ...os.Signal) {
		signalNotifyCalled++
	}
	listenAndServeFuncExpected = 0
	listenAndServeFuncCalled = 0
	listenAndServeFunc = func(session *session, server *http.Server, serveHTTPS bool) error {
		listenAndServeFuncCalled++
		return nil
	}
	contextWithTimeoutExpected = 0
	contextWithTimeoutCalled = 0
	contextWithTimeout = func(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
		contextWithTimeoutCalled++
		return nil, nil
	}
	contextBackgroundExpected = 0
	contextBackgroundCalled = 0
	contextBackground = func() context.Context {
		contextBackgroundCalled++
		return nil
	}
	shutdownServerFuncExpected = 0
	shutdownServerFuncCalled = 0
	shutdownServerFunc = func(runtimeContext context.Context, server *http.Server) error {
		shutdownServerFuncCalled++
		return nil
	}
	evaluateServerErrorsFuncExpected = 0
	evaluateServerErrorsFuncCalled = 0
	evaluateServerErrorsFunc = func(session *session, hostError error, shutdownError error) bool {
		evaluateServerErrorsFuncCalled++
		return false
	}
	getRequestBodyFuncExpected = 0
	getRequestBodyFuncCalled = 0
	getRequestBodyFunc = func(httpRequest *http.Request) string {
		getRequestBodyFuncCalled++
		return ""
	}
	logEndpointRequestFuncExpected = 0
	logEndpointRequestFuncCalled = 0
	logEndpointRequestFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logEndpointRequestFuncCalled++
	}
	tryUnmarshalFuncExpected = 0
	tryUnmarshalFuncCalled = 0
	tryUnmarshalFunc = func(value string, dataTemplate interface{}) error {
		tryUnmarshalFuncCalled++
		return nil
	}
	muxVarsExpected = 0
	muxVarsCalled = 0
	muxVars = func(r *http.Request) map[string]string {
		muxVarsCalled++
		return nil
	}
	getAllQueriesFuncExpected = 0
	getAllQueriesFuncCalled = 0
	getAllQueriesFunc = func(session *session, name string) []string {
		getAllQueriesFuncCalled++
		return nil
	}
	textprotoCanonicalMIMEHeaderKeyExpected = 0
	textprotoCanonicalMIMEHeaderKeyCalled = 0
	textprotoCanonicalMIMEHeaderKey = func(s string) string {
		textprotoCanonicalMIMEHeaderKeyCalled++
		return ""
	}
	getAllHeadersFuncExpected = 0
	getAllHeadersFuncCalled = 0
	getAllHeadersFunc = func(session *session, name string) []string {
		getAllHeadersFuncCalled++
		return nil
	}
	jsonMarshalExpected = 0
	jsonMarshalCalled = 0
	jsonMarshal = func(v interface{}) ([]byte, error) {
		jsonMarshalCalled++
		return nil, nil
	}
	runtimeCallerExpected = 0
	runtimeCallerCalled = 0
	runtimeCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		runtimeCallerCalled++
		return 0, "", 0, false
	}
	runtimeFuncForPCExpected = 0
	runtimeFuncForPCCalled = 0
	runtimeFuncForPC = func(pc uintptr) *runtime.Func {
		runtimeFuncForPCCalled++
		return nil
	}
	getMethodNameFuncExpected = 0
	getMethodNameFuncCalled = 0
	getMethodNameFunc = func() string {
		getMethodNameFuncCalled++
		return ""
	}
	logMethodEnterFuncExpected = 0
	logMethodEnterFuncCalled = 0
	logMethodEnterFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logMethodEnterFuncCalled++
	}
	logMethodParameterFuncExpected = 0
	logMethodParameterFuncCalled = 0
	logMethodParameterFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logMethodParameterFuncCalled++
	}
	logMethodLogicFuncExpected = 0
	logMethodLogicFuncCalled = 0
	logMethodLogicFunc = func(session *session, logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logMethodLogicFuncCalled++
	}
	logMethodReturnFuncExpected = 0
	logMethodReturnFuncCalled = 0
	logMethodReturnFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logMethodReturnFuncCalled++
	}
	logMethodExitFuncExpected = 0
	logMethodExitFuncCalled = 0
	logMethodExitFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logMethodExitFuncCalled++
	}
	timeNowExpected = 0
	timeNowCalled = 0
	timeNow = func() time.Time {
		timeNowCalled++
		return time.Time{}
	}
	clientDoFuncExpected = 0
	clientDoFuncCalled = 0
	clientDoFunc = func(httpClient *http.Client, httpRequest *http.Request) (*http.Response, error) {
		clientDoFuncCalled++
		return nil, nil
	}
	timeSleepExpected = 0
	timeSleepCalled = 0
	timeSleep = func(time.Duration) {
		timeSleepCalled++
	}
	getHTTPTransportFuncExpected = 0
	getHTTPTransportFuncCalled = 0
	getHTTPTransportFunc = func(skipServerCertVerification bool, clientCertificate *tls.Certificate, roundTripperWrapper func(originalTransport http.RoundTripper) http.RoundTripper) http.RoundTripper {
		getHTTPTransportFuncCalled++
		return nil
	}
	urlQueryEscapeExpected = 0
	urlQueryEscapeCalled = 0
	urlQueryEscape = func(s string) string {
		urlQueryEscapeCalled++
		return ""
	}
	createQueryStringFuncExpected = 0
	createQueryStringFuncCalled = 0
	createQueryStringFunc = func(query map[string][]string) string {
		createQueryStringFuncCalled++
		return ""
	}
	generateRequestURLFuncExpected = 0
	generateRequestURLFuncCalled = 0
	generateRequestURLFunc = func(baseURL string, query map[string][]string) string {
		generateRequestURLFuncCalled++
		return ""
	}
	stringsNewReaderExpected = 0
	stringsNewReaderCalled = 0
	stringsNewReader = func(s string) *strings.Reader {
		stringsNewReaderCalled++
		return nil
	}
	httpNewRequestExpected = 0
	httpNewRequestCalled = 0
	httpNewRequest = func(method, url string, body io.Reader) (*http.Request, error) {
		httpNewRequestCalled++
		return nil, nil
	}
	logWebcallStartFuncExpected = 0
	logWebcallStartFuncCalled = 0
	logWebcallStartFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logWebcallStartFuncCalled++
	}
	logWebcallRequestFuncExpected = 0
	logWebcallRequestFuncCalled = 0
	logWebcallRequestFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logWebcallRequestFuncCalled++
	}
	logWebcallResponseFuncExpected = 0
	logWebcallResponseFuncCalled = 0
	logWebcallResponseFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logWebcallResponseFuncCalled++
	}
	logWebcallFinishFuncExpected = 0
	logWebcallFinishFuncCalled = 0
	logWebcallFinishFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logWebcallFinishFuncCalled++
	}
	createHTTPRequestFuncExpected = 0
	createHTTPRequestFuncCalled = 0
	createHTTPRequestFunc = func(webRequest *webRequest) (*http.Request, error) {
		createHTTPRequestFuncCalled++
		return nil, nil
	}
	getClientForRequestFuncExpected = 0
	getClientForRequestFuncCalled = 0
	getClientForRequestFunc = func(sendClientCert bool) *http.Client {
		getClientForRequestFuncCalled++
		return nil
	}
	clientDoWithRetryFuncExpected = 0
	clientDoWithRetryFuncCalled = 0
	clientDoWithRetryFunc = func(httpClient *http.Client, httpRequest *http.Request, connectivityRetryCount int, httpStatusRetryCount map[int]int, retryDelay time.Duration) (*http.Response, error) {
		clientDoWithRetryFuncCalled++
		return nil, nil
	}
	logErrorResponseFuncExpected = 0
	logErrorResponseFuncCalled = 0
	logErrorResponseFunc = func(session *session, responseError error, startTime time.Time) {
		logErrorResponseFuncCalled++
	}
	logSuccessResponseFuncExpected = 0
	logSuccessResponseFuncCalled = 0
	logSuccessResponseFunc = func(session *session, response *http.Response, startTime time.Time) {
		logSuccessResponseFuncCalled++
	}
	doRequestProcessingFuncExpected = 0
	doRequestProcessingFuncCalled = 0
	doRequestProcessingFunc = func(webRequest *webRequest) (*http.Response, error) {
		doRequestProcessingFuncCalled++
		return nil, nil
	}
	getDataTemplateFuncExpected = 0
	getDataTemplateFuncCalled = 0
	getDataTemplateFunc = func(session *session, statusCode int, dataReceivers []dataReceiver) interface{} {
		getDataTemplateFuncCalled++
		return nil
	}
	parseResponseFuncExpected = 0
	parseResponseFuncCalled = 0
	parseResponseFunc = func(session *session, body io.ReadCloser, dataTemplate interface{}) error {
		parseResponseFuncCalled++
		return nil
	}
}

func verifyAll(t *testing.T) {
	fmtSprint = fmt.Sprint
	assert.Equal(t, fmtSprintExpected, fmtSprintCalled, "Unexpected number of calls to fmtSprint")
	getErrorMessageFunc = getErrorMessage
	assert.Equal(t, getErrorMessageFuncExpected, getErrorMessageFuncCalled, "Unexpected number of calls to getErrorMessageFunc")
	printInnerErrorsFunc = printInnerErrors
	assert.Equal(t, printInnerErrorsFuncExpected, printInnerErrorsFuncCalled, "Unexpected number of calls to printInnerErrorsFunc")
	errorsIs = errors.Is
	assert.Equal(t, errorsIsExpected, errorsIsCalled, "Unexpected number of calls to errorsIs")
	equalsErrorFunc = equalsError
	assert.Equal(t, equalsErrorFuncExpected, equalsErrorFuncCalled, "Unexpected number of calls to equalsErrorFunc")
	appErrorContainsFunc = appErrorContains
	assert.Equal(t, appErrorContainsFuncExpected, appErrorContainsFuncCalled, "Unexpected number of calls to appErrorContainsFunc")
	innerErrorContainsFunc = innerErrorContains
	assert.Equal(t, innerErrorContainsFuncExpected, innerErrorContainsFuncCalled, "Unexpected number of calls to innerErrorContainsFunc")
	cleanupInnerErrorsFunc = cleanupInnerErrors
	assert.Equal(t, cleanupInnerErrorsFuncExpected, cleanupInnerErrorsFuncCalled, "Unexpected number of calls to cleanupInnerErrorsFunc")
	newAppErrorFunc = newAppError
	assert.Equal(t, newAppErrorFuncExpected, newAppErrorFuncCalled, "Unexpected number of calls to newAppErrorFunc")
	isInterfaceValueNilFunc = isInterfaceValueNil
	assert.Equal(t, isInterfaceValueNilFuncExpected, isInterfaceValueNilFuncCalled, "Unexpected number of calls to method isInterfaceValueNilFunc")
	uuidNew = uuid.New
	assert.Equal(t, uuidNewExpected, uuidNewCalled, "Unexpected number of calls to method uuidNew")
	startApplicationFunc = startApplication
	assert.Equal(t, startApplicationFuncExpected, startApplicationFuncCalled, "Unexpected number of calls to method startApplicationFunc")
	haltServerFunc = haltServer
	assert.Equal(t, haltServerFuncExpected, haltServerFuncCalled, "Unexpected number of calls to method haltServerFunc")
	preBootstrapingFunc = preBootstraping
	assert.Equal(t, preBootstrapingFuncExpected, preBootstrapingFuncCalled, "Unexpected number of calls to method preBootstrapingFunc")
	bootstrapFunc = bootstrap
	assert.Equal(t, bootstrapFuncExpected, bootstrapFuncCalled, "Unexpected number of calls to method bootstrapFunc")
	postBootstrapingFunc = postBootstraping
	assert.Equal(t, postBootstrapingFuncExpected, postBootstrapingFuncCalled, "Unexpected number of calls to method postBootstrapingFunc")
	endApplicationFunc = endApplication
	assert.Equal(t, endApplicationFuncExpected, endApplicationFuncCalled, "Unexpected number of calls to method endApplicationFunc")
	beginApplicationFunc = beginApplication
	assert.Equal(t, beginApplicationFuncExpected, beginApplicationFuncCalled, "Unexpected number of calls to method beginApplicationFunc")
	logAppRootFunc = logAppRoot
	assert.Equal(t, logAppRootFuncExpected, logAppRootFuncCalled, "Unexpected number of calls to method logAppRootFunc")
	initializeHTTPClientsFunc = initializeHTTPClients
	assert.Equal(t, initializeHTTPClientsFuncExpected, initializeHTTPClientsFuncCalled, "Unexpected number of calls to method initializeHTTPClientsFunc")
	hostServerFunc = hostServer
	assert.Equal(t, hostServerFuncExpected, hostServerFuncCalled, "Unexpected number of calls to method hostServerFunc")
	fmtPrintf = fmt.Printf
	assert.Equal(t, fmtPrintfExpected, fmtPrintfCalled, "Unexpected number of calls to method fmtPrintf")
	fmtSprintf = fmt.Sprintf
	assert.Equal(t, fmtSprintfExpected, fmtSprintfCalled, "Unexpected number of calls to method fmtSprintf")
	marshalIgnoreErrorFunc = marshalIgnoreError
	assert.Equal(t, marshalIgnoreErrorFuncExpected, marshalIgnoreErrorFuncCalled, "Unexpected number of calls to method marshalIgnoreErrorFunc")
	debugStack = debug.Stack
	assert.Equal(t, debugStackExpected, debugStackCalled, "Unexpected number of calls to debugStack")
	getRecoverErrorFunc = getRecoverError
	assert.Equal(t, getRecoverErrorFuncExpected, getRecoverErrorFuncCalled, "Unexpected number of calls to getRecoverErrorFunc")
	stringsSplit = strings.Split
	assert.Equal(t, stringsSplitExpected, stringsSplitCalled, "Unexpected number of calls to method stringsSplit")
	strconvAtoi = strconv.Atoi
	assert.Equal(t, strconvAtoiExpected, strconvAtoiCalled, "Unexpected number of calls to method strconvAtoi")
	getRouteInfoFunc = getRouteInfo
	assert.Equal(t, getRouteInfoFuncExpected, getRouteInfoFuncCalled, "Unexpected number of calls to method getRouteInfoFunc")
	initiateSessionFunc = initiateSession
	assert.Equal(t, initiateSessionFuncExpected, initiateSessionFuncCalled, "Unexpected number of calls to method initiateSessionFunc")
	getTimeNowUTCFunc = getTimeNowUTC
	assert.Equal(t, getTimeNowUTCFuncExpected, getTimeNowUTCFuncCalled, "Unexpected number of calls to method getTimeNowUTCFunc")
	finalizeSessionFunc = finalizeSession
	assert.Equal(t, finalizeSessionFuncExpected, finalizeSessionFuncCalled, "Unexpected number of calls to method finalizeSessionFunc")
	logEndpointEnterFunc = logEndpointEnter
	assert.Equal(t, logEndpointEnterFuncExpected, logEndpointEnterFuncCalled, "Unexpected number of calls to method logEndpointEnterFunc")
	logEndpointExitFunc = logEndpointExit
	assert.Equal(t, logEndpointExitFuncExpected, logEndpointExitFuncCalled, "Unexpected number of calls to method logEndpointExitFunc")
	timeSince = time.Since
	assert.Equal(t, timeSinceExpected, timeSinceCalled, "Unexpected number of calls to method timeSince")
	handlePanicFunc = handlePanic
	assert.Equal(t, handlePanicFuncExpected, handlePanicFuncCalled, "Unexpected number of calls to method handlePanicFunc")
	writeResponseFunc = writeResponse
	assert.Equal(t, writeResponseFuncExpected, writeResponseFuncCalled, "Unexpected number of calls to method writeResponseFunc")
	handleActionFunc = handleAction
	assert.Equal(t, handleActionFuncExpected, handleActionFuncCalled, "Unexpected number of calls to method handleActionFunc")
	jsonNewEncoder = json.NewEncoder
	assert.Equal(t, jsonNewEncoderExpected, jsonNewEncoderCalled, "Unexpected number of calls to jsonNewEncoder")
	stringsTrimRight = strings.TrimRight
	assert.Equal(t, stringsTrimRightExpected, stringsTrimRightCalled, "Unexpected number of calls to stringsTrimRight")
	jsonUnmarshal = json.Unmarshal
	assert.Equal(t, jsonUnmarshalExpected, jsonUnmarshalCalled, "Unexpected number of calls to jsonUnmarshal")
	fmtErrorf = fmt.Errorf
	assert.Equal(t, fmtErrorfExpected, fmtErrorfCalled, "Unexpected number of calls to fmtErrorf")
	reflectTypeOf = reflect.TypeOf
	assert.Equal(t, reflectTypeOfExpected, reflectTypeOfCalled, "Unexpected number of calls to reflectTypeOf")
	stringsToLower = strings.ToLower
	assert.Equal(t, stringsToLowerExpected, stringsToLowerCalled, "Unexpected number of calls to stringsToLower")
	strconvParseBool = strconv.ParseBool
	assert.Equal(t, strconvParseBoolExpected, strconvParseBoolCalled, "Unexpected number of calls to strconvParseBool")
	strconvParseInt = strconv.ParseInt
	assert.Equal(t, strconvParseIntExpected, strconvParseIntCalled, "Unexpected number of calls to strconvParseInt")
	strconvParseFloat = strconv.ParseFloat
	assert.Equal(t, strconvParseFloatExpected, strconvParseFloatCalled, "Unexpected number of calls to strconvParseFloat")
	strconvParseUint = strconv.ParseUint
	assert.Equal(t, strconvParseUintExpected, strconvParseUintCalled, "Unexpected number of calls to strconvParseUint")
	tryUnmarshalPrimitiveTypesFunc = tryUnmarshalPrimitiveTypes
	assert.Equal(t, tryUnmarshalPrimitiveTypesFuncExpected, tryUnmarshalPrimitiveTypesFuncCalled, "Unexpected number of calls to tryUnmarshalPrimitiveTypesFunc")
	prepareLoggingFunc = prepareLogging
	assert.Equal(t, prepareLoggingFuncExpected, prepareLoggingFuncCalled, "Unexpected number of calls to prepareLoggingFunc")
	sortStrings = sort.Strings
	assert.Equal(t, sortStringsExpected, sortStringsCalled, "Unexpected number of calls to sortStrings")
	stringsJoin = strings.Join
	assert.Equal(t, stringsJoinExpected, stringsJoinCalled, "Unexpected number of calls to stringsJoin")
	regexpMatchString = regexp.MatchString
	assert.Equal(t, regexpMatchStringExpected, regexpMatchStringCalled, "Unexpected number of calls to regexpMatchString")
	reflectValueOf = reflect.ValueOf
	assert.Equal(t, reflectValueOfExpected, reflectValueOfCalled, "Unexpected number of calls to reflectValueOf")
	stringsReplace = strings.Replace
	assert.Equal(t, stringsReplaceExpected, stringsReplaceCalled, "Unexpected number of calls to method stringsReplace")
	doParameterReplacementFunc = doParameterReplacement
	assert.Equal(t, doParameterReplacementFuncExpected, doParameterReplacementFuncCalled, "Unexpected number of calls to method doParameterReplacementFunc")
	evaluatePathWithParametersFunc = evaluatePathWithParameters
	assert.Equal(t, evaluatePathWithParametersFuncExpected, evaluatePathWithParametersFuncCalled, "Unexpected number of calls to method evaluatePathWithParametersFunc")
	evaluateQueriesFunc = evaluateQueries
	assert.Equal(t, evaluateQueriesFuncExpected, evaluateQueriesFuncCalled, "Unexpected number of calls to method evaluateQueriesFunc")
	registerRouteFunc = registerRoute
	assert.Equal(t, registerRouteFuncExpected, registerRouteFuncCalled, "Unexpected number of calls to method registerRouteFunc")
	registerStaticFunc = registerStatic
	assert.Equal(t, registerStaticFuncExpected, registerStaticFuncCalled, "Unexpected number of calls to method registerStaticFunc")
	addMiddlewareFunc = addMiddleware
	assert.Equal(t, addMiddlewareFuncExpected, addMiddlewareFuncCalled, "Unexpected number of calls to method addMiddlewareFunc")
	muxNewRouter = mux.NewRouter
	assert.Equal(t, muxNewRouterExpected, muxNewRouterCalled, "Unexpected number of calls to method muxNewRouter")
	registerRoutesFunc = registerRoutes
	assert.Equal(t, registerRoutesFuncExpected, registerRoutesFuncCalled, "Unexpected number of calls to method registerRoutesFunc")
	registerStaticsFunc = registerStatics
	assert.Equal(t, registerStaticsFuncExpected, registerStaticsFuncCalled, "Unexpected number of calls to method registerStaticsFunc")
	registerMiddlewaresFunc = registerMiddlewares
	assert.Equal(t, registerMiddlewaresFuncExpected, registerMiddlewaresFuncCalled, "Unexpected number of calls to method registerMiddlewaresFunc")
	walkRegisteredRoutesFunc = walkRegisteredRoutes
	assert.Equal(t, walkRegisteredRoutesFuncExpected, walkRegisteredRoutesFuncCalled, "Unexpected number of calls to method walkRegisteredRoutesFunc")
	registerErrorHandlersFunc = registerErrorHandlers
	assert.Equal(t, registerErrorHandlersFuncExpected, registerErrorHandlersFuncCalled, "Unexpected number of calls to method registerErrorHandlersFunc")
	ioutilReadAll = ioutil.ReadAll
	assert.Equal(t, ioutilReadAllExpected, ioutilReadAllCalled, "Unexpected number of calls to ioutilReadAll")
	ioutilNopCloser = ioutil.NopCloser
	assert.Equal(t, ioutilNopCloserExpected, ioutilNopCloserCalled, "Unexpected number of calls to ioutilNopCloser")
	bytesNewBuffer = bytes.NewBuffer
	assert.Equal(t, bytesNewBufferExpected, bytesNewBufferCalled, "Unexpected number of calls to bytesNewBuffer")
	shouldSkipHandlingFunc = shouldSkipHandling
	assert.Equal(t, skipResponseHandlingFuncExpected, skipResponseHandlingFuncCalled, "Unexpected number of calls to shouldSkipHandlingFunc")
	constructResponseFunc = constructResponse
	assert.Equal(t, constructResponseFuncExpected, constructResponseFuncCalled, "Unexpected number of calls to method constructResponseFunc")
	logEndpointResponseFunc = logEndpointResponse
	assert.Equal(t, logEndpointResponseFuncExpected, logEndpointResponseFuncCalled, "Unexpected number of calls to method logEndpointResponseFunc")
	httpStatusText = http.StatusText
	assert.Equal(t, httpStatusTextExpected, httpStatusTextCalled, "Unexpected number of calls to method httpStatusText")
	strconvItoa = strconv.Itoa
	assert.Equal(t, strconvItoaExpected, strconvItoaCalled, "Unexpected number of calls to method strconvItoa")
	getPathTemplateFunc = getPathTemplate
	assert.Equal(t, getPathTemplateFuncExpected, getPathTemplateFuncCalled, "Unexpected number of calls to method getPathTemplateFunc")
	getPathRegexpFunc = getPathRegexp
	assert.Equal(t, getPathRegexpFuncExpected, getPathRegexpFuncCalled, "Unexpected number of calls to method getPathRegexpFunc")
	evaluateRouteFunc = evaluateRoute
	assert.Equal(t, evaluateRouteFuncExpected, evaluateRouteFuncCalled, "Unexpected number of calls to method evaluateRouteFunc")
	muxCurrentRoute = mux.CurrentRoute
	assert.Equal(t, muxCurrentRouteExpected, muxCurrentRouteCalled, "Unexpected number of calls to method muxCurrentRoute")
	getNameFunc = getName
	assert.Equal(t, getNameFuncExpected, getNameFuncCalled, "Unexpected number of calls to method getNameFunc")
	getEndpointByNameFunc = getEndpointByName
	assert.Equal(t, getEndpointByNameFuncExpected, getEndpointByNameFuncCalled, "Unexpected number of calls to method getEndpointByNameFunc")
	instantiateRouterFunc = instantiateRouter
	assert.Equal(t, instantiateRouterFuncExpected, instantiateRouterFuncCalled, "Unexpected number of calls to method instantiateRouterFunc")
	runServerFunc = runServer
	assert.Equal(t, runServerFuncExpected, runServerFuncCalled, "Unexpected number of calls to method runServerFunc")
	createServerFunc = createServer
	assert.Equal(t, createServerFuncExpected, createServerFuncCalled, "Unexpected number of calls to method createServerFunc")
	signalNotify = signal.Notify
	assert.Equal(t, signalNotifyExpected, signalNotifyCalled, "Unexpected number of calls to method signalNotify")
	listenAndServeFunc = listenAndServe
	assert.Equal(t, listenAndServeFuncExpected, listenAndServeFuncCalled, "Unexpected number of calls to method listenAndServeFunc")
	contextWithTimeout = context.WithTimeout
	assert.Equal(t, contextWithTimeoutExpected, contextWithTimeoutCalled, "Unexpected number of calls to method contextWithTimeout")
	contextBackground = context.Background
	assert.Equal(t, contextBackgroundExpected, contextBackgroundCalled, "Unexpected number of calls to method contextBackground")
	shutdownServerFunc = shutdownServer
	assert.Equal(t, shutdownServerFuncExpected, shutdownServerFuncCalled, "Unexpected number of calls to method shutdownServerFunc")
	evaluateServerErrorsFunc = evaluateServerErrors
	assert.Equal(t, evaluateServerErrorsFuncExpected, evaluateServerErrorsFuncCalled, "Unexpected number of calls to method evaluateServerErrorsFunc")
	getRequestBodyFunc = getRequestBody
	assert.Equal(t, getRequestBodyFuncExpected, getRequestBodyFuncCalled, "Unexpected number of calls to method getRequestBodyFunc")
	logEndpointRequestFunc = logEndpointRequest
	assert.Equal(t, logEndpointRequestFuncExpected, logEndpointRequestFuncCalled, "Unexpected number of calls to method logEndpointRequestFunc")
	tryUnmarshalFunc = tryUnmarshal
	assert.Equal(t, tryUnmarshalFuncExpected, tryUnmarshalFuncCalled, "Unexpected number of calls to method tryUnmarshalFunc")
	muxVars = mux.Vars
	assert.Equal(t, muxVarsExpected, muxVarsCalled, "Unexpected number of calls to method muxVars")
	getAllQueriesFunc = getAllQueries
	assert.Equal(t, getAllQueriesFuncExpected, getAllQueriesFuncCalled, "Unexpected number of calls to method getAllQueriesFunc")
	textprotoCanonicalMIMEHeaderKey = textproto.CanonicalMIMEHeaderKey
	assert.Equal(t, textprotoCanonicalMIMEHeaderKeyExpected, textprotoCanonicalMIMEHeaderKeyCalled, "Unexpected number of calls to method textprotoCanonicalMIMEHeaderKey")
	getAllHeadersFunc = getAllHeaders
	assert.Equal(t, getAllHeadersFuncExpected, getAllHeadersFuncCalled, "Unexpected number of calls to method getAllHeadersFunc")
	jsonMarshal = json.Marshal
	assert.Equal(t, jsonMarshalExpected, jsonMarshalCalled, "Unexpected number of calls to method jsonMarshal")
	runtimeCaller = runtime.Caller
	assert.Equal(t, runtimeCallerExpected, runtimeCallerCalled, "Unexpected number of calls to method runtimeCaller")
	runtimeFuncForPC = runtime.FuncForPC
	assert.Equal(t, runtimeFuncForPCExpected, runtimeFuncForPCCalled, "Unexpected number of calls to method runtimeFuncForPC")
	getMethodNameFunc = getMethodName
	assert.Equal(t, getMethodNameFuncExpected, getMethodNameFuncCalled, "Unexpected number of calls to method getMethodNameFunc")
	logMethodEnterFunc = logMethodEnter
	assert.Equal(t, logMethodEnterFuncExpected, logMethodEnterFuncCalled, "Unexpected number of calls to method logMethodEnterFunc")
	logMethodParameterFunc = logMethodParameter
	assert.Equal(t, logMethodParameterFuncExpected, logMethodParameterFuncCalled, "Unexpected number of calls to method logMethodParameterFunc")
	logMethodLogicFunc = logMethodLogic
	assert.Equal(t, logMethodLogicFuncExpected, logMethodLogicFuncCalled, "Unexpected number of calls to method logMethodLogicFunc")
	logMethodReturnFunc = logMethodReturn
	assert.Equal(t, logMethodReturnFuncExpected, logMethodReturnFuncCalled, "Unexpected number of calls to method logMethodReturnFunc")
	logMethodExitFunc = logMethodExit
	assert.Equal(t, logMethodExitFuncExpected, logMethodExitFuncCalled, "Unexpected number of calls to method logMethodExitFunc")
	timeNow = time.Now
	assert.Equal(t, timeNowExpected, timeNowCalled, "Unexpected number of calls to timeNow")
	clientDoFunc = clientDo
	assert.Equal(t, clientDoFuncExpected, clientDoFuncCalled, "Unexpected number of calls to method clientDoFunc")
	timeSleep = time.Sleep
	assert.Equal(t, timeSleepExpected, timeSleepCalled, "Unexpected number of calls to method timeSleep")
	getHTTPTransportFunc = getHTTPTransport
	assert.Equal(t, getHTTPTransportFuncExpected, getHTTPTransportFuncCalled, "Unexpected number of calls to method getHTTPTransportFunc")
	urlQueryEscape = url.QueryEscape
	assert.Equal(t, urlQueryEscapeExpected, urlQueryEscapeCalled, "Unexpected number of calls to method urlQueryEscape")
	createQueryStringFunc = createQueryString
	assert.Equal(t, createQueryStringFuncExpected, createQueryStringFuncCalled, "Unexpected number of calls to method createQueryStringFunc")
	generateRequestURLFunc = generateRequestURL
	assert.Equal(t, generateRequestURLFuncExpected, generateRequestURLFuncCalled, "Unexpected number of calls to method generateRequestURLFunc")
	stringsNewReader = strings.NewReader
	assert.Equal(t, stringsNewReaderExpected, stringsNewReaderCalled, "Unexpected number of calls to method stringsNewReader")
	httpNewRequest = http.NewRequest
	assert.Equal(t, httpNewRequestExpected, httpNewRequestCalled, "Unexpected number of calls to method httpNewRequest")
	logWebcallStartFunc = logWebcallStart
	assert.Equal(t, logWebcallStartFuncExpected, logWebcallStartFuncCalled, "Unexpected number of calls to method logWebcallStartFunc")
	logWebcallRequestFunc = logWebcallRequest
	assert.Equal(t, logWebcallRequestFuncExpected, logWebcallRequestFuncCalled, "Unexpected number of calls to method logWebcallRequestFunc")
	logWebcallResponseFunc = logWebcallResponse
	assert.Equal(t, logWebcallResponseFuncExpected, logWebcallResponseFuncCalled, "Unexpected number of calls to method logWebcallResponseFunc")
	logWebcallFinishFunc = logWebcallFinish
	assert.Equal(t, logWebcallFinishFuncExpected, logWebcallFinishFuncCalled, "Unexpected number of calls to method logWebcallFinishFunc")
	createHTTPRequestFunc = createHTTPRequest
	assert.Equal(t, createHTTPRequestFuncExpected, createHTTPRequestFuncCalled, "Unexpected number of calls to method createHTTPRequestFunc")
	getClientForRequestFunc = getClientForRequest
	assert.Equal(t, getClientForRequestFuncExpected, getClientForRequestFuncCalled, "Unexpected number of calls to method getClientForRequestFunc")
	clientDoWithRetryFunc = clientDoWithRetry
	assert.Equal(t, clientDoWithRetryFuncExpected, clientDoWithRetryFuncCalled, "Unexpected number of calls to method clientDoWithRetryFunc")
	logErrorResponseFunc = logErrorResponse
	assert.Equal(t, logErrorResponseFuncExpected, logErrorResponseFuncCalled, "Unexpected number of calls to method logErrorResponseFunc")
	logSuccessResponseFunc = logSuccessResponse
	assert.Equal(t, logSuccessResponseFuncExpected, logSuccessResponseFuncCalled, "Unexpected number of calls to method logSuccessResponseFunc")
	doRequestProcessingFunc = doRequestProcessing
	assert.Equal(t, doRequestProcessingFuncExpected, doRequestProcessingFuncCalled, "Unexpected number of calls to method doRequestProcessingFunc")
	getDataTemplateFunc = getDataTemplate
	assert.Equal(t, getDataTemplateFuncExpected, getDataTemplateFuncCalled, "Unexpected number of calls to method getDataTemplateFunc")
	parseResponseFunc = parseResponse
	assert.Equal(t, parseResponseFuncExpected, parseResponseFuncCalled, "Unexpected number of calls to method parseResponseFunc")
}

func functionPointerEquals(t *testing.T, expectFunc interface{}, actualFunc interface{}) {
	var expectValue = fmt.Sprintf("%v", reflect.ValueOf(expectFunc))
	var actualValue = fmt.Sprintf("%v", reflect.ValueOf(actualFunc))
	assert.Equal(t, expectValue, actualValue)
}

// mock structs
type dummyAppError struct {
	t *testing.T
}

func (dummyAppError *dummyAppError) Error() string {
	assert.Fail(dummyAppError.t, "Unexpected call to method AppError.Error")
	return ""
}

func (dummyAppError *dummyAppError) ErrorCode() string {
	assert.Fail(dummyAppError.t, "Unexpected call to method AppError.ErrorCode")
	return ""
}

func (dummyAppError *dummyAppError) HTTPStatusCode() int {
	assert.Fail(dummyAppError.t, "Unexpected call to method AppError.HTTPStatusCode")
	return 0
}

func (dummyAppError *dummyAppError) HTTPResponseMessage() string {
	assert.Fail(dummyAppError.t, "Unexpected call to method AppError.HTTPResponseMessage")
	return ""
}

func (dummyAppError *dummyAppError) Contains(err error) bool {
	assert.Fail(dummyAppError.t, "Unexpected call to method AppError.Contains")
	return false
}

func (dummyAppError *dummyAppError) Wrap(innerErrors ...error) AppError {
	assert.Fail(dummyAppError.t, "Unexpected call to method AppError.Wrap")
	return nil
}

type dummyCustomization struct {
	t *testing.T
}

func (customization *dummyCustomization) PreBootstrap() error {
	assert.Fail(customization.t, "Unexpected call to PreBootstrap")
	return nil
}

func (customization *dummyCustomization) PostBootstrap() error {
	assert.Fail(customization.t, "Unexpected call to PostBootstrap")
	return nil
}

func (customization *dummyCustomization) AppClosing() error {
	assert.Fail(customization.t, "Unexpected call to AppClosing")
	return nil
}

func (customization *dummyCustomization) Log(session Session, logType LogType, logLevel LogLevel, category, subcategory, description string) {
	assert.Fail(customization.t, "Unexpected call to Log")
}

func (customization *dummyCustomization) ServerCert() *tls.Certificate {
	assert.Fail(customization.t, "Unexpected call to ServerCert")
	return nil
}

func (customization *dummyCustomization) CaCertPool() *x509.CertPool {
	assert.Fail(customization.t, "Unexpected call to CaCertPool")
	return nil
}

func (customization *dummyCustomization) GraceShutdownWaitTime() time.Duration {
	assert.Fail(customization.t, "Unexpected call to GraceShutdownWaitTime")
	return 0
}

func (customization *dummyCustomization) Routes() []Route {
	assert.Fail(customization.t, "Unexpected call to Routes")
	return nil
}

func (customization *dummyCustomization) Statics() []Static {
	assert.Fail(customization.t, "Unexpected call to Statics")
	return nil
}

func (customization *dummyCustomization) Middlewares() []MiddlewareFunc {
	assert.Fail(customization.t, "Unexpected call to Middlewares")
	return nil
}

func (customization *dummyCustomization) InstrumentRouter(router *mux.Router) *mux.Router {
	assert.Fail(customization.t, "Unexpected call to InstrumentRouter")
	return nil
}

func (customization *dummyCustomization) WrapHandler(handler http.Handler) http.Handler {
	assert.Fail(customization.t, "Unexpected call to WrapHandler")
	return nil
}

func (customization *dummyCustomization) Listener() net.Listener {
	assert.Fail(customization.t, "Unexpected call to Listener")
	return nil
}

func (customization *dummyCustomization) PreAction(session Session) error {
	assert.Fail(customization.t, "Unexpected call to PreAction")
	return nil
}

func (customization *dummyCustomization) PostAction(session Session) error {
	assert.Fail(customization.t, "Unexpected call to PostAction")
	return nil
}

func (customization *dummyCustomization) InterpretSuccess(responseContent interface{}) (int, string) {
	assert.Fail(customization.t, "Unexpected call to InterpretSuccess")
	return 0, ""
}

func (customization *dummyCustomization) InterpretError(err error) (int, string) {
	assert.Fail(customization.t, "Unexpected call to InterpretError")
	return 0, ""
}

func (customization *dummyCustomization) RecoverPanic(session Session, recoverResult interface{}) (interface{}, error) {
	assert.Fail(customization.t, "Unexpected call to RecoverPanic")
	return nil, nil
}

func (customization *dummyCustomization) NotFoundHandler() http.Handler {
	assert.Fail(customization.t, "Unexpected call to NotFoundHandler")
	return nil
}

func (customization *dummyCustomization) MethodNotAllowedHandler() http.Handler {
	assert.Fail(customization.t, "Unexpected call to MethodNotAllowedHandler")
	return nil
}

func (customization *dummyCustomization) ClientCert() *tls.Certificate {
	assert.Fail(customization.t, "Unexpected call to ClientCert")
	return nil
}

func (customization *dummyCustomization) DefaultTimeout() time.Duration {
	assert.Fail(customization.t, "Unexpected call to DefaultTimeout")
	return 0
}

func (customization *dummyCustomization) SkipServerCertVerification() bool {
	assert.Fail(customization.t, "Unexpected call to SkipServerCertVerification")
	return false
}

func (customization *dummyCustomization) RoundTripper(originalTransport http.RoundTripper) http.RoundTripper {
	assert.Fail(customization.t, "Unexpected call to RoundTripper")
	return nil
}

func (customization *dummyCustomization) WrapRequest(session Session, httpRequest *http.Request) *http.Request {
	assert.Fail(customization.t, "Unexpected call to WrapRequest")
	return nil
}

type dummySession struct {
	t *testing.T
}

func (session *dummySession) GetID() uuid.UUID {
	assert.Fail(session.t, "Unexpected call to GetID")
	return uuid.Nil
}

func (session *dummySession) GetName() string {
	assert.Fail(session.t, "Unexpected call to GetName")
	return ""
}

func (session *dummySession) GetRequest() *http.Request {
	assert.Fail(session.t, "Unexpected call to GetRequest")
	return nil
}

func (session *dummySession) GetResponseWriter() http.ResponseWriter {
	assert.Fail(session.t, "Unexpected call to GetResponseWriter")
	return nil
}

func (session *dummySession) GetRequestBody(dataTemplate interface{}) error {
	assert.Fail(session.t, "Unexpected call to GetRequestBody")
	return nil
}

func (session *dummySession) GetRequestParameter(name string, dataTemplate interface{}) error {
	assert.Fail(session.t, "Unexpected call to GetRequestParameter")
	return nil
}

func (session *dummySession) GetRequestQuery(name string, index int, dataTemplate interface{}) error {
	assert.Fail(session.t, "Unexpected call to GetRequestQuery")
	return nil
}

func (session *dummySession) GetRequestHeader(name string, index int, dataTemplate interface{}) error {
	assert.Fail(session.t, "Unexpected call to GetRequestHeader")
	return nil
}

func (session *dummySession) Attach(name string, value interface{}) bool {
	assert.Fail(session.t, "Unexpected call to Attach")
	return false
}

func (session *dummySession) Detach(name string) bool {
	assert.Fail(session.t, "Unexpected call to Detach")
	return false
}

func (session *dummySession) GetRawAttachment(name string) (interface{}, bool) {
	assert.Fail(session.t, "Unexpected call to GetRawAttachment")
	return nil, false
}

func (session *dummySession) GetAttachment(name string, dataTemplate interface{}) bool {
	assert.Fail(session.t, "Unexpected call to GetAttachment")
	return false
}

func (session *dummySession) LogMethodEnter() {
	assert.Fail(session.t, "Unexpected call to LogMethodEnter")
}

func (session *dummySession) LogMethodParameter(parameters ...interface{}) {
	assert.Fail(session.t, "Unexpected call to LogMethodParameter")
}

func (session *dummySession) LogMethodLogic(logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	assert.Fail(session.t, "Unexpected call to LogMethodLogic")
}

func (session *dummySession) LogMethodReturn(returns ...interface{}) {
	assert.Fail(session.t, "Unexpected call to LogMethodReturn")
}

func (session *dummySession) LogMethodExit() {
	assert.Fail(session.t, "Unexpected call to LogMethodExit")
}

func (session *dummySession) CreateWebcallRequest(method string, url string, payload string, sendClientCert bool) WebRequest {
	assert.Fail(session.t, "Unexpected call to CreateWebcallRequest")
	return nil
}

type dummyTransport struct {
	t *testing.T
}

func (transport *dummyTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	assert.Fail(transport.t, "Unexpected call to RoundTrip")
	return nil, nil
}

type dummyResponseWriter struct {
	t *testing.T
}

func (drw *dummyResponseWriter) Header() http.Header {
	assert.Fail(drw.t, "Unexpected number of calls to Header")
	return nil
}

func (drw *dummyResponseWriter) Write([]byte) (int, error) {
	assert.Fail(drw.t, "Unexpected number of calls to Write")
	return 0, nil
}

func (drw *dummyResponseWriter) WriteHeader(statusCode int) {
	assert.Fail(drw.t, "Unexpected number of calls to WriteHeader")
}

type dummyHandler struct {
	t *testing.T
}

func (dh *dummyHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	assert.Fail(dh.t, "Unexpected number of calls to ServeHTTP")
}
