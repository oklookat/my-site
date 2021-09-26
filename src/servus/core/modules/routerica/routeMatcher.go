package routerica

import (
	"net/http"
)

// TODO: pass params map to request context.
// TODO: remove big case converting (replace to strings.EqualFold)
// TODO: remove adding prefix to route group routes(?)
// TODO: write normal tests and benchmarks
// depending on request method and path choose right routes array, executes middleware (if exists) and then run handler (if middleware not returned response).
func routeMatcher(r *Routerica, response http.ResponseWriter, request *http.Request) {
	var requestMethod = request.Method
	var isValidMethod = requestMethod == methodGET || requestMethod == methodPOST || requestMethod == methodPUT || requestMethod == methodDELETE
	// go back if method like OPTIONS.
	if !isValidMethod {
		return
	}
	//////// route matching.
	var routeLocal *RouteLocal
	var requestURISlice = uriSplitter(formatPath(request.URL.RequestURI()))
	// route groups.
	var routeGroup = parseRouteGroup(r.routeGroups, requestURISlice)
	// getting routeLocal depending on request URI.
	switch routeGroup {
	default:
		routeGroup.middlewareChain.ServeHTTP(response, request)
		routeLocal = parseRouteLocal(routeGroup.localRoutes, requestMethod, requestURISlice)
		break
	case nil:
		routeLocal = parseRouteLocal(r.localRoutes, requestMethod, requestURISlice)
		break
	}
	// execute routeLocal middleware or handler. Or go to 404.
	switch routeLocal {
	default:
		switch routeLocal.middlewareChain {
		default:
			// if middleware exists, at the end it executes handler.
			routeLocal.middlewareChain.ServeHTTP(response, request)
			break
		case nil:
			// Or middleware not exists we run handler directly.
			routeLocal.handler(response, request)
			break
		}
		break
	case nil:
		// TODO: add 404 handler here
		println("404")
		break
	}
}