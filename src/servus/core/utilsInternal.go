package core

import (
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)


func convertCookieSameSite(sameSite string) (http.SameSite, error) {
	sameSite = strings.ToUpper(sameSite)
	switch sameSite {
	case "DEFAULT":
		return http.SameSiteDefaultMode, nil
	case "LAX":
		return http.SameSiteLaxMode, nil
	case "STRICT":
		return http.SameSiteStrictMode, nil
	case "NONE":
		return http.SameSiteNoneMode, nil
	default:
		return http.SameSiteDefaultMode, errors.New("Wrong sameSite string.")
	}
}

func parseCorsConfig(writer http.ResponseWriter, request *http.Request){
	//////// allow origin parse
	var allowOrigin = servus.Config.Security.CORS.AllowOrigin
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
	var allowMethods = servus.Config.Security.CORS.AllowMethods
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
	var allowHeaders = servus.Config.Security.CORS.AllowHeaders
	// get headers like 'header-1, header-2, header-3'
	var clientHeaders = request.Header.Get("Access-Control-Request-Headers")
	// remove all spaces from headers
	clientHeaders = servus.Utils.RemoveSpaces(clientHeaders)
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
	var exposeHeaders = servus.Config.Security.CORS.ExposeHeaders
	var exposeHeadersFinally = strings.Join(exposeHeaders, ", ")
	writer.Header().Add("Access-Control-Expose-Headers", exposeHeadersFinally)
	//////// allow credentials parse
	var allowCredentials = servus.Config.Security.CORS.AllowCredentials
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials#directives
	if allowCredentials {
		writer.Header().Add("Access-Control-Allow-Credentials", "true")
	}
	//////// max age parse
	var maxAge = servus.Config.Security.CORS.MaxAge
	writer.Header().Add("Access-Control-Max-Age", strconv.FormatInt(maxAge, 10))
}