package auth

import (
	"net/http"
	"servus/apps/elven/foundation"
	"servus/core"

	"github.com/gorilla/mux"
)

var call *core.Instance

type Instance struct {
	validator authValidator
	Route     route
}

type route struct {
	validate   *authValidator
	middleware foundation.MiddlewareAuthorizedOnly
	pipe       foundation.TokenPiper
}

type authValidator struct {
}

func (a *Instance) Boot(
	_core *core.Instance,
	_middleware foundation.MiddlewareAuthorizedOnly,
	_pipe foundation.TokenPiper,
) {
	call = _core
	a.validator = authValidator{}
	a.Route = route{&a.validator, _middleware, _pipe}
}

func (a *route) Boot(router *mux.Router) {
	var all = router.PathPrefix("/auth").Subrouter()
	all.HandleFunc("/login", a.login).Methods(http.MethodPost)
	var authOnly = router.PathPrefix("/auth").Subrouter()
	authOnly.Use(a.middleware.AuthorizedOnly)
	authOnly.HandleFunc("/logout", a.logout).Methods(http.MethodPost)
}
