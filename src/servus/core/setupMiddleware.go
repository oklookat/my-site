package core

import (
	"net/http"
	"servus/core/external/way"
	"servus/core/internal/cors"
	"servus/core/internal/limiter"
	"servus/core/internal/middleware"
)

func (i *Instance) setupMiddleware() error {
	var err error

	// cors.
	var corsInstance = cors.New(i.Config.Security.CORS)

	// limiter.
	var bodyLimiter = limiter.NewBody(i.Config.Security.Limiter.Body)

	// http.
	var httpHelp = &httpHelper{}
	var variablesGetter httpParamsGetter = func(r *http.Request) map[string]string {
		return way.Vars(r)
	}
	httpHelp.new(i.Logger, i.Control, i.Config.Security.Cookie, variablesGetter)

	// middleware.
	var md = &middleware.Instance{}
	if err = md.New(corsInstance.Middleware, bodyLimiter.Middleware, httpHelp.middleware); err != nil {
		return err
	}
	i.Middleware = md
	return err
}
