package elven

import (
	"github.com/gorilla/mux"
	"net/http"
)

type articleRoutes struct {
}

func (a *articleRoutes) boot(router *mux.Router) {
	var ar = router.PathPrefix("/articles").Subrouter()
	ar.Use(eBase.middlewareReadOnly)
	ar.HandleFunc("", eArticle.getAll).Methods(http.MethodGet)
	ar.HandleFunc("/{id}", eArticle.getOne).Methods(http.MethodGet)
	ar.HandleFunc("", eArticle.create).Methods(http.MethodPost)
	ar.HandleFunc("/{id}", eArticle.update).Methods(http.MethodPut, http.MethodPatch)
	ar.HandleFunc("/{id}", eArticle.delete).Methods(http.MethodDelete)
}
