package webserver

// These are the constants used by the HTTP modules
const (
	ContentTypeJSON = "application/json; charset=utf-8"
)

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
	var statusCode, responseMessage = constructResponseFunc(
		session,
		responseObject,
		responseError,
	)
	logEndpointResponseFunc(
		session,
		httpStatusText(statusCode),
		strconvItoa(statusCode),
		responseMessage,
	)
	var responseWriter = session.GetResponseWriter()
	responseWriter.Header().Set("Content-Type", ContentTypeJSON)
	responseWriter.WriteHeader(statusCode)
	responseWriter.Write([]byte(responseMessage))
}
