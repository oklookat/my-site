package core

import (
	"io"
	"net/http"
)

type MiddlewareFunc func(http.Handler) http.Handler

type HTTP interface {
	Send(body string, statusCode int, err error)
	SetCookie(name string, value string) error
}

// Controller - sends information/controls server via 3rd party services, like Telegram bot.
//
// Use controller interface to provide service.
type Controller interface {
	GetEnabled() bool
	SendMessage(message string)
	SendFile(caption *string, filename string, reader io.Reader)
}