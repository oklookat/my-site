package auth

import (
	"net/http"
	"servus/apps/elven/base"
	"servus/core"
	"servus/core/external/way"
)

var call *core.Instance

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
) {
	call = _core
	a.middleware = _middleware
	a.pipe = _pipe
	a.throw = _throw
}

func (a *Instance) BootRoutes(router *way.Router) {
	// login
	router.Route("/auth", a.login).Methods(http.MethodPost)

	// logout
	var logout = router.Route("/auth", a.logout).Methods(http.MethodPost)
	logout.Use(a.middleware.AuthorizedOnly)
}
