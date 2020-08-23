package webserver

import (
	"net/http"
	"strconv"
)

// These are the constants used by the HTTP modules
const (
	ContentTypeJSON = "application/json; charset=utf-8"
)

func createOkResponse(
	responseContent interface{},
) (int, string) {
	if responseContent == nil {
		return http.StatusNoContent, ""
	}
	var responseMessage = marshalIgnoreError(responseContent)
	if responseMessage == "" {
		return http.StatusNoContent, ""
	}
	return http.StatusOK, responseMessage
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
	return createOkResponse(
		responseObject,
	)
}

// writeResponse responds to the consumer with corresponding HTTP status code and response body
func writeResponse(
	session *session,
	responseObject interface{},
	responseError error,
) {
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
