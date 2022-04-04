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

	// provide HTTP helper.
	elven.Use(call.Http.Middleware)

	// set content-type.
	elven.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			request.Header.Add("Content-Type", "application/json")
			next.ServeHTTP(response, request)
		})
	})

	// limit body.
	elven.Use(call.Limiter.Body.Middleware)

	// ip banner.
	elven.Use(call.Banhammer.Middleware)

	// user token.
	elven.Use(a.Middleware.ProvideTokenPipe)

	// user.
	elven.Use(a.Middleware.ProvideUserPipe)

	// start routes.
	a.Auth.Routes(elven)
	a.Article.Routes(elven)
	a.File.Routes(elven)
	a.User.Routes(elven)

	// get before router global middlewares.
	var useBeforeRouter = call.Cors.GetMiddleware(root)
	http.Handle("/", useBeforeRouter)
}
