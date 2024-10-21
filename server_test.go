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

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/zhongjie-cai/gomocker"
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
	var dummyRouter = &mux.Router{}
	var dummyError = errors.New("some error")

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(instantiateRouter, 1, func(app *application, session *session) (*mux.Router, error) {
		assert.Equal(t, dummyApplication, app)
		assert.Equal(t, dummySession, session)
		return dummyRouter, dummyError
	})

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
	var dummyRouter = &mux.Router{}
	var dummyAppError = &appError{Message: "some error message"}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(instantiateRouter, 1, func(app *application, session *session) (*mux.Router, error) {
		assert.Equal(t, dummyApplication, app)
		assert.Equal(t, dummySession, session)
		return dummyRouter, nil
	})
	m.ExpectFunc(logAppRoot, 2, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "server", category)
		assert.Equal(t, "hostServer", subcategory)
		if m.FuncCalledCount(logAppRoot) == 1 {
			assert.Equal(t, "Targeting address [%v]", messageFormat)
			assert.Equal(t, 1, len(parameters))
			assert.Equal(t, dummyAddress, parameters[0])
		} else if m.FuncCalledCount(logAppRoot) == 2 {
			assert.Equal(t, "Server terminated", messageFormat)
			assert.Empty(t, parameters)
		}
	})
	m.ExpectFunc(runServer, 1, func(address string, session *session, router *mux.Router, shutdownSignal chan os.Signal, started *bool) bool {
		assert.Equal(t, dummyAddress, address)
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
		assert.Equal(t, dummyShutdownSignal, shutdownSignal)
		assert.Equal(t, &dummyStarted, started)
		return false
	})
	m.ExpectFunc(newAppError, 1, func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageHostServer, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	})

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
	var dummyRouter = &mux.Router{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(instantiateRouter, 1, func(app *application, session *session) (*mux.Router, error) {
		assert.Equal(t, dummyApplication, app)
		assert.Equal(t, dummySession, session)
		return dummyRouter, nil
	})
	m.ExpectFunc(logAppRoot, 2, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "server", category)
		assert.Equal(t, "hostServer", subcategory)
		if m.FuncCalledCount(logAppRoot) == 1 {
			assert.Equal(t, "Targeting address [%v]", messageFormat)
			assert.Equal(t, 1, len(parameters))
			assert.Equal(t, dummyAddress, parameters[0])
		} else if m.FuncCalledCount(logAppRoot) == 2 {
			assert.Equal(t, "Server closed", messageFormat)
			assert.Empty(t, parameters)
		}
	})
	m.ExpectFunc(runServer, 1, func(address string, session *session, router *mux.Router, shutdownSignal chan os.Signal, started *bool) bool {
		assert.Equal(t, dummyAddress, address)
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
		assert.Equal(t, dummyShutdownSignal, shutdownSignal)
		assert.Equal(t, &dummyStarted, started)
		return true
	})

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
	var dummyRouter = &mux.Router{KeepContext: rand.Intn(100) > 50}
	var dummyServerCert *tls.Certificate

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "ServerCert", 1, func() *tls.Certificate {
		return dummyServerCert
	})
	m.ExpectMethod(dummyCustomization, "WrapHandler", 1, func(self *DefaultCustomization, handler http.Handler) http.Handler {
		assert.Equal(t, dummyCustomization, self)
		assert.Equal(t, dummyRouter, handler)
		return dummyRouter
	})

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
	var dummyRouter = &mux.Router{KeepContext: rand.Intn(100) > 50}
	var dummyServerCert = &tls.Certificate{}
	var dummyCaCertPool *x509.CertPool

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "ServerCert", 1, func() *tls.Certificate {
		return dummyServerCert
	})
	m.ExpectMethod(dummyCustomization, "CaCertPool", 1, func() *x509.CertPool {
		return dummyCaCertPool
	})
	m.ExpectMethod(dummyCustomization, "WrapHandler", 1, func(self *DefaultCustomization, handler http.Handler) http.Handler {
		assert.Equal(t, dummyCustomization, self)
		assert.Equal(t, dummyRouter, handler)
		return dummyRouter
	})

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
	var dummyRouter = &mux.Router{KeepContext: rand.Intn(100) > 50}
	var dummyServerCert = &tls.Certificate{}
	var dummyCaCertPool = &x509.CertPool{}

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectMethod(dummyCustomization, "ServerCert", 1, func() *tls.Certificate {
		return dummyServerCert
	})
	m.ExpectMethod(dummyCustomization, "CaCertPool", 1, func() *x509.CertPool {
		return dummyCaCertPool
	})
	m.ExpectMethod(dummyCustomization, "WrapHandler", 1, func(self *DefaultCustomization, handler http.Handler) http.Handler {
		assert.Equal(t, dummyCustomization, self)
		assert.Equal(t, dummyRouter, handler)
		return dummyRouter
	})

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
	m.ExpectMethod(dummyCustomization, "Listener", 1, func(self *DefaultCustomization) net.Listener {
		assert.Equal(t, dummyCustomization, self)
		return nil
	})
	m.ExpectMethod(dummyServer, "ListenAndServeTLS", 1, func(self *http.Server, certFile string, keyFile string) error {
		assert.Equal(t, dummyServer, self)
		assert.Equal(t, "", certFile)
		assert.Equal(t, "", keyFile)
		return dummyError
	})

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
	m.ExpectMethod(dummyCustomization, "Listener", 1, func(self *DefaultCustomization) net.Listener {
		assert.Equal(t, dummyCustomization, self)
		return nil
	})
	m.ExpectMethod(dummyServer, "ListenAndServe", 1, func(self *http.Server) error {
		assert.Equal(t, dummyServer, self)
		return dummyError
	})

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
	m.ExpectMethod(dummyCustomization, "Listener", 1, func(self *DefaultCustomization) net.Listener {
		assert.Equal(t, dummyCustomization, self)
		return dummyListener
	})
	m.ExpectMethod(dummyServer, "ServeTLS", 1, func(self *http.Server, listener net.Listener, certFile string, keyFile string) error {
		assert.Equal(t, dummyServer, self)
		assert.Equal(t, dummyListener, listener)
		assert.Equal(t, "", certFile)
		assert.Equal(t, "", keyFile)
		return dummyError
	})

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
	m.ExpectMethod(dummyCustomization, "Listener", 1, func(self *DefaultCustomization) net.Listener {
		assert.Equal(t, dummyCustomization, self)
		return dummyListener
	})
	m.ExpectMethod(dummyServer, "Serve", 1, func(self *http.Server, listener net.Listener) error {
		assert.Equal(t, dummyServer, self)
		assert.Equal(t, dummyListener, listener)
		return dummyError
	})

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
	m.ExpectFunc(logAppRoot, 2, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "server", category)
		assert.Equal(t, "runServer", subcategory)
		assert.Equal(t, 1, len(parameters))
		if m.FuncCalledCount(logAppRoot) == 1 {
			assert.Equal(t, "Host error found: %+v", messageFormat)
			assert.Equal(t, dummyHostError, parameters[0])
		} else if m.FuncCalledCount(logAppRoot) == 2 {
			assert.Equal(t, "Shutdown error found: %+v", messageFormat)
			assert.Equal(t, dummyShutDownError, parameters[0])
		}
	})

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
	var dummyRouter = &mux.Router{}
	var dummyShutdownSignal = make(chan os.Signal)
	var dummyStarted = false
	var cancelCallbackExpected = 1
	var cancelCallbackCalled = 0
	var dummyServer = &http.Server{}
	var dummyHTTPS = rand.Intn(100) > 50
	var dummyHostError = errors.New("some host error message")
	var dummyBackgroundContext = context.Background()
	var dummyRuntimeContext = context.TODO()
	var dummyGraceShutdownWaitTime = time.Duration(rand.Intn(100)) * time.Second
	var dummyShutDownError = errors.New("some shut down error message")
	var dummyResult = rand.Intn(100) > 50

	// mock
	var m = gomocker.NewMocker(t)

	// expect
	m.ExpectFunc(createServer, 1, func(address string, session *session, router *mux.Router) (*http.Server, bool) {
		assert.Equal(t, dummyAddress, address)
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
		return dummyServer, dummyHTTPS
	})
	m.ExpectFunc(signal.Notify, 1, func(c chan<- os.Signal, sig ...os.Signal) {
		assert.Equal(t, 2, len(sig))
		assert.Equal(t, os.Interrupt, sig[0])
		assert.Equal(t, syscall.SIGTERM, sig[1])
	})
	m.ExpectFunc(listenAndServe, 1, func(session *session, server *http.Server, serveHTTPS bool) error {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyServer, server)
		assert.Equal(t, dummyHTTPS, serveHTTPS)
		assert.True(t, dummyStarted)
		return dummyHostError
	})
	m.ExpectFunc(haltServer, 1, func(shutdownSignal chan os.Signal) {
		shutdownSignal <- os.Interrupt
	})
	m.ExpectFunc(logAppRoot, 1, func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "server", category)
		assert.Equal(t, "runServer", subcategory)
		assert.Equal(t, "Interrupt signal received: Terminating server", messageFormat)
		assert.Empty(t, parameters)
		assert.False(t, dummyStarted)
	})
	m.ExpectFunc(context.Background, 1, func() context.Context {
		return dummyBackgroundContext
	})
	var cancelCallback = func() {
		cancelCallbackCalled++
	}
	m.ExpectMethod(dummyCustomization, "GraceShutdownWaitTime", 1, func() time.Duration {
		return dummyGraceShutdownWaitTime
	})
	m.ExpectFunc(context.WithTimeout, 1, func(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
		assert.Equal(t, dummyBackgroundContext, parent)
		assert.Equal(t, dummyGraceShutdownWaitTime, timeout)
		return dummyRuntimeContext, cancelCallback
	})
	m.ExpectFunc(shutdownServer, 1, func(runtimeContext context.Context, server *http.Server) error {
		assert.Equal(t, dummyRuntimeContext, runtimeContext)
		assert.Equal(t, dummyServer, server)
		return dummyShutDownError
	})
	m.ExpectFunc(evaluateServerErrors, 1, func(session *session, hostError error, shutdownError error) bool {
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyHostError, hostError)
		assert.Equal(t, dummyShutDownError, shutdownError)
		return dummyResult
	})

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
	assert.Equal(t, cancelCallbackExpected, cancelCallbackCalled)
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
