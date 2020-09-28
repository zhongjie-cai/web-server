package webserver

import (
	"crypto/tls"
	"errors"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNilApplication(t *testing.T) {
	// mock
	createMock(t)

	// assert
	assert.NotNil(t, nilApplication)
	assert.Zero(t, nilApplication.name)
	assert.Zero(t, nilApplication.port)
	assert.Zero(t, nilApplication.version)
	assert.NotNil(t, nilApplication.session)
	assert.Equal(t, defaultRequest, nilApplication.session.request)
	assert.Equal(t, defaultResponseWriter, nilApplication.session.responseWriter)
	assert.Empty(t, nilApplication.session.attachment)
	assert.Equal(t, customizationDefault, nilApplication.session.customization)
	assert.Equal(t, customizationDefault, nilApplication.customization)
	assert.Empty(t, nilApplication.actionFuncMap)
	assert.NotZero(t, nilApplication.shutdownSignal)

	// verify
	verifyAll(t)
}

func TestNewApplication_PortUsed(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyPort = rand.Intn(65536)
	var dummyVersion = "some version"
	var dummyCustomization Customization
	var dummyError = errors.New("some error")

	// stub
	applicationMap[dummyPort] = &application{}

	// mock
	createMock(t)

	// expect
	fmtErrorfExpected = 1
	fmtErrorf = func(format string, a ...interface{}) error {
		fmtErrorfCalled++
		assert.Equal(t, "An existing application was already registered with given port: %v", format)
		assert.Equal(t, 1, len(a))
		assert.Equal(t, dummyPort, a[0])
		return dummyError
	}

	// assert
	assert.PanicsWithValue(t, dummyError, func() {
		NewApplication(
			dummyName,
			dummyPort,
			dummyVersion,
			dummyCustomization,
		)
	})

	// verify
	verifyAll(t)
}

func TestNewApplication_NilCustomization(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyPort = rand.Intn(65536)
	var dummyVersion = "some version"
	var dummyCustomization Customization
	var dummySessionID = uuid.New()

	// mock
	createMock(t)

	// expect
	isInterfaceValueNilFuncExpected = 1
	isInterfaceValueNilFunc = func(i interface{}) bool {
		isInterfaceValueNilFuncCalled++
		assert.Equal(t, dummyCustomization, i)
		return true
	}
	uuidNewExpected = 1
	uuidNew = func() uuid.UUID {
		uuidNewCalled++
		return dummySessionID
	}

	// SUT
	var result = NewApplication(
		dummyName,
		dummyPort,
		dummyVersion,
		dummyCustomization,
	)

	// act
	var value, ok = result.(*application)

	// assert
	assert.True(t, ok)
	assert.NotNil(t, value)
	assert.Equal(t, dummyName, value.name)
	assert.Equal(t, dummyPort, value.port)
	assert.Equal(t, dummyVersion, value.version)
	assert.NotNil(t, value.session)
	assert.Equal(t, dummySessionID, value.session.id)
	assert.Equal(t, dummyName, value.session.name)
	assert.Equal(t, defaultRequest, value.session.request)
	assert.Equal(t, defaultResponseWriter, value.session.responseWriter)
	assert.Empty(t, value.session.attachment)
	assert.Equal(t, customizationDefault, value.session.customization)
	assert.Equal(t, customizationDefault, value.customization)
	assert.Empty(t, value.actionFuncMap)
	assert.NotZero(t, value.shutdownSignal)
	assert.Equal(t, 1, len(applicationMap))
	assert.Equal(t, value, applicationMap[dummyPort])

	// verify
	verifyAll(t)
}

func TestNewApplication_HasCustomization(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyPort = rand.Intn(65536)
	var dummyVersion = "some version"
	var dummyCustomization = &dummyCustomization{t: t}
	var dummySessionID = uuid.New()

	// mock
	createMock(t)

	// expect
	isInterfaceValueNilFuncExpected = 1
	isInterfaceValueNilFunc = func(i interface{}) bool {
		isInterfaceValueNilFuncCalled++
		assert.Equal(t, dummyCustomization, i)
		return false
	}
	uuidNewExpected = 1
	uuidNew = func() uuid.UUID {
		uuidNewCalled++
		return dummySessionID
	}

	// SUT
	var result = NewApplication(
		dummyName,
		dummyPort,
		dummyVersion,
		dummyCustomization,
	)

	// act
	var value, ok = result.(*application)

	// assert
	assert.True(t, ok)
	assert.NotNil(t, value)
	assert.Equal(t, dummyName, value.name)
	assert.Equal(t, dummyPort, value.port)
	assert.Equal(t, dummyVersion, value.version)
	assert.NotNil(t, value.session)
	assert.Equal(t, dummySessionID, value.session.id)
	assert.Equal(t, dummyName, value.session.name)
	assert.Equal(t, defaultRequest, value.session.request)
	assert.Equal(t, defaultResponseWriter, value.session.responseWriter)
	assert.Empty(t, value.session.attachment)
	assert.Equal(t, dummyCustomization, value.session.customization)
	assert.Equal(t, dummyCustomization, value.customization)
	assert.Empty(t, value.actionFuncMap)
	assert.NotZero(t, value.shutdownSignal)
	assert.Equal(t, 1, len(applicationMap))
	assert.Equal(t, value, applicationMap[dummyPort])

	// verify
	verifyAll(t)
}

func TestGetApplication_NotFound(t *testing.T) {
	// arrange
	var dummyPort = rand.Intn(65536)

	// stub
	applicationMap = map[int]*application{}

	// mock
	createMock(t)

	// expect

	// SUT + act
	var result = getApplication(
		dummyPort,
	)

	// assert
	assert.Equal(t, nilApplication, result)

	// verify
	verifyAll(t)
}

func TestGetApplication_NilFound(t *testing.T) {
	// arrange
	var dummyPort = rand.Intn(65536)
	var dummyApplication *application

	// stub
	applicationMap = map[int]*application{
		dummyPort: dummyApplication,
	}

	// mock
	createMock(t)

	// expect
	isInterfaceValueNilFuncExpected = 1
	isInterfaceValueNilFunc = func(i interface{}) bool {
		isInterfaceValueNilFuncCalled++
		assert.Equal(t, dummyApplication, i)
		return true
	}

	// SUT + act
	var result = getApplication(
		dummyPort,
	)

	// assert
	assert.Equal(t, nilApplication, result)

	// verify
	verifyAll(t)
}

func TestGetApplication_ValidFound(t *testing.T) {
	// arrange
	var dummyPort = rand.Intn(65536)
	var dummyApplication = &application{
		name: "some name",
	}

	// stub
	applicationMap = map[int]*application{
		dummyPort: dummyApplication,
	}

	// mock
	createMock(t)

	// expect
	isInterfaceValueNilFuncExpected = 1
	isInterfaceValueNilFunc = func(i interface{}) bool {
		isInterfaceValueNilFuncCalled++
		assert.Equal(t, dummyApplication, i)
		return false
	}

	// SUT + act
	var result = getApplication(
		dummyPort,
	)

	// assert
	assert.Equal(t, dummyApplication, result)

	// verify
	verifyAll(t)
}

func TestApplication_Start(t *testing.T) {
	// arrange
	var dummyApplication = &application{
		name: "some name",
	}

	// mock
	createMock(t)

	// expect
	startApplicationFuncExpected = 1
	startApplicationFunc = func(app *application) {
		startApplicationFuncCalled++
		assert.Equal(t, dummyApplication, app)
	}

	// SUT + act
	dummyApplication.Start()

	// verify
	verifyAll(t)
}

func TestApplication_StartAsync_NilWaitGroup(t *testing.T) {
	// arrange
	var dummyApplication = &application{
		name: "some name",
	}

	// mock
	createMock(t)

	// expect
	startApplicationFuncExpected = 1
	startApplicationFunc = func(app *application) {
		startApplicationFuncCalled++
		assert.Equal(t, dummyApplication, app)
	}

	// SUT + act
	var result = dummyApplication.StartAsync(nil)

	// assert
	assert.NotNil(t, result)

	// wait
	result.Wait()

	// verify
	verifyAll(t)
}

func TestApplication_StartAsync_GivenWaitGroup(t *testing.T) {
	// arrange
	var dummyApplication = &application{
		name: "some name",
	}
	var dummyWaitGroup = &sync.WaitGroup{}

	// mock
	createMock(t)

	// expect
	startApplicationFuncExpected = 1
	startApplicationFunc = func(app *application) {
		startApplicationFuncCalled++
		assert.Equal(t, dummyApplication, app)
	}

	// SUT + act
	var result = dummyApplication.StartAsync(dummyWaitGroup)

	// assert
	assert.NotNil(t, result)
	assert.Equal(t, dummyWaitGroup, result)

	// wait
	result.Wait()

	// verify
	verifyAll(t)
}

func TestApplication_Session(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyApplication = &application{
		name:    "some name",
		session: dummySession,
	}

	// mock
	createMock(t)

	// SUT + act
	var session = dummyApplication.Session()

	// assert
	assert.Equal(t, dummySession, session)

	// verify
	verifyAll(t)
}

func TestApplication_Stop_NotStarted(t *testing.T) {
	// arrange
	var dummyApplication = &application{
		name:    "some name",
		started: false,
	}

	// mock
	createMock(t)

	// SUT + act
	dummyApplication.Stop()

	// verify
	verifyAll(t)
}

func TestApplication_Stop_HasStarted(t *testing.T) {
	// arrange
	var dummyShutdownSignal = make(chan os.Signal, 1)
	var dummyApplication = &application{
		name:           "some name",
		shutdownSignal: dummyShutdownSignal,
		started:        true,
	}

	// mock
	createMock(t)

	// expect
	haltServerFuncExpected = 1
	haltServerFunc = func(shutdownSignal chan os.Signal) {
		haltServerFuncCalled++
		assert.Equal(t, dummyShutdownSignal, shutdownSignal)
	}

	// SUT + act
	dummyApplication.Stop()

	// verify
	verifyAll(t)
}

func TestStartApplication_AlreadyStarted(t *testing.T) {
	// arrange
	var dummyApplication = &application{
		name:    "some name",
		started: true,
	}

	// mock
	createMock(t)

	// SUT + act
	startApplication(dummyApplication)

	// verify
	verifyAll(t)
}

func TestStartApplication_PreBootstrapingFailure(t *testing.T) {
	// arrange
	var dummyApplication = &application{
		name: "some name",
	}

	// mock
	createMock(t)

	// expect
	preBootstrapingFuncExpected = 1
	preBootstrapingFunc = func(app *application) bool {
		preBootstrapingFuncCalled++
		assert.Equal(t, dummyApplication, app)
		return false
	}

	// SUT + act
	startApplication(dummyApplication)

	// verify
	verifyAll(t)
}

func TestStartApplication_PostBootstrapingFailure(t *testing.T) {
	// arrange
	var dummyApplication = &application{
		name: "some name",
	}

	// mock
	createMock(t)

	// expect
	preBootstrapingFuncExpected = 1
	preBootstrapingFunc = func(app *application) bool {
		preBootstrapingFuncCalled++
		assert.Equal(t, dummyApplication, app)
		return true
	}
	bootstrapFuncExpected = 1
	bootstrapFunc = func(app *application) {
		bootstrapFuncCalled++
		assert.Equal(t, dummyApplication, app)
	}
	postBootstrapingFuncExpected = 1
	postBootstrapingFunc = func(app *application) bool {
		postBootstrapingFuncCalled++
		assert.Equal(t, dummyApplication, app)
		return false
	}

	// SUT + act
	startApplication(dummyApplication)

	// verify
	verifyAll(t)
}

func TestStartApplication_HappyPath(t *testing.T) {
	// arrange
	var dummyApplication = &application{
		name: "some name",
	}

	// mock
	createMock(t)

	// expect
	preBootstrapingFuncExpected = 1
	preBootstrapingFunc = func(app *application) bool {
		preBootstrapingFuncCalled++
		assert.Equal(t, dummyApplication, app)
		return true
	}
	bootstrapFuncExpected = 1
	bootstrapFunc = func(app *application) {
		bootstrapFuncCalled++
		assert.Equal(t, dummyApplication, app)
	}
	postBootstrapingFuncExpected = 1
	postBootstrapingFunc = func(app *application) bool {
		postBootstrapingFuncCalled++
		assert.Equal(t, dummyApplication, app)
		return true
	}
	beginApplicationFuncExpected = 1
	beginApplicationFunc = func(app *application) {
		beginApplicationFuncCalled++
		assert.Equal(t, dummyApplication, app)
	}
	endApplicationFuncExpected = 1
	endApplicationFunc = func(app *application) {
		endApplicationFuncCalled++
		assert.Equal(t, dummyApplication, app)
	}

	// SUT + act
	startApplication(dummyApplication)

	// verify
	verifyAll(t)
}

type dummyCustomizationPreBootstrapping struct {
	dummyCustomization
	preBootstrap func() error
}

func (customization *dummyCustomizationPreBootstrapping) PreBootstrap() error {
	if customization.preBootstrap != nil {
		return customization.preBootstrap()
	}
	assert.Fail(customization.t, "Unexpected call to PreBootstrap")
	return nil
}

func TestPreBootstraping_Error(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCustomization = &dummyCustomizationPreBootstrapping{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummyApplication = &application{
		session:       dummySession,
		customization: dummyCustomization,
	}
	var customizationPreBootstrapExpected int
	var customizationPreBootstrapCalled int
	var dummyError = errors.New("some error")
	var dummyMessageFormat = "Failed to execute customization.PreBootstrap. Error: %+v"

	// mock
	createMock(t)

	// expect
	customizationPreBootstrapExpected = 1
	dummyCustomization.preBootstrap = func() error {
		customizationPreBootstrapCalled++
		return dummyError
	}
	logAppRootFuncExpected = 1
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "preBootstraping", subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyError, parameters[0])
	}

	// SUT + act
	var result = preBootstraping(
		dummyApplication,
	)

	// assert
	assert.False(t, result)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationPreBootstrapExpected, customizationPreBootstrapCalled, "Unexpected number of calls to method customization.PreBootstrap")
}

