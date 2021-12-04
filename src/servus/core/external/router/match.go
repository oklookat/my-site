package router

import (
	"strings"
)

// matchGroup - get route group from request path. If success, returns routeGroup and requestURI without routeGroup prefix.
func matchGroup(groups []Group, path []string) (matched *Group, params map[string]string) {
	for index := range groups {
		matched = &groups[index]
		isMatched, params := matchPath(matched.prefix, path)
		if !isMatched {
			continue
		}
		return matched, params
	}
	return nil, nil
}

// matchRoute - get route local from routes, request path and method.
func matchRoute(routes routes, path []string) (matched *route, params map[string]string) {
	for defPath := range routes {
		isMatched, params := matchPath(pathToSlice(defPath), path)
		if !isMatched {
			continue
		}
		matched = routes[defPath]
		return matched, params
	}
	return nil, nil
}

// matchPath - if path equals requestPath, returns true and map like [name: value] if {param}.
//
// path with params like: [hello, {name}]
//
// requestPath like: [hello, john]
func matchPath(path []string, requestPath []string) (isMatch bool, params map[string]string) {
	if isEmpty(path) || isEmpty(requestPath) || moreThan(path, requestPath) {
		return false, nil
	}
	// used for comparing.
	var verified = 0
	params = make(map[string]string, 0)
	for pathIndex := range path {
		var pathPart = path[pathIndex]
		var requestPathPart = requestPath[pathIndex]
		var sameParts = strings.EqualFold(requestPathPart, pathPart)
		if sameParts {
			verified++
			continue
		}
		// check {param}.
		var hasParam = strings.HasPrefix(pathPart, paramOpen) && strings.HasSuffix(pathPart, paramClose)
		if hasParam {
			// get param name without { }. Ex: {user} = user.
			pathPart = strings.ReplaceAll(pathPart, paramOpen, "")
			pathPart = strings.ReplaceAll(pathPart, paramClose, "")
			// paste param and value to map. Ex: params[user] = 1.
			params[pathPart] = requestPathPart
			verified++
			continue
		}
		// parts not same / not have param, we don't need to continue.
		break
	}
	var pathLen = len(path) - 1
	return verified != pathLen, params
}
