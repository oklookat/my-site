package routerica

import (
	"net/http"
	"strings"
)

func New() *Routerica {
	return &Routerica{}
}

// Use middleware. Create one big middleware from middlewares (chain).
func (g *GetRoute) Use(middlewares ...TheMiddleware) {
	// https://gist.github.com/husobee/fd23681261a39699ee37#gistcomment-3111569
	var chainer = func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
	g.middlewareChain = chainer(g.handler)
}

// when a request comes in, it first appears here
func (r *Routerica) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// the function is called twice. One for /favicon, and one for route
	println("i am in entry point")
	var requestMethod = strings.ToUpper(request.Method)
	var requestURL = request.URL
	var requestPath = strings.ToUpper(requestURL.Path)
	switch requestMethod {
	case methodGet:
		parseGet(r, requestPath, response, request)
	default:
		break
	}

	//for index, middleware := range r.middlewares {
	//	if index > len(r.middlewares) - 1 {
	//		break
	//	}
	//	middleware(r.middlewares[index + 1])
	//}
}

// Get request
func (r *Routerica) Get(path string, handler func(http.ResponseWriter, *http.Request)) *GetRoute {
	var theRoute = theRoute{path: path, method: methodGet, handler: handler}
	var getRoute = &GetRoute{theRoute: theRoute}
	r.getRoutes = append(r.getRoutes, getRoute)
	return getRoute
}
