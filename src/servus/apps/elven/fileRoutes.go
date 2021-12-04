package elven

import (
	"github.com/gorilla/mux"
	"net/http"
)

type fileRoutes struct {
}

func (f *fileRoutes) boot(router *mux.Router) {
	var fr = router.PathPrefix("/files").Subrouter()
	fr.Use(eBase.middlewareAdminOnly)
	fr.HandleFunc("", eFile.controllerGetAll).Methods(http.MethodGet)
	fr.HandleFunc("", eFile.controllerCreateOne).Methods(http.MethodPost)
	fr.HandleFunc("/{id}", eFile.controllerDeleteOne).Methods(http.MethodDelete)
}
