package core

import (
	"fmt"
	"os"
	"servus/core/internal/limiter"
)

type Core struct {
	Utils  *Utils
	Config *ConfigFile
	Logger Logger
	Middleware *Middleware
	Encryption *Encryption
	DB         *Database
	Control    Controller
}

// Boot - boot Core.
func (c *Core) Boot() {

	// utils.
	c.Utils = &Utils{}

	// config.
	var config = ConfigFile{}
	var path = fmt.Sprintf("%v/settings/config.json", c.Utils.GetExecutionDir())
	err := config.new(path)
	if err != nil {
		panic(err)
	}
	c.Config = &config

	// logger.
	var log = &Log{}
	log.new(c.Config.Logger, c.Utils.GetExecutionDir())
	c.Logger = log

	// control.
	var controlTG = ControlTelegram{}
	controlTG.new(c.Config.Control.Telegram, c.Logger)
	var ctrl = &control{}
	ctrl.addController(&controlTG)
	c.Control = ctrl

	// cors.
	var cors = &Cors{}
	cors.new(c.Config.Security.CORS)

	// http.
	var http = &httpHelper{}
	http.new(c.Logger, c.Control, c.Config.Security.Cookie)

	// middleware.
	_limiter := limiter.New(c.Config.Security.Limiter.Body.Active, c.Config.Security.Limiter.Body.MaxSize, c.Config.Security.Limiter.Body.Except)
	c.Middleware = &Middleware{}
	c.Middleware.new(cors.middleware(), _limiter.Middleware, http.middleware)

	// encryption.
	c.Encryption = &Encryption{}
	c.Encryption.new(c.Config.Security.Encryption)

	// database.
	c.DB = &Database{}
	err = c.DB.new(c.Config, c.Logger)
	if err != nil {
		c.Logger.Panic(err)
		os.Exit(1)
	}

}
