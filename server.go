package webserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/zhongjie-cai/WebServiceTemplate/customization"
)

// hostServer hosts the service entries and starts HTTPS server
func hostServer(
	port int,
	session *session,
	shutdownSignal chan os.Signal,
) error {
	var router, routerError = instantiateRouter(
		port,
		session,
	)
	if routerError != nil {
		return errRouteRegistration
	}
	logAppRoot(
		session,
		"server",
		"Host",
		"Targeting port [%v]",
		port,
	)
	var hostError = runServer(
		port,
		session,
		router,
		shutdownSignal,
	)
	if hostError != nil {
		logAppRoot(
			session,
			"server",
			"Host",
			"Server failure: %+v",
			hostError,
		)
		return errHostServer
	}
	logAppRoot(
		session,
		"server",
		"Host",
		"Server terminated",
	)
	return nil
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
			if caCertPool != nil {
				tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
				tlsConfig.ClientCAs = caCertPool
			} else {
				tlsConfig.ClientAuth = tls.RequireAnyClientCert
			}
		} else {
			tlsConfig.ClientAuth = tls.RequestClientCert
		}
	}
	var address = fmt.Sprintf(
		":%v",
		port,
	)
	return &http.Server{
		Addr:      address,
		TLSConfig: tlsConfig,
		Handler:   router,
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

func shutDown(
	runtimeContext context.Context,
	server *http.Server,
) error {
	return server.Shutdown(
		runtimeContext,
	)
}

func runServer(
	port int,
	session *session,
	router *mux.Router,
	shutdownSignal chan os.Signal,
) error {
	var server, https = createServer(
		port,
		session,
		router,
	)

	signal.Notify(
		shutdownSignal,
		os.Interrupt,
		os.Kill,
	)

	var hostError error
	go func() {
		hostError = listenAndServe(
			server,
			https,
		)
		haltServer(
			shutdownSignal,
		)
	}()

	<-shutdownSignal

	logAppRoot(
		session,
		"server",
		"Host",
		"Interrupt signal received: Terminating server",
	)

	if hostError != http.ErrServerClosed {
		logAppRoot(
			session,
			"server",
			"Host",
			"Host error found: %+v",
			hostError,
		)
	}

	var runtimeContext, cancelCallback = context.WithTimeout(
		context.Background(),
		customization.GraceShutdownWaitTime(),
	)
	defer cancelCallback()

	return shutDown(
		runtimeContext,
		server,
	)
}

func haltServer(shutdownSignal chan os.Signal) {
	shutdownSignal <- os.Interrupt
}
