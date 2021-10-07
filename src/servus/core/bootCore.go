package core

import (
	"servus/core/modules/corsParse"
	"servus/core/modules/database"
	"servus/core/modules/logger"
)

var Middleware = BasicMiddleware{}
var Utils = BasicUtils{}
var Config *ConfigFile
var Logger *logger.Logger
var Database *database.DB
var corsParser *corsParse.CorsParse

func Boot() {
	Config = bootConfig()
	Logger = bootLogger()
	Database = bootDB(Config, Logger)
	var corsConfig = corsParse.Config{
		AllowCredentials: Config.Security.CORS.AllowCredentials,
		AllowOrigin:      Config.Security.CORS.AllowOrigin,
		AllowMethods:     Config.Security.CORS.AllowMethods,
		AllowHeaders:     Config.Security.CORS.AllowHeaders,
		ExposeHeaders:    Config.Security.CORS.ExposeHeaders,
		MaxAge:           Config.Security.CORS.MaxAge,
	}
	var corsParseInstance = corsParse.New(corsConfig)
	corsParser = &corsParseInstance
}
