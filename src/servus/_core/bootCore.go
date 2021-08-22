package _core

import (
	"log"
)


// тут будить CMD, если есть args. И результат передавать через context
func Boot() *Servus{
	log.Printf("booting utils...")
	servus.Utils = Utils{}
	log.Printf("booting config...")
	servus.Config = bootConfig()
	log.Printf("booting logger...")
	servus.Logger = bootLogger()
	servus.Logger.Info("booting database...")
	servus.DB = bootDB(servus.Config, &servus.Logger)
	servus.Logger.Info("core: booted")
	return &servus
}