package routerica

import "strings"

func parseRouteLocal(routes  map[string][]RouteLocal, requestMethod string, requestURI string) *RouteLocal{
	// routes without group
	requestMethod = strings.ToUpper(requestMethod)
	requestURI = strings.ToUpper(requestURI) + "/"
	for index := range routes[requestMethod] {
		var currentRoute = &routes[requestMethod][index]
		var routePath = strings.ToUpper(currentRoute.path)
		if routePath != requestURI {
			continue
		}
		// now we in correct route
		return currentRoute
	}
	return nil
}

type routeLocalI interface {
	// Use - add local middleware for a specific route
	Use(middlewares ...Middleware)
}

type RouteLocal struct {
	RouteBase
	routeLocalI
	path   string
	method string
}

// Use middleware (request methods). That method executes when you're adding middleware to request method
func (r *RouteLocal) Use(middlewares ...Middleware) {
	r.middlewareChain = middlewareChainer(middlewares, r.handler)
}
