package routerica

func New() *Routerica {
	var routerica = &Routerica{localRoutes: map[string][]TheRoute{}}
	routerica.globalRoute = &TheGlobalRoute{}
	routerica.globalRoute.middlewareChain = routerica.globalRoute
	return routerica
}
