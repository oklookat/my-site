package file

import (
	"net/http"
	"servus/apps/elven/base"
	"servus/core"
	"servus/core/external/way"
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

func (f *Instance) BootRoutes(router *way.Router) {
	var root = router.Group("/files")
	root.Use(f.middleware.AdminOnly)
	root.Route("", f.getAll).Methods(http.MethodGet)
	root.Route("", f.upload).Methods(http.MethodPost)
	root.Route("/{id}", f.deleteOne).Methods(http.MethodDelete)
}
