package webserver

// handlePanic prevents the application from halting when service handler panics unexpectedly
func handlePanic(
	session *session,
	recoverResult interface{},
) {
	if recoverResult == nil {
		return
	}
	var responseObject, responseError = session.customization.RecoverPanic(
		session,
		recoverResult,
	)
	writeResponseFunc(
		session,
		responseObject,
		responseError,
	)
}
