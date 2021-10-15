package core

import (
	"fmt"
	"os"
	"servus/core/modules/corsParse"
	"servus/core/modules/logger"
)

// TODO: convert other modules to interfaces, or include it directly here.
// TODO: add argon2 hashing.
type Core struct {
	Utils      *Utils
	Config     *ConfigFile
	Logger     Logger
	Middleware *Middleware
	DB         *Database
}

// Boot - boot Core.
func (c *Core) Boot() {
	c.bootUtils()
	c.bootConfig()
	c.bootLogger()
	c.bootMiddleware()
	c.bootDatabase()
}

// bootUtils - boot Utils.
func (c *Core) bootUtils() {
	c.Utils = &Utils{}
}

// bootConfig - boot config file.
func (c *Core) bootConfig() {
	var config = &ConfigFile{}
	var path = fmt.Sprintf("%v/settings/config.json", c.Utils.GetExecuteDir())
	err := config.boot(path)
	if err != nil {
		panic(err)
	}
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
		var wtfDir = fmt.Sprintf("%v/logs/", c.Utils.GetExecuteDir())
		var wtfMaxLogFiles = c.Config.Logger.WriteToFile.MaxLogFiles
		var wtfMaxLogSize = c.Config.Logger.WriteToFile.MaxLogSize
		loggerConfig.WriteToFile.Dir = wtfDir
		loggerConfig.WriteToFile.MaxLogFiles = wtfMaxLogFiles
		loggerConfig.WriteToFile.MaxLogSize = wtfMaxLogSize
	}
	var theLogger = logger.New(loggerConfig)
	c.Logger = &theLogger
}

// bootMiddleware - boot middleware. Use this after booting the config.
func (c *Core) bootMiddleware() {
	var corsConfig = corsParse.Config{
		AllowCredentials: c.Config.Security.CORS.AllowCredentials,
		AllowOrigin:      c.Config.Security.CORS.AllowOrigin,
		AllowMethods:     c.Config.Security.CORS.AllowMethods,
		AllowHeaders:     c.Config.Security.CORS.AllowHeaders,
		ExposeHeaders:    c.Config.Security.CORS.ExposeHeaders,
		MaxAge:           c.Config.Security.CORS.MaxAge,
	}
	var corsInstance = corsParse.New(corsConfig)
	*c.Middleware = Middleware{config: c.Config, cors: &corsInstance}
}

// bootDatabase - boot database. Use this after booting the config.
func (c *Core) bootDatabase() {
	var database = Database{config: c.Config, logger: c.Logger}
	err := database.boot()
	if err != nil {
		c.Logger.Panic(err)
		os.Exit(1)
	}
}
