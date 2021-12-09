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
	middleware *middleware
	validate *articleValidator
}

type articleValidator struct {
}


func (a *article) boot(m *middleware) {
	a.validator = articleValidator{}
	a.route = articleRoute{m, &a.validator}
}

func (a *articleRoute) boot(router *mux.Router) {
	var readOnly = router.PathPrefix("/articles").Subrouter()
	readOnly.Use(a.middleware.safeMethodsOnly)
	readOnly.HandleFunc("", a.getAll).Methods(http.MethodGet)
	readOnly.HandleFunc("/{id}", a.getOne).Methods(http.MethodGet)
	readOnly.HandleFunc("", a.create).Methods(http.MethodPost)
	readOnly.HandleFunc("/{id}", a.update).Methods(http.MethodPut, http.MethodPatch)
	readOnly.HandleFunc("/{id}", a.delete).Methods(http.MethodDelete)
}