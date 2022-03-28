package way

import (
	"net/http"
	"strings"
)

type routeMatcher struct {
	request          *http.Request
	requestPath      string
	requestPathSlice []string
	method           string
}

func (r *routeMatcher) New(req *http.Request) {
	r.request = req
	r.method = r.request.Method

	// convert request path to standart path, like we do it with route paths.
	r.requestPath = pathToStandart(req.URL.Path)
	r.requestPathSlice = strings.Split(removeSlashStartEnd(r.requestPath), "/")
}

func (r *routeMatcher) Groups(routers []*Router) (matched *Router, statusCode int) {
	statusCode = 200
	for i := range routers {

		// example: request /hello, prefix /hello/world. Not our group.
		if len(r.requestPathSlice) < len(routers[i].prefixSlice) {
			continue
		}

		var isPieceMatched = r.matchPathPieces(i, routers[i].prefixSlice, r.requestPathSlice)
		if isPieceMatched {
			matched = routers[i]
		}

		var isLastGroup = i == len(routers)-1

		// if last group and no match, set 404 code.
		if matched == nil && isLastGroup {
			// if we found group before, but with not allowed method
			// then not change status code
			if statusCode != 405 {
				statusCode = 404
			}
			matched = nil
			return
		} else if matched != nil {
			// if we found group.
			// check is method allowed.
			var isAllowed = isMethodAllowed(routers[i].allowedMethods, r.method)
			if isLastGroup && !isAllowed {
				// if it last group and method not allowed, set 405 code.
				statusCode = 405
				matched = nil
				return
			} else if !isAllowed {
				// not last group, but method not allowed, go next.
				continue
			}
			// method allowed, its our group.
			return
		}
	}
	return
}

func (r *routeMatcher) Routes(routes []*Route) (matched *Route, statusCode int) {
	statusCode = 200
	for i := range routes {

		var requestPath = r.requestPath
		var requestPathSlice = r.requestPathSlice

		// check is route under group (has prefix).
		if routes[i].isHasPrefix {
			// remove prefix from request path (its already computed in group).
			requestPath = strings.TrimPrefix(requestPath, routes[i].prefix)
			requestPathSlice = strings.Split(removeSlashStartEnd(requestPath), "/")
		}

		// example: request /hello, route /hello/world. Not our route.
		if len(requestPathSlice) != len(routes[i].pathSlice) {
			continue
		}

		var isPiecesMatched = r.matchPathPieces(i, routes[i].pathSlice, requestPathSlice)
		if isPiecesMatched {
			matched = routes[i]
		}

		var isLastRoute = i == len(routes)-1

		// if last route and no match, set 404 code.
		if matched == nil && isLastRoute {
			// if we found route before, but with not allowed method
			// then not change status code
			if statusCode != 405 {
				statusCode = 404
			}
			matched = nil
			return
		} else if matched != nil {
			// if we found route.
			// check is method allowed.
			var isAllowed = isMethodAllowed(routes[i].allowedMethods, r.method)
			if isLastRoute && !isAllowed {
				// if it last route and method not allowed, set 405 code.
				statusCode = 405
				matched = nil
				return
			} else if !isAllowed {
				// not last route, but method not allowed, go next.
				continue
			}
			// all checks passed, its our group.
			return
		}
	}
	return
}

func (r *routeMatcher) matchPathPieces(counter int, pathSlice []string, requestPathSlice []string) (matched bool) {
	// check path pieces.
	for j := range pathSlice {
		// example: is req piece /hello, and prefix /hello
		var isSame = pathSlice[j] == requestPathSlice[j]

		// not same? maybe it's route variable?
		if !isSame {
			// example: is req piece /users/12, and prefix /users/{id}
			var isVar, name = isRouteVar(pathSlice[j])
			if !isVar {
				// summary: pieces not same, and it's not var. Not our group.
				matched = false
				break
			}

			// add var to request context.
			addVarToContext(r.request, name, requestPathSlice[counter])
		}

		// check is last piece.
		var isLast = j == len(pathSlice)-1
		if isLast {
			// maybe our group.
			matched = true
		}
	}
	return
}
