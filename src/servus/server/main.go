package main

import (
	"fmt"
	"log"
	"net/http"
	"servus/apps/elven"
	"servus/core"
)

var servus *core.Servus

func main() {
	log.Println("core: booting")
	servus = core.Boot()
	servus.Logger.Info("elven: booting")
	elven.Boot(servus)
	var host = servus.Config.Host
	var port = servus.Config.Port
	var listen = fmt.Sprintf("%v:%v", host, port)
	servus.Logger.Info(fmt.Sprintf("servus: listening on http://%v", listen))
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		servus.Logger.Panic(err)
	}
}
