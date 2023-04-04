package webserver

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route holds the registration information of a dynamic route hosting
type Route struct {
	Endpoint   string
	Method     string
	Path       string
	Parameters map[string]ParameterType
	Queries    map[string]ParameterType
	ActionFunc ActionFunc
}

const (
	stringSeparator string = "|"
)

func getName(route *mux.Route) string {
	return route.GetName()
}

func getPathTemplate(route *mux.Route) (string, error) {
	return route.GetPathTemplate()
}

func getPathRegexp(route *mux.Route) (string, error) {
	return route.GetPathRegexp()
}

func getQueriesTemplates(route *mux.Route) string {
	var queriesTemplates, _ = route.GetQueriesTemplates()
	return stringsJoin(queriesTemplates, stringSeparator)
}

func getQueriesRegexp(route *mux.Route) string {
	var queriesRegexps, _ = route.GetQueriesRegexp()
	return stringsJoin(queriesRegexps, stringSeparator)
}

func getMethods(route *mux.Route) string {
	var methods, _ = route.GetMethods()
	return stringsJoin(methods, stringSeparator)
}

func evaluateRoute(
	route *mux.Route,
	router *mux.Router,
	ancestors []*mux.Route,
) error {
	var (
		_, pathTemplateError = getPathTemplateFunc(route)
		_, pathRegexpError   = getPathRegexpFunc(route)
	)
	if pathTemplateError != nil {
		return pathTemplateError
	}
	if pathRegexpError != nil {
		return pathRegexpError
	}
	return nil
}

// walkRegisteredRoutes examines the registered router for errors
func walkRegisteredRoutes(
	session *session,
	router *mux.Router,
) error {
	var walkError = router.Walk(
		evaluateRouteFunc,
	)
	if walkError != nil {
		logAppRootFunc(
			session,
			"route",
			"walkRegisteredRoutes",
			"Failure: %+v",
			walkError,
		)
		return newAppErrorFunc(
			errorCodeGeneralFailure,
			errorMessageRouteRegistration,
			[]error{walkError},
		)
	}
	return nil
}

// registerRoute wraps the mux route handler
func registerRoute(
	router *mux.Router,
	endpoint string,
	method string,
	path string,
	queries []string,
	handleFunc func(http.ResponseWriter, *http.Request),
	actionFunc ActionFunc,
) (string, *mux.Route) {
	var name = fmtSprintf(
		"%v:%v",
		endpoint,
		method,
	)
	var route = router.HandleFunc(
		path,
		handleFunc,
	).Methods(
		method,
	).Queries(
		queries...,
	).Name(
		name,
	)
	return name, route
}

func defaultActionFunc(session Session) (interface{}, error) {
	return nil,
		newAppErrorFunc(
			errorCodeNotImplemented,
			"No corresponding action function configured; falling back to default",
			[]error{},
		)
}

func getEndpointByName(name string) string {
	var splitSubs = stringsSplit(
		name,
		":",
	)
	if len(splitSubs) < 2 {
		return name
	}
	return splitSubs[0]
}

// getRouteInfo retrieves the registered name and action for the given route
func getRouteInfo(httpRequest *http.Request, actionFuncMap map[string]ActionFunc) (string, ActionFunc, error) {
	var route = muxCurrentRoute(httpRequest)
	if route == nil {
		return "",
			nil,
			newAppErrorFunc(
				errorCodeNotFound,
				"No corresponding route configured for path",
				[]error{},
			)
	}
	var name = getNameFunc(route)
	var endpoint = getEndpointByNameFunc(name)
	var action, found = actionFuncMap[name]
	if !found {
		action = defaultActionFunc
	}
	return endpoint, action, nil
}
