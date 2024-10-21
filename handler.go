package webserver

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func initiateSession(
	app *application,
	responseWriter http.ResponseWriter,
	httpRequest *http.Request,
) (*session, ActionFunc, error) {
	var endpoint, action, routeError = getRouteInfo(
		httpRequest,
		app.actionFuncMap,
	)
	return &session{
		uuid.New(),
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
	handlePanic(
		session,
		recoverResult,
	)
	logEndpointExit(
		session,
		session.name,
		session.request.Method,
		"%s",
		time.Since(startTime),
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
		writeResponse(
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
		writeResponse(
			session,
			responseObject,
			responseError,
		)
		return
	}
	var postActionError = session.customization.PostAction(
		session,
	)
	writeResponse(
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
	var session, action, routeError = initiateSession(
		app,
		responseWriter,
		httpRequest,
	)
	logEndpointEnter(
		session,
		session.name,
		httpRequest.Method,
		"",
	)
	defer finalizeSession(
		session,
		getTimeNowUTC(),
		recover(),
	)
	if routeError != nil {
		writeResponse(
			session,
			nil,
			routeError,
		)
		return
	}
	handleAction(
		session,
		action,
	)
}
