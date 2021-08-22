package main

import (
	"fmt"
	"log"
	"net/http"
	"servus/_app"
	"servus/_cmd"
	"servus/_core"
)

func main(){
	log.Println("core: booting")
	var servus = _core.Boot()
	servus.Logger.Info("cmd: booting")
	_cmd.BootCmd(servus)
	servus.Logger.Info("cmd: booted")
	servus.Logger.Info("app: booting")
	_app.Boot(servus)
	servus.Logger.Info("app: booted")
	var host = servus.Config.Host
	var port = servus.Config.Port
	var listen = fmt.Sprintf("%v:%v", host, port)
	servus.Logger.Info(fmt.Sprintf("servus: listen on http://%v", listen))
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		servus.Logger.Panic(err)
	}
}
