package elven

import (
	"net/http"
	"servus/apps/elven/article"
	"servus/apps/elven/auth"
	"servus/apps/elven/file"
	"servus/apps/elven/model"
	"servus/apps/elven/pipe"
	"servus/apps/elven/user"
	"servus/core"

	"github.com/gorilla/mux"
)

var call *core.Instance

type App struct {
	middleware *middleware
	auth       *auth.Instance
	article    *article.Instance
	file       *file.Instance
	user       *user.Instance
}

func (a *App) Boot(c *core.Instance) {
	call = c
	c.Logger.Info("elven: booting")
	var _cmd = &cmd{}
	_cmd.boot()
	a.middleware = &middleware{}
	// models.
	model.Boot(c)
	// pipe.
	pipe.Boot(c)
	var pipeToken = &pipe.Token{}
	var pipeUser = &pipe.User{}
	// auth.
	a.auth = &auth.Instance{}
	a.auth.Boot(call, a.middleware, pipeToken)
	// article.
	a.article = &article.Instance{}
	a.article.Boot(call, a.middleware, pipeUser)
	// file.
	a.file = &file.Instance{}
	a.file.Boot(call, a.middleware, pipeUser)
	// user.
	a.user = &user.Instance{}
	a.user.Boot(call, a.middleware, pipeUser)
	//
	a.bootRoutes()
}

func (a *App) bootRoutes() {
	router := mux.NewRouter().PathPrefix("/elven").Subrouter()
	router.Use(call.Middleware.ProvideHTTP())
	router.Use(call.Middleware.AsJson())
	router.Use(a.middleware.ProvideTokenPipe)
	router.Use(a.middleware.ProvideUserPipe)
	//
	a.auth.Route.Boot(router)
	a.article.Route.Boot(router)
	a.file.Route.Boot(router)
	a.user.Route.Boot(router)
	//
	var useBeforeRouter = call.Middleware.CORS()(
		call.Middleware.LimitBody()(router))
	http.Handle("/", useBeforeRouter)
}
