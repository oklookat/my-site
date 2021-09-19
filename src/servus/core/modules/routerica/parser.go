package routerica

import (
	"net/http"
	"strings"
)

func parseGet(r *Routerica, requestPath string, response http.ResponseWriter, request *http.Request) {
	for index := range r.getRoutes {
		var currentRoute = r.getRoutes[index]
		var routePath = strings.ToUpper(currentRoute.path)
		if routePath != requestPath {
			continue
		}
		// now we in correct route
		// run middleware
		currentRoute.middlewareChain.ServeHTTP(response, request)
		currentRoute.handler(response, request)
	}
}