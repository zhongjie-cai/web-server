package webserver

import (
	"crypto/tls"
	"errors"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zhongjie-cai/gomocker"
)

func TestNewApplication_NilCustomization(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyAddress = "some address"
	var dummyVersion = "some version"
	var dummyCustomization Customization
	var dummySessionID = uuid.New()

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(isInterfaceValueNil, 1, func(i interface{}) bool {
		assert.Equal(t, dummyCustomization, i)
		return true
	})
	m.ExpectFunc(uuid.New, 1, func() uuid.UUID {
		return dummySessionID
	})

	// SUT
	var result = NewApplication(
		dummyName,
		dummyAddress,
		dummyVersion,
		dummyCustomization,
	)

	// act
	var value, ok = result.(*application)

	// assert
	assert.True(t, ok)
	assert.NotNil(t, value)
	assert.Equal(t, dummyName, value.name)
	assert.Equal(t, dummyAddress, value.address)
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
}

func TestNewApplication_HasCustomization(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyAddress = "some address"
	var dummyVersion = "some version"
	var dummyCustomization = &DefaultCustomization{}
	var dummySessionID = uuid.New()

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(isInterfaceValueNil, 1, func(i interface{}) bool {
		assert.Equal(t, dummyCustomization, i)
		return false
	})
	m.ExpectFunc(uuid.New, 1, func() uuid.UUID {
		return dummySessionID
	})

	// SUT
	var result = NewApplication(
		dummyName,
		dummyAddress,
		dummyVersion,
		dummyCustomization,
	)

	// act
	var value, ok = result.(*application)

	// assert
	assert.True(t, ok)
	assert.NotNil(t, value)
	assert.Equal(t, dummyName, value.name)
	assert.Equal(t, dummyAddress, value.address)
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
}

func TestApplication_Start(t *testing.T) {
	// arrange
	var dummyApplication = &application{
		name: "some name",
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(startApplication, 1, func(app *application) {
		assert.Equal(t, dummyApplication, app)
	})

	// SUT + act
	dummyApplication.Start()
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

	// SUT + act
	var session = dummyApplication.Session()

	// assert
	assert.Equal(t, dummySession, session)
}

func TestApplication_IsRunning(t *testing.T) {
	// arrange
	var dummyApplication = &application{
		name:    "some name",
		started: rand.Intn(100) > 50,
	}

	// SUT + act
	var result = dummyApplication.IsRunning()

	// assert
	assert.Equal(t, dummyApplication.started, result)
}

func TestApplication_Stop_NotStarted(t *testing.T) {
	// arrange
	var dummyApplication = &application{
		name:    "some name",
		started: false,
	}

	// SUT + act
	dummyApplication.Stop()
}

func TestApplication_Stop_HasStarted(t *testing.T) {
	// arrange
	var dummyShutdownSignal = make(chan os.Signal)
	var dummyApplication = &application{
		name:           "some name",
		shutdownSignal: dummyShutdownSignal,
		started:        true,
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(haltServer, 1, func(shutdownSignal chan os.Signal) {
		assert.Equal(t, dummyShutdownSignal, shutdownSignal)
	})

	// SUT + act
	dummyApplication.Stop()
}

func TestStartApplication_AlreadyStarted(t *testing.T) {
	// arrange
	var dummyApplication = &application{
		name:    "some name",
		started: true,
	}

	// SUT + act
	startApplication(dummyApplication)
}

func TestStartApplication_PreBootstrapingFailure(t *testing.T) {
	// arrange
	var dummyApplication = &application{
		name: "some name",
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(preBootstraping, 1, func(app *application) bool {
		assert.Equal(t, dummyApplication, app)
		return false
	})

	// SUT + act
	startApplication(dummyApplication)
}

func TestStartApplication_PostBootstrapingFailure(t *testing.T) {
	// arrange
	var dummyApplication = &application{
		name: "some name",
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(preBootstraping, 1, func(app *application) bool {
		assert.Equal(t, dummyApplication, app)
		return true
	})
	m.ExpectFunc(bootstrap, 1, func(app *application) {
		assert.Equal(t, dummyApplication, app)
	})
	m.ExpectFunc(postBootstraping, 1, func(app *application) bool {
		assert.Equal(t, dummyApplication, app)
		return false
	})

	// SUT + act
	startApplication(dummyApplication)
}

func TestStartApplication_HappyPath(t *testing.T) {
	// arrange
	var dummyApplication = &application{
		name: "some name",
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(preBootstraping, 1, func(app *application) bool {
		assert.Equal(t, dummyApplication, app)
		return true
	})
	m.ExpectFunc(bootstrap, 1, func(app *application) {
		assert.Equal(t, dummyApplication, app)
	})
	m.ExpectFunc(postBootstraping, 1, func(app *application) bool {
		assert.Equal(t, dummyApplication, app)
		return true
	})
	m.ExpectFunc(beginApplication, 1, func(app *application) {
		assert.Equal(t, dummyApplication, app)
	})
	m.ExpectFunc(endApplication, 1, func(app *application) {
		assert.Equal(t, dummyApplication, app)
	})

	// SUT + act
	startApplication(dummyApplication)
}

func TestPreBootstraping_Error(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCustomization = &DefaultCustomization{}
	var dummyApplication = &application{
		session:       dummySession,
		customization: dummyCustomization,
	}
	var dummyError = errors.New("some error")
	var dummyMessageFormat = "Failed to execute customization.PreBootstrap. Error: %+v"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "PreBootstrap", 1, func() error {
		return dummyError
	})
	m.ExpectFunc(logAppRoot, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "preBootstraping", subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyError, parameters[0])
	})

	// SUT + act
	var result = preBootstraping(
		dummyApplication,
	)

	// assert
	assert.False(t, result)
}

func TestPreBootstraping_Success(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCustomization = &DefaultCustomization{}
	var dummyApplication = &application{
		session:       dummySession,
		customization: dummyCustomization,
	}
	var dummyMessageFormat = "customization.PreBootstrap executed successfully"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "PreBootstrap", 1, func() error {
		return nil
	})
	m.ExpectFunc(logAppRoot, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "preBootstraping", subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Empty(t, parameters)
	})

	// SUT + act
	var result = preBootstraping(
		dummyApplication,
	)

	// assert
	assert.True(t, result)
}

