package file

import (
	"net/http"
	"servus/apps/elven/base"
	"servus/core"

	"github.com/gorilla/mux"
)

var call *core.Instance

type Instance struct {
	middleware base.MiddlewareAdminOnly
	pipe       base.UserPiper
	throw      base.RequestError
}

func (f *Instance) Boot(
	_core *core.Instance,
	_middleware base.MiddlewareAdminOnly,
	_pipe base.UserPiper,
	_throw base.RequestError,
) {
	call = _core
	f.middleware = _middleware
	f.pipe = _pipe
	f.throw = _throw
}

func (f *Instance) BootRoutes(router *mux.Router) {
	var fr = router.PathPrefix("/files").Subrouter()
	fr.Use(f.middleware.AdminOnly)
	fr.HandleFunc("", f.getAll).Methods(http.MethodGet)
	fr.HandleFunc("", f.upload).Methods(http.MethodPost)
	fr.HandleFunc("/{id}", f.deleteOne).Methods(http.MethodDelete)
}
