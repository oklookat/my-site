package _app

import (
	"servus/_app/controllers"
	"servus/_core"
	"servus/_core/database"
)

var servus *_core.Servus

func Boot(_servus *_core.Servus){
	servus = _servus
	servus.Logger.Info("booting controllers...")
	controllers.BootControllers(servus)
	servus.Logger.Info("booting routes...")
	bootRoutes()
	var user = database.User{Username: "username1", Password: "53153f13515", RegIP: "123"}
	servus.DB.User.Create(user)
}
