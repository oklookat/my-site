package routerica

import (
	"net/http"
)

// New - create new Routerica instance
func New() *Routerica {
	var localRoutes = make(map[string][]RouteLocal, 0)
	var routerica = &Routerica{localRoutes: localRoutes}
	routerica.globalMiddleware = &GlobalMiddleware{}
	routerica.routeGroups = make([]RouteGroup, 0)
	routerica.globalMiddleware.middlewareChain = routerica.globalMiddleware
	return routerica
}

// newRouteGroup - create route group instance (internal)
func newRouteGroup(prefix string) RouteGroup {
	var localRoutes = map[string][]RouteLocal{}
	prefixSlice := uriSplitter(prefix)
	var routeGroup = RouteGroup{prefix: prefixSlice, localRoutes: localRoutes}
	routeGroup.middlewareChain = &routeGroup.RouteBase
	return routeGroup
}

// newRouteLocal - create route local instance (internal)
func newRouteLocal(path string, method string, handler http.HandlerFunc) RouteLocal {
	pathSlice := uriSplitter(path)
	return RouteLocal{path: pathSlice, method: method, RouteBase: RouteBase{handler: handler}}
}
