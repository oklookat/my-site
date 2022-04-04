package core

import (
	"io"
	"net"
	"net/http"
)

// Instance - servus kernel. Provides cool things.
type Instance struct {
	Dirs      Directories
	Config    *Config
	Logger    Logger
	Banhammer Banhammer
	Encryptor *Encryptor
	Control   Controller
	Cors      Corser
	Limiter   *Limiter
	Http      httpHelper
}

// provides functions to get HTTP.
type HttpHelper interface {
	// provide HTTP funcs.
	Middleware(next http.Handler) http.Handler

	// get HTTP helper from request context.
	Get(request *http.Request) HTTP
}

// request/response manipulations helper. Available after using middleware.
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

// sends information/controls server via 3rd party services like Telegram bot.
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

// encrypt/hash values.
type Encryptor struct {
	AES    EncryptorCryptor
	BCrypt EncryptorHasher
	Argon  EncryptorHasher
}

// hash/compare value.
type EncryptorHasher interface {
	// make hash.
	Hash(password string) (hash string, err error)

	// compare password with hash.
	Compare(password, hash string) (match bool, err error)

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

	// when IP unbanned.
	OnUnbanned(hook func(ip string))

	// add warn. 3 warns = ban.
	Warn(ip string) error

	// remove warn.
	Unwarn(ip string) error

	// when IP warned.
	OnWarned(hook func(ip string))

	// get ban checking middleware.
	Middleware(next http.Handler) http.Handler

	// get IP by request by X-REAL-IP / X-FORWARDED-FOR
	//
	// returns nil if failed to get IP.
	GetIpByRequest(request *http.Request) net.IP
}

// CORS funcs.
type Corser interface {
	GetMiddleware(next http.Handler) http.Handler
}
