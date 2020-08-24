package webserver

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var (
	isInterfaceValueNilFuncExpected   int
	isInterfaceValueNilFuncCalled     int
	uuidNewExpected                   int
	uuidNewCalled                     int
	startApplicationFuncExpected      int
	startApplicationFuncCalled        int
	haltServerFuncExpected            int
	haltServerFuncCalled              int
	preBootstrapingFuncExpected       int
	preBootstrapingFuncCalled         int
	bootstrapFuncExpected             int
	bootstrapFuncCalled               int
	postBootstrapingFuncExpected      int
	postBootstrapingFuncCalled        int
	endApplicationFuncExpected        int
	endApplicationFuncCalled          int
	beginApplicationFuncExpected      int
	beginApplicationFuncCalled        int
	logAppRootFuncExpected            int
	logAppRootFuncCalled              int
	initializeHTTPClientsFuncExpected int
	initializeHTTPClientsFuncCalled   int
	hostServerFuncExpected            int
	hostServerFuncCalled              int
	fmtPrintfExpected                 int
	fmtPrintfCalled                   int
	fmtSprintfExpected                int
	fmtSprintfCalled                  int
	marshalIgnoreErrorFuncExpected    int
	marshalIgnoreErrorFuncCalled      int
)

func createMock(t *testing.T) {
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
	hostServerFunc = func(port int, session *session, customization Customization, shutdownSignal chan os.Signal) error {
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
}

func verifyAll(t *testing.T) {
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

	applicationLock = sync.RWMutex{}
	applicationMap = map[int]*application{}
}

func functionPointerEquals(t *testing.T, expectFunc interface{}, actualFunc interface{}) {
	var expectValue = fmt.Sprintf("%v", reflect.ValueOf(expectFunc))
	var actualValue = fmt.Sprintf("%v", reflect.ValueOf(actualFunc))
	assert.Equal(t, expectValue, actualValue)
}

// mock structs
type dummyApplication struct {
	t *testing.T
}

func (application *dummyApplication) Start() {
	assert.Fail(application.t, "Unexpected call to Start")
}

func (application *dummyApplication) StartAsync(*sync.WaitGroup) *sync.WaitGroup {
	assert.Fail(application.t, "Unexpected call to StartAsync")
	return nil
}

func (application *dummyApplication) Stop() {
	assert.Fail(application.t, "Unexpected call to Stop")
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

// LogMethodEnter sends a logging entry of MethodEnter log type for the given session associated to the session ID
func (session *dummySession) LogMethodEnter() {
	assert.Fail(session.t, "Unexpected call to LogMethodEnter")
}

// LogMethodParameter sends a logging entry of MethodParameter log type for the given session associated to the session ID
func (session *dummySession) LogMethodParameter(parameters ...interface{}) {
	assert.Fail(session.t, "Unexpected call to LogMethodParameter")
}

// LogMethodLogic sends a logging entry of MethodLogic log type for the given session associated to the session ID
func (session *dummySession) LogMethodLogic(logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{}) {
	assert.Fail(session.t, "Unexpected call to LogMethodLogic")
}

// LogMethodReturn sends a logging entry of MethodReturn log type for the given session associated to the session ID
func (session *dummySession) LogMethodReturn(returns ...interface{}) {
	assert.Fail(session.t, "Unexpected call to LogMethodReturn")
}

// LogMethodExit sends a logging entry of MethodExit log type for the given session associated to the session ID
func (session *dummySession) LogMethodExit() {
	assert.Fail(session.t, "Unexpected call to LogMethodExit")
}

// CreateWebcallRequest generates a webcall request object to the targeted external web service for the given session associated to the session ID
func (session *dummySession) CreateWebcallRequest(method string, url string, payload string, header map[string]string, sendClientCert bool) WebRequest {
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
