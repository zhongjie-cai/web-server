package webserver

import (
	"fmt"

	"github.com/google/uuid"
)

// func pointers for injection / testing: application.go
var (
	isInterfaceValueNilFunc   = isInterfaceValueNil
	uuidNew                   = uuid.New
	startApplicationFunc      = startApplication
	haltServerFunc            = haltServer
	preBootstrapingFunc       = preBootstraping
	bootstrapFunc             = bootstrap
	postBootstrapingFunc      = postBootstraping
	endApplicationFunc        = endApplication
	beginApplicationFunc      = beginApplication
	logAppRootFunc            = logAppRoot
	initializeHTTPClientsFunc = initializeHTTPClients
	hostServerFunc            = hostServer
)

// func pointers for injection / testing: customization.go
var (
	fmtPrintf              = fmt.Printf
	fmtSprintf             = fmt.Sprintf
	marshalIgnoreErrorFunc = marshalIgnoreError
)
