package routerica

import "net/http"

// GET request
func (r *Routerica) GET(path string, handler func(http.ResponseWriter, *http.Request)) *TheRoute {
	var theRoute = TheRoute{path: path, method: methodGET, BaseRoute: BaseRoute{handler: handler}}
	r.localRoutes[methodGET] = append(r.localRoutes[methodGET], theRoute)
	var lastRoute = len(r.localRoutes[methodGET]) - 1
	return &r.localRoutes[methodGET][lastRoute]
}
