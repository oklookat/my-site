package routerica

import (
	"net/http"
)

type routericaI interface {
	// Use - add global middleware to all routes
	Use(middlewares ...Middleware)
	Group(prefix string) *RouteGroup
}

type Routerica struct {
	requestsI
	routericaI
	globalMiddleware *GlobalMiddleware
	routeGroups      []RouteGroup
	localRoutes      map[string][]RouteLocal
}

// Group - create group of routes.
func (r *Routerica) Group(prefix string) *RouteGroup {
	prefix = formatPath(prefix)
	var routeGroup = newRouteGroup(prefix)
	r.routeGroups = append(r.routeGroups, routeGroup)
	var lastGroup = len(r.routeGroups) - 1
	return &r.routeGroups[lastGroup]
}

// Use - add global middlewares. Any request will call these middlewares.
func (r *Routerica) Use(middlewares ...Middleware) {
	r.globalMiddleware.middlewareChain = middlewareChainer(middlewares, r.globalMiddleware)
}

// GET - add GET request.
func (r *Routerica) GET(path string, handler func(http.ResponseWriter, *http.Request)) *RouteLocal {
	path = formatPath(path)
	var theRoute = newRouteLocal(path, methodGET, handler)
	r.localRoutes[methodGET] = append(r.localRoutes[methodGET], theRoute)
	var lastRoute = len(r.localRoutes[methodGET]) - 1
	return &r.localRoutes[methodGET][lastRoute]
}
