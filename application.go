package webserver

import (
	"os"
	"sync"
)

// Application is the interface for web server application
type Application interface {
	// Start starts the web server hosting in the current running thread, causing the thread to be blocked until a Stop is called or an interrupt signal is received
	Start()
	// StartAsync starts the web server hosting in a new Goroutine (thread); if a sync wait group is provided, the wait group would be used for async management, otherwise a new wait group would be created and returned
	StartAsync(*sync.WaitGroup) *sync.WaitGroup
	// Stop interrupts the web server hosting, causing the web server to gracefully shutdown; a synchronous Start would then return, or an asynchronous StartSync would mark its wait group done and then return
	Stop()
}

var (
	applicationLock = sync.RWMutex{}
	applicationMap  = map[int]*application{}
	nilApplication  = &application{
		session: &session{
			request:        defaultRequest,
			responseWriter: defaultResponseWriter,
			attachment:     map[string]interface{}{},
			customization:  customizationDefault,
		},
		customization:  customizationDefault,
		actionFuncMap:  map[string]ActionFunc{},
		shutdownSignal: make(chan os.Signal, 1),
	}
)

type application struct {
	name           string
	port           int
	version        string
	session        *session
	customization  Customization
	actionFuncMap  map[string]ActionFunc
	shutdownSignal chan os.Signal
}

// NewApplication creates a new application for web server hosting
func NewApplication(
	name string,
	port int,
	version string,
	customization Customization,
) Application {
	applicationLock.Lock()
	defer applicationLock.Unlock()
	if isInterfaceValueNilFunc(customization) {
		customization = customizationDefault
	}
	var application = &application{
		name,
		port,
		version,
		&session{
			uuidNew(),
			name,
			defaultRequest,
			defaultResponseWriter,
			map[string]interface{}{},
			customization,
		},
		customization,
		map[string]ActionFunc{},
		make(chan os.Signal, 1),
	}
	applicationMap[port] = application
	return application
}

func getApplication(
	port int,
) *application {
	applicationLock.RLock()
	defer applicationLock.RUnlock()
	var application, found = applicationMap[port]
	if !found ||
		isInterfaceValueNilFunc(application) {
		return nilApplication
	}
	return application
}

func (app *application) Start() {
	startApplicationFunc(
		app,
	)
}

func (app *application) StartAsync(waitGroup *sync.WaitGroup) *sync.WaitGroup {
	if waitGroup == nil {
		waitGroup = &sync.WaitGroup{}
	}
	waitGroup.Add(1)
	go func() {
		if waitGroup != nil {
			defer waitGroup.Done()
		}
		startApplicationFunc(
			app,
		)
	}()
	return waitGroup
}

func (app *application) Stop() {
	haltServerFunc(
		app.shutdownSignal,
	)
}

func startApplication(app *application) {
	if !preBootstrapingFunc(app) {
		return
	}
	bootstrapFunc(app)
	if !postBootstrapingFunc(app) {
		return
	}
	defer endApplicationFunc(app)
	beginApplicationFunc(app)
}

func preBootstraping(app *application) bool {
	var preBootstrapError = app.customization.PreBootstrap()
	if preBootstrapError != nil {
		logAppRootFunc(
			app.session,
			"application",
			"preBootstraping",
			"Failed to execute customization.PreBootstrap. Error: %+v",
			preBootstrapError,
		)
		return false
	}
	logAppRootFunc(
		app.session,
		"application",
		"preBootstraping",
		"customization.PreBootstrap executed successfully",
	)
	return true
}

func bootstrap(app *application) {
	initializeHTTPClientsFunc(
		app.customization.DefaultTimeout(),
		app.customization.SkipServerCertVerification(),
		app.customization.ClientCert(),
		app.customization.RoundTripper,
	)
	logAppRootFunc(
		app.session,
		"application",
		"bootstrap",
		"Application bootstrapped successfully",
	)
}

func postBootstraping(app *application) bool {
	var postBootstrapError = app.customization.PostBootstrap()
	if postBootstrapError != nil {
		logAppRootFunc(
			app.session,
			"application",
			"postBootstraping",
			"Failed to execute customization.PostBootstrap. Error: %+v",
			postBootstrapError,
		)
		return false
	}
	logAppRootFunc(
		app.session,
		"application",
		"postBootstraping",
		"customization.PostBootstrap executed successfully",
	)
	return true
}

func beginApplication(app *application) {
	logAppRootFunc(
		app.session,
		"application",
		"beginApplication",
		"Trying to start server [%v] (v-%v)",
		app.name,
		app.version,
	)
	var serverHostError = hostServerFunc(
		app.port,
		app.session,
		app.customization,
		app.shutdownSignal,
	)
	if serverHostError != nil {
		logAppRootFunc(
			app.session,
			"application",
			"beginApplication",
			"Failed to host server. Error: %+v",
			serverHostError,
		)
	} else {
		logAppRootFunc(
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
		logAppRootFunc(
			app.session,
			"application",
			"endApplication",
			"Failed to execute customization.AppClosing. Error: %+v",
			appClosingError,
		)
	} else {
		logAppRootFunc(
			app.session,
			"application",
			"endApplication",
			"customization.AppClosing executed successfully",
		)
	}
}
