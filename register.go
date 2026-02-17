package webserver

import (
	"fmt"
	"strings"

	"github.com/go-chi/chi/v5"
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

func registerRoutes(
	app *application,
	session *session,
	router chi.Router,
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
		var name = generateRouteName(
			configuredRoute.Method,
			evaluatedPath,
		)
		router.MethodFunc(
			configuredRoute.Method,
			evaluatedPath,
			app.handleSession,
		)
		app.actionFuncMap[name] = configuredRoute.ActionFunc
	}
}

func registerStatics(
	session *session,
	router chi.Router,
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
		router.Handle(
			static.PathPrefix,
			static.Handler,
		)
	}
}

func registerMiddlewares(
	session *session,
	router chi.Router,
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
		router.Use(middleware)
	}
}

func registerErrorHandlers(
	customization Customization,
	router chi.Router,
) {
	var notAllowedHandler = customization.MethodNotAllowedHandler()
	if notAllowedHandler != nil {
		router.MethodNotAllowed(notAllowedHandler)
	}
	var notFoundHandler = customization.NotFoundHandler()
	if notFoundHandler != nil {
		router.NotFound(notFoundHandler)
	}
}

// instantiateRouter instantiates and registers the given routes according to custom specification
func instantiateRouter(
	app *application,
	session *session,
) (chi.Router, error) {
	var router chi.Router
	router = chi.NewRouter()
	registerMiddlewares(
		session,
		router,
	)
	registerRoutes(
		app,
		session,
		router,
	)
	registerStatics(
		session,
		router,
	)
	registerErrorHandlers(
		session.customization,
		router,
	)
	router = session.customization.InstrumentRouter(
		router,
	)
	var err = walkRegisteredRoutes(
		session,
		router,
	)
	if err != nil {
		logAppRoot(
			session,
			"register",
			"instantiateRouter",
			"%+v",
			err,
		)
		return router, newAppError(
			errorCodeGeneralFailure,
			errorMessageRouteRegistration,
			err,
		)
	}
	return router, err
}
