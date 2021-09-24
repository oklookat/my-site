package routerica

import (
	"net/http"
	"regexp"
)

// depending on request method and path choose right routes array, executes middleware (if exists) and then run handler (if middleware not returned response)
func routeMatcher(r *Routerica, response http.ResponseWriter, request *http.Request) {
	var requestMethod = request.Method
	var requestURL = request.URL
	var requestURI = requestURL.RequestURI()
	var isValidMethod = requestMethod == methodGET || requestMethod == methodPOST || requestMethod == methodPUT || requestMethod == methodDELETE
	// if method like OPTIONS we go back
	if !isValidMethod {
		return
	}
	//////// route matching
	// route groups
	var routeLocal *RouteLocal
	var routeGroup = parseRouteGroup(r.routeGroups, requestURI)
	if routeGroup != nil {

		routeGroup.middlewareChain.ServeHTTP(response, request)
		routeLocal = parseRouteLocal(routeGroup.localRoutes, requestMethod, requestURI)
	} else {
		routeLocal = parseRouteLocal(r.localRoutes, requestMethod, requestURI)
	}
	if routeLocal != nil {
		if routeLocal.middlewareChain != nil {
			// at the end middleware automatically executes handler
			routeLocal.middlewareChain.ServeHTTP(response, request)
		} else {
			// otherwise - run handler
			routeLocal.handler(response, request)
		}
	}
	// TODO: add 404 handler here
}

// create one big middleware from middlewares (chain)
func middlewareChainer(middlewares []Middleware, next http.Handler) http.Handler {
	// https://gist.github.com/husobee/fd23681261a39699ee37#gistcomment-3111569
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i](next)
	}
	return next
}

// ex: //1////2/3 to /1/2/3
func cleanupSlashes(data *string) {
	regex := regexp.MustCompile(`\/\/+`)
	*data = string(regex.ReplaceAll([]byte(*data), []byte("/")))
}
