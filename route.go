package webserver

import (
	"fmt"
	"net/http"
	"strings"

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
	return strings.Join(queriesTemplates, stringSeparator)
}

func getQueriesRegexp(route *mux.Route) string {
	var queriesRegexps, _ = route.GetQueriesRegexp()
	return strings.Join(queriesRegexps, stringSeparator)
}

func getMethods(route *mux.Route) string {
	var methods, _ = route.GetMethods()
	return strings.Join(methods, stringSeparator)
}

func printRegisteredRouteDetails(
	route *mux.Route,
	router *mux.Router,
	ancestors []*mux.Route,
) error {
	var (
		_, pathTemplateError = getPathTemplate(route)
		_, pathRegexpError   = getPathRegexp(route)
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
		printRegisteredRouteDetails,
	)
	if walkError != nil {
		appRoot(
			session,
			"route",
			"walkRegisteredRoutes",
			"Failure: %+v",
			walkError,
		)
		return errRouteRegistion
	}
	return nil
}

// createRouter initializes a router for route registrations
func createRouter() *mux.Router {
	return mux.NewRouter()
}

// handleFunc wraps the mux route handler
func handleFunc(
	router *mux.Router,
	endpoint string,
	method string,
	path string,
	queries []string,
	handleFunc func(http.ResponseWriter, *http.Request),
	actionFunc ActionFunc,
	port int,
) *mux.Route {
	var name = fmt.Sprintf(
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
	var application = getApplication(port)
	application.actionFuncMap[name] = actionFunc
	return route
}

func defaultActionFunc(session Session) (interface{}, error) {
	return nil, nil
}

func getEndpointByName(name string) string {
	var splitSubs = strings.Split(
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
	var route = mux.CurrentRoute(httpRequest)
	if route == nil {
		return "", nil, errRouteNotFound
	}
	var name = getName(route)
	var endpoint = getEndpointByName(name)
	var action, found = actionFuncMap[name]
	if !found {
		action = defaultActionFunc
	}
	return endpoint, action, nil
}