func TestPreBootstraping_Success(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCustomization = &dummyCustomizationPreBootstrapping{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummyApplication = &application{
		session:       dummySession,
		customization: dummyCustomization,
	}
	var customizationPreBootstrapExpected int
	var customizationPreBootstrapCalled int
	var dummyMessageFormat = "customization.PreBootstrap executed successfully"

	// mock
	createMock(t)

	// expect
	customizationPreBootstrapExpected = 1
	dummyCustomization.preBootstrap = func() error {
		customizationPreBootstrapCalled++
		return nil
	}
	logAppRootFuncExpected = 1
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "preBootstraping", subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Empty(t, parameters)
	}

	// SUT + act
	var result = preBootstraping(
		dummyApplication,
	)

	// assert
	assert.True(t, result)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationPreBootstrapExpected, customizationPreBootstrapCalled, "Unexpected number of calls to method customization.PreBootstrap")
}

type dummyCustomizationBootstrap struct {
	dummyCustomization
	clientCert                 func() *tls.Certificate
	defaultTimeout             func() time.Duration
	skipServerCertVerification func() bool
	roundTripper               func(http.RoundTripper) http.RoundTripper
}

func (customization *dummyCustomizationBootstrap) ClientCert() *tls.Certificate {
	if customization.clientCert != nil {
		return customization.clientCert()
	}
	assert.Fail(customization.t, "Unexpected call to ClientCert")
	return nil
}

