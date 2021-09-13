package elUser

import "servus/core"

var servus *core.Servus

func Boot(_servus *core.Servus){
	servus = _servus
	cmdSuperuser()
}
