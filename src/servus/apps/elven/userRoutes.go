package elven

import (
	"github.com/gorilla/mux"
	"net/http"
)

type userRoutes struct {
}

func (u *userRoutes) boot(router *mux.Router) {
	var ur = router.PathPrefix("/users").Subrouter()
	ur.Use(eBase.middlewareAuthorizedOnly)
	ur.HandleFunc("/me", eUser.controllerGetMe).Methods(http.MethodGet)
	ur.HandleFunc("/me/change", eUser.controllerChange).Methods(http.MethodPost)
}
