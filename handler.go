package webserver

import (
	"net/http"
	"time"
)

func initiateSession(
	app *application,
	responseWriter http.ResponseWriter,
	httpRequest *http.Request,
) (*session, ActionFunc, error) {
	var endpoint, action, routeError = getRouteInfoFunc(
		httpRequest,
		app.actionFuncMap,
	)
	return &session{
		uuidNew(),
		endpoint,
		httpRequest,
		responseWriter,
		map[string]interface{}{},
		app.customization,
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
func (app *application) handleSession(
	responseWriter http.ResponseWriter,
	httpRequest *http.Request,
) {
	var session, action, routeError = initiateSessionFunc(
		app,
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
