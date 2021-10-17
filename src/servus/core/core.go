package core

import (
	"fmt"
	"os"
	"servus/core/internal/corsParse"
	"servus/core/internal/logger"
)

type Core struct {
	Utils      *Utils
	Config     *ConfigFile
	Logger     Logger
	HTTP       *HTTP
	Encryption *Encryption
	DB         *Database
}

// Boot - boot Core.
func (c *Core) Boot() {
	c.bootUtils()
	c.bootConfig()
	c.bootLogger()
	c.bootHTTP()
	c.bootEncryption()
	c.bootDatabase()
}

// bootUtils - boot Utils.
func (c *Core) bootUtils() {
	c.Utils = &Utils{}
}

// bootConfig - boot config file.
func (c *Core) bootConfig() {
	var config = ConfigFile{}
	var path = fmt.Sprintf("%v/settings/config.json", c.Utils.GetExecutionDir())
	err := config.boot(path)
	if err != nil {
		panic(err)
	}
	c.Config = &config
}

// bootLogger - boot Logger. Use this after booting the config.
func (c *Core) bootLogger() {
	// create log file
	var writeToConsole = c.Config.Logger.WriteToConsole
	loggerConfig := logger.Config{
		LogLevel:       logger.LevelDebug,
		WriteToConsole: writeToConsole,
	}
	// write to file
	var wtfActive = c.Config.Logger.WriteToFile.Active
	loggerConfig.WriteToFile.Activated = wtfActive
	if wtfActive {
		var wtfDir = fmt.Sprintf("%v/logs/", c.Utils.GetExecutionDir())
		var wtfMaxLogFiles = c.Config.Logger.WriteToFile.MaxLogFiles
		var wtfMaxLogSize = c.Config.Logger.WriteToFile.MaxLogSize
		loggerConfig.WriteToFile.Dir = wtfDir
		loggerConfig.WriteToFile.MaxLogFiles = wtfMaxLogFiles
		loggerConfig.WriteToFile.MaxLogSize = wtfMaxLogSize
	}
	var theLogger = logger.New(loggerConfig)
	c.Logger = &theLogger
}

// bootHTTP - boot http utilities. Use this after booting the config.
func (c *Core) bootHTTP() {
	var corsConfig = corsParse.Config{
		AllowCredentials: c.Config.Security.CORS.AllowCredentials,
		AllowOrigin:      c.Config.Security.CORS.AllowOrigin,
		AllowMethods:     c.Config.Security.CORS.AllowMethods,
		AllowHeaders:     c.Config.Security.CORS.AllowHeaders,
		ExposeHeaders:    c.Config.Security.CORS.ExposeHeaders,
		MaxAge:           c.Config.Security.CORS.MaxAge,
	}
	var corsInstance = corsParse.New(corsConfig)
	middleware := Middleware{config: c.Config, cors: &corsInstance}
	c.HTTP = &HTTP{config: c.Config, Middleware: &middleware}
}

// bootEncryption - boot encryption. Use this after booting the config.
func (c *Core) bootEncryption() {
	enc := &Encryption{config: c.Config}
	enc.boot()
	c.Encryption = enc
}

// bootDatabase - boot database. Use this after booting the config.
func (c *Core) bootDatabase() {
	var database = Database{config: c.Config, logger: c.Logger}
	err := database.boot()
	if err != nil {
		c.Logger.Panic(err)
		os.Exit(1)
	}
	c.DB = &database
}
