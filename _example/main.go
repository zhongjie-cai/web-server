package main

import (
	"fmt"
	"net/http"

	webserver "github.com/zhongjie-cai/web-server"
)

// This is a sample of how to setup application for running the server
func main() {
	var application = webserver.NewApplication(
		"my web server",
		":18605",
		"0.0.1",
		&myCustomization{},
	)
	defer application.Stop()
	application.Start()
}

// myCustomization inherits from the default customization so you can skip setting up all customization methods
//
//	alternatively, you could bring in your own struct that instantiate the webserver.Customization interface to have a verbosed control over what to customize
type myCustomization struct {
	webserver.DefaultCustomization
}

func (customization *myCustomization) Log(session webserver.Session, logType webserver.LogType, logLevel webserver.LogLevel, category, subcategory, description string) {
	fmt.Printf("[%v|%v] <%v|%v> %v\n", logType, logLevel, category, subcategory, description)
}

func (customization *myCustomization) Middlewares() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		loggingRequestURIMiddleware,
	}
}

func (customization *myCustomization) Routes() []webserver.Route {
	return []webserver.Route{
		webserver.Route{
			Method:     http.MethodGet,
			Path:       "/health",
			ActionFunc: getHealth,
		},
		webserver.Route{
			Method:     http.MethodGet,
			Path:       "/health/{data}",
			ActionFunc: getHealth,
			Parameters: map[string]webserver.ParameterType{
				"data": webserver.ParameterTypeAnything,
			},
		},
		webserver.Route{
			Method:     http.MethodGet,
			Path:       "/health/simple",
			ActionFunc: getHealthSimple,
		},
	}
}

func getHealth(
	session webserver.Session,
) (any, error) {
	var value string
	var err = session.GetRequestParameter("data", &value)
	if err != nil {
		return nil, err
	}
	session.LogMethodLogic(
		webserver.LogLevelWarn,
		"Health",
		"Summary",
		"Value = %v",
		value,
	)
	return value, nil
}

func getHealthSimple(
	session webserver.Session,
) (any, error) {
	session.LogMethodLogic(
		webserver.LogLevelWarn,
		"Health",
		"Summary",
		"Simple",
	)
	return "--", nil
}

// loggingRequestURIMiddleware is an example of how a middleware function is written with this library
func loggingRequestURIMiddleware(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(
			responseWriter http.ResponseWriter,
			httpRequest *http.Request,
		) {
			// middleware logic & processing
			fmt.Println(httpRequest.RequestURI)
			// hand over to next handler in the chain
			nextHandler.ServeHTTP(responseWriter, httpRequest)
		},
	)
}
