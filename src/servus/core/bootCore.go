package core

import (
	"os"
	"servus/core/ancientUI"
)

func Boot() *Servus{
	var item = ancientUI.SelectItem{Title: "hello", Items: []string{"Friday", "Saturday", "Sunday"}}
	var selectI = ancientUI.AddSelect(item)
	var selected = selectI.Start()
	println(selected)
	os.Exit(1)


	servus.Utils = Utils{}
	servus.Config = bootConfig()
	servus.Logger = bootLogger()
	servus.DB = bootDB(servus.Config, &servus.Logger)
	return &servus
}