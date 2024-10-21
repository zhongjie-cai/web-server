package webserver

import (
	"net/http"
	"reflect"
	"strconv"
)

// These are the constants used by the HTTP modules
const (
	ContentTypeJSON = "application/json; charset=utf-8"
)

type skipResponseHandlingDummy struct{}

var typeOfSkipResponseHandling = reflect.TypeOf(skipResponseHandlingDummy{})

// SkipResponseHandling indicates to the library to skip operating on the HTTP response writer
func SkipResponseHandling() (interface{}, error) {
	return skipResponseHandlingDummy{}, nil
}

func shouldSkipHandling(
	responseObject interface{},
	responseError error,
) bool {
	if responseError != nil {
		return false
	}
	var responseType = reflect.TypeOf(responseObject)
	return responseType == typeOfSkipResponseHandling
}

func constructResponse(
	session *session,
	responseObject interface{},
	responseError error,
) (int, string) {
	if responseError != nil {
		return session.customization.InterpretError(
			responseError,
		)
	}
	return session.customization.InterpretSuccess(
		responseObject,
	)
}

// writeResponse responds to the consumer with corresponding HTTP status code and response body
func writeResponse(
	session *session,
	responseObject interface{},
	responseError error,
) {
	if shouldSkipHandling(
		responseObject,
		responseError,
	) {
		logEndpointResponse(
			session,
			"None",
			"-1",
			"Skipped response handling",
		)
		return
	}
	var statusCode, responseMessage = constructResponse(
		session,
		responseObject,
		responseError,
	)
	logEndpointResponse(
		session,
		http.StatusText(statusCode),
		strconv.Itoa(statusCode),
		responseMessage,
	)
	var responseWriter = session.GetResponseWriter()
	responseWriter.Header().Set("Content-Type", ContentTypeJSON)
	responseWriter.WriteHeader(statusCode)
	responseWriter.Write([]byte(responseMessage))
}
