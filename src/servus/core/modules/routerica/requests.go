package routerica

import "net/http"


// GET - add GET request to router.
func (r *Routerica) GET(path string, endpoint http.HandlerFunc) *RouteLocal {
	return r.addRequest(path, methodGET, endpoint)
}

// POST - add POST request to router.
func (r *Routerica) POST(path string, endpoint http.HandlerFunc) *RouteLocal {
	return r.addRequest(path, methodPOST, endpoint)
}

// PUT - add PUT request to router.
func (r *Routerica) PUT(path string, endpoint http.HandlerFunc) *RouteLocal {
	return r.addRequest(path, methodPUT, endpoint)
}

// DELETE - add DELETE request to router.
func (r *Routerica) DELETE(path string, endpoint http.HandlerFunc) *RouteLocal {
	return r.addRequest(path, methodDELETE, endpoint)
}

// GET - add GET request to the route group.
func (g *RouteGroup) GET(path string, endpoint http.HandlerFunc) *RouteLocal {
	return g.addRequest(path, methodGET, endpoint)
}

// POST - add POST request to the route group.
func (g *RouteGroup) POST(path string, endpoint http.HandlerFunc) *RouteLocal {
	return g.addRequest(path, methodPOST, endpoint)
}

// PUT - add PUT request to the route group.
func (g *RouteGroup) PUT(path string, endpoint http.HandlerFunc) *RouteLocal {
	return g.addRequest(path, methodPUT, endpoint)
}

// DELETE - add DELETE request to the route group.
func (g *RouteGroup) DELETE(path string, endpoint http.HandlerFunc) *RouteLocal {
	return g.addRequest(path, methodDELETE, endpoint)
}
