package webserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Route holds the registration information of a dynamic route hosting
type Route struct {
	Method     string
	Path       string
	Parameters map[string]ParameterType
	ActionFunc ActionFunc
}

func evaluateRoute(
	method string,
	route string,
	handler http.Handler,
	middlewares ...func(http.Handler) http.Handler,
) error {
	// TODO: how and what to walk
	return nil
}

// walkRegisteredRoutes examines the registered router for errors
func walkRegisteredRoutes(
	session *session,
	router chi.Router,
) error {
	var walkError = chi.Walk(
		router,
		evaluateRoute,
	)
	if walkError != nil {
		logAppRoot(
			session,
			"route",
			"walkRegisteredRoutes",
			"Failure: %+v",
			walkError,
		)
		return newAppError(
			errorCodeGeneralFailure,
			errorMessageRouteRegistration,
			walkError,
		)
	}
	return nil
}

func generateRouteName(method string, pattern string) string {
	return fmt.Sprintf(
		"%v:%v",
		method,
		pattern,
	)
}

func defaultActionFunc(session Session) (any, error) {
	return nil,
		newAppError(
			errorCodeNotImplemented,
			"No corresponding action function configured; falling back to default",
		)
}

// getRouteInfo retrieves the registered name and action for the given route
func getRouteInfo(httpRequest *http.Request, actionFuncMap map[string]ActionFunc) (string, ActionFunc, error) {
	var ctx = chi.RouteContext(httpRequest.Context())
	if ctx == nil {
		return "",
			nil,
			newAppError(
				errorCodeDataCorruption,
				"No go-chi context found in HTTP request",
			)
	}
	for _, routePattern := range ctx.RoutePatterns {
		var name = generateRouteName(
			ctx.RouteMethod,
			routePattern,
		)
		var action, found = actionFuncMap[name]
		if found {
			return name, action, nil
		}
	}
	return "",
		nil,
		newAppError(
			errorCodeNotFound,
			fmt.Sprintf(
				"No corresponding route configured for path: %v",
				httpRequest.URL.String(),
			),
		)
}
