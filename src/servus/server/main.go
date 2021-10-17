package main

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"servus/apps/elven"
	"servus/core"
)

var instance = core.Core{}

func main() {
	instance.Boot()
	// boot elven.
	instance.Logger.Info("elven: booting")
	elven.BootApp(&instance)
	// serve.
	var host = instance.Config.Host
	var port = instance.Config.Port
	var hostAndPort = fmt.Sprintf("%v:%v", host, port)
	if !instance.Config.Security.HTTPS.Active {
		serveHttp(hostAndPort)
	} else {
		serveHttps(hostAndPort)
	}
}

func serveHttp(hostAndPort string) {
	var listeningOn = fmt.Sprintf("servus: listening on http://%v", hostAndPort)
	instance.Logger.Info(listeningOn)
	err := http.ListenAndServe(hostAndPort, nil)
	if err != nil {
		var prettyErr = errors.Wrap(err, "HTTP serve error: ")
		instance.Logger.Panic(prettyErr)
	}
}

func serveHttps(hostAndPort string) {
	var listeningOn = fmt.Sprintf("servus: listening on https://%v", hostAndPort)
	instance.Logger.Info(listeningOn)
	var certPath = instance.Config.Security.HTTPS.CertPath
	var keyPath = instance.Config.Security.HTTPS.KeyPath
	err := http.ListenAndServeTLS(hostAndPort, certPath, keyPath, nil)
	if err != nil {
		var prettyErr = errors.Wrap(err, "HTTPS serve error: ")
		instance.Logger.Panic(prettyErr)
	}
}
