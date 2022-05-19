package elven

import (
	"fmt"
	"net/http"
	"servus/apps/elven/article"
	"servus/apps/elven/auth"
	"servus/apps/elven/cmd"
	"servus/apps/elven/file"
	"servus/apps/elven/token"
	"servus/apps/elven/user"
	"servus/core"

	"github.com/oklookat/goway"
)

var call *core.Instance
var appMiddleware = &middleware{}

type App struct {
}

func (a *App) Boot(c *core.Instance) {
	call = c
	c.Logger.Info("elven: booting")
	cmd.Boot(call.Logger)

	token.Start(&token.Starter{
		Core: call,
	})
	var tokenPipe = &token.Pipe{}

	user.Start(&user.Starter{
		Core:       call,
		Middleware: appMiddleware,
	})
	var userPipe = &user.Pipe{}

	auth.Start(&auth.Starter{
		Core:       call,
		Middleware: appMiddleware,
		Pipe:       tokenPipe,
	})

	article.Start(&article.Starter{
		Core:       call,
		Middleware: appMiddleware,
		Pipe:       userPipe,
	})

	file.Start(&file.Starter{
		Core:       call,
		Middleware: appMiddleware,
		Pipe:       userPipe,
	})

	a.setupHooks()

	// routes.
	a.setupRoutes()
}

func (a *App) setupHooks() {
	call.Banhammer.OnBanned(func(ip string) {
		var msg = fmt.Sprintf("[#BAN] %v", ip)
		call.Control.SendMessage(msg)
	})
}

func (a *App) setupRoutes() {
	var root = goway.New()
	var elven = root.Group("/elven")

	// provide HTTP helper.
	elven.Use(call.Http.Middleware)

	// limit body.
	elven.Use(call.Limiter.Body.Middleware)

	// ip ban checking.
	elven.Use(call.Banhammer.Middleware)

	// user token.
	elven.Use(appMiddleware.ProvideTokenPipe)

	// user.
	elven.Use(appMiddleware.ProvideUserPipe)

	// start routes.
	auth.StartRoutes(elven)
	article.StartRoutes(elven)
	file.StartRoutes(elven)
	user.StartRoutes(elven)

	// get before router global middlewares.
	var useBeforeRouter = call.Cors.GetMiddleware(root)
	http.Handle("/", useBeforeRouter)
}
