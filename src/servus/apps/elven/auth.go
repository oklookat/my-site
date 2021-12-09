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
	middleware *middleware
	validate *authValidator
}

type authValidator struct {
}

func (a *auth) boot(m *middleware) {
	a.validator = authValidator{}
	a.route = authRoute{m, &a.validator}
}

func (a *authRoute) boot(router *mux.Router) {
	var all = router.PathPrefix("/auth").Subrouter()
	all.HandleFunc("/login", a.login).Methods(http.MethodPost)
	var authOnly = router.PathPrefix("/auth").Subrouter()
	authOnly.Use(a.middleware.authorizedOnly)
	authOnly.HandleFunc("/logout", a.logout).Methods(http.MethodPost)
}
