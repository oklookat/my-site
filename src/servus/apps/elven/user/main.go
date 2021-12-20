package user

import (
	"net/http"
	"servus/apps/elven/foundation"
	"servus/core"

	"github.com/gorilla/mux"
)

var call *core.Instance

type Instance struct {
	validator validator
	Route     route
}

type route struct {
	middleware foundation.MiddlewareAuthorizedOnly
	pipe       foundation.UserPiper
}

type validator struct {
}

func (u *Instance) Boot(
	_core *core.Instance,
	_authOnly foundation.MiddlewareAuthorizedOnly,
	_userPipe foundation.UserPiper,
) {
	call = _core
	u.validator = validator{}
	u.Route = route{
		middleware: _authOnly,
		pipe:       _userPipe,
	}
}

func (u *route) Boot(router *mux.Router) {
	var authOnly = router.PathPrefix("/users").Subrouter()
	authOnly.Use(u.middleware.AuthorizedOnly)
	authOnly.HandleFunc("/me", u.getMe).Methods(http.MethodGet)
	authOnly.HandleFunc("/me/change", u.change).Methods(http.MethodPost)
}
