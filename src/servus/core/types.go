package core

import "net/http"

// Logger - logger.
type Logger interface {
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(message string)
	Panic(err error)
}

type cors interface {
	SetHeaders(response http.ResponseWriter, request *http.Request) (isPreflight bool)
}
