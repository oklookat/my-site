package core

import (
	"net/http"
	icors "servus/core/internal/cors"
)

type cors struct {
	config   *CorsConfig
	instance *icors.Instance
}

func (c *cors) new(config *CorsConfig) {
	if config == nil {
		panic("[core/cors]: config nil pointer")
	}
	c.config = config
	var corsConfig = icors.Config{
		AllowCredentials: c.config.AllowCredentials,
		AllowOrigin:      c.config.AllowOrigin,
		AllowMethods:     c.config.AllowMethods,
		AllowHeaders:     c.config.AllowHeaders,
		ExposeHeaders:    c.config.ExposeHeaders,
		MaxAge:           c.config.MaxAge,
	}
	var corsInstance = icors.New(&corsConfig)
	c.instance = corsInstance
}

func (c *cors) Enabled() bool {
	return c.config.Active
}

func (c *cors) middleware() func(next http.Handler) http.Handler {
	return c.instance.Middleware
}
