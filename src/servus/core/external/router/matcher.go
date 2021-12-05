package router

type matchedRoute struct {
	notFound         bool
	methodNotAllowed bool
	group            *Group
	route            *route
	method           *RouteMethod
	params           map[string]string
}

// match - match route by request.
func (r *Router) match(requestMethod string, requestPath string) (matched *matchedRoute) {
	matched = &matchedRoute{}
	// TODO: improve {params} matching. Use regex like {([a-zA-Z]*)} ?
	routeGroup, params := matchGroup(r.groups, requestPath)
	var _routes routes
	if routeGroup != nil {
		// we have route group, need to find route inside this group
		matched.group = routeGroup
		_routes = routeGroup.routes
	} else {
		// we don't have route group, need to find route without group.
		_routes = r.routes
	}
	// find route in routes
	_route, paramsTemp := matchRoute(_routes, requestPath)
	params = mapsToMap(params, paramsTemp)
	matched.notFound = _route == nil
	if matched.notFound {
		return
	}
	var hasMethod = _route.isHasMethod(requestMethod)
	if hasMethod {
		matched.method = _route.methods[requestMethod]
	}
	matched.methodNotAllowed = !hasMethod
	matched.route = _route
	matched.params = params
	return
}
