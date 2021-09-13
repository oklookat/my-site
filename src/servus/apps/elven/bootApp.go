package elven

import (
	"servus/apps/elven/elUser"
	"servus/core"
)

var servus *core.Servus

func Boot(_servus *core.Servus){
	servus = _servus
	elUser.Boot(servus)
	bootRoutes()
}
