package webserver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Static holds the registration information of a static content hosting
type Static struct {
	PathPrefix string
	Handler    http.Handler
}

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
	if handler == nil {
		return fmt.Errorf(
			"Invalid handler for %v:%v",
			method,
			route,
		)
	}
	for index, middleware := range middlewares {
		if middleware == nil {
			return fmt.Errorf(
				"Invalid middleware for %v:%v @ #%d",
				method,
				route,
				index+1,
			)
		}
	}
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

func extractRouteMethodAndPattern(name string) (string, string) {
	var parts = strings.Split(name, ":")
	if len(parts) < 2 {
		return "??", name
	}
	if len(parts) > 2 {
		return parts[0], strings.Join(parts[1:], ":")
	}
	return parts[0], parts[1]
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
				httpRequest.RequestURI,
			),
		)
}
