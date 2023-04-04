package webserver

import (
	"github.com/gorilla/mux"
)

func doParameterReplacement(
	session *session,
	originalPath string,
	parameterName string,
	parameterType ParameterType,
) string {
	if parameterType == "" {
		logAppRootFunc(
			session,
			"register",
			"doParameterReplacement",
			"Path parameter [%v] in path [%v] has no type specification; fallback to default.",
			parameterName,
			originalPath,
		)
		return originalPath
	}
	var oldParameter = fmtSprintf(
		"{%v}",
		parameterName,
	)
	var newParameter = fmtSprintf(
		"{%v:%v}",
		parameterName,
		parameterType,
	)
	return stringsReplace(
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
		updatedPath = doParameterReplacementFunc(
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
			queryParameter = fmtSprintf(
				"{%v}",
				key,
			)
		} else {
			queryParameter = fmtSprintf(
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
		logAppRootFunc(
			session,
			"register",
			"registerRoutes",
			"customization.Routes function empty: no routes returned!",
		)
		return
	}
	for _, configuredRoute := range configuredRoutes {
		var evaluatedPath = evaluatePathWithParametersFunc(
			session,
			configuredRoute.Path,
			configuredRoute.Parameters,
		)
		var queries = evaluateQueriesFunc(
			configuredRoute.Queries,
		)
		var name, _ = registerRouteFunc(
			router,
			configuredRoute.Endpoint,
			configuredRoute.Method,
			evaluatedPath,
			queries,
			app.handleSession,
			configuredRoute.ActionFunc,
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
		logAppRootFunc(
			session,
			"register",
			"registerStatics",
			"customization.Statics function empty: no static content returned!",
		)
		return
	}
	for _, static := range statics {
		registerStaticFunc(
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
		logAppRootFunc(
			session,
			"register",
			"registerMiddlewares",
			"customization.Middlewares function empty: no middleware returned!",
		)
		return
	}
	for _, middleware := range middlewares {
		addMiddlewareFunc(
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
	var router = muxNewRouter()
	registerRoutesFunc(
		app,
		session,
		router,
	)
	registerStaticsFunc(
		session,
		router,
	)
	registerMiddlewaresFunc(
		session,
		router,
	)
	var routerError = walkRegisteredRoutesFunc(
		session,
		router,
	)
	if routerError != nil {
		logAppRootFunc(
			session,
			"register",
			"instantiateRouter",
			"%+v",
			routerError,
		)
		return router,
			newAppErrorFunc(
				errorCodeGeneralFailure,
				errorMessageRouteRegistration,
				[]error{routerError},
			)
	}
	registerErrorHandlersFunc(
		session.customization,
		router,
	)
	return session.customization.InstrumentRouter(
		router,
	), nil
}
