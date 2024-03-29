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
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
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
	createMock(t)

	// expect
	instantiateRouterFuncExpected = 1
	instantiateRouterFunc = func(app *application, session *session) (*mux.Router, error) {
		instantiateRouterFuncCalled++
		assert.Equal(t, dummyApplication, app)
		assert.Equal(t, dummySession, session)
		return dummyRouter, dummyError
	}

	// SUT + act
	var err = hostServer(
		dummyApplication,
		dummySession,
		dummyShutdownSignal,
		&dummyStarted,
	)

	// assert
	assert.Equal(t, dummyError, err)

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	instantiateRouterFuncExpected = 1
	instantiateRouterFunc = func(app *application, session *session) (*mux.Router, error) {
		instantiateRouterFuncCalled++
		assert.Equal(t, dummyApplication, app)
		assert.Equal(t, dummySession, session)
		return dummyRouter, nil
	}
	logAppRootFuncExpected = 2
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "server", category)
		assert.Equal(t, "hostServer", subcategory)
		if logAppRootFuncCalled == 1 {
			assert.Equal(t, "Targeting address [%v]", messageFormat)
			assert.Equal(t, 1, len(parameters))
			assert.Equal(t, dummyAddress, parameters[0])
		} else if logAppRootFuncCalled == 2 {
			assert.Equal(t, "Server terminated", messageFormat)
			assert.Empty(t, parameters)
		}
	}
	runServerFuncExpected = 1
	runServerFunc = func(address string, session *session, router *mux.Router, shutdownSignal chan os.Signal, started *bool) bool {
		runServerFuncCalled++
		assert.Equal(t, dummyAddress, address)
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
		assert.Equal(t, dummyShutdownSignal, shutdownSignal)
		assert.Equal(t, &dummyStarted, started)
		return false
	}
	newAppErrorFuncExpected = 1
	newAppErrorFunc = func(errorCode errorCode, errorMessage string, innerErrors []error) *appError {
		newAppErrorFuncCalled++
		assert.Equal(t, errorCodeGeneralFailure, errorCode)
		assert.Equal(t, errorMessageHostServer, errorMessage)
		assert.Empty(t, innerErrors)
		return dummyAppError
	}

	// SUT + act
	var err = hostServer(
		dummyApplication,
		dummySession,
		dummyShutdownSignal,
		&dummyStarted,
	)

	// assert
	assert.Equal(t, dummyAppError, err)

	// verify
	verifyAll(t)
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
	createMock(t)

	// expect
	instantiateRouterFuncExpected = 1
	instantiateRouterFunc = func(app *application, session *session) (*mux.Router, error) {
		instantiateRouterFuncCalled++
		assert.Equal(t, dummyApplication, app)
		assert.Equal(t, dummySession, session)
		return dummyRouter, nil
	}
	logAppRootFuncExpected = 2
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "server", category)
		assert.Equal(t, "hostServer", subcategory)
		if logAppRootFuncCalled == 1 {
			assert.Equal(t, "Targeting address [%v]", messageFormat)
			assert.Equal(t, 1, len(parameters))
			assert.Equal(t, dummyAddress, parameters[0])
		} else if logAppRootFuncCalled == 2 {
			assert.Equal(t, "Server closed", messageFormat)
			assert.Empty(t, parameters)
		}
	}
	runServerFuncExpected = 1
	runServerFunc = func(address string, session *session, router *mux.Router, shutdownSignal chan os.Signal, started *bool) bool {
		runServerFuncCalled++
		assert.Equal(t, dummyAddress, address)
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
		assert.Equal(t, dummyShutdownSignal, shutdownSignal)
		assert.Equal(t, &dummyStarted, started)
		return true
	}

	// SUT + act
	var err = hostServer(
		dummyApplication,
		dummySession,
		dummyShutdownSignal,
		&dummyStarted,
	)

	// assert
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

type dummyCustomizationCreateServer struct {
	dummyCustomization
	serverCert  func() *tls.Certificate
	caCertPool  func() *x509.CertPool
	wrapHandler func(http.Handler) http.Handler
	listener    func() net.Listener
}

func (customization *dummyCustomizationCreateServer) ServerCert() *tls.Certificate {
	if customization.serverCert != nil {
		return customization.serverCert()
	}
	assert.Fail(customization.t, "Unexpected call to ServerCert")
	return nil
}

func (customization *dummyCustomizationCreateServer) CaCertPool() *x509.CertPool {
	if customization.caCertPool != nil {
		return customization.caCertPool()
	}
	assert.Fail(customization.t, "Unexpected call to CaCertPool")
	return nil
}

