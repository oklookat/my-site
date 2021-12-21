package article

import (
	"net/http"
	"servus/apps/elven/base"
	"servus/core"

	"github.com/gorilla/mux"
)

var call *core.Instance
var validate base.Validate

type Instance struct {
	middleware base.MiddlewareSafeMethodsOnly
	pipe       base.UserPiper
	throw      base.RequestError
}

// Boot - initial setup.
func (a *Instance) Boot(
	_core *core.Instance,
	_middleware base.MiddlewareSafeMethodsOnly,
	_pipe base.UserPiper,
	_throw base.RequestError,
	_validate base.Validate,

) {
	call = _core
	a.middleware = _middleware
	a.pipe = _pipe
	a.throw = _throw
	validate = _validate
}

// BootRoutes - add routes to router.
func (a *Instance) BootRoutes(router *mux.Router) {
	var readOnly = router.PathPrefix("/articles").Subrouter()
	readOnly.Use(a.middleware.SafeMethodsOnly)
	readOnly.HandleFunc("", a.getAll).Methods(http.MethodGet)
	readOnly.HandleFunc("/{id}", a.getOne).Methods(http.MethodGet)
	readOnly.HandleFunc("", a.create).Methods(http.MethodPost)
	readOnly.HandleFunc("/{id}", a.update).Methods(http.MethodPut, http.MethodPatch)
	readOnly.HandleFunc("/{id}", a.delete).Methods(http.MethodDelete)
}
