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
	var hostAndPort = fmt.Sprintf("%v:%v", host, port)
	// check http or https
	if !servus.Config.Security.HTTPS.Active {
		serveHttp(hostAndPort)
	} else {
		serveHttps(hostAndPort)
	}
}

func serveHttp(hostAndPort string){
	var listeningOn = fmt.Sprintf("servus: listening on http://%v", hostAndPort)
	servus.Logger.Info(listeningOn)
	err := http.ListenAndServe(hostAndPort, nil)
	if err != nil {
		servus.Logger.Panic(err)
	}
}

func serveHttps(hostAndPort string){
	var listeningOn = fmt.Sprintf("servus: listening on https://%v", hostAndPort)
	servus.Logger.Info(listeningOn)
	var certPath = servus.Config.Security.HTTPS.CertPath
	var keyPath = servus.Config.Security.HTTPS.KeyPath
	err := http.ListenAndServeTLS(hostAndPort, certPath, keyPath, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}