func (customization *dummyCustomizationCreateServer) WrapHandler(handler http.Handler) http.Handler {
	if customization.wrapHandler != nil {
		return customization.wrapHandler(handler)
	}
	assert.Fail(customization.t, "Unexpected call to WrapHandler")
	return nil
}

func (customization *dummyCustomizationCreateServer) Listener() net.Listener {
	if customization.listener != nil {
		return customization.listener()
	}
	assert.Fail(customization.t, "Unexpected call to Listener")
	return nil
}

func TestCreateServer_NoServerCert(t *testing.T) {
	// arrange
	var dummyAddress = "some address"
	var dummyCustomizationCreateServer = &dummyCustomizationCreateServer{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		id:            uuid.New(),
		customization: dummyCustomizationCreateServer,
	}
	var dummyRouter = &mux.Router{KeepContext: rand.Intn(100) > 50}
	var customizationServerCertExpected int
	var customizationServerCertCalled int
	var customizationCaCertPoolExpected int
	var customizationCaCertPoolCalled int
	var customizationWrapHandlerExpected int
	var customizationWrapHandlerCalled int
	var dummyServerCert *tls.Certificate

	// mock
	createMock(t)

	// expect
	customizationServerCertExpected = 1
	dummyCustomizationCreateServer.serverCert = func() *tls.Certificate {
		customizationServerCertCalled++
		return dummyServerCert
	}
	customizationWrapHandlerExpected = 1
	dummyCustomizationCreateServer.wrapHandler = func(handler http.Handler) http.Handler {
		customizationWrapHandlerCalled++
		assert.Equal(t, dummyRouter, handler)
		return dummyRouter
	}

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
	assert.Equal(t, true, server.TLSConfig.PreferServerCipherSuites)
	assert.Equal(t, uint16(tls.VersionTLS12), server.TLSConfig.MinVersion)
	assert.Zero(t, server.WriteTimeout)
	assert.Zero(t, server.ReadTimeout)
	assert.Zero(t, server.IdleTimeout)
	assert.False(t, https)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationServerCertExpected, customizationServerCertCalled, "Unexpected number of calls to customization.ServerCert")
	assert.Equal(t, customizationCaCertPoolExpected, customizationCaCertPoolCalled, "Unexpected number of calls to customization.CaCertPool")
	assert.Equal(t, customizationWrapHandlerExpected, customizationWrapHandlerCalled, "Unexpected number of calls to customization.WrapHandler")
}

func TestCreateServer_WithServerCert_NoCaCertPool(t *testing.T) {
	// arrange
	var dummyAddress = "some address"
	var dummyCustomizationCreateServer = &dummyCustomizationCreateServer{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		id:            uuid.New(),
		customization: dummyCustomizationCreateServer,
	}
	var dummyRouter = &mux.Router{KeepContext: rand.Intn(100) > 50}
	var customizationServerCertExpected int
	var customizationServerCertCalled int
	var customizationCaCertPoolExpected int
	var customizationCaCertPoolCalled int
	var customizationWrapHandlerExpected int
	var customizationWrapHandlerCalled int
	var dummyServerCert = &tls.Certificate{}
	var dummyCaCertPool *x509.CertPool

	// mock
	createMock(t)

	// expect
	customizationServerCertExpected = 1
	dummyCustomizationCreateServer.serverCert = func() *tls.Certificate {
		customizationServerCertCalled++
		return dummyServerCert
	}
	customizationCaCertPoolExpected = 1
	dummyCustomizationCreateServer.caCertPool = func() *x509.CertPool {
		customizationCaCertPoolCalled++
		return dummyCaCertPool
	}
	customizationWrapHandlerExpected = 1
	dummyCustomizationCreateServer.wrapHandler = func(handler http.Handler) http.Handler {
		customizationWrapHandlerCalled++
		assert.Equal(t, dummyRouter, handler)
		return dummyRouter
	}

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
	assert.Equal(t, true, server.TLSConfig.PreferServerCipherSuites)
	assert.Equal(t, uint16(tls.VersionTLS12), server.TLSConfig.MinVersion)
	assert.Zero(t, server.WriteTimeout)
	assert.Zero(t, server.ReadTimeout)
	assert.Zero(t, server.IdleTimeout)
	assert.True(t, https)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationServerCertExpected, customizationServerCertCalled, "Unexpected number of calls to customization.ServerCert")
	assert.Equal(t, customizationCaCertPoolExpected, customizationCaCertPoolCalled, "Unexpected number of calls to customization.CaCertPool")
	assert.Equal(t, customizationWrapHandlerExpected, customizationWrapHandlerCalled, "Unexpected number of calls to customization.WrapHandler")
}

