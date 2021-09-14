package elven

import (
	"servus/apps/elven/elUser"
)

func Boot(){
	elUser.Boot()
	bootRoutes()
}
