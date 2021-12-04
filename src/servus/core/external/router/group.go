package router

import (
	"fmt"
	"net/http"
	"strings"
)

// Group - group of routes.
type Group struct {
	middleware groupMiddleware
	prefix     []string
	routes     routes
}

// middlewareGroupGlobal - used only for route group (global middleware).
type groupMiddleware struct {
	resume   bool
	chain    MiddlewareChain
}

// new - create route group.
func (g *Group) new(prefix string) {
	prefixSlice := pathToSlice(prefix)
	g.prefix = prefixSlice
	g.routes = make(routes, 0)
	g.middleware = groupMiddleware{}
	g.middleware.chain = &g.middleware
}

// addRoute - add route to route group.
func (g *Group) addRoute(path string, method string, endpoint http.HandlerFunc) *RouteMethod {
	var prefix = strings.Join(g.prefix, "/")
	var pathFormatted = fmt.Sprintf("/%v/%v/", prefix, path)
	var rMethod = addRoute(g.routes, pathFormatted, method, endpoint)
	return rMethod
}

// GET - add GET request to the route group.
func (g *Group) GET(path string, endpoint http.HandlerFunc) *RouteMethod {
	return g.addRoute(path, methodGET, endpoint)
}

// HEAD - add HEAD request to the route group.
func (g *Group) HEAD(path string, endpoint http.HandlerFunc) *RouteMethod {
	return g.addRoute(path, methodHEAD, endpoint)
}

// POST - add POST request to the route group.
func (g *Group) POST(path string, endpoint http.HandlerFunc) *RouteMethod {
	return g.addRoute(path, methodPOST, endpoint)
}

// PUT - add PUT request to the route group.
func (g *Group) PUT(path string, endpoint http.HandlerFunc) *RouteMethod {
	return g.addRoute(path, methodPUT, endpoint)
}

// DELETE - add DELETE request to the route group.
func (g *Group) DELETE(path string, endpoint http.HandlerFunc) *RouteMethod {
	return g.addRoute(path, methodDELETE, endpoint)
}

// OPTIONS - add OPTIONS request to the route group.
func (g *Group) OPTIONS(path string, endpoint http.HandlerFunc) *RouteMethod {
	return g.addRoute(path, methodOPTIONS, endpoint)
}

// PATCH - add PATCH request to the route group.
func (g *Group) PATCH(path string, endpoint http.HandlerFunc) *RouteMethod {
	return g.addRoute(path, methodPATCH, endpoint)
}

// Use - add group global middlewares. Any request on this route group will call these middlewares.
func (g *Group) Use(middlewares ...Middleware) *Group {
	g.middleware.chain = middlewareChainer(middlewares, &g.middleware)
	return g
}

func (g *groupMiddleware) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	g.resume = true
}

// run - run middleware if exists. Returns true if middleware called next().
func (g *Group) run(response http.ResponseWriter, request *http.Request) (resume bool) {
	if g.middleware.chain != nil {
		g.middleware.chain.ServeHTTP(response, request)
		return g.middleware.resume
	}
	return true
}
