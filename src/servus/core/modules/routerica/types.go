package routerica

import "net/http"

type ctxMiddlewarePipe string

const (
	methodGET                       = "GET"
	methodPOST                      = "POST"
	methodPUT                       = "PUT"
	methodDELETE                    = "DELETE"
	ctxMiddleware ctxMiddlewarePipe = "PIPE_ROUTERICA"
)

type requestsI interface {
	GET(path string, handler func(http.ResponseWriter, *http.Request)) *RouteLocal
	POST(path string, handler func(http.ResponseWriter, *http.Request)) *RouteLocal
	PUT(path string, handler func(http.ResponseWriter, *http.Request)) *RouteLocal
	DELETE(path string, handler func(http.ResponseWriter, *http.Request)) *RouteLocal
}

type MiddlewareChain http.Handler
type Middleware func(next http.Handler) http.Handler

// RouteBase using for typical routes
type RouteBase struct {
	handler         http.HandlerFunc
	middlewareChain MiddlewareChain
}

// GlobalMiddleware used only for global middleware (router instance)
type GlobalMiddleware struct {
	handler         http.HandlerFunc
	middlewareChain MiddlewareChain
}

