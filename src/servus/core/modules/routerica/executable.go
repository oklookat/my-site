package routerica

import (
	"fmt"
	"net/http"
	"strings"
)

// NotFound - add 404 page. When no suitable route is found, the endpoint is called.
func (r *Routerica) NotFound(endpoint http.HandlerFunc) {
	r.routeNotFound.endpoint = endpoint
}

// Group - create group of routes.
func (r *Routerica) Group(prefix string) *RouteGroup {
	var routeGroup = newRouteGroup(prefix)
	r.routeGroups = append(r.routeGroups, routeGroup)
	var lastGroup = len(r.routeGroups) - 1
	return &r.routeGroups[lastGroup]
}

// addRequest - add request to router.
func (r *Routerica) addRequest(path string, method string, endpoint http.HandlerFunc) *RouteLocal {
	path = formatPath(path)
	var theRoute = newRouteLocal(path, method, endpoint)
	r.routeLocals[method] = append(r.routeLocals[method], theRoute)
	var lastRoute = len(r.routeLocals[method]) - 1
	return &r.routeLocals[method][lastRoute]
}

// addRequest - add request to route group.
func (g *RouteGroup) addRequest(path string, method string, endpoint http.HandlerFunc) *RouteLocal {
	var pathFormatted = fmt.Sprintf("/%v/%v/", strings.Join(g.prefix, "/"), path)
	pathFormatted = formatPath(pathFormatted)
	var theRoute = newRouteLocal(pathFormatted, method, endpoint)
	g.routeLocals[method] = append(g.routeLocals[method], theRoute)
	var lastRoute = len(g.routeLocals[method]) - 1
	return &g.routeLocals[method][lastRoute]
}
