package file

import (
	"errors"
	"net/http"
	"servus/apps/elven/base"
	"servus/core"

	"github.com/oklookat/goway"
)

// get paginated files by params.
type GetParams struct {
	// number of page.
	Page int

	// newest (DESC); oldest (ASC).
	Start string

	// created (by creation date).
	By string

	// extensions without dot.
	Extensions []string

	// filename.
	Filename *string
}

var isBooted = false
var call *core.Instance
var middleware base.MiddlewareAdminOnly
var pipe base.UserPipe

type Starter struct {
	Core       *core.Instance
	Middleware base.MiddlewareAdminOnly
	Pipe       base.UserPipe
}

func Start(s *Starter) error {
	// check.
	if s == nil {
		return errors.New("starter nil pointer")
	}
	if s.Core == nil {
		return errors.New("core nil pointer")
	}
	if s.Middleware == nil {
		return errors.New("middleware nil pointer")
	}
	if s.Pipe == nil {
		return errors.New("pipe nil pointer")
	}

	// set.
	call = s.Core
	middleware = s.Middleware
	pipe = s.Pipe

	// ok.
	isBooted = true
	return nil
}

func StartRoutes(router *goway.Router) error {
	if !isBooted {
		return errors.New("you must call Starter.Start() before Starter.Routes()")
	}

	var root = router.Group("/files")
	root.Use(middleware.AdminOnly)
	root.Route("", getAll).Methods(http.MethodGet)
	root.Route("", upload).Methods(http.MethodPost)
	root.Route("/{id}", deleteOne).Methods(http.MethodDelete)

	return nil
}
