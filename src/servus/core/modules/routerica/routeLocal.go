package routerica

type routeLocalI interface {
	// Use - add local middleware for a specific route
	Use(middlewares ...Middleware)
}

type RouteLocal struct {
	RouteBase
	routeLocalI
	path   []string
	method string
}

// parseRouteLocal - get route local from routes, request URI and method.
func parseRouteLocal(routes map[string][]RouteLocal, requestMethod string, requestURI []string) *RouteLocal {
	for index := range routes[requestMethod] {
		var currentRoute = &routes[requestMethod][index]
		var isMatched, _ = uriParser(currentRoute.path, requestURI)
		if !isMatched {
			continue
		}
		// now we in correct route
		return currentRoute
	}
	return nil
}

// Use middleware (request methods). That method executes when you're adding middleware to request method
func (r *RouteLocal) Use(middlewares ...Middleware) {
	r.middlewareChain = middlewareChainer(middlewares, r.handler)
}
