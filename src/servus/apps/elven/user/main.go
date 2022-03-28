package user

import (
	"net/http"
	"servus/apps/elven/base"
	"servus/core"
	"servus/core/external/way"
)

var call *core.Instance

type Instance struct {
	middleware base.MiddlewareAuthorizedOnly
	pipe       base.UserPiper
	throw      base.RequestError
}

func (u *Instance) Boot(
	_core *core.Instance,
	_middleware base.MiddlewareAuthorizedOnly,
	_pipe base.UserPiper,
	_throw base.RequestError,
) {
	call = _core
	u.middleware = _middleware
	u.pipe = _pipe
	u.throw = _throw
}

func (u *Instance) BootRoutes(router *way.Router) {
	// current user
	var currentUser = router.Group("/users/me")
	currentUser.Use(u.middleware.AuthorizedOnly)
	currentUser.Route("", u.getMe).Methods(http.MethodGet)
	currentUser.Route("/change", u.change).Methods(http.MethodPost)
}
