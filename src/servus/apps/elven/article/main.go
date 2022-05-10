package article

import (
	"errors"
	"net/http"
	"servus/apps/elven/base"
	"servus/core"

	"github.com/oklookat/goway"
)

// get paginated articles by params.
type GetParams struct {
	// number of page.
	Page int

	// show published or drafts?
	Drafts bool

	// start from newest? true == DESC; false == ASC.
	Newest bool

	// created; updated; published.
	By string

	// search by title.
	Title *string
}

// article request body that user should send to create/update article.
type Body struct {
	CoverID     *string `json:"cover_id"`
	IsPublished *bool   `json:"is_published"`
	Title       *string `json:"title"`
	Content     *string `json:"content"`
}

var isBooted = false
var call *core.Instance
var middleware base.MiddlewareSafeMethodsOnly
var pipe base.UserPipe

type Starter struct {
	Core       *core.Instance
	Middleware base.MiddlewareSafeMethodsOnly
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
		return errors.New("you must call Start() before StartRoutes()")
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
