package webserver

import (
	"os"

	"github.com/google/uuid"
)

// Application is the interface for web server application
type Application interface {
	Start()
	Stop()
}

type application struct {
	name           string
	port           int
	version        string
	session        *session
	customization  Customization
	shutdownSignal chan os.Signal
}

// NewApplication creates a new application for web server hosting
func NewApplication(
	name string,
	port int,
	version string,
	customization Customization,
) Application {
	registerCustomization(
		port,
		customization,
	)
	return &application{
		name,
		port,
		version,
		&session{
			uuid.New(),
			name,
			defaultRequest,
			defaultResponseWriter,
			map[string]interface{}{},
			customization,
		},
		customization,
		make(chan os.Signal, 1),
	}
}

func (app *application) Start() {
	if !app.preBootstraping() {
		return
	}
	if !app.bootstrap() {
		return
	}
	if !app.postBootstraping() {
		return
	}
	defer app.end()
	app.begin()
}

func (app *application) Stop() {
	haltServer(
		app.shutdownSignal,
	)
}

func (app *application) preBootstraping() bool {
	var preBootstrapError = app.customization.PreBootstrap()
	if preBootstrapError != nil {
		appRoot(
			app.session,
			"application",
			"preBootstraping",
			"Failed to execute customization.PreBootstrap. Error: %+v",
			preBootstrapError,
		)
		return false
	}
	appRoot(
		app.session,
		"application",
		"preBootstraping",
		"customization.PreBootstrap executed successfully",
	)
	return true
}

func (app *application) bootstrap() bool {
	initializeHTTPClients(
		app.customization.DefaultTimeout(),
		app.customization.SkipServerCertVerification(),
		app.customization.ClientCert(),
		app.customization.RoundTripper,
	)
	appRoot(
		app.session,
		"application",
		"bootstrap",
		"Application bootstrapped successfully",
	)
	return true
}

func (app *application) postBootstraping() bool {
	var postBootstrapError = app.customization.PostBootstrap()
	if postBootstrapError != nil {
		appRoot(
			app.session,
			"application",
			"postBootstraping",
			"Failed to execute customization.PostBootstrap. Error: %+v",
			postBootstrapError,
		)
		return false
	}
	appRoot(
		app.session,
		"application",
		"postBootstraping",
		"customization.PostBootstrap executed successfully",
	)
	return true
}

func (app *application) begin() {
	appRoot(
		app.session,
		"application",
		"begin",
		"Trying to start server (v-%v)",
		app.version,
	)
	var serverHostError = hostServer(
		app.port,
		app.session,
		app.customization,
		app.shutdownSignal,
	)
	if serverHostError != nil {
		appRoot(
			app.session,
			"application",
			"begin",
			"Failed to host server. Error: %+v",
			serverHostError,
		)
	} else {
		appRoot(
			app.session,
			"application",
			"begin",
			"Server hosting terminated",
		)
	}
}

func (app *application) end() {
	var appClosingError = app.customization.AppClosing()
	if appClosingError != nil {
		appRoot(
			app.session,
			"application",
			"end",
			"Failed to execute customization.AppClosing. Error: %+v",
			appClosingError,
		)
	} else {
		appRoot(
			app.session,
			"application",
			"end",
			"customization.AppClosing executed successfully",
		)
	}
}
