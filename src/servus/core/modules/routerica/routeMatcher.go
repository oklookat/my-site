package routerica

import (
	"context"
	"net/http"
)


// routeMatcher - depending on request method and path choose right routes array, executes middleware (if exists) and then run handler (if middleware not returned response).
func routeMatcher(r *Routerica, response http.ResponseWriter, request *http.Request) {
	var requestMethod = request.Method
	var isValidMethod = requestMethod == methodGET || requestMethod == methodPOST || requestMethod == methodPUT || requestMethod == methodDELETE
	// go back if method like OPTIONS.
	if !isValidMethod {
		return
	}
	//////// route matching.
	var routeLocal *RouteLocal
	// get uri without query string.
	var requestPath = uriSplitter(request.URL.Path)
	// define params (if {delimiter} exists, it will be replaced request value and added to context)
	// ex: /{id}/change => /1337/change => map[id: 1337].
	var paramsMap = make(map[string]string, 0)
	var paramsMapTemp = make(map[string]string, 0)
	// route groups.
	routeGroup, paramsMap := parserRouteGroups(r.routeGroups, requestPath)
	request = request.WithContext(context.WithValue(context.Background(), ctxValueParams, paramsMap))
	// getting routeLocal depending on request URI (group or not).
	switch routeGroup {
	default:
		// we have route group, execute group middleware and parse routes inside group.
		routeGroup.middlewareGroupGlobal.chain.ServeHTTP(response, request)
		routeLocal, paramsMapTemp = parserRouteLocals(routeGroup.routeLocals, requestMethod, requestPath)
		break
	case nil:
		// we don't have route group, switching to parse routes without group.
		routeLocal, paramsMapTemp = parserRouteLocals(r.routeLocals, requestMethod, requestPath)
		break
	}
	paramsMap = mapConcat(paramsMap, paramsMapTemp)
	request = request.WithContext(context.WithValue(context.Background(), ctxValueParams, paramsMap))
	// execute routeLocal middleware or handler. Or go to 404.
	switch routeLocal {
	default:
		// we have routeLocal.
		switch routeLocal.middlewareLocal.chain {
		default:
			// if middleware exists, at the end it executes endpoint.
			routeLocal.middlewareLocal.chain.ServeHTTP(response, request)
			break
		case nil:
			// middleware not exists - we run endpoint directly.
			routeLocal.endpoint(response, request)
			break
		}
		break
	case nil:
		// we don't have routeLocal. 404 moment.
		r.routeNotFound.endpoint.ServeHTTP(response, request)
		break
	}
}
