package core

import (
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
	httpHelp.new(i.Logger, i.Control, i.Config.Security.Cookie)

	// middleware.
	var md = &middleware.Instance{}
	if err = md.New(corsInstance.Middleware, bodyLimiter.Middleware, httpHelp.middleware); err != nil {
		return err
	}
	i.Middleware = md
	return err
}
