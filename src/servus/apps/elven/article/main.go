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

// add routes to router.
func (a *Instance) BootRoutes(router *mux.Router) {
	var root = router.PathPrefix("/article").Subrouter()
	// articles
	var articles = root.PathPrefix("/articles").Subrouter()
	articles.Use(a.middleware.SafeMethodsOnly)
	articles.HandleFunc("", a.getArticles).Methods(http.MethodGet)
	articles.HandleFunc("", a.createArticle).Methods(http.MethodPost)
	articles.HandleFunc("/{id}", a.getArticle).Methods(http.MethodGet)
	articles.HandleFunc("/{id}", a.updateArticle).Methods(http.MethodPut, http.MethodPatch)
	articles.HandleFunc("/{id}", a.deleteArticle).Methods(http.MethodDelete)
	// categories
	var categories = root.PathPrefix("/categories").Subrouter()
	categories.Use(a.middleware.SafeMethodsOnly)
	categories.HandleFunc("", a.getCategories).Methods(http.MethodGet)
	categories.HandleFunc("", a.addCategory).Methods(http.MethodPost)
	categories.HandleFunc("/{id}", a.getCategory).Methods(http.MethodGet)
	categories.HandleFunc("/{id}", a.renameCategory).Methods(http.MethodPut, http.MethodPatch)
	categories.HandleFunc("/{id}", a.deleteCategory).Methods(http.MethodDelete)
}
