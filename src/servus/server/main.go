package main

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"servus/apps/elven"
	"servus/core"
)

func main() {
	core.Boot()
	core.Logger.Info("elven: booting")
	elven.BootApp()
	var host = core.Config.Host
	var port = core.Config.Port
	var hostAndPort = fmt.Sprintf("%v:%v", host, port)
	if !core.Config.Security.HTTPS.Active {
		serveHttp(hostAndPort)
	} else {
		serveHttps(hostAndPort)
	}
}

func serveHttp(hostAndPort string) {
	var listeningOn = fmt.Sprintf("servus: listening on http://%v", hostAndPort)
	core.Logger.Info(listeningOn)
	err := http.ListenAndServe(hostAndPort, nil)
	if err != nil {
		var prettyErr = errors.Wrap(err, "HTTP serve error: ")
		core.Logger.Panic(prettyErr)
	}
}

func serveHttps(hostAndPort string) {
	var listeningOn = fmt.Sprintf("servus: listening on https://%v", hostAndPort)
	core.Logger.Info(listeningOn)
	var certPath = core.Config.Security.HTTPS.CertPath
	var keyPath = core.Config.Security.HTTPS.KeyPath
	err := http.ListenAndServeTLS(hostAndPort, certPath, keyPath, nil)
	if err != nil {
		var prettyErr = errors.Wrap(err, "HTTPS serve error: ")
		core.Logger.Panic(prettyErr)
	}
}
