package article

import (
	"net/http"
	"servus/apps/elven/foundation"
	"servus/core"

	"github.com/gorilla/mux"
)

var call *core.Instance

type Instance struct {
	validator validator
	Route     Route
}

type Route struct {
	validate   *validator
	middleware foundation.MiddlewareSafeMethodsOnly
	pipe       foundation.UserPiper
}

type validator struct {
}

func (a *Instance) Boot(
	_core *core.Instance,
	_middleware foundation.MiddlewareSafeMethodsOnly,
	pipe foundation.UserPiper,
) {
	call = _core
	a.validator = validator{}
	a.Route = Route{&a.validator, _middleware, pipe}
}

func (a *Route) Boot(router *mux.Router) {
	var readOnly = router.PathPrefix("/articles").Subrouter()
	readOnly.Use(a.middleware.SafeMethodsOnly)
	readOnly.HandleFunc("", a.getAll).Methods(http.MethodGet)
	readOnly.HandleFunc("/{id}", a.getOne).Methods(http.MethodGet)
	readOnly.HandleFunc("", a.create).Methods(http.MethodPost)
	readOnly.HandleFunc("/{id}", a.update).Methods(http.MethodPut, http.MethodPatch)
	readOnly.HandleFunc("/{id}", a.delete).Methods(http.MethodDelete)
}
