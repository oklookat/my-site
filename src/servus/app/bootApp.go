package app

import (
	"servus/app/controllers"
	"servus/core"
	"servus/core/database"
)

var servus *core.Servus

func Boot(_servus *core.Servus){
	servus = _servus
	servus.Logger.Info("booting controllers...")
	controllers.BootControllers(servus)
	servus.Logger.Info("booting routes...")
	bootRoutes()
	var user = database.User{Username: "username1", Password: "53153f13515", RegIP: "123"}
	servus.DB.User.Create(user)
}
