package core

import (
	"net/http"
	"strings"
)

// Middleware - basic middleware.
type Middleware struct {
	config *ConfigFile
	cors   cors
}

// AsJSON - set application/json header.
func (m *Middleware) AsJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(writer, request)
	})
}

// CORS - CORS headers depending on config file.
func (m *Middleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if m.config.Security.CORS.Active {
			var isPreflight = m.cors.SetHeaders(writer, request)
			if isPreflight {
				return
			}
		}
		next.ServeHTTP(writer, request)
	})
}

// Security - basic security.
func (m *Middleware) Security(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var method = request.Method
		method = strings.ToUpper(method)
		// request body size check
		if m.config.Security.Limiter.Body.Active {
			var isNotReadOnlyMethod = !(method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions)
			if isNotReadOnlyMethod {
				var requestURI = request.RequestURI
				var exceptSlice = m.config.Security.Limiter.Body.Except
				var isBypassed = false
				for _, exceptOne := range exceptSlice {
					if exceptOne == requestURI {
						isBypassed = true
					}
				}
				if !isBypassed {
					// TODO: fix max body length
					request.Body = http.MaxBytesReader(writer, request.Body, m.config.Security.Limiter.Body.MaxSize)
					//payload, err := ioutil.ReadAll(request.Body)
					//if err != nil {
					//	if strings.Contains(err.Error(), "request body too large"){
					//		http.Error(writer, "Request body too large.", 413)
					//		return
					//	}
					//}
				}
			}

		}
		next.ServeHTTP(writer, request)
	})
}
