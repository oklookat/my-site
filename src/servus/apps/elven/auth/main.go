package auth

import (
	"net/http"
	"servus/apps/elven/base"
	"servus/core"

	"github.com/gorilla/mux"
)

var call *core.Instance
var validate base.Validate

type Instance struct {
	middleware base.MiddlewareAuthorizedOnly
	pipe       base.TokenPiper
	throw      base.RequestError
}

func (a *Instance) Boot(
	_core *core.Instance,
	_middleware base.MiddlewareAuthorizedOnly,
	_pipe base.TokenPiper,
	_throw base.RequestError,
	_validate base.Validate,
) {
	call = _core
	a.middleware = _middleware
	a.pipe = _pipe
	a.throw = _throw
	validate = _validate
}

func (a *Instance) BootRoutes(router *mux.Router) {
	var all = router.PathPrefix("/auth").Subrouter()
	all.HandleFunc("/login", a.login).Methods(http.MethodPost)
	var authOnly = router.PathPrefix("/auth").Subrouter()
	authOnly.Use(a.middleware.AuthorizedOnly)
	authOnly.HandleFunc("/logout", a.logout).Methods(http.MethodPost)
}