func (customization *dummyCustomizationBootstrap) DefaultTimeout() time.Duration {
	if customization.defaultTimeout != nil {
		return customization.defaultTimeout()
	}
	assert.Fail(customization.t, "Unexpected call to DefaultTimeout")
	return 0
}

func (customization *dummyCustomizationBootstrap) SkipServerCertVerification() bool {
	if customization.skipServerCertVerification != nil {
		return customization.skipServerCertVerification()
	}
	assert.Fail(customization.t, "Unexpected call to SkipServerCertVerification")
	return false
}

func (customization *dummyCustomizationBootstrap) RoundTripper(originalTransport http.RoundTripper) http.RoundTripper {
	if customization.roundTripper != nil {
		return customization.roundTripper(originalTransport)
	}
	assert.Fail(customization.t, "Unexpected call to RoundTripper")
	return nil
}

func TestBootstrap_HappyPath(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCustomization = &dummyCustomizationBootstrap{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummyApplication = &application{
		session:       dummySession,
		customization: dummyCustomization,
	}
	var customizationDefaultTimeoutExpected int
	var customizationDefaultTimeoutCalled int
	var customizationSkipServerCertVerificationExpected int
	var customizationSkipServerCertVerificationCalled int
	var customizationClientCertExpected int
	var customizationClientCertCalled int
	var customizationRoundTripperExpected int
	var customizationRoundTripperCalled int
	var dummyWebcallTimeout = time.Duration(rand.Intn(100))
	var dummySkipCertVerification = rand.Intn(100) > 50
	var dummyClientCertificate = &tls.Certificate{Certificate: [][]byte{{0}}}
	var dummyOriginalTransport = &dummyTransport{t: t}
	var dummyMessageFormat = "Application bootstrapped successfully"

	// mock
	createMock(t)

	// expect
	initializeHTTPClientsFuncExpected = 1
	initializeHTTPClientsFunc = func(webcallTimeout time.Duration, skipServerCertVerification bool, clientCertificate *tls.Certificate, roundTripperWrapper func(originalTransport http.RoundTripper) http.RoundTripper) {
		initializeHTTPClientsFuncCalled++
		assert.Equal(t, dummyWebcallTimeout, webcallTimeout)
		assert.Equal(t, dummySkipCertVerification, skipServerCertVerification)
		assert.Equal(t, dummyClientCertificate, clientCertificate)
		roundTripperWrapper(dummyOriginalTransport)
	}
	customizationDefaultTimeoutExpected = 1
	dummyCustomization.defaultTimeout = func() time.Duration {
		customizationDefaultTimeoutCalled++
		return dummyWebcallTimeout
	}
	customizationSkipServerCertVerificationExpected = 1
	dummyCustomization.skipServerCertVerification = func() bool {
		customizationSkipServerCertVerificationCalled++
		return dummySkipCertVerification
	}
	customizationClientCertExpected = 1
	dummyCustomization.clientCert = func() *tls.Certificate {
		customizationClientCertCalled++
		return dummyClientCertificate
	}
	customizationRoundTripperExpected = 1
	dummyCustomization.roundTripper = func(originalTransport http.RoundTripper) http.RoundTripper {
		customizationRoundTripperCalled++
		assert.Equal(t, dummyOriginalTransport, originalTransport)
		return originalTransport
	}
	logAppRootFuncExpected = 1
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "bootstrap", subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Empty(t, parameters)
	}

	// SUT + act
	bootstrap(
		dummyApplication,
	)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationDefaultTimeoutExpected, customizationDefaultTimeoutCalled, "Unexpected number of calls to method customization.DefaultTimeout")
	assert.Equal(t, customizationSkipServerCertVerificationExpected, customizationSkipServerCertVerificationCalled, "Unexpected number of calls to method customization.SkipServerCertVerification")
	assert.Equal(t, customizationClientCertExpected, customizationClientCertCalled, "Unexpected number of calls to method customization.ClientCert")
	assert.Equal(t, customizationRoundTripperExpected, customizationRoundTripperCalled, "Unexpected number of calls to method customization.RoundTripper")
}

