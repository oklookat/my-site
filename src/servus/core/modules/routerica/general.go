package routerica

import (
	"net/http"
)

// New - create new Routerica instance
func New() *Routerica {
	var routeLocals = make(map[string][]RouteLocal, 0)
	var routerica = &Routerica{routeLocals: routeLocals}
	routerica.routeGroups = make([]RouteGroup, 0)
	routerica.middlewareGlobal = &middlewareGlobal{}
	routerica.middlewareGlobal.chain = routerica.middlewareGlobal
	routerica.routeNotFound.endpoint = defaultsEndpointNotFound
	return routerica
}

// newRouteGroup - create route group instance (internal)
func newRouteGroup(prefix string) RouteGroup {
	var routeLocals = map[string][]RouteLocal{}
	prefixSlice := uriSplitter(prefix)
	var routeGroup = RouteGroup{prefix: prefixSlice, routeLocals: routeLocals}
	routeGroup.middlewareGroupGlobal = &middlewareGroupGlobal{}
	// its works because middlewareGroupGlobal implements ServeHTTP
	routeGroup.middlewareGroupGlobal.chain = routeGroup.middlewareGroupGlobal
	return routeGroup
}

// newRouteLocal - create route local instance (internal)
func newRouteLocal(path string, method string, handler http.HandlerFunc) RouteLocal {
	pathSlice := uriSplitter(path)
	return RouteLocal{path: pathSlice, method: method, endpoint: handler, middlewareLocal: &middlewareLocal{}}
}
