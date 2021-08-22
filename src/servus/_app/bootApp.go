package _app

import (
	"servus/_core"
)

var servus *_core.Servus

func Boot(_servus *_core.Servus){
	servus = _servus
	servus.Logger.Info("booting routes...")
	bootRoutes()
}
