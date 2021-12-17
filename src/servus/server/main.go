package main

import (
	"fmt"
	"net/http"
	"servus/apps/elven"
	"servus/core"

	"github.com/pkg/errors"
)

var call = core.Instance{}

func main() {
	// boot.
	call.Boot()
	var _elven = elven.App{}
	_elven.Boot(&call)
	// serve.
	var host = call.Config.Host
	var port = call.Config.Port
	var hostAndPort = fmt.Sprintf("%v:%v", host, port)
	if !call.Config.Security.HTTPS.Active {
		serveHttp(hostAndPort)
	} else {
		serveHttps(hostAndPort)
	}
}

func serveHttp(hostAndPort string) {
	var listeningOn = fmt.Sprintf("servus: listening on http://%v", hostAndPort)
	call.Logger.Info(listeningOn)
	err := http.ListenAndServe(hostAndPort, nil)
	if err != nil {
		var prettyErr = errors.Wrap(err, "HTTP serve error: ")
		call.Logger.Panic(prettyErr)
	}
}

func serveHttps(hostAndPort string) {
	var listeningOn = fmt.Sprintf("servus: listening on https://%v", hostAndPort)
	call.Logger.Info(listeningOn)
	var certPath = call.Config.Security.HTTPS.CertPath
	var keyPath = call.Config.Security.HTTPS.KeyPath
	err := http.ListenAndServeTLS(hostAndPort, certPath, keyPath, nil)
	if err != nil {
		var prettyErr = errors.Wrap(err, "HTTPS serve error: ")
		call.Logger.Panic(prettyErr)
	}
}
