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
	controller fileController
}

type fileController struct {
	*baseController
	validate *fileValidator
}

type fileValidator struct {

}

func (f *file) boot(b *baseController) {
	f.validator = fileValidator{}
	var controller = fileController{b, &f.validator}
	f.route = fileRoute{controller}
}

func (f *fileRoute) boot(router *mux.Router) {
	var fr = router.PathPrefix("/files").Subrouter()
	fr.Use(f.controller.middlewareAdminOnly)
	fr.HandleFunc("", f.controller.getAll).Methods(http.MethodGet)
	fr.HandleFunc("", f.controller.createOne).Methods(http.MethodPost)
	fr.HandleFunc("/{id}", f.controller.deleteOne).Methods(http.MethodDelete)
}