func TestCreateServer_WithServerCert_WithCaCertPool(t *testing.T) {
	// arrange
	var dummyAddress = "some address"
	var dummyCustomizationCreateServer = &dummyCustomizationCreateServer{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		id:            uuid.New(),
		customization: dummyCustomizationCreateServer,
	}
	var dummyRouter = &mux.Router{KeepContext: rand.Intn(100) > 50}
	var customizationServerCertExpected int
	var customizationServerCertCalled int
	var customizationCaCertPoolExpected int
	var customizationCaCertPoolCalled int
	var customizationWrapHandlerExpected int
	var customizationWrapHandlerCalled int
	var dummyServerCert = &tls.Certificate{}
	var dummyCaCertPool = &x509.CertPool{}

	// mock
	createMock(t)

	// expect
	customizationServerCertExpected = 1
	dummyCustomizationCreateServer.serverCert = func() *tls.Certificate {
		customizationServerCertCalled++
		return dummyServerCert
	}
	customizationCaCertPoolExpected = 1
	dummyCustomizationCreateServer.caCertPool = func() *x509.CertPool {
		customizationCaCertPoolCalled++
		return dummyCaCertPool
	}
	customizationWrapHandlerExpected = 1
	dummyCustomizationCreateServer.wrapHandler = func(handler http.Handler) http.Handler {
		customizationWrapHandlerCalled++
		assert.Equal(t, dummyRouter, handler)
		return dummyRouter
	}

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
	assert.Equal(t, true, server.TLSConfig.PreferServerCipherSuites)
	assert.Equal(t, uint16(tls.VersionTLS12), server.TLSConfig.MinVersion)
	assert.Zero(t, server.WriteTimeout)
	assert.Zero(t, server.ReadTimeout)
	assert.Zero(t, server.IdleTimeout)
	assert.True(t, https)

	// verify
	verifyAll(t)
	assert.Equal(t, customizationServerCertExpected, customizationServerCertCalled, "Unexpected number of calls to customization.ServerCert")
	assert.Equal(t, customizationCaCertPoolExpected, customizationCaCertPoolCalled, "Unexpected number of calls to customization.CaCertPool")
	assert.Equal(t, customizationWrapHandlerExpected, customizationWrapHandlerCalled, "Unexpected number of calls to customization.WrapHandler")
}

func TestListenAndServe_NoListener_HTTPS(t *testing.T) {
	// arrange
	var dummyCustomizationCreateServer = &dummyCustomizationCreateServer{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		id:            uuid.New(),
		customization: dummyCustomizationCreateServer,
	}
	var dummyServer = &http.Server{}
	var dummyServeHTTPS = true
	
	// stub
	dummyCustomizationCreateServer.listener = func() net.Listener {
		return nil
	}

	// mock
	createMock(t)

	// SUT + act
	assert.NotPanics(
		t,
		func() {
			go listenAndServe(
				dummySession,
				dummyServer,
				dummyServeHTTPS,
			)
			var err = dummyServer.Close()

			// assert
			assert.NoError(t, err)
		},
	)

	// verify
	verifyAll(t)
}

func TestListenAndServe_NoListener_HTTP(t *testing.T) {
	// arrange
	var dummyCustomizationCreateServer = &dummyCustomizationCreateServer{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		id:            uuid.New(),
		customization: dummyCustomizationCreateServer,
	}
	var dummyServer = &http.Server{}
	var dummyServeHTTPS = false

	// stub
	dummyCustomizationCreateServer.listener = func() net.Listener {
		return nil
	}

	// mock
	createMock(t)

	// SUT + act
	assert.NotPanics(
		t,
		func() {
			go listenAndServe(
				dummySession,
				dummyServer,
				dummyServeHTTPS,
			)
			var err = dummyServer.Close()

			// assert
			assert.NoError(t, err)
		},
	)

	// verify
	verifyAll(t)
}

