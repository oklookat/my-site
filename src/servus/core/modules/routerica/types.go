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

type routericaI interface {
	// Use (G) - add global middleware
	Use(middlewares ...TheMiddleware) *TheGlobalRoute
	POST(path string, handler func(http.ResponseWriter, *http.Request)) *TheRoute
	GET(path string, handler func(http.ResponseWriter, *http.Request)) *TheRoute
	DELETE(path string, handler func(http.ResponseWriter, *http.Request)) *TheRoute
	Group(prefix string)
}

type Routerica struct {
	routericaI
	globalRoute *TheGlobalRoute
	localRoutes      map[string][]TheRoute
}

type TheMiddlewareChain http.Handler
type TheMiddleware func(next http.Handler) http.Handler

type baseMethodI interface {
	// Use (R) - add local middleware
	Use(...TheMiddleware)
}

type BaseMethod struct {
	baseMethodI
	middlewares []TheMiddleware
}

type BaseRoute struct {
	baseMethodI
	handler         http.HandlerFunc
	middlewareChain TheMiddlewareChain
}

type TheGlobalRoute struct {
	BaseRoute
}

type TheRoute struct {
	BaseRoute
	path   string
	method string
}

type TheGroup func()
type RouterGroup struct {
	Prefix string
}
