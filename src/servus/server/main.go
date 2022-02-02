package main

import (
	"fmt"
	"net/http"
	"servus/apps/elven"
	"servus/core"

	"github.com/pkg/errors"
)

var call = &core.Instance{}

func main() {
	// boot.
	call.Boot()
	// elven.
	var _elven = elven.App{}
	_elven.Boot(call)
	// serve.
	serve()
}

func serve() {
	var host = call.Config.Host
	var port = call.Config.Port
	var hostAndPort = fmt.Sprintf("%s:%s", host, port)
	// http(s)
	var isHTTPS = call.Config.Security.HTTPS.Active
	if isHTTPS {
		serveHttps(hostAndPort)
	} else {
		serveHttp(hostAndPort)
	}
}

func serveHttp(hostAndPort string) {
	var listeningOn = fmt.Sprintf("servus: listening on http://%s", hostAndPort)
	call.Logger.Info(listeningOn)
	err := http.ListenAndServe(hostAndPort, nil)
	if err != nil {
		var prettyErr = errors.Wrap(err, "HTTP serve error: ")
		call.Logger.Panic(prettyErr)
	}
}

func serveHttps(hostAndPort string) {
	var listeningOn = fmt.Sprintf("servus: listening on https://%s", hostAndPort)
	call.Logger.Info(listeningOn)
	var certPath = call.Config.Security.HTTPS.CertPath
	var keyPath = call.Config.Security.HTTPS.KeyPath
	err := http.ListenAndServeTLS(hostAndPort, certPath, keyPath, nil)
	if err != nil {
		var prettyErr = errors.Wrap(err, "HTTPS serve error: ")
		call.Logger.Panic(prettyErr)
	}
}
