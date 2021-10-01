package routerica

import "net/http"

type CtxInternalPipe string
type CtxValueParamsPipe string

const (
	methodGET                         = "GET"
	methodPOST                        = "POST"
	methodPUT                         = "PUT"
	methodDELETE                      = "DELETE"
	paramOpen                         = "{"
	paramClose                        = "}"
	CtxInternal CtxInternalPipe    = "ROUTERICA_PIPE_INTERNAL"
	CtxValueParams CtxValueParamsPipe = "ROUTERICA_PIPE_PARAMS"
)

type requestsI interface {
	GET(path string, endpoint http.HandlerFunc) *RouteLocal
	POST(path string, endpoint http.HandlerFunc) *RouteLocal
	PUT(path string, endpoint http.HandlerFunc) *RouteLocal
	DELETE(path string, endpoint http.HandlerFunc) *RouteLocal
}

type MiddlewareChain http.Handler
type Middleware func(next http.Handler) http.Handler
