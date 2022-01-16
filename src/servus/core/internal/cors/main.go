package cors

import (
	"net/http"
	"strconv"
	"strings"
)

type Instance struct {
	config *Config
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

func New(config *Config) *Instance {
	return &Instance{config: config}
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

// add CORS headers to response. Returns CorsResult with information about method (is preflight).
func (i *Instance) setHeaders(writer http.ResponseWriter, request *http.Request) (isPreflight bool) {
	preflight := i.isPreflight(request.Method)
	if preflight {
		i.allowMethodsParse(writer)
		i.allowHeadersParse(writer, request)
		i.maxAgeParse(writer)
	}
	i.allowOriginParse(writer, request)
	i.exposeHeadersParse(writer)
	i.allowCredentialsParse(writer)
	return preflight
}

func (i *Instance) isPreflight(method string) bool {
	return method == http.MethodOptions
}

func (i *Instance) allowOriginParse(writer http.ResponseWriter, request *http.Request) {
	var allowOrigin = i.config.AllowOrigin
	// bypass.
	if allowOrigin[0] == "*" {
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		return
	}
	// https://stackoverflow.com/a/1850482
	var clientOrigin = request.Header.Get("Origin")
	for index := range allowOrigin {
		if allowOrigin[index] != clientOrigin {
			continue
		}
		writer.Header().Add("Access-Control-Allow-Origin", allowOrigin[index])
		break
	}
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin#cors_and_caching
	writer.Header().Add("Vary", "Origin")
}

func (i *Instance) allowMethodsParse(writer http.ResponseWriter) {
	var allowMethods = i.config.AllowMethods
	if allowMethods[0] == "*" {
		writer.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, HEAD")
		return
	}
	writer.Header().Add("Access-Control-Allow-Methods", strings.Join(allowMethods, ", "))
}

func (i *Instance) allowHeadersParse(writer http.ResponseWriter, request *http.Request) {
	var allowHeaders = i.config.AllowHeaders
	var allowed string
	// allow all?
	if allowHeaders[0] == "*" {
		allowed = request.Header.Get("Access-Control-Request-Headers")
	} else {
		// make allowed headers string like 'header-1, header-2, header-3'
		allowed = strings.Join(allowHeaders, ", ")
	}
	writer.Header().Add("Access-Control-Allow-Headers", allowed)
}

func (i *Instance) exposeHeadersParse(writer http.ResponseWriter) {
	var conf = i.config.ExposeHeaders
	var exposeHeaders = strings.Join(conf, ", ")
	writer.Header().Add("Access-Control-Expose-Headers", exposeHeaders)
}

func (i *Instance) allowCredentialsParse(writer http.ResponseWriter) {
	var allowCredentials = i.config.AllowCredentials
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials#directives
	if allowCredentials {
		writer.Header().Add("Access-Control-Allow-Credentials", "true")
	}
}

func (i *Instance) maxAgeParse(writer http.ResponseWriter) {
	var maxAge = i.config.MaxAge
	writer.Header().Add("Access-Control-Max-Age", strconv.FormatInt(maxAge, 10))
}
