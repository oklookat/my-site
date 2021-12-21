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
	"servus/core/external/errorMan"

	"github.com/gorilla/mux"
)

var call *core.Instance
var requestErrors = errorMan.RequestError{}
var validate = &Validate{}

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
	a.auth.Boot(call, a.middleware, pipeToken, requestErrors, validate)
	// article.
	a.article = &article.Instance{}
	a.article.Boot(call, a.middleware, pipeUser, requestErrors, validate)
	// file.
	a.file = &file.Instance{}
	a.file.Boot(call, a.middleware, pipeUser, requestErrors, validate)
	// user.
	a.user = &user.Instance{}
	a.user.Boot(call, a.middleware, pipeUser, requestErrors, validate)
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
	a.auth.BootRoutes(router)
	a.article.BootRoutes(router)
	a.file.BootRoutes(router)
	a.user.BootRoutes(router)
	//
	var useBeforeRouter = call.Middleware.CORS()(
		call.Middleware.LimitBody()(router))
	http.Handle("/", useBeforeRouter)
}
