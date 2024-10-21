package webserver

import (
	"os"

	"github.com/google/uuid"
)

// Application is the interface for web server application
type Application interface {
	// Start starts the web server hosting in the current running thread, causing the thread to be blocked until a Stop is called or an interrupt signal is received
	Start()
	// Session retrieves the application-level session instance for logging or any other necessary operations
	Session() Session
	// IsRunning returns true if the server has been successfully started and is currently running
	IsRunning() bool
	// Stop interrupts the web server hosting, causing the web server to gracefully shutdown; a synchronous Start would then return, or an asynchronous StartSync would mark its wait group done and then return
	Stop()
}

type application struct {
	name           string
	address        string
	version        string
	session        *session
	customization  Customization
	actionFuncMap  map[string]ActionFunc
	shutdownSignal chan os.Signal
	started        bool
}

// NewApplication creates a new application for web server hosting
func NewApplication(
	name string,
	address string,
	version string,
	customization Customization,
) Application {
	if isInterfaceValueNil(customization) {
		customization = customizationDefault
	}
	var application = &application{
		name,
		address,
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
		map[string]ActionFunc{},
		make(chan os.Signal),
		false,
	}
	return application
}

func (app *application) Start() {
	startApplication(
		app,
	)
}

func (app *application) Session() Session {
	return app.session
}

func (app *application) IsRunning() bool {
	return app.started
}

func (app *application) Stop() {
	if !app.started {
		return
	}
	haltServer(
		app.shutdownSignal,
	)
}

func startApplication(app *application) {
	if app.started {
		return
	}
	if !preBootstraping(app) {
		return
	}
	bootstrap(app)
	if !postBootstraping(app) {
		return
	}
	defer endApplication(app)
	beginApplication(app)
}

func preBootstraping(app *application) bool {
	var preBootstrapError = app.customization.PreBootstrap()
	if preBootstrapError != nil {
		logAppRoot(
			app.session,
			"application",
			"preBootstraping",
			"Failed to execute customization.PreBootstrap. Error: %+v",
			preBootstrapError,
		)
		return false
	}
	logAppRoot(
		app.session,
		"application",
		"preBootstraping",
		"customization.PreBootstrap executed successfully",
	)
	return true
}

func bootstrap(app *application) {
	initializeHTTPClients(
		app.customization.DefaultTimeout(),
		app.customization.SkipServerCertVerification(),
		app.customization.ClientCert(),
		app.customization.RoundTripper,
	)
	logAppRoot(
		app.session,
		"application",
		"bootstrap",
		"Application bootstrapped successfully",
	)
}

func postBootstraping(app *application) bool {
	var postBootstrapError = app.customization.PostBootstrap()
	if postBootstrapError != nil {
		logAppRoot(
			app.session,
			"application",
			"postBootstraping",
			"Failed to execute customization.PostBootstrap. Error: %+v",
			postBootstrapError,
		)
		return false
	}
	logAppRoot(
		app.session,
		"application",
		"postBootstraping",
		"customization.PostBootstrap executed successfully",
	)
	return true
}

func beginApplication(app *application) {
	logAppRoot(
		app.session,
		"application",
		"beginApplication",
		"Trying to start server [%v] (v-%v)",
		app.name,
		app.version,
	)
	var serverHostError = hostServer(
		app,
		app.session,
		app.shutdownSignal,
		&app.started,
	)
	if serverHostError != nil {
		logAppRoot(
			app.session,
			"application",
			"beginApplication",
			"Failed to host server. Error: %+v",
			serverHostError,
		)
	} else {
		logAppRoot(
			app.session,
			"application",
			"beginApplication",
			"Server hosting terminated",
		)
	}
}

func endApplication(app *application) {
	var appClosingError = app.customization.AppClosing()
	if appClosingError != nil {
		logAppRoot(
			app.session,
			"application",
			"endApplication",
			"Failed to execute customization.AppClosing. Error: %+v",
			appClosingError,
		)
	} else {
		logAppRoot(
			app.session,
			"application",
			"endApplication",
			"customization.AppClosing executed successfully",
		)
	}
}
