package article

import (
	"net/http"
	"servus/apps/elven/base"
	"servus/core"
	"servus/core/external/way"
)

var call *core.Instance

type Instance struct {
	middleware base.MiddlewareSafeMethodsOnly
	pipe       base.UserPipe
	throw      base.RequestError
}

func (a *Instance) Boot(
	_core *core.Instance,
	_middleware base.MiddlewareSafeMethodsOnly,
	_pipe base.UserPipe,
	_throw base.RequestError,

) {
	call = _core
	a.middleware = _middleware
	a.pipe = _pipe
	a.throw = _throw
}

// add routes to router.
func (a *Instance) BootRoutes(router *way.Router) {
	var root = router.Group("/article")

	// articles | /article/articles
	var articles = root.Group("/articles")
	articles.Use(a.middleware.SafeMethodsOnly)
	articles.Route("", a.getArticles).Methods(http.MethodGet)
	articles.Route("", a.createArticle).Methods(http.MethodPost)
	articles.Route("/{id}", a.getArticle).Methods(http.MethodGet)
	articles.Route("/{id}", a.updateArticle).Methods(http.MethodPut, http.MethodPatch)
	articles.Route("/{id}", a.deleteArticle).Methods(http.MethodDelete)

	// categories | /article/categories
	var categories = root.Group("/categories")
	categories.Use(a.middleware.SafeMethodsOnly)
	categories.Route("", a.getCategories).Methods(http.MethodGet)
	categories.Route("", a.addCategory).Methods(http.MethodPost)
	categories.Route("/{id}", a.getCategory).Methods(http.MethodGet)
	categories.Route("/{id}", a.renameCategory).Methods(http.MethodPut, http.MethodPatch)
	categories.Route("/{id}", a.deleteCategory).Methods(http.MethodDelete)
}
