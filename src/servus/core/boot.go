package core

import (
	"flag"
	"fmt"
	"os"
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
	i.bootMiddleware()
	i.bootEncryptor()
	i.bootDatabase()
	i.bootControl()
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
	pathArg := flag.String("config", "nil", "path to config file.")
	flag.Parse()
	// check is path provided in args.
	if pathArg == nil || *pathArg == "nil" {
		getFromDir()
	} else {
		get(*pathArg)
	}
	i.Config = &config
}

func (i *Instance) bootLogger() {
	var log = logger.New(i.Config.Logger.Level)
	i.Logger = log
}

func (i *Instance) bootControl() {
	var ctrl = &control{}
	// telegram.
	if i.Config.Control.Telegram.Enabled {
		var controlTG = controlTelegram{}
		controlTG.new(i.Config.Control.Telegram, i.Logger)
		ctrl.add(&controlTG)
	}
	i.Control = ctrl
}

func (i *Instance) bootMiddleware() {
	// cors.
	var cors = &cors{}
	cors.new(i.Config.Security.CORS)
	// limiter.
	_limiter := limiter.New(i.Config.Security.Limiter.Body.Active, i.Config.Security.Limiter.Body.MaxSize, i.Config.Security.Limiter.Body.Except)
	// http
	var http = &httpHelper{}
	http.new(i.Logger, i.Control, i.Config.Security.Cookie)
	// middleware.
	var md = &middleware.Instance{}
	md.New(cors.middleware(), _limiter.Middleware, http.middleware)
	i.Middleware = md
}

func (i *Instance) bootEncryptor() {
	var encryptor = &encryptor{}
	encryptor.new(i.Config.Security.Encryption)
	i.Encryptor = &Encryptor{}
	i.Encryptor.AES = encryptor.AES
	i.Encryptor.Argon = encryptor.Argon
	i.Encryptor.BCrypt = encryptor.BCrypt
}

func (i *Instance) bootDatabase() {
	i.DB = &Database{}
	err := i.DB.new(i.Config, i.Logger)
	if err != nil {
		i.Logger.Panic(err)
		os.Exit(1)
	}
}
