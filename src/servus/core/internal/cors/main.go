package cors

import (
	"net/http"
	"strconv"
	"strings"
)

type Instance struct {
	config  Config
	writer  http.ResponseWriter
	request *http.Request
}

type Config struct {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
	AllowOrigin []string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods
	AllowMethods []string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers
	AllowHeaders []string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers
	ExposeHeaders []string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials
	AllowCredentials bool
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age
	MaxAge int64
}

// New - create new instance of CorsParse.
func New(config Config) Instance {
	return Instance{config: config}
}

func (i *Instance) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var isPreflight = i.setHeaders(writer, request)
		if isPreflight {
			writer.WriteHeader(204)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

// setHeaders - add CORS headers to response. Returns CorsResult with information about method (is preflight).
func (i *Instance) setHeaders(writer http.ResponseWriter, request *http.Request) (isPreflight bool) {
	i.writer = writer
	i.request = request
	preflight := i.isPreflight()
	if preflight {
		i.allowMethodsParse()
		i.allowHeadersParse()
		i.maxAgeParse()
	}
	i.allowOriginParse()
	i.exposeHeadersParse()
	i.allowCredentialsParse()
	return preflight
}

func (i *Instance) isPreflight() bool {
	var method = i.request.Method
	method = strings.ToUpper(method)
	return method == http.MethodOptions
}

func (i *Instance) allowOriginParse() {
	var allowOrigin = i.config.AllowOrigin
	// bypass.
	if allowOrigin[0] == "*" {
		i.writer.Header().Add("Access-Control-Allow-Origin", "*")
		return
	}
	// https://stackoverflow.com/a/1850482
	var clientOrigin = i.request.Header.Get("Origin")
	for index := range allowOrigin {
		if allowOrigin[index] != clientOrigin {
			continue
		}
		i.writer.Header().Add("Access-Control-Allow-Origin", allowOrigin[index])
		break
	}
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin#cors_and_caching
	i.writer.Header().Add("Vary", "Origin")
}

func (i *Instance) allowMethodsParse() {
	var allowMethods = i.config.AllowMethods
	if allowMethods[0] == "*" {
		i.writer.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, HEAD")
		return
	}
	i.writer.Header().Add("Access-Control-Allow-Methods", strings.Join(allowMethods, ", "))
}

func (i *Instance) allowHeadersParse() {
	var allowHeaders = i.config.AllowHeaders
	var allowed string
	// allow all?
	if allowHeaders[0] == "*" {
		allowed = i.request.Header.Get("Access-Control-Request-Headers")
	} else {
		// make allowed headers string like 'header-1, header-2, header-3'
		allowed = strings.Join(allowHeaders, ", ")
	}
	i.writer.Header().Add("Access-Control-Allow-Headers", allowed)
}

func (i *Instance) exposeHeadersParse() {
	var conf = i.config.ExposeHeaders
	var exposeHeaders = strings.Join(conf, ", ")
	i.writer.Header().Add("Access-Control-Expose-Headers", exposeHeaders)
}

func (i *Instance) allowCredentialsParse() {
	var allowCredentials = i.config.AllowCredentials
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials#directives
	if allowCredentials {
		i.writer.Header().Add("Access-Control-Allow-Credentials", "true")
	}
}

func (i *Instance) maxAgeParse() {
	var maxAge = i.config.MaxAge
	i.writer.Header().Add("Access-Control-Max-Age", strconv.FormatInt(maxAge, 10))
}
