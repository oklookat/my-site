package core

import (
	"io"
	"net/http"
)

type _ctxHTTP string

const ctxHTTP _ctxHTTP = "CORE_HTTP_PIPE"

// helper for request/response manipulations.
type HTTP interface {
	// send response.
	Send(body string, statusCode int, err error)
	// set cookie.
	SetCookie(name string, value string) error
	// get route args [name: value].
	//
	// ----example:----
	//
	// route: /api/users/{username}/{id}
	//
	// request: /api/users/oklookat/1
	//
	// map: = [username: oklookat, id: 1]
	GetRouteArgs() map[string]string
}

// sends information/controls server via 3rd party services, like Telegram bot.
type Controller interface {
	SendMessage(message string)
	SendFile(caption *string, filename string, data io.Reader)
}

// writes information.
type Logger interface {
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(message string)
	Panic(err error)
}

// useful utilities.
type Utils interface {
	// remove spaces from string.
	RemoveSpaces(str string) string
	// get dir where binary started.
	GetExecutionDir() (string, error)
	// format path to system specific slashes.
	FormatPath(path string) string
	// get HTTP from request context.
	GetHTTP(request *http.Request) HTTP
	// returns unique string like 1GFGVSSRTHYWW52GVXZ.
	GenerateULID() (ul string, err error)
	// convert string to rune slice -> get len() of this slice
	//
	// example:
	//
	// with LenRune() "hello" and "вечер" will have the same length in 5
	//
	// with len(), result was be 5 and 10.
	LenRune(val string) int
}

// basic middlewares.
type Middlewarer interface {
	// set application/json header.
	AsJson() func(http.Handler) http.Handler
	// set CORS headers depending on config.
	CORS() func(http.Handler) http.Handler
	// limit request body size.
	LimitBody() func(http.Handler) http.Handler
	// get HTTP helper.
	ProvideHTTP() func(http.Handler) http.Handler
}

// encrypt/hash values.
type Encryptor struct {
	AES    EncryptorCryptor
	BCrypt EncryptorHasher
	Argon  EncryptorHasher
}

// hash/compare value.
type EncryptorHasher interface {
	// make hash.
	Hash(data string) (hash string, err error)
	// compare value with hash.
	Compare(what, with string) (match bool, err error)
	// check is value a hash.
	IsHash(data string) bool
}

// encrypt/decrypt value.
type EncryptorCryptor interface {
	// encrypt data.
	Encrypt(data string) (encrypted string, err error)
	// decrypt data.
	Decrypt(encrypted string) (data string, err error)
}
