# web-server

[![Build Status](https://travis-ci.com/zhongjie-cai/web-server.svg?branch=master)](https://travis-ci.com/zhongjie-cai/web-server)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/zhongjie-cai/web-server)
[![Go Report Card](https://goreportcard.com/badge/github.com/zhongjie-cai/web-server)](https://goreportcard.com/report/github.com/zhongjie-cai/web-server)
[![Coverage](http://gocover.io/_badge/github.com/zhongjie-cai/web-server)](http://gocover.io/github.com/zhongjie-cai/web-server)

This library is provided as a wrapper utility for quickly create and host your web services.

Original source: https://github.com/zhongjie-cai/web-server

Library dependencies (must be present in vendor folder or in Go path):
* [UUID](https://github.com/google/uuid): `go get -u github.com/google/uuid`
* [MUX](https://github.com/gorilla/mux): `go get -u github.com/gorilla/mux`
* [Testify](https://github.com/stretchr/testify): `go get -u github.com/stretchr/testify` (For tests only)

A sample application is shown below:

# main.go
```golang
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
		18605,
		"0.0.1",
		&myCustomization{},
	)
	defer application.Stop()
	application.Start()
}

// myCustomization inherits from the default customization so you can skip setting up all customization methods
//   alternatively, you could bring in your own struct that instantiate the webserver.Customization interface to have a verbosed control over what to customize
type myCustomization struct {
	webserver.DefaultCustomization
}

func (customization *myCustomization) Log(session webserver.Session, logType webserver.LogType, logLevel webserver.LogLevel, category, subcategory, description string) {
	fmt.Printf("[%v|%v] <%v|%v> %v\n", logType, logLevel, category, subcategory, description)
}

func (customization *myCustomization) Middlewares() []webserver.MiddlewareFunc {
	return []webserver.MiddlewareFunc{
		loggingRequestURIMiddleware,
	}
}

func (customization *myCustomization) Statics() []webserver.Static {
	return []webserver.Static{
		webserver.Static{
			Name:       "SwaggerUI",
			PathPrefix: "/docs/",
			Handler:    swaggerHandler(),
		},
		webserver.Static{
			Name:       "SwaggerRedirect",
			PathPrefix: "/docs",
			Handler:    swaggerRedirect(),
		},
	}
}
func (customization *myCustomization) Routes() []webserver.Route {
	return []webserver.Route{
		webserver.Route{
			Endpoint:   "Health",
			Method:     http.MethodGet,
			Path:       "/health",
			ActionFunc: getHealth,
		},
	}
}

// getHealth is an example of how a normal HTTP handling method is written with this library
func getHealth(
	session webserver.Session,
) (interface{}, error) {
	var appVersion = "some application version"
	session.LogMethodLogic(
		webserver.LogLevelWarn,
		"Health",
		"Summary",
		"AppVersion = %v",
		appVersion,
	)
	return appVersion, nil
}

// swaggerRedirect is an example of how an HTTP redirection could be managed with this library
func swaggerRedirect() http.Handler {
	return http.RedirectHandler(
		"/docs/",
		http.StatusPermanentRedirect,
	)
}

// swaggerHandler is an example of how a normal HTTP static content hosting is written with this library
func swaggerHandler() http.Handler {
	return http.StripPrefix(
		"/docs/",
		http.FileServer(
			http.Dir("./docs"),
		),
	)
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
```

# Request & Response

The registered handler could retrieve request body, parameters and query strings through session methods, thus it is normally not necessary to load request from session:

```golang
// request body: {"foo":"bar","test":123}
var body struct {
	Foo string `json:"foo"`
	Test int `json:"test"`
}
var bodyError = session.GetRequestBody(&body)

// parameters: "id"=456
var id int
var idError = session.GetRequestParameter("id", &id)

// query strigns: "uuid"="123456-1234-1234-1234-123456789abc"
var uuid uuid.UUID
var uuidError = session.GetRequestQuery("uuid", 0, &uuid)
```

However, if specific data is needed from request, one could always retrieve request from session through following function call using session object:

```golang
var httpRequest = session.GetRequest()
```

The response functions accept the session ID and internally load the response writer accordingly, thus it is normally not necessary to load response writer from session.

However, if specific operation is needed for response, one could always retrieve response writer through following function call using session object:

```golang
var responseWriter = session.GetResponseWriter()
```

After your customized operations on the response writer, use the following as return type in the handler function to enforce the library to skip any unnecessary handling of the HTTP response writer:

```golang
return webserver.SkipResponseHandling()
```

# Error Handling

To simplify the error handling, one could utilize the built-in error interface `AppError`, which provides support to many basic types of errors that are mapped to corresponding HTTP status codes:

* BadRequest       => BadRequest (400)
* Unauthorized     => Unauthorized (401)
* CircuitBreak     => Forbidden (403)
* AccessForbidden  => Forbidden (403)
* NotFound         => NotFound (404)
* InvalidOperation => MethodNotAllowed (405)
* DataCorruption   => Conflict (409)
* OperationLock    => Locked (423)
* GeneralFailure   => InternalServerError (500)
* NotImplemented   => NotImplemented (501)

If you bring in your own implementation for the error interface `AppHTTPError`, the web server engine could automatically utilise the corresponding methods to translate an `AppHTTPError` into HTTP status code and response message:

```golang
HTTPStatusCode() int
HTTPResponseMessage() string
```

However, if specific operation is needed for response, one could always customize the error interpretation by customizing the `InterpretError` function:

```golang
func (customization *myCustomization) InterpretError(err error) (statusCode int, responseMessage string) {
	return 500, err.Error()
}
```

# Logging

The library allows the user to customize its logging function by customizing the `Log` method. 
The logging is split into two management areas: log type and log level. 

## Log Type

The log type definitions can be found under the `logType.go` file. 
Apart from all `Method`-prefixed log types, all remainig log types are managed by the library internally and should not be worried by the consumer. 

## Log Level

The log level definitions can be found under the `logLevel.go` file. 
Log level only affects all `Method`-prefixed log types; for all other log types, the log level is default to `Info`. 

## Session Logging

The registered session allows the user to add manual logging to its codebase, through several listed methods as
```golang
session.LogMethodEnter()
session.LogMethodParameter(parameters ...interface{})
session.LogMethodLogic(logLevel LogLevel, category string, subcategory string, messageFormat string, parameters ...interface{})
session.LogMethodReturn(returns ...interface{})
session.LogMethodExit()
```

The `Enter`, `Parameter`, `Return` and `Exit` are limited to the scope of method boundary area loggings. 
The `Logic` is the normal logging that can be used in any place at any level in the codebase to enforce the user's customized logging entries.

# Session Attachment

The registered session contains an attachment dictionary, which allows the user to attach any object which is JSON serializable into the given session associated to a session ID.

```golang
var myAttachmentName = "my attachment name"
var myAttachmentObject = anyJSONSerializableStruct {
	...
}
var success = session.Attach(myAttachmentName, myAttachmentObject)
if !success {
	// failed to attach an object: add your customized logic here if needed
} else {
	// succeeded to attach an object: add your customized logic here if needed
}
```

To retrieve a previously attached object from session, simply use the following sample logic.

```golang
var myAttachmentName = "my attachment name"
var retrievedAttachment anyJSONSerializableStruct
var success = session.GetAttachment(myAttachmentName, &retrievedAttachment)
if !success {
	// failed to retrieve an attachment: add your customized logic here if needed
} else {
	// succeeded to retrieve an attachment: add your customized logic here if needed
}
```

In some situations, it is good to detach a certain attachment, especially if it is a big object consuming large memory, which can be done as following.

```golang
var myAttachmentName = "my attachment name"
var success = session.Detach(myAttachmentName)
if !success {
	// failed to detach an attachment: add your customized logic here if needed
} else {
	// succeeded to detach an attachment: add your customized logic here if needed
}
```

# External Webcall Requests

The library provides a way to send out HTTP/HTTPS requests to external web services based on current session. 
Using this provided feature ensures the logging of the web service requests into corresponding log type for the given session. 

You can reuse a same struct for multiple HTTP status codes, as long as the structures in JSON format are compatible.
If there is no receiver entry registered for a particular HTTP status code, the corresponding response body is ignored for deserialization when that HTTP status code is received.

```golang
...

var webcallRequest = session.CreateWebcallRequest(
	HTTP.POST,                       // Method
	"https://www.example.com/tests", // URL
	"{\"foo\":\"bar\"}",             // Payload
	true,                            // SendClientCert
)
var responseOn200 responseOn200Struct
var responseOn400 responseOn400Struct
var responseOn500 responseOn500Struct
webcallRequest.AddHeader(
	"Content-Type",
	"application/json",
).AddHeader(
	"Accept",
	"application/json",
).Anticipate(
	http.StatusOK,
	http.StatusBadRequest,
	&responseOn200,
).Anticipate(
	http.StatusBadRequest,
	http.StatusInternalServerError,
	&responseOn400,
).Anticipate(
	http.StatusInternalServerError,
	999,
	&responseOn500,
)
var statusCode, responseHeader, responseError = webcallRequest.Process()

...
```

Webcall requests would send out client certificate for mTLS communications if the following customization is in place.

```golang
func (customization *myCustomization) ClientCert() *tls.Certificate {
    return ... // replace with however you would load the client certificate
}
```

Webcall requests could also be customized for：

## HTTP Client's HTTP Transport (http.RoundTripper)

This is to enable the 3rd party monitoring libraries, e.g. new relic, to wrap the HTTP transport for better handling of webcall communications. 

```golang
func (customization *myCustomization) RoundTripper(originalTransport http.RoundTripper) http.RoundTripper {
	return ... // replace with whatever round trip wrapper logic you would like to have
}
```

## HTTP Request (http.Request)

This is to enable the 3rd party monitoring libraries, e.g. new relic, to wrap individual HTTP request for better handling of web requests.

```golang
func (customization *myCustomization) WrapRequest(session Session, httpRequest *http.Request) *http.Request {
	return ... // replace with whatever HTTP request wrapper logic you would like to have
}
```

## Webcall Timeout

This is to provide the default HTTP request timeouts for HTTP Client over all webcall communications.

```golang
func (customization *myCustomization) DefaultTimeout() time.Duration {
	return 3 * time.Minute // replace with whatever timeout duration you would like to have
}
```
