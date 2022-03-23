package core

import (
	"net/http"
	"servus/core/internal/cors"
	"servus/core/internal/limiter"
	"servus/core/internal/middleware"

	"github.com/gorilla/mux"
)

func (i *Instance) setupMiddleware() {
	// cors.
	var corsInstance = cors.New(i.Config.Security.CORS)

	// limiter.
	var bodyLimiter = limiter.NewBody(i.Config.Security.Limiter.Body)

	// http.
	var httpHelp = &httpHelper{}
	var variablesGetter httpParamsGetter = func(r *http.Request) map[string]string {
		return mux.Vars(r)
	}
	httpHelp.new(i.Logger, i.Control, i.Config.Security.Cookie, variablesGetter)

	// middleware.
	var md = &middleware.Instance{}
	md.New(corsInstance.Middleware, bodyLimiter.Middleware, httpHelp.middleware)
	i.Middleware = md
}
