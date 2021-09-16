package corsParse

import "net/http"

type CorsResult struct {
	IsPreflight bool
}

type CorsParse struct {
	corsParseI
	// internal
	config Config
	isPreflight bool
	writer http.ResponseWriter
	request *http.Request
}
type corsParseI interface {
	SetHeaders(writer http.ResponseWriter, request *http.Request)
	preflightParse()
	allowOriginParse()
	allowMethodsParse()
	allowHeadersParse()
	exposeHeadersParse()
	allowCredentialsParse()
	maxAgeParse()
}


type Config struct {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
	AllowOrigin      []string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods
	AllowMethods     []string
 	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers
	AllowHeaders     []string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers
	ExposeHeaders    []string
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials
	AllowCredentials bool
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age
	MaxAge           int64
}
