package core

func Boot() *Servus{
	servus.Utils = Utils{}
	servus.Config = bootConfig()
	servus.Logger = bootLogger()
	servus.DB = bootDB(servus.Config, &servus.Logger)
	return &servus
}