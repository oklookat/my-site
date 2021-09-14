package corsParse

import (
	"net/http"
	"strconv"
	"strings"
)

func New(config Config) CorsParse{
	return CorsParse{config: config}
}

func (c *CorsParse) SetHeaders(writer http.ResponseWriter, request *http.Request){
	//////// allow origin parse
	var allowOrigin = c.config.AllowOrigin
	var isAllowOriginBypass = true
	switch allowOrigin[0] {
	case "*":
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		break
	case "null":
		writer.Header().Add("Access-Control-Allow-Origin", "null")
		break
	default:
		isAllowOriginBypass = false
		break
	}
	if !isAllowOriginBypass {
		// https://stackoverflow.com/a/1850482
		var clientOrigin = request.Header.Get("Origin")
		for _, origin := range allowOrigin {
			if origin == clientOrigin {
				writer.Header().Add("Access-Control-Allow-Origin", origin)
				break
			}
		}
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin#cors_and_caching
		writer.Header().Add("Vary", "Origin")
	}
	//////// allow methods parse
	var allowMethods = c.config.AllowMethods
	var isAllowMethodsBypass = true
	if allowMethods[0] == "*"{
		writer.Header().Add("Access-Control-Allow-Methods", "OPTIONS, GET, HEAD, POST, PUT, DELETE")
	} else {
		isAllowMethodsBypass = false
	}
	if !isAllowMethodsBypass {
		var clientMethod = request.Header.Get("Access-Control-Request-Method")
		for _, method := range allowMethods {
			if method == clientMethod {
				writer.Header().Add("Access-Control-Allow-Methods", method)
				break
			}
		}
	}
	//////// allow headers parse
	var allowHeaders = c.config.AllowHeaders
	// get headers like 'header-1, header-2, header-3'
	var clientHeaders = request.Header.Get("Access-Control-Request-Headers")
	// remove all spaces from headers
	clientHeaders = removeSpaces(clientHeaders)
	// split headers string to slice of headers
	var clientHeadersSlice = strings.Split(clientHeaders, ",")
	// here we store finally allowed headers
	var responseAllowedHeadersSlice []string
	// get allowed headers from config
	var responseAllowedHeaders string
	// if wildcard - allow all headers
	if allowHeaders[0] == "*"{
		responseAllowedHeaders = strings.Join(clientHeadersSlice, ", ")
	} else {
		for _, header := range allowHeaders {
			// get client headers
			for _, clientHeader := range clientHeadersSlice {
				// if allowed header and client header same
				if header == clientHeader {
					// allow this header
					responseAllowedHeadersSlice = append(responseAllowedHeadersSlice, header)
					break
				}
			}
		}
		responseAllowedHeaders = strings.Join(responseAllowedHeadersSlice, ", ")
	}
	// make allowed headers string like 'header-1, header-2, header-3'
	writer.Header().Add("Access-Control-Allow-Headers", responseAllowedHeaders)
	//////// expose headers parse
	var exposeHeaders = c.config.ExposeHeaders
	var exposeHeadersFinally = strings.Join(exposeHeaders, ", ")
	writer.Header().Add("Access-Control-Expose-Headers", exposeHeadersFinally)
	//////// allow credentials parse
	var allowCredentials = c.config.AllowCredentials
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials#directives
	if allowCredentials {
		writer.Header().Add("Access-Control-Allow-Credentials", "true")
	}
	//////// max age parse
	var maxAge = c.config.MaxAge
	writer.Header().Add("Access-Control-Max-Age", strconv.FormatInt(maxAge, 10))

	//////// bypass any auth if simple request
	var method = request.Method
	method = strings.ToUpper(method)
	var isSimpleRequest = method == "OPTIONS"
	if isSimpleRequest {
		writer.WriteHeader(200)
		writer.Write([]byte(""))
		return
	}
}
