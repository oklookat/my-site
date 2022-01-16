package way

import (
	"net/http"
	"strings"
)

// group of routes.
type Group struct {
	prefix     string
	middleware *middleware
	routes     []*route
}

// create group of routes.
func (r *Router) Group(prefix string) *Group {
	var group = &Group{}
	group.new(prefix)
	r.groups = append(r.groups, group)
	return group
}

// create route group.
func (g *Group) new(prefix string) {
	g.prefix = prefix
	g.routes = make([]*route, 0)
	g.middleware = &middleware{}
}

// add group global middlewares. Any request on this route group will call these middlewares.
func (g *Group) Use(middlewares ...Middleware) *Group {
	g.middleware.add(middlewares)
	return g
}

// run middleware if exists. Returns true if middleware called next().
func (g *Group) runMiddleware(response http.ResponseWriter, request *http.Request) (executed bool) {
	if g.middleware.chain != nil {
		g.middleware.chain.ServeHTTP(response, request)
		return g.middleware.executed
	}
	return true
}

// this group? If it is, get params.
func (g *Group) match(requestPath string) (matched bool, params map[string]string) {
	matched = false
	if !strings.HasPrefix(requestPath, g.prefix) {
		return
	}
	invalid, differ, params := verifyPaths(g.prefix, requestPath, false)
	if invalid || differ {
		return
	}
	matched = true
	return
}
