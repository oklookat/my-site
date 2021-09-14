package core

import "servus/core/modules/corsParse"

// global user-usable vars

var Middleware = BasicMiddleware{}
var Utils = BasicUtils{}
var Config = bootConfig()
var Logger = bootLogger()
var Database = bootDB(Config, &Logger)

// internal core-usable vars

var corsConfig = corsParse.Config{
	AllowCredentials: Config.Security.CORS.AllowCredentials,
	AllowOrigin:   Config.Security.CORS.AllowOrigin,
	AllowMethods:  Config.Security.CORS.AllowMethods,
	AllowHeaders:  Config.Security.CORS.AllowHeaders,
	ExposeHeaders: Config.Security.CORS.ExposeHeaders,
	MaxAge:        Config.Security.CORS.MaxAge,
}
var corsParser = corsParse.New(corsConfig)
