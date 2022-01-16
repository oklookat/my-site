package way

import (
	"fmt"
	"net/http"
)

// handler with path, methods, and limiter.
type Point struct {
	methods    []string
	middleware *middleware
	handler    http.HandlerFunc
}

// add endpoint.
func (r *Router) Endpoint(path string, handler http.HandlerFunc) *Point {
	var _route = route{}
	var notExists, point = _route.new(path, handler, r.routes)
	if notExists {
		r.routes = append(r.routes, &_route)
	}
	return point
}

// add endpoint.
func (g *Group) Endpoint(path string, handler http.HandlerFunc) *Point {
	var pathFormatted = fmt.Sprintf("/%v/%v/", g.prefix, path)
	var _route = route{}
	var notExists, point = _route.new(pathFormatted, handler, g.routes)
	if notExists {
		g.routes = append(g.routes, &_route)
	}
	return point
}

// create new Point.
func (p *Point) new(handler http.HandlerFunc) {
	p.middleware = &middleware{}
	p.handler = handler
}

// add methods to route.
func (p *Point) Methods(methods ...string) *Point {
	var add = func(method string) {
		for index := range p.methods {
			var _method = p.methods[index]
			// if method exists.
			if _method == method {
				return
			}
		}
		// add method.
		p.methods = append(p.methods, method)
	}
	for index := range methods {
		add(methods[index])
	}
	return p
}

// add route middlewares. Request on this route will call these middlewares.
func (p *Point) Use(middlewares ...Middleware) {
	p.middleware.add(middlewares)
}

// run middleware if exists and run endpoint.
func (p *Point) runMiddlewareAndEndpoint(response http.ResponseWriter, request *http.Request) (incorrectMethod bool) {
	// check is method exists in point.
	var method = request.Method
	var correctMethod = false
	for _, allowedMethod := range p.methods {
		if method == allowedMethod {
			correctMethod = true
			break
		}
	}
	// 405 if method not exists.
	if !correctMethod {
		incorrectMethod = true
		return
	}
	// run limiter.
	if p.middleware.chain != nil {
		p.middleware.chain.ServeHTTP(response, request)
		// if limiter not called next().
		if !p.middleware.executed {
			return
		}
	}
	// run endpoint.
	p.handler.ServeHTTP(response, request)
	return
}
