package webserver

import (
	"net/http"
	"time"
)

func getRequestedPort(
	httpRequest *http.Request,
) int {
	if httpRequest == nil {
		return 0
	}
	var hostAddress = httpRequest.Host
	var hostParts = stringsSplit(hostAddress, ":")
	if len(hostParts) < 2 {
		return 0
	}
	var portNumber, parseError = strconvAtoi(hostParts[1])
	if parseError != nil {
		return 0
	}
	return portNumber
}

func initiateSession(
	responseWriter http.ResponseWriter,
	httpRequest *http.Request,
) (*session, ActionFunc, error) {
	var port = getRequestedPortFunc(
		httpRequest,
	)
	var application = getApplicationFunc(
		port,
	)
	var endpoint, action, routeError = getRouteInfoFunc(
		httpRequest,
		application.actionFuncMap,
	)
	return &session{
		uuidNew(),
		endpoint,
		httpRequest,
		responseWriter,
		map[string]interface{}{},
		application.customization,
	}, action, routeError
}

func finalizeSession(
	session *session,
	startTime time.Time,
	recoverResult interface{},
) {
	handlePanicFunc(
		session,
		recoverResult,
	)
	logEndpointExitFunc(
		session,
		session.name,
		session.request.Method,
		"%s",
		timeSince(startTime),
	)
}

func handleAction(
	session *session,
	action ActionFunc,
) {
	var preActionError = session.customization.PreAction(
		session,
	)
	if preActionError != nil {
		writeResponseFunc(
			session,
			nil,
			preActionError,
		)
		return
	}
	var responseObject, responseError = action(
		session,
	)
	if responseError != nil {
		writeResponseFunc(
			session,
			responseObject,
			responseError,
		)
		return
	}
	var postActionError = session.customization.PostAction(
		session,
	)
	writeResponseFunc(
		session,
		responseObject,
		postActionError,
	)
}

// handleSession wraps the HTTP handler with session related operations
func handleSession(
	responseWriter http.ResponseWriter,
	httpRequest *http.Request,
) {
	var session, action, routeError = initiateSessionFunc(
		responseWriter,
		httpRequest,
	)
	logEndpointEnterFunc(
		session,
		session.name,
		httpRequest.Method,
		"",
	)
	defer finalizeSessionFunc(
		session,
		getTimeNowUTCFunc(),
		recover(),
	)
	if routeError != nil {
		writeResponseFunc(
			session,
			nil,
			routeError,
		)
		return
	}
	handleActionFunc(
		session,
		action,
	)
}
