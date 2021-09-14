package main

import (
	"fmt"
	"log"
	"net/http"
	"servus/apps/elven"
	"servus/core"
)


func main() {
	core.Logger.Info("elven: booting")
	elven.Boot()
	var host = core.Config.Host
	var port = core.Config.Port
	var hostAndPort = fmt.Sprintf("%v:%v", host, port)
	// check http or https
	if !core.Config.Security.HTTPS.Active {
		serveHttp(hostAndPort)
	} else {
		serveHttps(hostAndPort)
	}
}

func serveHttp(hostAndPort string){
	var listeningOn = fmt.Sprintf("servus: listening on http://%v", hostAndPort)
	core.Logger.Info(listeningOn)
	err := http.ListenAndServe(hostAndPort, nil)
	if err != nil {
		core.Logger.Panic(err)
	}
}

func serveHttps(hostAndPort string){
	var listeningOn = fmt.Sprintf("servus: listening on https://%v", hostAndPort)
	core.Logger.Info(listeningOn)
	var certPath = core.Config.Security.HTTPS.CertPath
	var keyPath = core.Config.Security.HTTPS.KeyPath
	err := http.ListenAndServeTLS(hostAndPort, certPath, keyPath, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}