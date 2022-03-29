package way

import (
	"net/http"
	"strings"
)

/*
way router
some ideas from gorilla/mux router
https://github.com/gorilla/mux/blob/master/LICENSE
*/

// middleware.
type MiddlewareFunc func(http.Handler) http.Handler

// route endpoint.
type RouteHandler func(http.ResponseWriter, *http.Request)

// vars for request context.
type CTX_VAL string

const CTX_VARS_NAME CTX_VAL = "WAY_ROUTER_VARS"

func New() *Router {
	return &Router{}
}

type Router struct {

	// exclude this prefix from request path, because it's computed in parent group.
	excludePrefix string

	// group prefix like: /hello/world.
	prefix string

	// group prefix slice like: [hello, world].
	prefixSlice []string

	// allowed methods.
	allowedMethods []string

	// middleware chain.
	middleware MiddlewareFunc

	// route groups.
	groups []*Router

	// routes.
	routes []*Route

	// 404 handler.
	notFoundHandler RouteHandler
}

// any parents (routes or groups) should remove this exclude prefix from
// request path to match request.
func (r *Router) getExcludePrefix() string {
	var exclude string
	if len(r.excludePrefix) > 0 {
		exclude = "/" + r.excludePrefix + "/" + r.prefix
	} else if len(r.prefix) > 0 {
		exclude = "/" + r.prefix
	}
	exclude = pathToStandart(exclude)
	return exclude
}

// when request coming.
func (r *Router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// run middleware.
	var isResponseSended = executeMiddleware(response, request, r.middleware)
	if isResponseSended {
		return
	}

	// check is method allowed.
	var isMethodAllowed = isMethodAllowed(r.allowedMethods, request.Method)
	if !isMethodAllowed {
		send405(response)
		return
	}

	var matcher = routeMatcher{}
	matcher.New(request)

	// try to match groups first.
	if r.groups != nil {
		var matched, code = matcher.Groups(r.groups)
		if matched != nil {
			matched.ServeHTTP(response, request)
			return
		}
		if code == 405 {
			send405(response)
			return
		}
	}

	// try to match routes.
	if r.routes != nil {
		var matched, code = matcher.Routes(r.routes)
		if matched != nil {
			matched.ServeHTTP(response, request)
			return
		}
		if code == 405 {
			send405(response)
			return
		}
	}

	// 404.
	if r.notFoundHandler != nil {
		r.notFoundHandler(response, request)
		return
	}
	send404(response)
	return
}

// add route.
func (r *Router) Route(to string, handler RouteHandler) *Route {
	if r.routes == nil {
		r.routes = make([]*Route, 0)
	}

	// new route.
	var newRoute = &Route{}
	newRoute.new(r.getExcludePrefix(), to, handler)

	// add to routes.
	r.routes = append(r.routes, newRoute)
	return newRoute
}

// add route group.
func (r *Router) Group(prefix string) (group *Router) {
	// make groups if not.
	if r.groups == nil {
		r.groups = make([]*Router, 0)
	}

	// create new router.
	var newRouter = New()
	newRouter.excludePrefix = r.getExcludePrefix()
	newRouter.prefix = pathToStandart(prefix)
	newRouter.prefixSlice = strings.Split(removeSlashStartEnd(newRouter.prefix), "/")

	// add to groups.
	r.groups = append(r.groups, newRouter)
	return newRouter
}

// provide middleware.
func (r *Router) Use(middleware ...MiddlewareFunc) {
	var chained = make([]MiddlewareFunc, 0)
	if r.middleware != nil {
		chained = append(chained, r.middleware)
	}
	for _, m := range middleware {
		if m == nil {
			continue
		}
		chained = append(chained, m)
	}
	var finalChain = middlewareChain(chained...)
	r.middleware = finalChain
}

// 404 handler.
func (r *Router) NotFound(handler RouteHandler) {
	r.notFoundHandler = handler
}

// add allowed methods.
func (r *Router) Methods(methods ...string) *Router {
	r.allowedMethods = processAllowedMethods(r.allowedMethods, methods...)
	return r
}