type dummyCustomizationPostBootstrapping struct {
	dummyCustomization
	postBootstrap func() error
}

func (customization *dummyCustomizationPostBootstrapping) PostBootstrap() error {
	if customization.postBootstrap != nil {
		return customization.postBootstrap()
	}
	assert.Fail(customization.t, "Unexpected call to PostBootstrap")
	return nil
}

func TestPostBootstraping_Error(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCustomization = &dummyCustomizationPostBootstrapping{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummyApplication = &application{
		session:       dummySession,
		customization: dummyCustomization,
	}
	var customizationPostBootstrapExpected int
	var customizationPostBootstrapCalled int
	var dummyError = errors.New("some error")
	var dummyMessageFormat = "Failed to execute customization.PostBootstrap. Error: %+v"

	// mock
	createMock(t)

	// expect
	customizationPostBootstrapExpected = 1
	dummyCustomization.postBootstrap = func() error {
		customizationPostBootstrapCalled++
		return dummyError
	}
	logAppRootFuncExpected = 1
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "postBootstraping", subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyError, parameters[0])
	}

	// SUT + act
	var result = postBootstraping(
		dummyApplication,
	)

	// assert
	assert.False(t, result)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationPostBootstrapExpected, customizationPostBootstrapCalled, "Unexpected number of calls to method customization.PostBootstrap")
}

