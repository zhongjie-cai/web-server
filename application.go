package webserver

import (
	"os"
	"sync"

	"github.com/google/uuid"
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
	if isInterfaceValueNil(customization) {
		customization = customizationDefault
	}
	var application = &application{
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
	if !found {
		return nilApplication
	}
	return application
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

func (app *application) StartAsync(waitGroup *sync.WaitGroup) *sync.WaitGroup {
	if waitGroup == nil {
		waitGroup = &sync.WaitGroup{}
	}
	waitGroup.Add(1)
	go func() {
		if waitGroup != nil {
			defer waitGroup.Done()
		}
		app.Start()
	}()
	return waitGroup
}

func (app *application) Stop() {
	haltServer(
		app.shutdownSignal,
	)
}

func (app *application) preBootstraping() bool {
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

func (app *application) bootstrap() bool {
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
	return true
}

func (app *application) postBootstraping() bool {
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

func (app *application) begin() {
	logAppRoot(
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
		logAppRoot(
			app.session,
			"application",
			"begin",
			"Failed to host server. Error: %+v",
			serverHostError,
		)
	} else {
		logAppRoot(
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
		logAppRoot(
			app.session,
			"application",
			"end",
			"Failed to execute customization.AppClosing. Error: %+v",
			appClosingError,
		)
	} else {
		logAppRoot(
			app.session,
			"application",
			"end",
			"customization.AppClosing executed successfully",
		)
	}
}
