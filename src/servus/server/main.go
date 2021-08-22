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
	log.Println("booting core...")
	var servus = _core.Boot()
	servus.Logger.Info("booting cmd...")
	_cmd.BootCmd(servus)
	servus.Logger.Info("booting app...")
	_app.Boot(servus)
	var host = servus.Config.Host
	var port = servus.Config.Port
	err := http.ListenAndServe(fmt.Sprintf("%v:%v", host, port), nil)
	if err != nil {
		servus.Logger.Panic(err.Error())
	}
}
