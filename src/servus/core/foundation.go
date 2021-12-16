package core

import (
	"io"
	"net/http"
)

type MiddlewareFunc func(http.Handler) http.Handler

// HTTP - helper for request/response manipulations.
type HTTP interface {
	Send(body string, statusCode int, err error)
	SetCookie(name string, value string) error
}

// Controller - sends information/controls server via 3rd party services, like Telegram bot.
type Controller interface {
	GetEnabled() bool
	SendMessage(message string)
	SendFile(caption *string, filename string, reader io.Reader)
}

// Logger - writes information.
type Logger interface {
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(message string)
	Panic(err error)
}

// Utils - useful utilities.
type Utils interface {
	// RemoveSpaces - remove spaces from string.
	RemoveSpaces(str string) string
	// GetExecutionDir - get server execution directory.
	GetExecutionDir() (string, error)
	// FormatPath - format path to system specific slashes.
	FormatPath(path string) string
}

// Middlewarer - basic middlewares.
type Middlewarer interface {
	// AsJson - set application/json header.
	AsJson() func(http.Handler) http.Handler
	// LimitBody - set CORS headers depending on config.
	CORS() func(http.Handler) http.Handler
	// LimitBody - limit request body size.
	LimitBody() func(http.Handler) http.Handler
	// ProvideHTTP - get HTTP helper.
	ProvideHTTP() func(http.Handler) http.Handler
	// GetHTTP - get HTTP from request context.
	GetHTTP(request *http.Request) (HTTP, error)
}
