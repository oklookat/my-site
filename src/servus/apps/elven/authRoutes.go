package elven

import (
	"github.com/gorilla/mux"
	"net/http"
)

type authRoutes struct {
}

func (a *authRoutes) boot(router *mux.Router) {
	var main = router.PathPrefix("/auth").Subrouter()
	main.HandleFunc("/login", eAuth.controllerLogin).Methods(http.MethodPost)
	var logout = main.PathPrefix("/logout").Subrouter()
	logout.Use(eBase.middlewareAuthorizedOnly)
	logout.HandleFunc("", eAuth.controllerLogout).Methods(http.MethodPost)
}
