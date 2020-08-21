package webserver

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func executeCustomizedFunction(
	session Session,
	customFunc func(Session) error,
) error {
	if customFunc == nil {
		return nil
	}
	return customFunc(
		session,
	)
}

// handleSession wraps the HTTP handler with session related operations
func handleSession(
	writeResponser http.ResponseWriter,
	httpRequest *http.Request,
) {
	var endpoint, action, routeError = getRouteInfo(
		httpRequest,
	)
	var customization = &defaultCustomization{} // TODO: where to get this from app??
	var session = &session{
		uuid.New(),
		endpoint,
		httpRequest,
		writeResponser,
		map[string]interface{}{},
		customization,
	}
	var startTime = getTimeNowUTC()
	apiEnter(
		session,
		endpoint,
		httpRequest.Method,
		"",
	)
	defer func() {
		handlePanic(
			session,
			customization,
			recover(),
		)
		apiExit(
			session,
			endpoint,
			httpRequest.Method,
			"%s",
			time.Since(startTime),
		)
	}()
	if routeError != nil {
		writeResponse(
			session,
			customization,
			nil,
			routeError,
		)
	} else {
		var preActionError = session.customization.PreAction(
			session,
		)
		if preActionError != nil {
			writeResponse(
				session,
				customization,
				nil,
				preActionError,
			)
		} else {
			var responseObject, responseError = action(
				session,
			)
			var postActionError = session.customization.PostAction(
				session,
			)
			if postActionError != nil {
				if responseError != nil {
					apiExit(
						session,
						endpoint,
						httpRequest.Method,
						"Post-action error: %v",
						postActionError,
					)
					writeResponse(
						session,
						nil,
						responseError,
					)
				} else {
					writeResponse(
						session,
						customization,
						nil,
						postActionError,
					)
				}
			} else {
				writeResponse(
					session,
					customization,
					responseObject,
					responseError,
				)
			}
		}
	}
}
