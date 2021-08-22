package _app

import (
	"servus/_app/controllers"
	"servus/_core"
)

var servus *_core.Servus

func Boot(_servus *_core.Servus){
	servus = _servus
	servus.Logger.Info("booting controllers...")
	controllers.BootControllers(servus)
	servus.Logger.Info("booting routes...")
	bootRoutes()
	data := map[string]string{"role": "admin", "username": "юзернейм6", "password": "passwordpassword"}
	servus.DB.User.Create(data)
}
