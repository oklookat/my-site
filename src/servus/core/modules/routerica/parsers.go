package routerica

import (
	"strings"
)

// parserRouteGroups - get route group from request URI. If success, returns routeGroup and requestURI without routeGroup prefix.
func parserRouteGroups(groups []RouteGroup, requestPath []string) (routeGroup *RouteGroup, paramsMap map[string]string) {
	for index := range groups {
		var currentGroup = &groups[index]
		var isGroupMatched, paramsMap = parserPath(currentGroup.prefix, requestPath)
		if !isGroupMatched {
			continue
		}
		return currentGroup, paramsMap
	}
	return nil, nil
}

// parserRouteLocals - get route local from routes, request URI and method.
func parserRouteLocals(routes map[string][]RouteLocal, requestMethod string, requestPath []string) (routeLocal *RouteLocal, paramsMap map[string]string) {
	for index := range routes[requestMethod] {
		var currentRoute = &routes[requestMethod][index]
		var isMatched, paramsMap = parserPath(currentRoute.path, requestPath)
		if !isMatched {
			continue
		}
		// now we in correct route.
		return currentRoute, paramsMap
	}
	return nil, nil
}

// parserPath - if url path equals requestPath, returns true and map like [param: value] if {param}.
func parserPath(path []string, requestPath []string) (isMatch bool, paramsMap map[string]string) {
	if len(path) > len(requestPath) {
		return false, nil
	}
	// equalCount - get size of uri to compare with request uri slice.
	var equalCount = len(path) - 1
	// verified - used for comparing.
	var verified = 0
	// requestUriSliceFormatted - contains requestUri items without matched.
	paramsMap = make(map[string]string, 0)
	for pathPartIndex := range path {
		var currentPathPart = path[pathPartIndex]
		var currentPathRequestPart = requestPath[pathPartIndex]
		var hasParam = strings.HasPrefix(currentPathPart, paramOpen) && strings.HasSuffix(currentPathPart, paramClose)
		// if part of request URI not equals to prefix part, but prefix part has {param}.
		if hasParam {
			// get prefix name without param. Ex: {user} = user.
			currentPathPart = strings.ReplaceAll(currentPathPart, paramOpen, "")
			currentPathPart = strings.ReplaceAll(currentPathPart, paramClose, "")
			// paste prefix and value to map. Ex: /{user}/hello to /69/hello.
			// in summary: paramsMap[user] = 69.
			paramsMap[currentPathPart] = currentPathRequestPart
			verified++
			continue
		}
		var isEquals = strings.EqualFold(currentPathRequestPart, currentPathPart)
		if isEquals {
			verified++
			continue
		} else {
			// if part of requestUriSlice not equals to uriSlice part, and has no {param}, we don't need to continue.
			return false, nil
		}
	}
	return equalCount != verified, paramsMap
}
