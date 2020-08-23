package webserver

import (
	"fmt"
	"runtime/debug"
)

func getRecoverError(recoverResult interface{}) error {
	var err, ok = recoverResult.(error)
	if !ok {
		err = fmt.Errorf("Endpoint panic: %v", recoverResult)
	}
	return err
}

func getDebugStack() string {
	return string(debug.Stack())
}

// handlePanic prevents the application from halting when service handler panics unexpectedly
func handlePanic(
	session *session,
	recoverResult interface{},
) {
	if recoverResult != nil {
		var appError = getRecoverError(
			recoverResult,
		)
		writeResponse(
			session,
			nil,
			appError,
		)
		logAppRoot(
			session,
			"panic",
			"Handle",
			"%+v\n%v",
			appError,
			getDebugStack(),
		)
	}
}