func TestListenAndServe_WithListener_HTTPS(t *testing.T) {
	// arrange
	var dummyCustomizationCreateServer = &dummyCustomizationCreateServer{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		id:            uuid.New(),
		customization: dummyCustomizationCreateServer,
	}
	var dummyServer = &http.Server{}
	var dummyServeHTTPS = true
	var dummyListener = &net.UnixListener{}

	// expect
	dummyCustomizationCreateServer.listener = func() net.Listener {
		return dummyListener
	}

	// mock
	createMock(t)

	// SUT + act
	assert.NotPanics(
		t,
		func() {
			go listenAndServe(
				dummySession,
				dummyServer,
				dummyServeHTTPS,
			)
			var err = dummyServer.Close()

			// assert
			assert.NoError(t, err)
		},
	)

	// verify
	verifyAll(t)
}

func TestListenAndServe_WithListener_HTTP(t *testing.T) {
	// arrange
	var dummyCustomizationCreateServer = &dummyCustomizationCreateServer{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		id:            uuid.New(),
		customization: dummyCustomizationCreateServer,
	}
	var dummyServer = &http.Server{}
	var dummyServeHTTPS = false
	var dummyListener = &net.UnixListener{}

	// expect
	dummyCustomizationCreateServer.listener = func() net.Listener {
		return dummyListener
	}

	// mock
	createMock(t)

	// SUT + act
	assert.NotPanics(
		t,
		func() {
			go listenAndServe(
				dummySession,
				dummyServer,
				dummyServeHTTPS,
			)
			var err = dummyServer.Close()

			// assert
			assert.NoError(t, err)
		},
	)

	// verify
	verifyAll(t)
}

func TestShutdownServer(t *testing.T) {
	// arrange
	var dummyContext = context.TODO()
	var dummyServer = &http.Server{}

	// mock
	createMock(t)

	// SUT + act
	var err = shutdownServer(
		dummyContext,
		dummyServer,
	)

	// assert
	assert.NoError(t, err)

	// verify
	verifyAll(t)
}

func TestEvaluateServerErrors_NoErrors(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyHostError error
	var dummyShutDownError error

	// mock
	createMock(t)

	// SUT + act
	var result = evaluateServerErrors(
		dummySession,
		dummyHostError,
		dummyShutDownError,
	)

	// assert
	assert.True(t, result)

	// verify
	verifyAll(t)
}

func TestEvaluateServerErrors_FilteredErrors(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyHostError = http.ErrServerClosed
	var dummyShutDownError = http.ErrServerClosed

	// mock
	createMock(t)

	// SUT + act
	var result = evaluateServerErrors(
		dummySession,
		dummyHostError,
		dummyShutDownError,
	)

	// assert
	assert.True(t, result)

	// verify
	verifyAll(t)
}

func TestEvaluateServerErrors_OtherErrors(t *testing.T) {
	// arrange
	var dummySession = &session{id: uuid.New()}
	var dummyHostError = errors.New("some host error")
	var dummyShutDownError = errors.New("some shutdown error")

	// mock
	createMock(t)

	// expect
	logAppRootFuncExpected = 2
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "server", category)
		assert.Equal(t, "runServer", subcategory)
		assert.Equal(t, 1, len(parameters))
		if logAppRootFuncCalled == 1 {
			assert.Equal(t, "Host error found: %+v", messageFormat)
			assert.Equal(t, dummyHostError, parameters[0])
		} else if logAppRootFuncCalled == 2 {
			assert.Equal(t, "Shutdown error found: %+v", messageFormat)
			assert.Equal(t, dummyShutDownError, parameters[0])
		}
	}

	// SUT + act
	var result = evaluateServerErrors(
		dummySession,
		dummyHostError,
		dummyShutDownError,
	)

	// assert
	assert.False(t, result)

	// verify
	verifyAll(t)
}

type dummyCustomizationRunServer struct {
	dummyCustomization
	graceShutdownWaitTime func() time.Duration
}

func (customization *dummyCustomizationRunServer) GraceShutdownWaitTime() time.Duration {
	if customization.graceShutdownWaitTime != nil {
		return customization.graceShutdownWaitTime()
	}
	assert.Fail(customization.t, "Unexpected call to GraceShutdownWaitTime")
	return 0
}

