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
	requestsI
	routericaI
	middlewareGlobal *middlewareGlobal
	routeNotFound    routeNotFound
	routeGroups      []RouteGroup
	routeLocals      map[string][]RouteLocal
}

type routeNotFound struct {
	endpoint http.HandlerFunc
}

type RouteGroup struct {
	requestsI
	routeGroupI
	middlewareGroupGlobal *middlewareGroupGlobal
	prefix                []string
	routeLocals           map[string][]RouteLocal
}

type RouteLocal struct {
	routeLocalI
	middlewareLocal *middlewareLocal
	endpoint        http.HandlerFunc
	path            []string
	method          string
}
