package router

import "net/http"

// route - one route.
type route struct {
	path []string
	allowedMethods []string
	// methods - key = method.
	methods routeMethods
}

// routes - key is path.
type routes map[string]*route

type routeMethods map[string]*RouteMethod

type RouteMethod struct {
	resume     bool
	middleware MiddlewareChain
	endpoint   http.HandlerFunc
}

// new - safely create new route instance.
func (r *route) new() {
	r.path = make([]string, 0)
	r.allowedMethods = make([]string, 0)
	r.methods = make(routeMethods, 0)
}

// addMethod - add method to route. If method exists - replaced.
func (r *route) addMethod(method string, m *RouteMethod) {
	var allowedFound = false
	for _, allowed := range r.allowedMethods {
		if allowed == method {
			allowedFound = true
			break
		}
	}
	if !allowedFound {
		r.allowedMethods = append(r.allowedMethods, method)
	}
	r.methods[method] = m
}

// isHasMethod - check is route has method.
func (r *route) isHasMethod(method string) bool {
	_, ok := r.methods[method]
	return ok
}

// run - run middleware if exists and run endpoint.
func (r *RouteMethod) run(response http.ResponseWriter, request *http.Request) {
	if r.middleware != nil {
		r.middleware.ServeHTTP(response, request)
		if !r.resume {
			return
		}
	}
	r.endpoint.ServeHTTP(response, request)
}

// newRouteLocal - create route local instance (internal).
func (r *RouteMethod) new(route *route, path string, handler http.HandlerFunc) {
	pathSlice := pathToSlice(path)
	route.path = pathSlice
	r.endpoint = handler
}

// Use - add route middlewares. Request on this route will call these middlewares.
func (r *RouteMethod) Use(middlewares ...Middleware) {
	r.middleware = middlewareChainer(middlewares, r.middleware)
}

func (r *RouteMethod) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	r.resume = true
}