func TestPostBootstraping_Success(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCustomization = &dummyCustomizationPostBootstrapping{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummyApplication = &application{
		session:       dummySession,
		customization: dummyCustomization,
	}
	var customizationPostBootstrapExpected int
	var customizationPostBootstrapCalled int
	var dummyMessageFormat = "customization.PostBootstrap executed successfully"

	// mock
	createMock(t)

	// expect
	customizationPostBootstrapExpected = 1
	dummyCustomization.postBootstrap = func() error {
		customizationPostBootstrapCalled++
		return nil
	}
	logAppRootFuncExpected = 1
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "postBootstraping", subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Empty(t, parameters)
	}

	// SUT + act
	var result = postBootstraping(
		dummyApplication,
	)

	// assert
	assert.True(t, result)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationPostBootstrapExpected, customizationPostBootstrapCalled, "Unexpected number of calls to method customization.PostBootstrap")
}

func TestBeginApplication_HostError(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyVersion = "some version"
	var dummyPort = rand.Intn(65536)
	var dummySession = &session{id: uuid.New()}
	var dummyCustomization = &dummyCustomization{t: t}
	var dummyShutdownSignal = make(chan os.Signal, 1)
	var dummyStarted = rand.Intn(100) > 50
	var dummyApplication = &application{
		name:           dummyName,
		version:        dummyVersion,
		port:           dummyPort,
		session:        dummySession,
		customization:  dummyCustomization,
		shutdownSignal: dummyShutdownSignal,
		started:        dummyStarted,
	}
	var dummyError = errors.New("some error")

	// mock
	createMock(t)

	// expect
	logAppRootFuncExpected = 2
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "beginApplication", subcategory)
		if logAppRootFuncCalled == 1 {
			assert.Equal(t, "Trying to start server [%v] (v-%v)", messageFormat)
			assert.Equal(t, 2, len(parameters))
			assert.Equal(t, dummyName, parameters[0])
			assert.Equal(t, dummyVersion, parameters[1])
		} else if logAppRootFuncCalled == 2 {
			assert.Equal(t, "Failed to host server. Error: %+v", messageFormat)
			assert.Equal(t, 1, len(parameters))
			assert.Equal(t, dummyError, parameters[0])
		}
	}
	hostServerFuncExpected = 1
	hostServerFunc = func(port int, session *session, shutdownSignal chan os.Signal, started *bool) error {
		hostServerFuncCalled++
		assert.Equal(t, dummyPort, port)
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyShutdownSignal, shutdownSignal)
		assert.Equal(t, &dummyStarted, started)
		return dummyError
	}

	// SUT + act
	beginApplication(
		dummyApplication,
	)

	// verify
	verifyAll(t)
}

