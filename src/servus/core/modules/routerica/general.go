package routerica

func New() *Routerica {
	var localRoutes = map[string][]RouteLocal{}
	var routerica = &Routerica{localRoutes: localRoutes}
	routerica.globalMiddleware = &GlobalMiddleware{}
	routerica.routeGroups = append(routerica.routeGroups, newRouteGroup(""))
	routerica.globalMiddleware.middlewareChain = routerica.globalMiddleware
	return routerica
}

func newRouteGroup(prefix string) RouteGroup {
	var localRoutes = map[string][]RouteLocal{}
	var routeGroup = RouteGroup{prefix: prefix, localRoutes: localRoutes}
	routeGroup.middlewareChain = &routeGroup.RouteBase
	return routeGroup
}
