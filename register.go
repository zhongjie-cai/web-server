package webserver

import (
	"fmt"
	"strings"

	"github.com/gorilla/mux"
)

func doParameterReplacement(
	session *session,
	originalPath string,
	parameterName string,
	parameterType ParameterType,
) string {
	if parameterType == "" {
		appRoot(
			session,
			"register",
			"doParameterReplacement",
			"Path parameter [%v] in path [%v] has no type specification; fallback to default.",
			parameterName,
			originalPath,
		)
		return originalPath
	}
	var oldParameter = fmt.Sprintf(
		"{%v}",
		parameterName,
	)
	var newParameter = fmt.Sprintf(
		"{%v:%v}",
		parameterName,
		parameterType,
	)
	return strings.Replace(
		originalPath,
		oldParameter,
		newParameter,
		-1,
	)
}

func evaluatePathWithParameters(
	session *session,
	path string,
	parameters map[string]ParameterType,
) string {
	var updatedPath = path
	for parameterName, parameterType := range parameters {
		updatedPath = doParameterReplacement(
			session,
			updatedPath,
			parameterName,
			parameterType,
		)
	}
	return updatedPath
}

func evaluateQueries(
	queries map[string]ParameterType,
) []string {
	var evaluatedQueries = []string{}
	for key, value := range queries {
		var queryParameter string
		if value == "" {
			queryParameter = fmt.Sprintf(
				"{%v}",
				key,
			)
		} else {
			queryParameter = fmt.Sprintf(
				"{%v:%v}",
				key,
				value,
			)
		}
		evaluatedQueries = append(
			evaluatedQueries,
			key,
			queryParameter,
		)
	}
	return evaluatedQueries
}

func registerRoutes(
	session *session,
	customization Customization,
	router *mux.Router,
) {
	var configuredRoutes = customization.Routes()
	if configuredRoutes == nil ||
		len(configuredRoutes) == 0 {
		appRoot(
			session,
			"register",
			"registerRoutes",
			"customization.Routes function empty: no routes returned!",
		)
		return
	}
	for _, configuredRoute := range configuredRoutes {
		var evaluatedPath = evaluatePathWithParameters(
			session,
			configuredRoute.Path,
			configuredRoute.Parameters,
		)
		var queries = evaluateQueries(
			configuredRoute.Queries,
		)
		handleFunc(
			router,
			configuredRoute.Endpoint,
			configuredRoute.Method,
			evaluatedPath,
			queries,
			handleSession,
			configuredRoute.ActionFunc,
		)
	}
}

func registerStatics(
	session *session,
	customization Customization,
	router *mux.Router,
) {
	var statics = customization.Statics()
	if statics == nil ||
		len(statics) == 0 {
		appRoot(
			session,
			"register",
			"registerStatics",
			"customization.Statics function empty: no static content returned!",
		)
		return
	}
	for _, static := range statics {
		hostStatic(
			router,
			static.Name,
			static.PathPrefix,
			static.Handler,
		)
	}
}

func registerMiddlewares(
	session *session,
	customization Customization,
	router *mux.Router,
) {
	var middlewares = customization.Middlewares()
	if middlewares == nil ||
		len(middlewares) == 0 {
		appRoot(
			session,
			"register",
			"registerMiddlewares",
			"customization.Middlewares function empty: no middleware returned!",
		)
		return
	}
	for _, middleware := range middlewares {
		addMiddleware(
			router,
			middleware,
		)
	}
}

func registerErrorHandlers(
	customization Customization,
	router *mux.Router,
) {
	router.MethodNotAllowedHandler = customization.MethodNotAllowedHandler()
	router.NotFoundHandler = customization.NotFoundHandler()
}

func instrumentRouter(
	customization Customization,
	router *mux.Router,
) *mux.Router {
	return customization.InstrumentRouter(router)
}

// instantiateRouter instantiates and registers the given routes according to custom specification
func instantiateRouter(
	session *session,
	customization Customization,
) (*mux.Router, error) {
	var router = createRouter()
	registerRoutes(
		session,
		customization,
		router,
	)
	registerStatics(
		session,
		customization,
		router,
	)
	registerMiddlewares(
		session,
		customization,
		router,
	)
	var routerError = walkRegisteredRoutes(
		session,
		router,
	)
	if routerError != nil {
		return nil, errRouteRegistion
	}
	registerErrorHandlers(
		customization,
		router,
	)
	return instrumentRouter(
		customization,
		router,
	), nil
}
