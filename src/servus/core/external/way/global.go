package way

import (
	"context"
	"net/http"
)

// simple router.
type Router struct {
	notFound         http.HandlerFunc
	methodNotAllowed http.HandlerFunc
	middleware       *middleware
	groups           []*Group
	routes           []*route
}

// create new Router instance.
func New() *Router {
	var router = &Router{}
	router.groups = make([]*Group, 0)
	router.routes = make([]*route, 0)
	router.middleware = &middleware{}
	router.notFound = defaultNotFound
	router.methodNotAllowed = defaultNotAllowed
	return router
}

// add 404 page. When no suitable route is found, the endpoint is called.
func (r *Router) NotFound(handler http.HandlerFunc) {
	r.notFound = handler
}

// add 405 page. When no method for route, the endpoint is called.
func (r *Router) MethodNotAllowed(handler http.HandlerFunc) {
	r.methodNotAllowed = handler
}

// add global middlewares. Any request will call these middlewares.
func (r *Router) Use(middlewares ...Middleware) {
	r.middleware.add(middlewares)
}

// run middleware if exists. Returns true if middleware called next().
func (r *Router) runMiddleware(response http.ResponseWriter, request *http.Request) (executed bool) {
	if r.middleware.chain != nil {
		r.middleware.chain.ServeHTTP(response, request)
		return r.middleware.executed
	}
	return true
}

// when a request comes in and goes out, it will be here.
func (r *Router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// global limiter
	executed := r.runMiddleware(response, request)
	if !executed {
		return
	}
	group, _route, params := r.match(request.URL.Path)
	// 404
	if _route == nil {
		r.notFound.ServeHTTP(response, request)
		return
	}
	// provide {params}
	request = request.WithContext(context.WithValue(context.Background(), ctxPathParams, params))
	// group
	if group != nil {
		executed = group.runMiddleware(response, request)
		if !executed {
			return
		}
	}
	// find route and execute limiter, endpoint.
	_route.findAndExecute(response, request, r.methodNotAllowed)
}
