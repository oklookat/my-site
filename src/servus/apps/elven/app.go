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

type App struct {
	Middleware *middleware
	Auth       *auth.Instance
	Article    *article.Instance
	File       *file.Instance
	User       *user.Instance
}

func (a *App) Boot(c *core.Instance) {
	call = c
	c.Logger.Info("elven: booting")

	// models.
	model.Boot(c)

	// cmd.
	var cmdArgs = &cmd{}
	cmdArgs.boot(a)

	// middleware.
	a.Middleware = &middleware{}

	// pipe.
	pipe.Boot(c)
	var pipeToken = &pipe.Token{}
	var pipeUser = &pipe.User{}

	// auth.
	a.Auth = &auth.Instance{}
	a.Auth.Boot(call, a.Middleware, pipeToken, requestErrors)

	// article.
	a.Article = &article.Instance{}
	a.Article.Boot(call, a.Middleware, pipeUser, requestErrors)

	// file.
	a.File = &file.Instance{}
	a.File.Boot(call, a.Middleware, pipeUser, requestErrors)

	// user.
	a.User = &user.Instance{}
	a.User.Boot(call, a.Middleware, pipeUser, requestErrors)

	// routes.
	a.bootRoutes()

	// TEST HOOKS.
	call.Banhammer.OnWarned(func(ip string) {
		println("[WARN] IP: " + ip)
	})
	call.Banhammer.OnBanned(func(ip string) {
		println("[BAN] IP: " + ip)
	})
}

func (a *App) bootRoutes() {
	router := mux.NewRouter().PathPrefix("/elven").Subrouter()
	router.Use(call.Middleware.ProvideHTTP())
	router.Use(call.Middleware.AsJson())
	router.Use(a.Middleware.ProvideTokenPipe)
	router.Use(a.Middleware.ProvideUserPipe)
	//
	a.Auth.BootRoutes(router)
	a.Article.BootRoutes(router)
	a.File.BootRoutes(router)
	a.User.BootRoutes(router)
	//
	call.Banhammer.GetMiddleware()

	var corsMiddleware = call.Middleware.CORS()
	var limitBodyMiddleware = call.Middleware.LimitBody()
	var banhammerMiddleware = call.Banhammer.GetMiddleware()

	var useBeforeRouter = corsMiddleware(
		limitBodyMiddleware(banhammerMiddleware(router)),
	)
	http.Handle("/", useBeforeRouter)
}
