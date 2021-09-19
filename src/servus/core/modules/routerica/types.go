package routerica

import "net/http"

const (
	methodGet    = "GET"
	methodPost   = "POST"
	methodPut    = "PUT"
	methodDelete = "DELETE"
)

type routericaI interface {
	Post(path string, handler func(http.ResponseWriter, *http.Request))
	Get(path string, handler func(http.ResponseWriter, *http.Request)) *GetRoute
	Delete(path string, handler func(http.ResponseWriter, *http.Request))
	Group(prefix string)
}

type Routerica struct {
	getRoutes []*GetRoute
	routericaI
}
// TODO: check when user has no middlewares. For now we get a panic, when no middlewares.
type TheMiddlewareChain http.Handler
type TheMiddleware func(next http.Handler) http.Handler

type baseMethodI interface {
	Use(...TheMiddleware)
}

type BaseMethod struct {
	middlewares []TheMiddleware
	baseMethodI
}

type theRoute struct {
	baseMethodI
	path            string
	method          string
	handler         http.HandlerFunc
	middlewareChain TheMiddlewareChain
}

type GetRoute struct {
	theRoute
}

type TheGroup func()
type RouterGroup struct {
	Prefix string
}
