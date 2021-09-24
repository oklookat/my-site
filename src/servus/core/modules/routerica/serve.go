package routerica

import (
	"context"
	"net/http"
)

// ServeHTTP (R) - when a request comes in and goes out, it will be here
func (r *Routerica) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// the function is called twice. One for /favicon, and one for the route
	// if favicon - we don't handle request
	if request.RequestURI == "/favicon.ico" {
		return
	}
	println("i am in entry point")
	// pass routerica instance to request
	request = request.WithContext(context.WithValue(context.Background(), ctxMiddleware, r))
	r.globalMiddleware.middlewareChain.ServeHTTP(response, request)
	println("exit from entry point")
}

// ServeHTTP (G) - when global middleware finished or if no global middlewares, it calls this method. If global middleware send response, this method will not be called.
func (g *GlobalMiddleware) ServeHTTP(response http.ResponseWriter, request *http.Request){
	var routerica = request.Context().Value(ctxMiddleware)
	// clear context
	request = request.WithContext(context.Background())
	routeMatcher(routerica.(*Routerica), response, request)
}


func (b *RouteBase) ServeHTTP(response http.ResponseWriter, request *http.Request){

}
