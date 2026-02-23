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
	"github.com/zhongjie-cai/gomocker/v2"
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
	m.Mock(isInterfaceValueNil).Expects(dummyCustomization).Returns(true).Once()
	m.Mock(uuid.New).Expects().Returns(dummySessionID).Once()

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
	m.Mock(isInterfaceValueNil).Expects(dummyCustomization).Returns(false).Once()
	m.Mock(uuid.New).Expects().Returns(dummySessionID).Once()

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
	m.Mock(startApplication).Expects(dummyApplication).Returns().Once()

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
	m.Mock(haltServer).Expects(dummyShutdownSignal).Returns().Once()

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
	m.Mock(preBootstraping).Expects(dummyApplication).Returns(false).Once()

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
	m.Mock(preBootstraping).Expects(dummyApplication).Returns(true).Once()
	m.Mock(bootstrap).Expects(dummyApplication).Returns().Once()
	m.Mock(postBootstraping).Expects(dummyApplication).Returns(false).Once()

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
	m.Mock(preBootstraping).Expects(dummyApplication).Returns(true).Once()
	m.Mock(bootstrap).Expects(dummyApplication).Returns().Once()
	m.Mock(postBootstraping).Expects(dummyApplication).Returns(true).Once()
	m.Mock(beginApplication).Expects(dummyApplication).Returns().Once()
	m.Mock(endApplication).Expects(dummyApplication).Returns().Once()

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
	m.Mock((*DefaultCustomization).PreBootstrap).Expects(dummyCustomization).Returns(dummyError).Once()
	m.Mock(logAppRoot).Expects(dummySession, "application", "preBootstraping", dummyMessageFormat, dummyError).Returns().Once()

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
	m.Mock((*DefaultCustomization).PreBootstrap).Expects(dummyCustomization).Returns(nil).Once()
	m.Mock(logAppRoot).Expects(dummySession, "application", "preBootstraping", dummyMessageFormat).Returns().Once()

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
	var dummyHttpClient = &http.Client{Timeout: time.Duration(rand.Intn(100))}
	var dummyWebcallTimeout = time.Duration(rand.Intn(100))
	var dummySkipCertVerification = rand.Intn(100) > 50
	var dummyClientCertificate = &tls.Certificate{Certificate: [][]byte{{0}}}
	var dummyMessageFormat = "Application bootstrapped successfully"

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(initializeHTTPClients).Expects(dummyHttpClient, dummyWebcallTimeout, dummySkipCertVerification, dummyClientCertificate, gomocker.Anything()).Returns().Once()
	m.Mock((*DefaultCustomization).HttpClient).Expects(dummyCustomization).Returns(dummyHttpClient).Once()
	m.Mock((*DefaultCustomization).DefaultTimeout).Expects(dummyCustomization).Returns(dummyWebcallTimeout).Once()
	m.Mock((*DefaultCustomization).SkipServerCertVerification).Expects(dummyCustomization).Returns(dummySkipCertVerification).Once()
	m.Mock((*DefaultCustomization).ClientCert).Expects(dummyCustomization).Returns(dummyClientCertificate).Once()
	m.Mock(logAppRoot).Expects(dummySession, "application", "bootstrap", dummyMessageFormat).Returns().Once()

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
	m.Mock((*DefaultCustomization).PostBootstrap).Expects(dummyCustomization).Returns(dummyError).Once()
	m.Mock(logAppRoot).Expects(dummySession, "application", "postBootstraping", dummyMessageFormat, dummyError).Returns().Once()

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
	m.Mock((*DefaultCustomization).PostBootstrap).Expects(dummyCustomization).Returns(nil).Once()
	m.Mock(logAppRoot).Expects(dummySession, "application", "postBootstraping", dummyMessageFormat).Returns().Once()

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
	m.Mock(logAppRoot).Expects(dummySession, "application", "beginApplication", "Trying to start server [%v] (v-%v)", dummyName, dummyVersion).Returns().Once()
	m.Mock(hostServer).Expects(dummyApplication, dummySession, dummyShutdownSignal, &dummyStarted).Returns(dummyError).Once()
	m.Mock(logAppRoot).Expects(dummySession, "application", "beginApplication", "Failed to host server. Error: %+v", dummyError).Returns().Once()

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
	m.Mock(logAppRoot).Expects(dummySession, "application", "beginApplication", "Trying to start server [%v] (v-%v)", dummyName, dummyVersion).Returns().Once()
	m.Mock(hostServer).Expects(dummyApplication, dummySession, dummyShutdownSignal, &dummyStarted).Returns(nil).Once()
	m.Mock(logAppRoot).Expects(dummySession, "application", "beginApplication", "Server hosting terminated").Returns().Once()

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
	m.Mock((*DefaultCustomization).AppClosing).Expects(dummyCustomization).Returns(dummyError).Once()
	m.Mock(logAppRoot).Expects(dummySession, "application", "endApplication", dummyMessageFormat, dummyError).Returns().Once()

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
	m.Mock((*DefaultCustomization).AppClosing).Expects(dummyCustomization).Returns(nil).Once()
	m.Mock(logAppRoot).Expects(dummySession, "application", "endApplication", dummyMessageFormat).Returns().Once()

	// SUT + act
	endApplication(
		dummyApplication,
	)
}
