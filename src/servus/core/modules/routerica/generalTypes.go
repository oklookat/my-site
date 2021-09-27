package routerica

import "net/http"

type ctxMainPipe string
type ctxPathParamsPipe string

const (
	methodGET                             = "GET"
	methodPOST                            = "POST"
	methodPUT                             = "PUT"
	methodDELETE                          = "DELETE"
	paramDelimiterOpen                    = "{"
	paramDelimiterClose                   = "}"
	ctxMain             ctxMainPipe       = "PIPE_MAIN"
	ctxPathParams       ctxPathParamsPipe = "PIPE_PATH_PARAMS"
)

type requestsI interface {
	GET(path string, endpoint http.HandlerFunc) *RouteLocal
	POST(path string, endpoint http.HandlerFunc) *RouteLocal
	PUT(path string, endpoint http.HandlerFunc) *RouteLocal
	DELETE(path string, endpoint http.HandlerFunc) *RouteLocal
}

type MiddlewareChain http.Handler
type Middleware func(next http.Handler) http.Handler
