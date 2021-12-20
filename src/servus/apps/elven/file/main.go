package file

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
	validate   *validator
	middleware foundation.MiddlewareAdminOnly
	pipe       foundation.UserPiper
}

type validator struct {
}

func (f *Instance) Boot(
	_core *core.Instance,
	_middleware foundation.MiddlewareAdminOnly,
	_pipe foundation.UserPiper,
) {
	call = _core
	f.validator = validator{}
	f.Route = route{&f.validator, _middleware, _pipe}
}

func (f *route) Boot(router *mux.Router) {
	var fr = router.PathPrefix("/files").Subrouter()
	fr.Use(f.middleware.AdminOnly)
	fr.HandleFunc("", f.getAll).Methods(http.MethodGet)
	fr.HandleFunc("", f.createOne).Methods(http.MethodPost)
	fr.HandleFunc("/{id}", f.deleteOne).Methods(http.MethodDelete)
}
