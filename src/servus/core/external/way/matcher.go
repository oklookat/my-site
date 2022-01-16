package way

import "strings"

// match route by request.
func (r *Router) match(requestPath string) (gr *Group, ro *route, params map[string]string) {
	// find group
	for index := range r.groups {
		yes, _params := r.groups[index].match(requestPath)
		if yes {
			params = _params
			gr = r.groups[index]
			break
		}
	}
	// choose points.
	var _routes []*route
	if gr != nil {
		// have route group, use group points.
		_routes = gr.routes
	} else {
		// don't have route group, use router points.
		_routes = r.routes
	}
	// find point
	for index := range _routes {
		yes, paramsTemp := _routes[index].match(requestPath)
		if yes {
			params = mapsToMap(params, paramsTemp)
			ro = _routes[index]
			break
		}
	}
	return
}

// compare two paths, get params.
//
// path: path in way like /hello/{username}.
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
		var inEnd = str[len(str)-1] == '/'
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
