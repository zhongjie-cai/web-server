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
		logAppRoot(
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
	app *application,
	session *session,
	router *mux.Router,
) {
	var configuredRoutes = session.customization.Routes()
	if len(configuredRoutes) == 0 {
		logAppRoot(
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
		var name, _ = registerRoute(
			router,
			configuredRoute.Endpoint,
			configuredRoute.Method,
			evaluatedPath,
			queries,
			app.handleSession,
		)
		app.actionFuncMap[name] = configuredRoute.ActionFunc
	}
}

func registerStatics(
	session *session,
	router *mux.Router,
) {
	var statics = session.customization.Statics()
	if len(statics) == 0 {
		logAppRoot(
			session,
			"register",
			"registerStatics",
			"customization.Statics function empty: no static content returned!",
		)
		return
	}
	for _, static := range statics {
		registerStatic(
			router,
			static.Name,
			static.PathPrefix,
			static.Handler,
		)
	}
}

func registerMiddlewares(
	session *session,
	router *mux.Router,
) {
	var middlewares = session.customization.Middlewares()
	if len(middlewares) == 0 {
		logAppRoot(
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

// instantiateRouter instantiates and registers the given routes according to custom specification
func instantiateRouter(
	app *application,
	session *session,
) (*mux.Router, error) {
	var router = mux.NewRouter()
	registerRoutes(
		app,
		session,
		router,
	)
	registerStatics(
		session,
		router,
	)
	registerMiddlewares(
		session,
		router,
	)
	var routerError = walkRegisteredRoutes(
		session,
		router,
	)
	if routerError != nil {
		logAppRoot(
			session,
			"register",
			"instantiateRouter",
			"%+v",
			routerError,
		)
		return router,
			newAppError(
				errorCodeGeneralFailure,
				errorMessageRouteRegistration,
				[]error{routerError},
			)
	}
	registerErrorHandlers(
		session.customization,
		router,
	)
	return session.customization.InstrumentRouter(
		router,
	), nil
}
