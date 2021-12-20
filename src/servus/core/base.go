package core

import (
	"io"
	"net/http"
)

type _ctxHTTP string

const ctxHTTP _ctxHTTP = "CORE_HTTP_PIPE"

// HTTP - helper for request/response manipulations.
type HTTP interface {
	Send(body string, statusCode int, err error)
	SetCookie(name string, value string) error
}

// Controller - sends information/controls server via 3rd party services, like Telegram bot.
type Controller interface {
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
	// GetHTTP - get HTTP from request context.
	GetHTTP(request *http.Request) HTTP
	// GenerateULID - returns unique string like 1GFGVSSRTHYWW52GVXZ.
	GenerateULID() (ul string, err error)
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
}

// Encryptor - encrypt/hash values.
type Encryptor struct {
	AES    EncryptorCryptor
	BCrypt EncryptorHasher
	Argon  EncryptorHasher
}

// EncryptorHasher - hash/compare value.
type EncryptorHasher interface {
	// Hash - make hash.
	Hash(data string) (hash string, err error)
	// Compare - compare value with hash.
	Compare(what, with string) (match bool, err error)
	// IsHash - check is value a hash.
	IsHash(data string) bool
}

// EncryptorCryptor - encrypt/decrypt value.
type EncryptorCryptor interface {
	// Encrypt - encrypt value.
	Encrypt(data string) (encrypted string, err error)
	// Decrypt - decrypt value.
	Decrypt(encrypted string) (data string, err error)
}
