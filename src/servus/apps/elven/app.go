package elven

import (
	"github.com/gorilla/mux"
	"net/http"
	"servus/core"
)

var call *core.Core

// ResponseContent - template for response.
type ResponseContent struct {
	Meta struct {
		PerPage     int `json:"per_page"`
		TotalPages  int `json:"total_pages"`
		CurrentPage int `json:"current_page"`
	} `json:"meta"`
	Data interface{} `json:"data"`
}

type App struct {
	middleware *middleware
	auth           *auth
	article        *article
	file           *file
	user           *user
}

func (a *App) Boot(c *core.Core) {
	call = c
	c.Logger.Info("elven: booting")
	var _cmd = &cmd{}
	_cmd.boot()
	a.middleware = &middleware{}
	//
	a.auth = &auth{}
	a.auth.boot(a.middleware)
	a.article = &article{}
	a.article.boot(a.middleware)
	a.file = &file{}
	a.file.boot(a.middleware)
	a.user = &user{}
	a.user.boot(a.middleware)
	//
	a.bootRoutes()
}

func (a *App) bootRoutes() {
	router := mux.NewRouter().PathPrefix("/elven").Subrouter()
	router.Use(call.Middleware.ProvideHTTP)
	router.Use(call.Middleware.AsJson)
	//
	a.auth.route.boot(router)
	a.article.route.boot(router)
	a.file.route.boot(router)
	a.user.route.boot(router)
	//
	var useBeforeRouter = call.Middleware.CORS(call.Middleware.LimitBody(router))
	http.Handle("/", useBeforeRouter)
}
