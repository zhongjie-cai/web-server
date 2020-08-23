package webserver

import (
	"net/http"
	"strconv"
	"strings"
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

func getRequestedPort(
	httpRequest *http.Request,
) int {
	if httpRequest == nil {
		return 0
	}
	var hostAddress = httpRequest.Host
	var hostParts = strings.Split(hostAddress, ":")
	if len(hostParts) < 2 {
		return 0
	}
	var portNumber, parseError = strconv.Atoi(hostParts[1])
	if parseError != nil {
		return 0
	}
	return portNumber
}

// handleSession wraps the HTTP handler with session related operations
func handleSession(
	writeResponser http.ResponseWriter,
	httpRequest *http.Request,
) {
	var port = getRequestedPort(
		httpRequest,
	)
	var application = getApplication(
		port,
	)
	var endpoint, action, routeError = getRouteInfo(
		httpRequest,
		application.actionFuncMap,
	)
	var customization = application.customization
	var session = &session{
		uuid.New(),
		endpoint,
		httpRequest,
		writeResponser,
		map[string]interface{}{},
		customization,
	}
	var startTime = getTimeNowUTC()
	logEndpointEnter(
		session,
		endpoint,
		httpRequest.Method,
		"",
	)
	defer func() {
		handlePanic(
			session,
			recover(),
		)
		logEndpointExit(
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
			nil,
			routeError,
		)
	} else {
		var preActionError = customization.PreAction(
			session,
		)
		if preActionError != nil {
			writeResponse(
				session,
				nil,
				preActionError,
			)
		} else {
			var responseObject, responseError = action(
				session,
			)
			var postActionError = customization.PostAction(
				session,
			)
			if postActionError != nil {
				if responseError != nil {
					logEndpointExit(
						session,
						endpoint,
						httpRequest.Method,
						"Post-action error: %+v",
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
						nil,
						postActionError,
					)
				}
			} else {
				writeResponse(
					session,
					responseObject,
					responseError,
				)
			}
		}
	}
}
