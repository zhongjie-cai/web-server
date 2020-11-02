package webserver

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// hostServer hosts the service entries and starts HTTPS server
func hostServer(
	port int,
	session *session,
	shutdownSignal chan os.Signal,
	started *bool,
) error {
	var router, routerError = instantiateRouterFunc(
		port,
		session,
	)
	if routerError != nil {
		return routerError
	}
	logAppRootFunc(
		session,
		"server",
		"hostServer",
		"Targeting port [%v]",
		port,
	)
	if runServerFunc(
		port,
		session,
		router,
		shutdownSignal,
		started,
	) {
		logAppRootFunc(
			session,
			"server",
			"hostServer",
			"Server closed",
		)
		return nil
	}
	logAppRootFunc(
		session,
		"server",
		"hostServer",
		"Server terminated",
	)
	return newAppErrorFunc(
		errorCodeGeneralFailure,
		errorMessageHostServer,
		[]error{},
	)
}

func createServer(
	port int,
	session *session,
	router *mux.Router,
) (*http.Server, bool) {
	var tlsConfig = &tls.Config{
		// Force it server side
		PreferServerCipherSuites: true,
		// TLS 1.2 as minimum requirement
		MinVersion: tls.VersionTLS12,
	}
	var https = false
	var serverCert = session.customization.ServerCert()
	if serverCert != nil {
		https = true
		tlsConfig.Certificates = []tls.Certificate{
			*serverCert,
		}
		var caCertPool = session.customization.CaCertPool()
		if caCertPool != nil {
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
			tlsConfig.ClientCAs = caCertPool
		} else {
			tlsConfig.ClientAuth = tls.RequireAnyClientCert
		}
	}
	var address = fmtSprintf(
		":%v",
		port,
	)
	return &http.Server{
		Addr:      address,
		TLSConfig: tlsConfig,
		Handler: session.customization.WrapHandler(
			router,
		),
	}, https
}

func listenAndServe(
	server *http.Server,
	serveHTTPS bool,
) error {
	if serveHTTPS {
		return server.ListenAndServeTLS("", "")
	}
	return server.ListenAndServe()
}

func shutdownServer(
	runtimeContext context.Context,
	server *http.Server,
) error {
	return server.Shutdown(
		runtimeContext,
	)
}

func evaluateServerErrors(
	session *session,
	hostError error,
	shutdownError error,
) bool {
	var result = true
	if hostError != nil &&
		hostError != http.ErrServerClosed {
		logAppRootFunc(
			session,
			"server",
			"runServer",
			"Host error found: %+v",
			hostError,
		)
		result = false
	}
	if shutdownError != nil &&
		shutdownError != http.ErrServerClosed {
		logAppRootFunc(
			session,
			"server",
			"runServer",
			"Shutdown error found: %+v",
			shutdownError,
		)
		result = false
	}
	return result
}

func runServer(
	port int,
	session *session,
	router *mux.Router,
	shutdownSignal chan os.Signal,
	started *bool,
) bool {
	var server, https = createServerFunc(
		port,
		session,
		router,
	)

	signalNotify(
		shutdownSignal,
		os.Interrupt,
		os.Kill,
	)

	*started = true

	var hostError error
	go func() {
		hostError = listenAndServeFunc(
			server,
			https,
		)
		haltServerFunc(
			shutdownSignal,
		)
	}()

	<-shutdownSignal

	*started = false

	logAppRootFunc(
		session,
		"server",
		"runServer",
		"Interrupt signal received: Terminating server",
	)

	var runtimeContext, cancelCallback = contextWithTimeout(
		contextBackground(),
		session.customization.GraceShutdownWaitTime(),
	)
	defer cancelCallback()

	var shutdownError = shutdownServerFunc(
		runtimeContext,
		server,
	)

	return evaluateServerErrorsFunc(
		session,
		hostError,
		shutdownError,
	)
}

func haltServer(shutdownSignal chan os.Signal) {
	shutdownSignal <- os.Interrupt
}
