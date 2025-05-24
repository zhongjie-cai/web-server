package webserver

// handlePanic prevents the application from halting when service handler panics unexpectedly
func handlePanic(
	session *session,
	recoverResult any,
) {
	if recoverResult == nil {
		return
	}
	var responseObject, responseError = session.customization.RecoverPanic(
		session,
		recoverResult,
	)
	writeResponse(
		session,
		responseObject,
		responseError,
	)
}
