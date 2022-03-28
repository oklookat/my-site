package way

import (
	"net/http"
	"strings"
)

type Route struct {
	// is route has prefix? (it means route under group)
	isHasPrefix bool

	// route prefix.
	prefix string

	// route path like: /hello/world.
	path string

	// route path slice like: [hello, world].
	pathSlice []string

	// allowed route methods.
	allowedMethods []string

	// route middleware chain.
	middleware MiddlewareFunc

	// route endpoint.
	handler RouteHandler
}

func (r *Route) new(prefix string, to string, handler RouteHandler) {
	// check prefix.
	if len(prefix) > 0 {
		r.isHasPrefix = true
		r.prefix = prefix
	}

	// clean.
	r.path = pathToStandart(to)

	// split.
	r.pathSlice = strings.Split(removeSlashStartEnd(r.path), "/")

	// set.
	r.handler = handler
}

func (r *Route) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	// run middleware.
	var isResponseSended = executeMiddleware(response, request, r.middleware)
	if isResponseSended {
		return
	}

	r.handler(response, request)
	return
}

// route trigger on this methods only.
func (r *Route) Methods(methods ...string) *Route {
	r.allowedMethods = processAllowedMethods(r.allowedMethods, methods...)
	return r
}

// provide middleware.
func (r *Route) Use(middleware ...MiddlewareFunc) {
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
