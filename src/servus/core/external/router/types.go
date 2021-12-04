package router

import "net/http"

type ctxPathParamsPipe string

const (
	methodGET                       = "GET"
	methodHEAD                      = "HEAD"
	methodPOST                      = "POST"
	methodPUT                       = "PUT"
	methodDELETE                    = "DELETE"
	methodOPTIONS                   = "OPTIONS"
	methodPATCH                     = "PATCH"
	paramOpen                       = "{"
	paramClose                      = "}"
	ctxPathParams ctxPathParamsPipe = "ROUTER_PIPE_PATH_PARAMS"
)

type MiddlewareChain http.Handler
type Middleware func(next http.Handler) http.Handler
