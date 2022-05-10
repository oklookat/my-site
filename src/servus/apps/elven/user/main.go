package user

import (
	"errors"
	"net/http"
	"servus/apps/elven/base"
	"servus/core"

	"github.com/oklookat/goway"
)

var isBooted = false
var call *core.Instance
var middleware base.MiddlewareAuthorizedOnly
var pipe = Pipe{}

type Starter struct {
	Core       *core.Instance
	Middleware base.MiddlewareAuthorizedOnly
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

	// set.
	call = s.Core
	middleware = s.Middleware

	// ok.
	isBooted = true
	return nil
}

func StartRoutes(router *goway.Router) error {
	if !isBooted {
		return errors.New("you must call Starter.Start() before Starter.Routes()")
	}

	var currentUser = router.Group("/users/me")
	currentUser.Use(middleware.AuthorizedOnly)
	currentUser.Route("", getMe).Methods(http.MethodGet)
	currentUser.Route("/change", change).Methods(http.MethodPost)

	return nil
}
