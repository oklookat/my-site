package router

import (
	"strings"
)

// matchGroup - get route group from request path. If success, returns routeGroup and requestURI without routeGroup prefix.
func matchGroup(groups []Group, path string) (group *Group, params map[string]string) {
	for index := range groups {
		group = &groups[index]
		// is not correct group?
		if !strings.HasPrefix(path, group.prefix) {
			continue
		}
		// correct group, get params
		invalid, differ, params := verifyPaths(group.prefix, path, false)
		if invalid || differ {
			continue
		}
		return group, params
	}
	return nil, nil
}

// matchRoute - get route local from routes, request path and method.
func matchRoute(routes routes, path string) (matched *route, params map[string]string) {
	for defPath := range routes {
		invalid, differ, params := verifyPaths(defPath, path, true)
		if invalid || differ {
			continue
		}
		matched = routes[defPath]
		return matched, params
	}
	return nil, nil
}

// verifyPaths - compare two paths, get params.
//
// path: path in router like /hello/{username}.
//
// requestPath: request path like /hello/john.
//
// onlySameLength: verify only if two paths has same length.
//
// returns -
//
// invalid: validation error like one of paths are empty
//
// differ: different paths, makes not sense to verify it
//
// params: like [username: john]. Or nil.
func verifyPaths(path string, requestPath string, onlySameLength bool) (invalid bool, differ bool, params map[string]string) {
	invalid = false
	differ = false
	var cutSlashes = func(str string) string {
		if len(str) < 1 {
			return str
		}
 		var inStart = str[0] == '/'
		var inEnd = str[len(str) - 1] == '/'
		if inStart {
			str = trimFirstRune(str)
		}
		if inEnd {
			str = strings.TrimSuffix(str, "/")
		}
		return str
	}
	path = cutSlashes(path)
	requestPath = cutSlashes(requestPath)
	var pathS = strings.Split(path, "/")
	var requestPathS = strings.Split(requestPath, "/")
	if isEmpty(pathS) || isEmpty(requestPathS) {
		invalid = true
		return
	}
	params = make(map[string]string, 0)
	// get correct length.
	iterator := 0
	pathLen := len(pathS) - 1
	requestPathLen := len(requestPathS) - 1
	if onlySameLength {
		if pathLen != requestPathLen {
			differ = true
			return
		}
		iterator = pathLen
	} else {
		// get lowest length.
		if pathLen > requestPathLen {
			iterator = requestPathLen
		} else {
			// pathLen < requestPathLen.
			iterator = pathLen
		}
	}
	// compare & get params.
	for i := 0; i <= iterator; i++ {
		var pathPart = pathS[i]
		var requestPathPart = requestPathS[i]
		var equals = pathPart == requestPathPart
		if equals {
			continue
		}
		exists, name := getParamName(pathPart)
		if !exists || len(requestPathPart) < 1 {
			// path part not same and no params - 404.
			differ = true
			return
		}
		params[name] = requestPathS[i]
	}
	return
}