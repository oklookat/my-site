package main

import (
	"fmt"
	"net/http"
	"servus/apps/elven"
	"servus/core"
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
	var listeningOn = fmt.Sprintf("servus: listening at http://%s", hostAndPort)
	call.Logger.Info(listeningOn)
	if err := http.ListenAndServe(hostAndPort, nil); err != nil {
		var prettyErr = fmt.Errorf("HTTP serve error: %w", err)
		call.Logger.Panic(prettyErr)
	}
}

func serveHttps(hostAndPort string) {
	var listeningOn = fmt.Sprintf("servus: listening at https://%s", hostAndPort)
	call.Logger.Info(listeningOn)
	var certPath = call.Config.Security.HTTPS.CertPath
	var keyPath = call.Config.Security.HTTPS.KeyPath
	if err := http.ListenAndServeTLS(hostAndPort, certPath, keyPath, nil); err != nil {
		var prettyErr = fmt.Errorf("HTTPS serve error: %w", err)
		call.Logger.Panic(prettyErr)
	}
}
