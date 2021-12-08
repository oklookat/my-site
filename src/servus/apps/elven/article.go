package elven

import (
	"github.com/gorilla/mux"
	"net/http"
)

type article struct {
	validator articleValidator
	route articleRoute
}

type articleRoute struct {
	controller articleController
}

type articleController struct {
	*baseController
	validate *articleValidator
}

type articleValidator struct {
}


func (a *article) boot(b *baseController) {
	a.validator = articleValidator{}
	var controller = articleController{b, &a.validator}
	a.route = articleRoute{controller}
}

func (a *articleRoute) boot(router *mux.Router) {
	var readOnly = router.PathPrefix("/articles").Subrouter()
	readOnly.Use(a.controller.middlewareReadOnly)
	readOnly.HandleFunc("", a.controller.getAll).Methods(http.MethodGet)
	readOnly.HandleFunc("/{id}", a.controller.getOne).Methods(http.MethodGet)
	var adminOnly = router.PathPrefix("/articles").Subrouter()
	adminOnly.Use(a.controller.middlewareAdminOnly)
	adminOnly.HandleFunc("", a.controller.create).Methods(http.MethodPost)
	adminOnly.HandleFunc("/{id}", a.controller.update).Methods(http.MethodPut, http.MethodPatch)
	adminOnly.HandleFunc("/{id}", a.controller.delete).Methods(http.MethodDelete)
}