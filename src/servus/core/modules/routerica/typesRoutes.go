package routerica

import "net/http"

type routericaI interface {
	NotFound(endpoint http.HandlerFunc)
	Use(middlewares ...Middleware)
	Group(prefix string) *RouteGroup
	addRequest(path string, method string, endpoint http.HandlerFunc) *RouteLocal
}

type routeGroupI interface {
	Use(middlewares ...Middleware)
	addRequest(path string, method string, endpoint http.HandlerFunc) *RouteLocal
}

type routeLocalI interface {
	Use(middlewares ...Middleware)
}

type Routerica struct {
	middlewareGlobal *middlewareGlobal
	routeNotFound    routeNotFound
	routeGroups      []RouteGroup
	routeLocals      map[string][]RouteLocal
}

type routeNotFound struct {
	endpoint http.HandlerFunc
}

type RouteGroup struct {
	middlewareGroupGlobal *middlewareGroupGlobal
	prefix                []string
	routeLocals           map[string][]RouteLocal
}

type RouteLocal struct {
	middlewareLocal *middlewareLocal
	endpoint        http.HandlerFunc
	path            []string
	method          string
}
