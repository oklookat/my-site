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
	pipe       base.TokenPipe
	throw      base.RequestError
}

func (a *Instance) Boot(
	_core *core.Instance,
	_middleware base.MiddlewareAuthorizedOnly,
	_pipe base.TokenPipe,
	_throw base.RequestError,
) {
	call = _core
	a.middleware = _middleware
	a.pipe = _pipe
	a.throw = _throw
}

func (a *Instance) BootRoutes(router *way.Router) {

	var authGroup = router.Group("/auth")

	// login
	authGroup.Route("/login", a.login).Methods(http.MethodPost)

	// logout
	var logout = authGroup.Route("/logout", a.logout).Methods(http.MethodPost)
	logout.Use(a.middleware.AuthorizedOnly)
}
