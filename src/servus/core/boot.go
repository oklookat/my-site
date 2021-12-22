package core

import (
	"fmt"
	"os"
	"servus/core/external/argument"
	"servus/core/internal/controlTelegram"
	"servus/core/internal/cors"
	"servus/core/internal/cryptor"
	"servus/core/internal/limiter"
	"servus/core/internal/logger"
	"servus/core/internal/middleware"
)

// TODO: test body limiter, continue refactoring / fixes.

// Instance - servus kernel. Provides cool things.
type Instance struct {
	Utils      Utils
	Config     *ConfigFile
	Logger     Logger
	Middleware Middlewarer
	Encryptor  *Encryptor
	DB         *Database
	Control    Controller
}

// Boot - boot Instance.
func (i *Instance) Boot() {
	i.bootUtils()
	i.bootConfig()
	i.bootLogger()
	i.bootControl()
	i.bootMiddleware()
	i.bootEncryptor()
	i.bootDatabase()
}

func (i *Instance) bootUtils() {
	i.Utils = &utils{}
}

func (i *Instance) bootConfig() {
	var config = ConfigFile{}
	var get = func(path string) {
		err := config.load(path)
		if err != nil {
			panic(err)
		}
	}
	var getFromDir = func() {
		var executionDir, err = i.Utils.GetExecutionDir()
		if err != nil {
			panic(err)
		}
		var path = fmt.Sprintf("%v/settings/config.json", executionDir)
		get(path)
	}
	// check is path provided in args.
	var configFlag = "-config"
	var arg = argument.Get(configFlag)
	if arg == nil {
		getFromDir()
	} else {
		if arg.Value == nil {
			panic("config flag cannot be empty")
		}
		get(*arg.Value)
	}
	i.Config = &config
}

func (i *Instance) bootLogger() {
	var log = logger.New(i.Config.Logger)
	i.Logger = log
}

func (i *Instance) bootControl() {
	var ctrl = &controller{}
	// telegram.
	if i.Config.Control.Telegram.Enabled {
		var controlTG = controlTelegram.Controller{}
		controlTG.New(i.Config.Control.Telegram, i.Logger)
		ctrl.add(&controlTG)
	}
	i.Control = ctrl
}

func (i *Instance) bootMiddleware() {
	// cors.
	var corsInstance = cors.New(i.Config.Security.CORS)
	// limiter.
	_limiter := limiter.New(i.Config.Security.Limiter.Body.Active, i.Config.Security.Limiter.Body.MaxSize, i.Config.Security.Limiter.Body.Except)
	// http
	var http = &httpHelper{}
	http.new(i.Logger, i.Control, i.Config.Security.Cookie)
	// middleware.
	var md = &middleware.Instance{}
	md.New(corsInstance.Middleware, _limiter.Middleware, http.middleware)
	i.Middleware = md
}

func (i *Instance) bootEncryptor() {
	var cr = &cryptor.Instance{}
	cr.New(i.Config.Security.Encryption)
	var en = &Encryptor{}
	en.AES = cr.AES
	en.Argon = cr.Argon
	en.BCrypt = cr.BCrypt
	i.Encryptor = en
}

func (i *Instance) bootDatabase() {
	i.DB = &Database{}
	err := i.DB.new(i.Config, i.Logger)
	if err != nil {
		i.Logger.Panic(err)
		os.Exit(1)
	}
}
