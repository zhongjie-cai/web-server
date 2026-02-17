package webserver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zhongjie-cai/gomocker/v2"
)

func TestHostServer_ErrorRegisterRoutes(t *testing.T) {
	// arrange
	var dummyAddress = "some address"
	var dummyApplication = &application{
		address: dummyAddress,
	}
	var dummySession = &session{id: uuid.New()}
	var dummyShutdownSignal = make(chan os.Signal)
	var dummyStarted = rand.Intn(100) > 50
	type router struct {
		chi.Router
	}
	var dummyRouter = &router{}
	var dummyError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(instantiateRouter).Expects(dummyApplication, dummySession).Returns(dummyRouter, dummyError).Once()

	// SUT + act
	var err = hostServer(
		dummyApplication,
		dummySession,
		dummyShutdownSignal,
		&dummyStarted,
	)

	// assert
	assert.Equal(t, dummyError, err)
}

func TestHostServer_RunServerFailure(t *testing.T) {
	// arrange
	var dummyAddress = "some address"
	var dummyApplication = &application{
		address: dummyAddress,
	}
	var dummySession = &session{id: uuid.New()}
	var dummyShutdownSignal = make(chan os.Signal)
	var dummyStarted = rand.Intn(100) > 50
	type router struct {
		chi.Router
	}
	var dummyRouter = &router{}
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(instantiateRouter).Expects(dummyApplication, dummySession).Returns(dummyRouter, nil).Once()
	m.Mock(logAppRoot).Expects(dummySession, "server", "hostServer", "Targeting address [%v]", dummyAddress).Returns().Once()
	m.Mock(logAppRoot).Expects(dummySession, "server", "hostServer", "Server terminated").Returns().Once()
	m.Mock(runServer).Expects(dummyAddress, dummySession, dummyRouter, dummyShutdownSignal, &dummyStarted).Returns(false).Once()
	m.Mock(newAppError).Expects(errorCodeGeneralFailure, errorMessageHostServer).Returns(dummyAppError).Once()

	// SUT + act
	var err = hostServer(
		dummyApplication,
		dummySession,
		dummyShutdownSignal,
		&dummyStarted,
	)

	// assert
	assert.Equal(t, dummyAppError, err)
}

func TestHostServer_RunServerSuccess(t *testing.T) {
	// arrange
	var dummyAddress = "some address"
	var dummyApplication = &application{
		address: dummyAddress,
	}
	var dummySession = &session{id: uuid.New()}
	var dummyShutdownSignal = make(chan os.Signal)
	var dummyStarted = rand.Intn(100) > 50
	type router struct {
		chi.Router
	}
	var dummyRouter = &router{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(instantiateRouter).Expects(dummyApplication, dummySession).Returns(dummyRouter, nil).Once()
	m.Mock(logAppRoot).Expects(dummySession, "server", "hostServer", "Targeting address [%v]", dummyAddress).Returns().Once()
	m.Mock(logAppRoot).Expects(dummySession, "server", "hostServer", "Server closed").Returns().Once()
	m.Mock(runServer).Expects(dummyAddress, dummySession, dummyRouter, dummyShutdownSignal, &dummyStarted).Returns(true).Once()

	// SUT + act
	var err = hostServer(
		dummyApplication,
		dummySession,
		dummyShutdownSignal,
		&dummyStarted,
	)

	// assert
	assert.NoError(t, err)
}

func TestCreateServer_NoServerCert(t *testing.T) {
	// arrange
	var dummyAddress = "some address"
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		id:            uuid.New(),
		customization: dummyCustomization,
	}
	type router struct {
		chi.Router
	}
	var dummyRouter = &router{}
	var dummyServerCert *tls.Certificate

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock((*DefaultCustomization).ServerCert).Expects(dummyCustomization).Returns(dummyServerCert).Once()
	m.Mock((*DefaultCustomization).WrapHandler).Expects(dummyCustomization, dummyRouter).Returns(dummyRouter).Once()

	// SUT + act
	var server, https = createServer(
		dummyAddress,
		dummySession,
		dummyRouter,
	)

	// assert
	assert.NotNil(t, server)
	assert.Equal(t, dummyAddress, server.Addr)
	assert.NotNil(t, server.TLSConfig)
	assert.Empty(t, server.TLSConfig.Certificates)
	assert.Equal(t, tls.NoClientCert, server.TLSConfig.ClientAuth)
	assert.Nil(t, server.TLSConfig.ClientCAs)
	assert.Empty(t, server.TLSConfig.CipherSuites)
	assert.Equal(t, uint16(tls.VersionTLS12), server.TLSConfig.MinVersion)
	assert.Zero(t, server.WriteTimeout)
	assert.Zero(t, server.ReadTimeout)
	assert.Zero(t, server.IdleTimeout)
	assert.False(t, https)
}

