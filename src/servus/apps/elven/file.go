package elven

import (
	"github.com/gorilla/mux"
	"net/http"
)

type file struct {
	validator fileValidator
	route fileRoute
}

type fileRoute struct {
	middleware *middleware
	validate *fileValidator
}

type fileValidator struct {
}

func (f *file) boot(m *middleware) {
	f.validator = fileValidator{}
	f.route = fileRoute{m, &f.validator}
}

func (f *fileRoute) boot(router *mux.Router) {
	var fr = router.PathPrefix("/files").Subrouter()
	fr.Use(f.middleware.adminOnly)
	fr.HandleFunc("", f.getAll).Methods(http.MethodGet)
	fr.HandleFunc("", f.createOne).Methods(http.MethodPost)
	fr.HandleFunc("/{id}", f.deleteOne).Methods(http.MethodDelete)
}
