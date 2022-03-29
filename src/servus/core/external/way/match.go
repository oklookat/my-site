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
	var requestPath = r.requestPath
	var requestPathSlice = r.requestPathSlice

	for i := range routers {
		// check is route has parent prefix.
		var isHasExcludePrefix = len(routers[i].excludePrefix) > 0
		if isHasExcludePrefix {
			// exclude parent prefix from request path (its already computed before).
			requestPath = strings.TrimPrefix(requestPath, routers[i].excludePrefix)
			requestPathSlice = strings.Split(removeSlashStartEnd(requestPath), "/")
		}

		// example: prefix /hello/world, request /hello. Not our group.
		if len(routers[i].prefixSlice) > len(r.requestPathSlice) {
			continue
		}

		var isPieceMatched = r.matchPathPieces(i, routers[i].prefixSlice, requestPathSlice)
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
	var requestPath = r.requestPath
	var requestPathSlice = r.requestPathSlice
	for i := range routes {
		// check is route under group (has prefix).
		if routes[i].isUnderGroup {
			// remove prefix from request path (its already computed in group).
			requestPath = strings.TrimPrefix(requestPath, routes[i].excludePrefix)
			requestPathSlice = strings.Split(removeSlashStartEnd(requestPath), "/")
		}

		// example: route /hello/world, request /hello. Not our route.
		if len(routes[i].pathSlice) > len(requestPathSlice) {
			continue
		}

		var isPiecesMatched = r.matchPathPieces(i, routes[i].pathSlice, requestPathSlice)
		if isPiecesMatched {
			matched = routes[i]
		}

		var isLastRoute = i == len(routes)-1

		// if we found route.
		if matched != nil {
			// check is method allowed.
			var isAllowed = isMethodAllowed(matched.allowedMethods, r.method)
			if isAllowed {
				// it's our route.
				return
			} else {
				// method not allowed.
				if isLastRoute {
					// method not allowed and no more routes.
					matched = nil
					statusCode = 405
					return
				}
				// method not allowed, but maybe try other paths?
				statusCode = 405
			}
		} else if isLastRoute {
			// if we not found route and it's last path.
			if statusCode != 405 {
				statusCode = 404
			}
			matched = nil
			return
		}

	}
	return
}

func (r *routeMatcher) matchPathPieces(counter int, pathSlice []string, requestPathSlice []string) (matched bool) {
	// compare paths.
	for pieceCounter := range pathSlice {
		var pathPiece = pathSlice[pieceCounter]
		var requestPathPiece = requestPathSlice[pieceCounter]

		var isSame = pathPiece == requestPathPiece
		// not same? maybe it's route variable?
		if !isSame {
			var isVar, name = isRouteVar(pathPiece)
			if !isVar {
				// summary: pieces not same, and it's not var. Not our path.
				matched = false
				break
			}

			// if it's var, add it to context
			addVarToContext(r.request, name, requestPathPiece)
		}

		// check is last piece.
		var isLast = pieceCounter == len(pathSlice)-1
		if isLast {
			// pieces same, maybe it's our path.
			matched = true
			return
		}
	}
	return
}
