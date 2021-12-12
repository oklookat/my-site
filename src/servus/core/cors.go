package core

import (
	"net/http"
	"servus/core/internal/cors"
)

type Cors struct {
	config   *CorsConfig
	instance *cors.Instance
}

type CorsConfig struct {
	Active           bool     `json:"active"`
	AllowCredentials bool     `json:"allowCredentials"`
	AllowOrigin      []string `json:"allowOrigin"`
	AllowMethods     []string `json:"allowMethods"`
	AllowHeaders     []string `json:"allowHeaders"`
	ExposeHeaders    []string `json:"exposeHeaders"`
	MaxAge           int64    `json:"maxAge"`
}

func (c *Cors) new(config *CorsConfig) {
	c.config = config
	var corsConfig = cors.Config{
		AllowCredentials: c.config.AllowCredentials,
		AllowOrigin:      c.config.AllowOrigin,
		AllowMethods:     c.config.AllowMethods,
		AllowHeaders:     c.config.AllowHeaders,
		ExposeHeaders:    c.config.ExposeHeaders,
		MaxAge:           c.config.MaxAge,
	}
	var corsInstance = cors.New(&corsConfig)
	c.instance = corsInstance
}

func (c *Cors) Enabled() bool {
	return c.config.Active
}

func (c *Cors) middleware() func(next http.Handler) http.Handler {
	return c.instance.Middleware
}
