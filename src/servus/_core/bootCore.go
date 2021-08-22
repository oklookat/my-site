package _core

import (
	"log"
)


func Boot() *Servus{
	log.Println("booting utils...")
	servus.Utils = Utils{}
	log.Println("booting config...")
	servus.Config = bootConfig()
	log.Println("booting logger...")
	servus.Logger = bootLogger()
	servus.Logger.Info("booting database...")
	servus.DB = bootDB(servus.Config, &servus.Logger)
	servus.Logger.Info("core: booted")
	return &servus
}