func TestBootstrap_HappyPath(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCustomization = &DefaultCustomization{}
	var dummyApplication = &application{
		session:       dummySession,
		customization: dummyCustomization,
	}
	var dummyWebcallTimeout = time.Duration(rand.Intn(100))
	var dummySkipCertVerification = rand.Intn(100) > 50
	var dummyClientCertificate = &tls.Certificate{Certificate: [][]byte{{0}}}
	var dummyMessageFormat = "Application bootstrapped successfully"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(initializeHTTPClients, 1, func(webcallTimeout time.Duration, skipServerCertVerification bool, clientCertificate *tls.Certificate, roundTripperWrapper func(originalTransport http.RoundTripper) http.RoundTripper) {
		assert.Equal(t, dummyWebcallTimeout, webcallTimeout)
		assert.Equal(t, dummySkipCertVerification, skipServerCertVerification)
		assert.Equal(t, dummyClientCertificate, clientCertificate)
	})
	m.ExpectMethod(dummyCustomization, "DefaultTimeout", 1, func() time.Duration {
		return dummyWebcallTimeout
	})
	m.ExpectMethod(dummyCustomization, "SkipServerCertVerification", 1, func() bool {
		return dummySkipCertVerification
	})
	m.ExpectMethod(dummyCustomization, "ClientCert", 1, func() *tls.Certificate {
		return dummyClientCertificate
	})
	m.ExpectFunc(logAppRoot, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "bootstrap", subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Empty(t, parameters)
	})

	// SUT + act
	bootstrap(
		dummyApplication,
	)
}

func TestPostBootstraping_Error(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCustomization = &DefaultCustomization{}
	var dummyApplication = &application{
		session:       dummySession,
		customization: dummyCustomization,
	}
	var dummyError = errors.New("some error")
	var dummyMessageFormat = "Failed to execute customization.PostBootstrap. Error: %+v"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "PostBootstrap", 1, func() error {
		return dummyError
	})
	m.ExpectFunc(logAppRoot, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "postBootstraping", subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyError, parameters[0])
	})

	// SUT + act
	var result = postBootstraping(
		dummyApplication,
	)

	// assert
	assert.False(t, result)
}

func TestPostBootstraping_Success(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCustomization = &DefaultCustomization{}
	var dummyApplication = &application{
		session:       dummySession,
		customization: dummyCustomization,
	}
	var dummyMessageFormat = "customization.PostBootstrap executed successfully"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "PostBootstrap", 1, func() error {
		return nil
	})
	m.ExpectFunc(logAppRoot, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "postBootstraping", subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Empty(t, parameters)
	})

	// SUT + act
	var result = postBootstraping(
		dummyApplication,
	)

	// assert
	assert.True(t, result)
}