func TestCreateServer_WithServerCert_NoCaCertPool(t *testing.T) {
	// arrange
	var dummyAddress = "some address"
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		id:            uuid.New(),
		customization: dummyCustomization,
	}
	type router struct {
		chi.Router
	}
	var dummyRouter = &router{}
	var dummyServerCert = &tls.Certificate{}
	var dummyCaCertPool *x509.CertPool

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock((*DefaultCustomization).ServerCert).Expects(dummyCustomization).Returns(dummyServerCert).Once()
	m.Mock((*DefaultCustomization).CaCertPool).Expects(dummyCustomization).Returns(dummyCaCertPool).Once()
	m.Mock((*DefaultCustomization).WrapHandler).Expects(dummyCustomization, dummyRouter).Returns(dummyRouter).Once()

	// SUT + act
	var server, https = createServer(
		dummyAddress,
		dummySession,
		dummyRouter,
	)

	// assert
	assert.NotNil(t, server)
	assert.Equal(t, dummyAddress, server.Addr)
	assert.NotNil(t, server.TLSConfig)
	assert.Equal(t, 1, len(server.TLSConfig.Certificates))
	assert.Equal(t, *dummyServerCert, server.TLSConfig.Certificates[0])
	assert.Equal(t, tls.RequireAnyClientCert, server.TLSConfig.ClientAuth)
	assert.Nil(t, server.TLSConfig.ClientCAs)
	assert.Empty(t, server.TLSConfig.CipherSuites)
	assert.Equal(t, uint16(tls.VersionTLS12), server.TLSConfig.MinVersion)
	assert.Zero(t, server.WriteTimeout)
	assert.Zero(t, server.ReadTimeout)
	assert.Zero(t, server.IdleTimeout)
	assert.True(t, https)
}

func TestCreateServer_WithServerCert_WithCaCertPool(t *testing.T) {
	// arrange
	var dummyAddress = "some address"
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		id:            uuid.New(),
		customization: dummyCustomization,
	}
	type router struct {
		chi.Router
	}
	var dummyRouter = &router{}
	var dummyServerCert = &tls.Certificate{}
	var dummyCaCertPool = &x509.CertPool{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock((*DefaultCustomization).ServerCert).Expects(dummyCustomization).Returns(dummyServerCert).Once()
	m.Mock((*DefaultCustomization).CaCertPool).Expects(dummyCustomization).Returns(dummyCaCertPool).Once()
	m.Mock((*DefaultCustomization).WrapHandler).Expects(dummyCustomization, dummyRouter).Returns(dummyRouter).Once()

	// SUT + act
	var server, https = createServer(
		dummyAddress,
		dummySession,
		dummyRouter,
	)

	// assert
	assert.NotNil(t, server)
	assert.Equal(t, dummyAddress, server.Addr)
	assert.NotNil(t, server.TLSConfig)
	assert.Equal(t, 1, len(server.TLSConfig.Certificates))
	assert.Equal(t, *dummyServerCert, server.TLSConfig.Certificates[0])
	assert.Equal(t, tls.RequireAndVerifyClientCert, server.TLSConfig.ClientAuth)
	assert.Equal(t, dummyCaCertPool, server.TLSConfig.ClientCAs)
	assert.Empty(t, server.TLSConfig.CipherSuites)
	assert.Equal(t, uint16(tls.VersionTLS12), server.TLSConfig.MinVersion)
	assert.Zero(t, server.WriteTimeout)
	assert.Zero(t, server.ReadTimeout)
	assert.Zero(t, server.IdleTimeout)
	assert.True(t, https)
}

