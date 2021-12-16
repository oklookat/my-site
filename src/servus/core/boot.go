package core

import (
	"flag"
	"fmt"
	"os"
	"servus/core/internal/limiter"
	coreutils "servus/core/internal/utils"
)

type Core struct {
	Utils      Utils
	Config     *ConfigFile
	Logger     Logger
	Middleware Middlewarer
	Encryption *Encryption
	DB         *Database
	Control    Controller
}

// Boot - boot Core.
func (c *Core) Boot() {
	c.bootUtils()
	c.bootConfig()
	c.bootLogger()
	c.bootMiddleware()
	c.bootEncryption()
	c.bootDatabase()
	c.bootControl()
}

func (c *Core) bootUtils() {
	c.Utils = &coreutils.Instance{}
}

func (c *Core) bootConfig() {
	var config = ConfigFile{}
	var get = func(path string) {
		err := config.get(path)
		if err != nil {
			panic(err)
		}
	}
	var getFromDir = func() {
		var executionDir, err = c.Utils.GetExecutionDir()
		if err != nil {
			panic(err)
		}
		var path = fmt.Sprintf("%v/settings/config.json", executionDir)
		get(path)
	}
	pathArg := flag.String("config", "nil", "path to config file.")
	flag.Parse()
	// check is path provided in args.
	if pathArg == nil || *pathArg == "nil" {
		getFromDir()
	} else {
		get(*pathArg)
	}
	c.Config = &config
}

func (c *Core) bootLogger() {
	var log = &Log{}
	var executionDir, err = c.Utils.GetExecutionDir()
	if err != nil {
		panic(err)
	}
	log.new(c.Config.Logger, executionDir)
	c.Logger = log
}

func (c *Core) bootControl() {
	var controlTG = ControlTelegram{}
	controlTG.new(c.Config.Control.Telegram, c.Logger)
	var ctrl = &control{}
	ctrl.addController(&controlTG)
	c.Control = ctrl
}

func (c *Core) bootMiddleware() {
	// cors.
	var cors = &Cors{}
	cors.new(c.Config.Security.CORS)
	// limiter.
	_limiter := limiter.New(c.Config.Security.Limiter.Body.Active, c.Config.Security.Limiter.Body.MaxSize, c.Config.Security.Limiter.Body.Except)
	// http
	var http = &httpHelper{}
	http.new(c.Logger, c.Control, c.Config.Security.Cookie)
	// middleware.
	var md = &middleware{}
	md.new(cors.middleware(), _limiter.Middleware, http.middleware)
	c.Middleware = md
}

func (c *Core) bootEncryption() {
	c.Encryption = &Encryption{}
	c.Encryption.new(c.Config.Security.Encryption)
}

func (c *Core) bootDatabase() {
	c.DB = &Database{}
	err := c.DB.new(c.Config, c.Logger)
	if err != nil {
		c.Logger.Panic(err)
		os.Exit(1)
	}
}
