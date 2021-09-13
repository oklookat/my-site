package core

var servus = Servus{}

func Boot() *Servus{
	servus.Utils = Utils{}
	servus.Config = bootConfig()
	servus.Logger = bootLogger()
	servus.DB = bootDB(servus.Config, &servus.Logger)
	servus.Middleware = Middleware{}
	return &servus
}