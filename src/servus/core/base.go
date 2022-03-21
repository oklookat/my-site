package core

import (
	"io"
	"net"
	"net/http"
	"time"
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
	AddCommand(command string, callback func(args []string))
}

// directories/files management.
type Directories interface {
	// get dir where bin executes.
	GetExecution() (path string, err error)
	// get dir where config and other files placed.
	GetData() (path string, err error)
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
	// create debouncer with interval
	//
	// when you call debouncer, callback can only be called once after *interval*.
	//
	// example: interval 1 second
	//
	// if debouncer will be called every 10 milliseconds, then the callback will not be called
	//
	// only after 1 second callback will be called once.
	Debounce(interval time.Duration) (debouncer func(callback func()))
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

// ban/unban/warn IP's.
type Banhammer interface {
	//// ban.
	Ban(ip string) error
	// unban.
	Unban(ip string) error
	// is IP banned?
	IsBanned(ip string) (bool, error)
	// when IP banned.
	OnBanned(hook func(ip string))

	//// warn.
	// add warn. 3 warns = ban.
	Warn(ip string) error
	// remove warn.
	Unwarn(ip string) error
	// when IP warned.
	OnWarned(hook func(ip string))

	//// service.
	// get ban checking middleware.
	GetMiddleware() func(http.Handler) http.Handler
	// get IP by request by X-REAL-IP / X-FORWARDED-FOR
	//
	// returns nil if failed to get IP.
	GetIpByRequest(request *http.Request) net.IP
}
