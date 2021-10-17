package corsParse

import (
	"net/http"
	"strconv"
	"strings"
)

// New - create new instance of CorsParse.
func New(config Config) CorsParse {
	return CorsParse{config: config}
}

// SetHeaders - add CORS headers to response. Returns CorsResult with information about method (is preflight).
func (c *CorsParse) SetHeaders(writer http.ResponseWriter, request *http.Request) (isPreflight bool) {
	c.writer = writer
	c.request = request
	c.preflightParse()
	c.allowOriginParse()
	c.allowMethodsParse()
	c.allowHeadersParse()
	c.exposeHeadersParse()
	c.allowCredentialsParse()
	c.maxAgeParse()
	// if it is preflight method, we need bypass any auth and return in middleware
	return c.isPreflight
}

func (c *CorsParse) preflightParse(){
	var method = c.request.Method
	method = strings.ToUpper(method)
	c.isPreflight = method == http.MethodOptions
}

func (c *CorsParse) allowOriginParse(){
	var allowOrigin = c.config.AllowOrigin
	var isAllowOriginBypass = true
	switch allowOrigin[0] {
	case "*":
		c.writer.Header().Add("Access-Control-Allow-Origin", "*")
		return
	default:
		isAllowOriginBypass = false
		break
	}
	if !isAllowOriginBypass {
		// https://stackoverflow.com/a/1850482
		var clientOrigin = c.request.Header.Get("Origin")
		for index := range allowOrigin {
			if allowOrigin[index] == clientOrigin {
				c.writer.Header().Add("Access-Control-Allow-Origin", allowOrigin[index])
				break
			}
		}
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin#cors_and_caching
		c.writer.Header().Add("Vary", "Origin")
	}
}

func (c *CorsParse) allowMethodsParse(){
	if !c.isPreflight {
		return
	}
	var allowMethods = c.config.AllowMethods
	if allowMethods[0] == "*" {
		c.writer.Header().Add("Access-Control-Allow-Methods", "OPTIONS, GET, HEAD, POST, PUT, DELETE")
		return
	}
	var clientMethod = c.request.Header.Get("Access-Control-Request-Method")
	for index := range allowMethods {
		if allowMethods[index] == clientMethod {
			c.writer.Header().Add("Access-Control-Allow-Methods", allowMethods[index])
			break
		}
	}
}

func (c *CorsParse) allowHeadersParse(){
	if !c.isPreflight {
		return
	}
	// get allowed headers
	var allowHeaders = c.config.AllowHeaders
	// get request headers string like 'header-1, header-2, header-3'
	var requestHeaders = c.request.Header.Get("Access-Control-Request-Headers")
	// remove all spaces from headers
	requestHeaders = removeSpaces(requestHeaders)
	// split headers string to slice of headers ['header-1', 'header-2', 'header-3']
	var requestHeadersSlice = strings.Split(requestHeaders, ",")
	// get allowed headers from config
	var responseAllowedHeaders string
	// if wildcard - allow all headers
	var allowAll = allowHeaders[0] == "*"
	if allowAll {
		responseAllowedHeaders = strings.Join(requestHeadersSlice, ", ")
	} else {
		// here we store finally allowed headers
		var responseAllowedHeadersSlice []string
		for headerIndex := range allowHeaders {
			// get client headers
			for clientHeaderIndex := range requestHeadersSlice {
				// if allowed header and client header same
				if allowHeaders[headerIndex] == requestHeadersSlice[clientHeaderIndex] {
					// allow this header
					responseAllowedHeadersSlice = append(responseAllowedHeadersSlice, allowHeaders[headerIndex])
					break
				}
			}
		}
		// make allowed headers string like 'header-1, header-2, header-3'
		responseAllowedHeaders = strings.Join(responseAllowedHeadersSlice, ", ")
		return
	}
	c.writer.Header().Add("Access-Control-Allow-Headers", responseAllowedHeaders)
}

func (c *CorsParse) exposeHeadersParse(){
	var exposeHeaders = c.config.ExposeHeaders
	var exposeHeadersFinally = strings.Join(exposeHeaders, ", ")
	c.writer.Header().Add("Access-Control-Expose-Headers", exposeHeadersFinally)
}

func (c *CorsParse) allowCredentialsParse(){
	var allowCredentials = c.config.AllowCredentials
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials#directives
	if allowCredentials {
		c.writer.Header().Add("Access-Control-Allow-Credentials", "true")
	}
}

func (c *CorsParse) maxAgeParse(){
	var maxAge = c.config.MaxAge
	c.writer.Header().Add("Access-Control-Max-Age", strconv.FormatInt(maxAge, 10))
}