func TestRunServer_HappyPath(t *testing.T) {
	// arrange
	var dummyAddress = "some address"
	var dummyCustomizationRunServer = &dummyCustomizationRunServer{
		dummyCustomization: dummyCustomization{t: t},
	}
	var dummySession = &session{
		id:            uuid.New(),
		customization: dummyCustomizationRunServer,
	}
	var dummyRouter = &mux.Router{}
	var dummyShutdownSignal = make(chan os.Signal)
	var dummyStarted = false
	var customizationGraceShutdownWaitTimeExpected int
	var customizationGraceShutdownWaitTimeCalled int
	var cancelCallbackExpected int
	var cancelCallbackCalled int
	var cancelCallback func()
	var dummyServer = &http.Server{}
	var dummyHTTPS = rand.Intn(100) > 50
	var dummyHostError = errors.New("some host error message")
	var dummyBackgroundContext = context.Background()
	var dummyRuntimeContext = context.TODO()
	var dummyGraceShutdownWaitTime = time.Duration(rand.Intn(100)) * time.Second
	var dummyShutDownError = errors.New("some shut down error message")
	var dummyResult = rand.Intn(100) > 50

	// mock
	createMock(t)

	// expect
	createServerFuncExpected = 1
	createServerFunc = func(address string, session *session, router *mux.Router) (*http.Server, bool) {
		createServerFuncCalled++
		assert.Equal(t, dummyAddress, address)
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyRouter, router)
		return dummyServer, dummyHTTPS
	}
	signalNotifyExpected = 1
	signalNotify = func(c chan<- os.Signal, sig ...os.Signal) {
		signalNotifyCalled++
		assert.Equal(t, 2, len(sig))
		assert.Equal(t, os.Interrupt, sig[0])
		assert.Equal(t, os.Kill, sig[1])
	}
	listenAndServeFuncExpected = 1
	listenAndServeFunc = func(session *session, server *http.Server, serveHTTPS bool) error {
		listenAndServeFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyServer, server)
		assert.Equal(t, dummyHTTPS, serveHTTPS)
		assert.True(t, dummyStarted)
		return dummyHostError
	}
	haltServerFuncExpected = 1
	haltServerFunc = func(shutdownSignal chan os.Signal) {
		haltServerFuncCalled++
		shutdownSignal <- os.Interrupt
	}
	logAppRootFuncExpected = 1
	logAppRootFunc = func(session *session, category string, subcategory string, messageFormat string, parameters ...interface{}) {
		logAppRootFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, "server", category)
		assert.Equal(t, "runServer", subcategory)
		assert.Equal(t, "Interrupt signal received: Terminating server", messageFormat)
		assert.Empty(t, parameters)
		assert.False(t, dummyStarted)
	}
	contextBackgroundExpected = 1
	contextBackground = func() context.Context {
		contextBackgroundCalled++
		return dummyBackgroundContext
	}
	cancelCallbackExpected = 1
	cancelCallback = func() {
		cancelCallbackCalled++
	}
	customizationGraceShutdownWaitTimeExpected = 1
	dummyCustomizationRunServer.graceShutdownWaitTime = func() time.Duration {
		customizationGraceShutdownWaitTimeCalled++
		return dummyGraceShutdownWaitTime
	}
	contextWithTimeoutExpected = 1
	contextWithTimeout = func(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
		contextWithTimeoutCalled++
		assert.Equal(t, dummyBackgroundContext, parent)
		assert.Equal(t, dummyGraceShutdownWaitTime, timeout)
		return dummyRuntimeContext, cancelCallback
	}
	shutdownServerFuncExpected = 1
	shutdownServerFunc = func(runtimeContext context.Context, server *http.Server) error {
		shutdownServerFuncCalled++
		assert.Equal(t, dummyRuntimeContext, runtimeContext)
		assert.Equal(t, dummyServer, server)
		return dummyShutDownError
	}
	evaluateServerErrorsFuncExpected = 1
	evaluateServerErrorsFunc = func(session *session, hostError error, shutdownError error) bool {
		evaluateServerErrorsFuncCalled++
		assert.Equal(t, dummySession, session)
		assert.Equal(t, dummyHostError, hostError)
		assert.Equal(t, dummyShutDownError, shutdownError)
		return dummyResult
	}

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

	// verify
	verifyAll(t)
	assert.Equal(t, cancelCallbackExpected, cancelCallbackCalled, "Unexpected number of calls to cancelCallback")
	assert.Equal(t, customizationGraceShutdownWaitTimeExpected, customizationGraceShutdownWaitTimeCalled, "Unexpected number of calls to customization.GraceShutdownWaitTime")
}

func TestHaltServer(t *testing.T) {
	// arrange
	var shutdownSignal = make(chan os.Signal)

	// mock
	createMock(t)

	// SUT
	go haltServer(shutdownSignal)

	// act
	var result, ok = <-shutdownSignal

	// assert
	assert.True(t, ok)
	assert.Equal(t, os.Interrupt, result)

	// verify
	verifyAll(t)
}
