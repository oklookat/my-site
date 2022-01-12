package user

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
	pipe       base.UserPiper
	throw      base.RequestError
}

func (u *Instance) Boot(
	_core *core.Instance,
	_middleware base.MiddlewareAuthorizedOnly,
	_pipe base.UserPiper,
	_throw base.RequestError,
	_validate base.Validate,
) {
	call = _core
	u.middleware = _middleware
	u.pipe = _pipe
	u.throw = _throw
	validate = _validate
}

func (u *Instance) BootRoutes(router *mux.Router) {
	var authOnly = router.PathPrefix("/users").Subrouter()
	authOnly.Use(u.middleware.AuthorizedOnly)
	authOnly.HandleFunc("/me", u.getMe).Methods(http.MethodGet)
	authOnly.HandleFunc("/me/change", u.change).Methods(http.MethodPost)
}