package core

import (
	"net/http"
	"servus/core/internal/limiter"
)

type LimiterBody interface {
	// is limiter active?
	IsActive() bool

	// get max body size in MB.
	GetMaxSize() int64

	// get bypassed paths.
	GetExcept() []string

	// get middleware.
	Middleware(next http.Handler) http.Handler
}

type Limiter struct {
	Body LimiterBody
}

func (i *Instance) setupLimiter() error {
	var module = &Limiter{}

	// body.
	var bodyConf = i.Config.Security.Limiter.Body
	var body, err = limiter.NewBody(bodyConf)
	if err != nil {
		return err
	}
	module.Body = body

	i.Limiter = module
	return err
}