func TestBeginApplication_HostError(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyVersion = "some version"
	var dummyAddress = "some address"
	var dummySession = &session{id: uuid.New()}
	var dummyCustomization = &DefaultCustomization{}
	var dummyShutdownSignal = make(chan os.Signal)
	var dummyStarted = rand.Intn(100) > 50
	var dummyApplication = &application{
		name:           dummyName,
		version:        dummyVersion,
		address:        dummyAddress,
		session:        dummySession,
		customization:  dummyCustomization,
		shutdownSignal: dummyShutdownSignal,
		started:        dummyStarted,
	}
	var dummyError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(logAppRoot, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "beginApplication", subcategory)
		assert.Equal(t, "Trying to start server [%v] (v-%v)", messageFormat)
		assert.Equal(t, 2, len(parameters))
		assert.Equal(t, dummyName, parameters[0])
		assert.Equal(t, dummyVersion, parameters[1])
	}).ExpectFunc(logAppRoot, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "beginApplication", subcategory)
		assert.Equal(t, "Failed to host server. Error: %+v", messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyError, parameters[0])
	})
	m.ExpectFunc(hostServer, 1, func(app *application, session *session, shutdownSignal chan os.Signal, started *bool) error {
		assert.Equal(t, dummyApplication, app)
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyShutdownSignal, shutdownSignal)
		assert.Equal(t, &dummyStarted, started)
		return dummyError
	})

	// SUT + act
	beginApplication(
		dummyApplication,
	)
}

func TestBeginApplication_HostSuccess(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyVersion = "some version"
	var dummyAddress = "some address"
	var dummySession = &session{id: uuid.New()}
	var dummyCustomization = &DefaultCustomization{}
	var dummyShutdownSignal = make(chan os.Signal)
	var dummyStarted = rand.Intn(100) > 50
	var dummyApplication = &application{
		name:           dummyName,
		version:        dummyVersion,
		address:        dummyAddress,
		session:        dummySession,
		customization:  dummyCustomization,
		shutdownSignal: dummyShutdownSignal,
		started:        dummyStarted,
	}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(logAppRoot, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "beginApplication", subcategory)
		assert.Equal(t, "Trying to start server [%v] (v-%v)", messageFormat)
		assert.Equal(t, 2, len(parameters))
		assert.Equal(t, dummyName, parameters[0])
		assert.Equal(t, dummyVersion, parameters[1])
	}).ExpectFunc(logAppRoot, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "beginApplication", subcategory)
		assert.Equal(t, "Server hosting terminated", messageFormat)
		assert.Empty(t, parameters)
	})
	m.ExpectFunc(hostServer, 1, func(app *application, session *session, shutdownSignal chan os.Signal, started *bool) error {
		assert.Equal(t, dummyApplication, app)
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyShutdownSignal, shutdownSignal)
		assert.Equal(t, &dummyStarted, started)
		return nil
	})

	// SUT + act
	beginApplication(
		dummyApplication,
	)
}

func TestEndApplication_Error(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCustomization = &DefaultCustomization{}
	var dummyApplication = &application{
		session:       dummySession,
		customization: dummyCustomization,
	}
	var dummyError = errors.New("some error")
	var dummyMessageFormat = "Failed to execute customization.AppClosing. Error: %+v"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "AppClosing", 1, func() error {
		return dummyError
	})
	m.ExpectFunc(logAppRoot, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "endApplication", subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Equal(t, 1, len(parameters))
		assert.Equal(t, dummyError, parameters[0])
	})

	// SUT + act
	endApplication(
		dummyApplication,
	)
}

func TestEndApplication_Success(t *testing.T) {
	// arrange
	var dummySession = &session{
		id: uuid.New(),
	}
	var dummyCustomization = &DefaultCustomization{}
	var dummyApplication = &application{
		session:       dummySession,
		customization: dummyCustomization,
	}
	var dummyMessageFormat = "customization.AppClosing executed successfully"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "AppClosing", 1, func() error {
		return nil
	})
	m.ExpectFunc(logAppRoot, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "application", category)
		assert.Equal(t, "endApplication", subcategory)
		assert.Equal(t, dummyMessageFormat, messageFormat)
		assert.Empty(t, parameters)
	})

	// SUT + act
	endApplication(
		dummyApplication,
	)
}
