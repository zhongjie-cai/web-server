package webserver

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

// hostServer hosts the service entries and starts HTTPS server
func hostServer(
	app *application,
	session *session,
	shutdownSignal chan os.Signal,
	started *bool,
) error {
	var router, routerError = instantiateRouter(
		app,
		session,
	)
	if routerError != nil {
		return routerError
	}
	logAppRoot(
		session,
		"server",
		"hostServer",
		"Targeting address [%v]",
		app.address,
	)
	if runServer(
		app.address,
		session,
		router,
		shutdownSignal,
		started,
	) {
		logAppRoot(
			session,
			"server",
			"hostServer",
			"Server closed",
		)
		return nil
	}
	logAppRoot(
		session,
		"server",
		"hostServer",
		"Server terminated",
	)
	return newAppError(
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
		logAppRoot(
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
		logAppRoot(
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
	var server, https = createServer(
		address,
		session,
		router,
	)

	signal.Notify(
		shutdownSignal,
		os.Interrupt,
		syscall.SIGTERM,
	)

	*started = true

	var hostError error
	go func() {
		hostError = listenAndServe(
			session,
			server,
			https,
		)
		haltServer(
			shutdownSignal,
		)
	}()

	<-shutdownSignal

	*started = false

	logAppRoot(
		session,
		"server",
		"runServer",
		"Interrupt signal received: Terminating server",
	)

	var runtimeContext, cancelCallback = context.WithTimeout(
		context.Background(),
		session.customization.GraceShutdownWaitTime(),
	)
	defer cancelCallback()

	var shutdownError = shutdownServer(
		runtimeContext,
		server,
	)

	return evaluateServerErrors(
		session,
		hostError,
		shutdownError,
	)
}

func haltServer(shutdownSignal chan os.Signal) {
	shutdownSignal <- os.Interrupt
}