func TestListenAndServe_NoListener_HTTPS(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		id:            uuid.New(),
		customization: dummyCustomization,
	}
	var dummyServer = &http.Server{}
	var dummyServeHTTPS = true
	var dummyError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock((*DefaultCustomization).Listener).Expects(dummyCustomization).Returns(nil).Once()
	m.Mock((*http.Server).ListenAndServeTLS).Expects(dummyServer, "", "").Returns(dummyError).Once()

	// SUT + act
	var err = listenAndServe(
		dummySession,
		dummyServer,
		dummyServeHTTPS,
	)

	// assert
	assert.Equal(t, dummyError, err)
}

func TestListenAndServe_NoListener_HTTP(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		id:            uuid.New(),
		customization: dummyCustomization,
	}
	var dummyServer = &http.Server{}
	var dummyServeHTTPS = false
	var dummyError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock((*DefaultCustomization).Listener).Expects(dummyCustomization).Returns(nil).Once()
	m.Mock((*http.Server).ListenAndServe).Expects(dummyServer).Returns(dummyError).Once()

	// SUT + act
	var err = listenAndServe(
		dummySession,
		dummyServer,
		dummyServeHTTPS,
	)

	// assert
	assert.Equal(t, dummyError, err)
}

func TestListenAndServe_WithListener_HTTPS(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		id:            uuid.New(),
		customization: dummyCustomization,
	}
	var dummyServer = &http.Server{}
	var dummyServeHTTPS = true
	var dummyListener = &net.UnixListener{}
	var dummyError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock((*DefaultCustomization).Listener).Expects(dummyCustomization).Returns(dummyListener).Once()
	m.Mock((*http.Server).ServeTLS).Expects(dummyServer, dummyListener, "", "").Returns(dummyError).Once()

	// SUT + act
	var err = listenAndServe(
		dummySession,
		dummyServer,
		dummyServeHTTPS,
	)

	// assert
	assert.Equal(t, dummyError, err)
}

func TestListenAndServe_WithListener_HTTP(t *testing.T) {
	// arrange
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		id:            uuid.New(),
		customization: dummyCustomization,
	}
	var dummyServer = &http.Server{}
	var dummyServeHTTPS = false
	var dummyListener = &net.UnixListener{}
	var dummyError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock((*DefaultCustomization).Listener).Expects(dummyCustomization).Returns(dummyListener).Once()
	m.Mock((*http.Server).Serve).Expects(dummyServer, dummyListener).Returns(dummyError).Once()

	// SUT + act
	var err = listenAndServe(
		dummySession,
		dummyServer,
		dummyServeHTTPS,
	)

	// assert
	assert.Equal(t, dummyError, err)
}

func TestShutdownServer(t *testing.T) {
	// arrange
	var dummyContext = context.TODO()
	var dummyServer = &http.Server{}

	// SUT + act
	var err = shutdownServer(
		dummyContext,
		dummyServer,
	)

	// assert
	assert.NoError(t, err)
}

func TestEvaluateServerErrors_NoErrors(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyHostError error
	var dummyShutDownError error

	// SUT + act
	var result = evaluateServerErrors(
		dummySession,
		dummyHostError,
		dummyShutDownError,
	)

	// assert
	assert.True(t, result)
}

