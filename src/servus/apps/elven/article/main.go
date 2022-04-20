package article

import (
	"errors"
	"net/http"
	"servus/apps/elven/base"
	"servus/core"

	"github.com/oklookat/goway"
)

var isBooted = false
var call *core.Instance
var middleware base.MiddlewareSafeMethodsOnly
var pipe base.UserPipe

type Starter struct {
	Core       *core.Instance
	Middleware base.MiddlewareSafeMethodsOnly
	Pipe       base.UserPipe
}

func (s *Starter) Start() error {
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

func (s *Starter) Routes(router *goway.Router) error {
	if !isBooted {
		return errors.New("you must call Starter.Boot() before Starter.Routes()")
	}

	var root = router.Group("/article")

	// articles | /article/articles
	var articles = root.Group("/articles")
	articles.Use(middleware.SafeMethodsOnly)
	articles.Route("", getArticles).Methods(http.MethodGet)
	articles.Route("", createArticle).Methods(http.MethodPost)
	articles.Route("/{id}", getArticle).Methods(http.MethodGet)
	articles.Route("/{id}", updateArticle).Methods(http.MethodPut, http.MethodPatch)
	articles.Route("/{id}", deleteArticle).Methods(http.MethodDelete)

	return nil
}
