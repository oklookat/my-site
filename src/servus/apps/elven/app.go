package elven

import (
	"net/http"
	"servus/apps/elven/article"
	"servus/apps/elven/auth"
	"servus/apps/elven/cmd"
	"servus/apps/elven/file"
	"servus/apps/elven/model"
	"servus/apps/elven/pipe"
	"servus/apps/elven/user"
	"servus/core"
	"servus/core/external/errorMan"

	"github.com/oklookat/goway"
)

var call *core.Instance
var requestErrors = errorMan.RequestError{}

type App struct {
	Middleware *middleware
	Auth       *auth.Starter
	Article    *article.Starter
	File       *file.Starter
	User       *user.Starter
}

func (a *App) Boot(c *core.Instance) {
	call = c
	c.Logger.Info("elven: booting")

	// models.
	model.Boot(c)

	// cmd.
	var cmdConfig = &cmd.Config{
		Logger: call.Logger,
	}
	cmd.Boot(cmdConfig)

	// middleware.
	a.Middleware = &middleware{}

	// pipe.
	pipe.Boot(c)
	var pipeToken = &pipe.Token{}
	var pipeUser = &pipe.User{}

	// auth.
	var auth = &auth.Starter{
		Core:       call,
		Middleware: a.Middleware,
		Pipe:       pipeToken,
		Throw:      requestErrors,
	}
	auth.Start()
	a.Auth = auth

	// article.
	var article = &article.Starter{
		Core:       call,
		Middleware: a.Middleware,
		Pipe:       pipeUser,
		Throw:      requestErrors,
	}
	article.Start()
	a.Article = article

	// file.
	var file = &file.Starter{
		Core:       call,
		Middleware: a.Middleware,
		Pipe:       pipeUser,
		Throw:      requestErrors,
	}
	file.Start()
	a.File = file

	// user.
	var user = &user.Starter{
		Core:       call,
		Middleware: a.Middleware,
		Pipe:       pipeUser,
		Throw:      requestErrors,
	}
	user.Start()
	a.User = user

	// routes.
	a.bootRoutes()
}

func (a *App) bootRoutes() {
	var root = goway.New()
	var elven = root.Group("/elven")

	// add global middleware.
	elven.Use(call.Middleware.ProvideHTTP())
	elven.Use(call.Middleware.AsJson())
	elven.Use(call.Middleware.LimitBody())
	elven.Use(call.Banhammer.GetMiddleware())
	elven.Use(a.Middleware.ProvideTokenPipe)
	elven.Use(a.Middleware.ProvideUserPipe)

	// start routes.
	a.Auth.Routes(elven)
	a.Article.Routes(elven)
	a.File.Routes(elven)
	a.User.Routes(elven)

	// get before router global middlewares.
	var corsMiddleware = call.Middleware.CORS()
	var useBeforeRouter = corsMiddleware(root)
	http.Handle("/", useBeforeRouter)
}