func TestEvaluateServerErrors_FilteredErrors(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyHostError = http.ErrServerClosed
	var dummyShutDownError = http.ErrServerClosed

	// SUT + act
	var result = evaluateServerErrors(
		dummySession,
		dummyHostError,
		dummyShutDownError,
	)

	// assert
	assert.True(t, result)
}

func TestEvaluateServerErrors_OtherErrors(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyHostError = errors.New("some host error")
	var dummyShutDownError = errors.New("some shutdown error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(logAppRoot).Expects(dummySession, "server", "runServer", "Host error found: %+v", dummyHostError).Returns().Once()
	m.Mock(logAppRoot).Expects(dummySession, "server", "runServer", "Shutdown error found: %+v", dummyShutDownError).Returns().Once()

	// SUT + act
	var result = evaluateServerErrors(
		dummySession,
		dummyHostError,
		dummyShutDownError,
	)

	// assert
	assert.False(t, result)
}

func TestRunServer_HappyPath(t *testing.T) {
	// arrange
	var dummyAddress = "some address"
	var dummyCustomization = &DefaultCustomization{}
	var dummySession = &session{
		id:            uuid.New(),
		customization: dummyCustomization,
	}
	type router struct {
		chi.Router
	}
	var dummyRouter = &router{}
	var dummyShutdownSignal = make(chan os.Signal)
	var dummyStarted = false
	var dummyServer = &http.Server{}
	var dummyHTTPS = rand.Intn(100) > 50
	var dummyHostError = errors.New("some host error message")
	var dummyBackgroundContext = context.Background()
	var dummyRuntimeContext = context.TODO()
	var dummyCallback = func() {}
	var dummyGraceShutdownWaitTime = time.Duration(rand.Intn(100)) * time.Second
	var dummyShutDownError = errors.New("some shut down error message")
	var dummyResult = rand.Intn(100) > 50

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.Mock(createServer).Expects(dummyAddress, dummySession, dummyRouter).Returns(dummyServer, dummyHTTPS).Once()
	m.Mock(signal.Notify).Expects(gomocker.Anything(), os.Interrupt, syscall.SIGTERM).Returns().Once()
	m.Mock(listenAndServe).Expects(dummySession, dummyServer, dummyHTTPS).Returns(dummyHostError).Once()
	m.Mock(haltServer).Expects(gomocker.Anything()).Returns().SideEffects(gomocker.GeneralSideEffect(
		0, func() { dummyShutdownSignal <- os.Interrupt })).Once()
	m.Mock(logAppRoot).Expects(dummySession, "server", "runServer", "Interrupt signal received: Terminating server").Returns().Once()
	m.Mock(context.Background).Expects().Returns(dummyBackgroundContext).Once()
	m.Mock((*DefaultCustomization).GraceShutdownWaitTime).Expects(dummyCustomization).Returns(dummyGraceShutdownWaitTime).Once()
	m.Mock(context.WithTimeout).Expects(dummyBackgroundContext, dummyGraceShutdownWaitTime).Returns(dummyRuntimeContext, dummyCallback).Once()
	m.Mock(shutdownServer).Expects(dummyRuntimeContext, dummyServer).Returns(dummyShutDownError).Once()
	m.Mock(evaluateServerErrors).Expects(dummySession, dummyHostError, dummyShutDownError).Returns(dummyResult).Once()
	m.Mock(dummyCallback).Expects().Returns().Once()

	// SUT + act
	var result = runServer(
		dummyAddress,
		dummySession,
		dummyRouter,
		dummyShutdownSignal,
		&dummyStarted,
	)

	// assert
	assert.Equal(t, dummyResult, result)
	assert.False(t, dummyStarted)
}

func TestHaltServer(t *testing.T) {
	// arrange
	var shutdownSignal = make(chan os.Signal)

	// SUT
	go haltServer(shutdownSignal)

	// act
	var result, ok = <-shutdownSignal

	// assert
	assert.True(t, ok)
	assert.Equal(t, os.Interrupt, result)
}
