package core

import (
	"fmt"
	"servus/core/external/argument"
	"servus/core/external/database"
	"servus/core/internal/controlTelegram"
	"servus/core/internal/cors"
	"servus/core/internal/cryptor"
	"servus/core/internal/limiter"
	"servus/core/internal/logger"
	"servus/core/internal/middleware"
)

// Instance - servus kernel. Provides cool things.
type Instance struct {
	Utils      Utils
	Config     *ConfigFile
	Logger     Logger
	Middleware Middlewarer
	Encryptor  *Encryptor
	Control    Controller
}

// boot Instance.
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
	// get from path.
	var get = func(path string) {
		err := config.load(path)
		if err != nil {
			panic(err)
		}
	}
	// get from current dir.
	var getFromDir = func() {
		var executionDir, err = i.Utils.GetExecutionDir()
		if err != nil {
			panic(err)
		}
		var path = fmt.Sprintf("%s/settings/config.json", executionDir)
		get(path)
	}
	// is config path in args?
	var configFlag = "-config"
	var arg = argument.Get(configFlag)
	var notInArgs = arg == nil || arg.Value == nil
	if notInArgs {
		// get from current dir.
		getFromDir()
	} else {
		// get by path in args.
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
	// Telegram bot.
	var tgEnabled = i.Config.Control.Telegram.Enabled
	if tgEnabled {
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
	var bodyLimiter = limiter.NewBody(i.Config.Security.Limiter.Body)
	// http.
	var http = &httpHelper{}
	http.new(i.Logger, i.Control, i.Config.Security.Cookie)
	// middleware.
	var md = &middleware.Instance{}
	md.New(corsInstance.Middleware, bodyLimiter.Middleware, http.middleware)
	i.Middleware = md
}

func (i *Instance) bootEncryptor() {
	var cr = &cryptor.Instance{}
	cr.New(i.Config.Security.Encryption)
	//
	var en = &Encryptor{}
	en.AES = cr.AES
	en.Argon = cr.Argon
	en.BCrypt = cr.BCrypt
	i.Encryptor = en
}

func (i *Instance) bootDatabase() {
	// connect to DB. Database be available via database.Adapter.
	var conn = database.Connector{}
	conn.New(i.Config.DB, i.Logger)
}
