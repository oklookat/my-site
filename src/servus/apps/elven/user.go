package elven

import (
	"github.com/gorilla/mux"
	"net/http"
)

type ResponseUser struct {
	IsAdmin  bool   `json:"is_admin"`
	Username string `json:"username"`
}

type user struct {
	validator userValidator
	route userRoute
}

type userRoute struct {
	middleware *middleware
	validate *userValidator
}

type userValidator struct {
}


func (u *user) boot(m *middleware) {
	u.validator = userValidator{}
	u.route = userRoute{m, &u.validator}
}

func (u *userRoute) boot(router *mux.Router) {
	var authOnly = router.PathPrefix("/users").Subrouter()
	authOnly.Use(u.middleware.authorizedOnly)
	authOnly.HandleFunc("/me", u.getMe).Methods(http.MethodGet)
	authOnly.HandleFunc("/me/change", u.change).Methods(http.MethodPost)
}

