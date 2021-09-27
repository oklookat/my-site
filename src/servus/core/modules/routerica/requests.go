package routerica

import "net/http"

// TODO: define all requests

// GET - add GET request to router.
func (r *Routerica) GET(path string, endpoint http.HandlerFunc) *RouteLocal {
	return r.addRequest(path, methodGET, endpoint)
}

// GET - add GET request to the route group.
func (g *RouteGroup) GET(path string, endpoint http.HandlerFunc) *RouteLocal {
	return g.addRequest(path, methodGET, endpoint)
}

