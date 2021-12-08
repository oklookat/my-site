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
	controller userController
}

type userController struct {
	*baseController
	validate *userValidator
}

type userValidator struct {
}


func (u *user) boot(b *baseController) {
	u.validator = userValidator{}
	var controller = userController{b, &u.validator}
	u.route = userRoute{controller}
}

func (u *userRoute) boot(router *mux.Router) {
	var authOnly = router.PathPrefix("/users").Subrouter()
	authOnly.Use(u.controller.middlewareAuthorizedOnly)
	authOnly.HandleFunc("/me", u.controller.getMe).Methods(http.MethodGet)
	authOnly.HandleFunc("/me/change", u.controller.change).Methods(http.MethodPost)
}

