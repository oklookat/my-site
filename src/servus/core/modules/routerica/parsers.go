package routerica

import "strings"

// parserPath - if url path equals requestPath returns true, map like [param: value] {if delimiter}, and requestPath without matched.
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
	for prefixIndex := range path {
		var currentPref = path[prefixIndex]
		var isEquals = strings.EqualFold(requestPath[prefixIndex], currentPref)
		if isEquals {
			verified++
			continue
		} else {
			var hasDelimiter = strings.HasPrefix(currentPref, paramDelimiterOpen) && strings.HasSuffix(currentPref, paramDelimiterClose)
			// if part of request URI not equals to prefix part, but prefix part has {delimiter}.
			if hasDelimiter {
				// get prefix name without delimiter. Ex: {user} = user.
				currentPref = strings.ReplaceAll(currentPref, paramDelimiterOpen, "")
				currentPref = strings.ReplaceAll(currentPref, paramDelimiterClose, "")
				// paste prefix and value to map. Ex: /{user}/hello to /69/hello.
				// in summary: paramsMap[user] = 69.
				paramsMap[currentPref] = requestPath[prefixIndex]
				verified++
				continue
			} else {
				// if part of requestUriSlice not equals to uriSlice part, and has no {delimiter}, we don't need to continue.
				break
			}
		}
	}
	return equalCount != verified, paramsMap
}

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
		// now we in correct route
		return currentRoute, paramsMap
	}
	return nil, nil
}

