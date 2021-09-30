package routerica

import (
	"context"
	"net/http"
)

// middlewareGlobal - used only for router instance (global middleware).
type middlewareGlobal struct {
	endpoint http.HandlerFunc
	chain    MiddlewareChain
}

// middlewareGroupGlobal - used only for route group (global middleware).
type middlewareGroupGlobal struct {
	endpoint http.HandlerFunc
	chain    MiddlewareChain
}

// middlewareLocal - used only for a specific route (local middleware)
type middlewareLocal struct {
	chain MiddlewareChain
}

// Use - add global middlewares. Any request will call these middlewares.
func (r *Routerica) Use(middlewares ...Middleware) {
	r.middlewareGlobal.chain = middlewareChainer(middlewares, r.middlewareGlobal)
}

// Use - add group global middlewares. Any request on this route group will call these middlewares.
func (g *RouteGroup) Use(middlewares ...Middleware) *RouteGroup {
	g.middlewareGroupGlobal.chain = middlewareChainer(middlewares, g.middlewareGroupGlobal)
	return g
}

// Use - add route middlewares. Request on this route will call these middlewares.
func (l *RouteLocal) Use(middlewares ...Middleware) {
	l.middlewareLocal.chain = middlewareChainer(middlewares, l.endpoint)
}

// ServeHTTP - when a request comes in and goes out, it will be here
func (r *Routerica) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// the function is called twice. One for /favicon, and one for the route
	// if favicon - we don't handle request
	if request.RequestURI == "/favicon.ico" {
		return
	}
	// pass routerica instance to request
	request = request.WithContext(context.WithValue(context.Background(), ctxInternal, r))
	r.middlewareGlobal.chain.ServeHTTP(response, request)
}

// ServeHTTP - when global middleware finished or if no global middlewares, it calls this method. If global middleware send response, this method will not be called.
func (g *middlewareGlobal) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	var routerica = request.Context().Value(ctxInternal)
	// clear context
	request = request.WithContext(context.Background())
	routeMatcher(routerica.(*Routerica), response, request)
}

// ServeHTTP - when group global middleware finished or if no global middlewares, it calls this method. If global middleware send response, this method will not be called.
func (g *middlewareGroupGlobal) ServeHTTP(response http.ResponseWriter, request *http.Request) {

}
