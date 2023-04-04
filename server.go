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
	app *application,
	session *session,
	shutdownSignal chan os.Signal,
	started *bool,
) error {
	var router, routerError = instantiateRouterFunc(
		app,
		session,
	)
	if routerError != nil {
		return routerError
	}
	logAppRootFunc(
		session,
		"server",
		"hostServer",
		"Targeting address [%v]",
		app.address,
	)
	if runServerFunc(
		app.address,
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
	address string,
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
	return &http.Server{
		Addr:      address,
		TLSConfig: tlsConfig,
		Handler: session.customization.WrapHandler(
			router,
		),
	}, https
}

func listenAndServe(
	session *session,
	server *http.Server,
	serveHTTPS bool,
) error {
	var listener = session.customization.Listener()
	if listener == nil {
		if serveHTTPS {
			return server.ListenAndServeTLS("", "")
		}
		return server.ListenAndServe()
	} else {
		if serveHTTPS {
			return server.ServeTLS(listener, "", "")
		}
		return server.Serve(listener)
	}
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
	address string,
	session *session,
	router *mux.Router,
	shutdownSignal chan os.Signal,
	started *bool,
) bool {
	var server, https = createServerFunc(
		address,
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
			session,
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
