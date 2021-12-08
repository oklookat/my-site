package elven

import (
	"github.com/gorilla/mux"
	"net/http"
)

type auth struct {
	validator authValidator
	route     authRoute
}

type authRoute struct {
	controller authController
}

type authController struct {
	*baseController
	validate *authValidator
}

type authValidator struct {
}

func (a *auth) boot(b *baseController) {
	a.validator = authValidator{}
	var controller = authController{b, &a.validator}
	a.route = authRoute{controller}
}

func (a *authRoute) boot(router *mux.Router) {
	var all = router.PathPrefix("/auth").Subrouter()
	all.HandleFunc("/login", a.controller.login).Methods(http.MethodPost)
	var authOnly = router.PathPrefix("/auth").Subrouter()
	authOnly.Use(a.controller.middlewareAuthorizedOnly)
	authOnly.HandleFunc("/logout", a.controller.logout).Methods(http.MethodPost)
}
