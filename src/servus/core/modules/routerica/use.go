package routerica


// Use middleware (global). That method executes when you're adding global middleware to the router instance
func (r *Routerica) Use(middlewares ...TheMiddleware) {
	r.globalRoute.middlewareChain = middlewareChainer(middlewares, r.globalRoute)
}


// Use middleware (request methods). That method executes when you're adding middleware to request method
func (r *TheRoute) Use(middlewares ...TheMiddleware) {
	r.middlewareChain = middlewareChainer(middlewares, r.handler)
}
