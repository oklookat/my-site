package routerica

import (
	"net/http"
	"strings"
)

// depending on request method choose right routes array, executes middleware (if exists) and handler
func urlParser(r *Routerica, response http.ResponseWriter, request *http.Request) {
	var requestMethod = strings.ToUpper(request.Method)
	var requestURL = request.URL
	var requestPath = strings.ToUpper(requestURL.Path)
	var isValidMethod = requestMethod == methodGET || requestMethod == methodPOST || requestMethod == methodPUT || requestMethod == methodDELETE
	// if method like OPTIONS we go back
	if !isValidMethod {
		return
	}
	for index := range r.localRoutes[requestMethod] {
		var currentRoute = r.localRoutes[requestMethod][index]
		var routePath = strings.ToUpper(currentRoute.path)
		if routePath != requestPath {
			continue
		}
		// now we in correct route
		// run middleware if exists
		if currentRoute.middlewareChain != nil {
			// at the end middleware automatically executes handler
			currentRoute.middlewareChain.ServeHTTP(response, request)
		} else {
			// otherwise - run handler
			currentRoute.handler(response, request)
		}
		break
	}
}

// create one big middleware from middlewares (chain)
func middlewareChainer(middlewares []TheMiddleware, next http.Handler) http.Handler {
	// https://gist.github.com/husobee/fd23681261a39699ee37#gistcomment-3111569
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i](next)
	}
	return next
}