func TestBeginApplication_HostSuccess(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyVersion = "some version"
	var dummyPort = rand.Intn(65536)
	var dummySession = &session{id: uuid.New()}
	var dummyCustomization = &dummyCustomization{t: t}
	var dummyShutdownSignal = make(chan os.Signal, 1)
	var dummyStarted = rand.Intn(100) > 50
	var dummyApplication = &application{
		name:           dummyName,
		version:        dummyVersion,
		port:           dummyPort,
		session:        dummySession,
		customization:  dummyCustomization,
		shutdownSignal: dummyShutdownSignal,
		started:        dummyStarted,
	}

	// mock
	createMock(t)

	// expect
	logAppRootFuncExpected = 2
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "beginApplication", subcategory)
		if logAppRootFuncCalled == 1 {
			assert.Equal(t, "Trying to start server [%v] (v-%v)", messageFormat)
			assert.Equal(t, 2, len(parameters))
			assert.Equal(t, dummyName, parameters[0])
			assert.Equal(t, dummyVersion, parameters[1])
		} else if logAppRootFuncCalled == 2 {
			assert.Equal(t, "Server hosting terminated", messageFormat)
			assert.Empty(t, parameters)
		}
	}
	hostServerFuncExpected = 1
	hostServerFunc = func(port int, session *session, shutdownSignal chan os.Signal, started *bool) error {
		hostServerFuncCalled++
		assert.Equal(t, dummyPort, port)
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyShutdownSignal, shutdownSignal)
		assert.Equal(t, &dummyStarted, started)
		return nil
	}

	// SUT + act
	beginApplication(
		dummyApplication,
	)

	// verify
	verifyAll(t)
}

