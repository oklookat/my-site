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
	"servus/core/external/way"
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
}

func (a *App) bootRoutes() {
	var router = way.New()
	var root = router.Group("/elven")
	// add global middleware.
	root.Use(call.Middleware.ProvideHTTP())
	root.Use(call.Middleware.AsJson())
	root.Use(a.Middleware.ProvideTokenPipe)
	root.Use(a.Middleware.ProvideUserPipe)
	// provide routes..
	a.Auth.BootRoutes(root)
	a.Article.BootRoutes(root)
	a.File.BootRoutes(root)
	a.User.BootRoutes(root)
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
