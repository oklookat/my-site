package router

import (
	"context"
	"net/http"
	"strings"
)

type Router struct {
	middleware routerMiddleware
	notFound   http.HandlerFunc
	methodNotAllowed http.HandlerFunc
	groups     []Group
	routes     routes
}

// routerMiddleware - used only for router instance (global middleware).
type routerMiddleware struct {
	// resume - need for check, is all middleware executed, or one of middleware send response.
	resume bool
	chain  MiddlewareChain
}

// New - create new Router instance.
func New() *Router {
	var router = &Router{}
	router.groups = make([]Group, 0)
	router.routes = make(routes, 0)
	router.middleware = routerMiddleware{}
	router.middleware.chain = &router.middleware
	router.notFound = defaultNotFound
	router.methodNotAllowed = defaultNotAllowed
	return router
}

// NotFound - add 404 page. When no suitable route is found, the endpoint is called.
func (r *Router) NotFound(endpoint http.HandlerFunc) {
	r.notFound = endpoint
}

// MethodNotAllowed - add 405 page. When no method for route, the endpoint is called.
func (r *Router) MethodNotAllowed(endpoint http.HandlerFunc) {
	r.methodNotAllowed = endpoint
}

// Group - create group of routes.
func (r *Router) Group(prefix string) *Group {
	var group = Group{}
	group.new(prefix)
	r.groups = append(r.groups, group)
	var lastGroup = len(r.groups) - 1
	return &r.groups[lastGroup]
}

// addRoute - add route to router.
func (r *Router) addRoute(path string, method string, endpoint http.HandlerFunc) *RouteMethod {
	return addRoute(r.routes, path, method, endpoint)
}

// GET - add GET request to router.
func (r *Router) GET(path string, endpoint http.HandlerFunc) *RouteMethod {
	return r.addRoute(path, methodGET, endpoint)
}

// HEAD - add HEAD request to router.
func (r *Router) HEAD(path string, endpoint http.HandlerFunc) *RouteMethod {
	return r.addRoute(path, methodHEAD, endpoint)
}

// POST - add POST request to router.
func (r *Router) POST(path string, endpoint http.HandlerFunc) *RouteMethod {
	return r.addRoute(path, methodPOST, endpoint)
}

// PUT - add PUT request to router.
func (r *Router) PUT(path string, endpoint http.HandlerFunc) *RouteMethod {
	return r.addRoute(path, methodPUT, endpoint)
}

// DELETE - add DELETE request to router.
func (r *Router) DELETE(path string, endpoint http.HandlerFunc) *RouteMethod {
	return r.addRoute(path, methodDELETE, endpoint)
}

// OPTIONS - add OPTIONS request to router.
func (r *Router) OPTIONS(path string, endpoint http.HandlerFunc) *RouteMethod {
	return r.addRoute(path, methodOPTIONS, endpoint)
}

// PATCH - add PATCH request to router.
func (r *Router) PATCH(path string, endpoint http.HandlerFunc) *RouteMethod {
	return r.addRoute(path, methodPATCH, endpoint)
}

// Use - add global middlewares. Any request will call these middlewares.
func (r *Router) Use(middlewares ...Middleware) {
	r.middleware.chain = middlewareChainer(middlewares, &r.middleware)
}

// ServeHTTP - when a request comes in and goes out, it will be here.
func (r *Router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// if favicon - don't handle request
	//if request.RequestURI == "/favicon.ico" {
	//	return
	//}
	// global middleware
	resume := r.run(response, request)
	if !resume {
		return
	}
	var path = pathToSlice(request.URL.Path)
	var matched = r.match(request.Method, path)
	// 404
	if matched.notFound {
		r.notFound.ServeHTTP(response, request)
		return
	}
	// 405
	if matched.methodNotAllowed {
		response.Header().Add("Allow", strings.Join(matched.route.allowedMethods, ", "))
		r.methodNotAllowed.ServeHTTP(response, request)
		return
	}
	// provide {params}
	request = request.WithContext(context.WithValue(context.Background(), ctxPathParams, matched.params))
	// group
	var group = matched.group
	if group != nil {
		resume = group.run(response, request)
		if !resume {
			return
		}
	}
	// route with method found, run middlewares and endpoint
	var _method = matched.method
	_method.run(response, request)
}

// run - run middleware if exists. Returns true if middleware called next().
func (r *Router) run(response http.ResponseWriter, request *http.Request) (resume bool) {
	if r.middleware.chain != nil {
		r.middleware.chain.ServeHTTP(response, request)
		return r.middleware.resume
	}
	return true
}

// ServeHTTP - if global middleware exists and next() called, this method executes.
func (r *routerMiddleware) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	r.resume = true
}