type dummyCustomizationEndApplication struct {
	dummyCustomization
	appClosing func() error
}

func (customization *dummyCustomizationEndApplication) AppClosing() error {
	if customization.appClosing != nil {
		return customization.appClosing()
	}
	assert.Fail(customization.t, "Unexpected call to AppClosing")
	return nil
}

func TestEndApplication_Error(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCustomization = &dummyCustomizationEndApplication{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummyApplication = &application{
		session:       dummySession,
		customization: dummyCustomization,
	}
	var customizationAppClosingExpected int
	var customizationAppClosingCalled int
	var dummyError = errors.New("some error")
	var dummyMessageFormat = "Failed to execute customization.AppClosing. Error: %+v"

	// mock
	createMock(t)

	// expect
	customizationAppClosingExpected = 1
	dummyCustomization.appClosing = func() error {
		customizationAppClosingCalled++
		return dummyError
	}
	logAppRootFuncExpected = 1
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "endApplication", subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyError, parameters[0])
	}

	// SUT + act
	endApplication(
		dummyApplication,
	)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationAppClosingExpected, customizationAppClosingCalled, "Unexpected number of calls to method customization.AppClosing")
}

func TestEndApplication_Success(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCustomization = &dummyCustomizationEndApplication{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummyApplication = &application{
		session:       dummySession,
		customization: dummyCustomization,
	}
	var customizationAppClosingExpected int
	var customizationAppClosingCalled int
	var dummyMessageFormat = "customization.AppClosing executed successfully"

	// mock
	createMock(t)

	// expect
	customizationAppClosingExpected = 1
	dummyCustomization.appClosing = func() error {
		customizationAppClosingCalled++
		return nil
	}
	logAppRootFuncExpected = 1
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "endApplication", subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Empty(t, parameters)
	}

	// SUT + act
	endApplication(
		dummyApplication,
	)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationAppClosingExpected, customizationAppClosingCalled, "Unexpected number of calls to method customization.AppClosing")
}
