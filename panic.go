package webserver

func getRecoverError(recoverResult interface{}) error {
	var err, ok = recoverResult.(error)
	if !ok {
		err = fmtErrorf("Endpoint panic: %v", recoverResult)
	}
	return err
}

func getDebugStack() string {
	return string(debugStack())
}

// handlePanic prevents the application from halting when service handler panics unexpectedly
func handlePanic(
	session *session,
	recoverResult interface{},
) {
	if recoverResult == nil {
		return
	}
	var err = getRecoverErrorFunc(
		recoverResult,
	)
	writeResponseFunc(
		session,
		nil,
		err,
	)
	logAppRootFunc(
		session,
		"panic",
		"Handle",
		"%+v\n%v",
		err,
		getDebugStackFunc(),
	)
}
