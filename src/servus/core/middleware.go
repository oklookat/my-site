package core

import (
	"net/http"
	"strings"
)

// MiddlewareAsJSON this middleware set application/json header
func (m *BasicMiddleware) MiddlewareAsJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})
}

// MiddlewareCORS CORS middleware depending on config file
func (m *BasicMiddleware) MiddlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if Config.Security.CORS.Active {
			corsParser.SetHeaders(writer, request)
		}
		next.ServeHTTP(writer, request)
	})
}

// MiddlewareSecurity basic security depending on config file
func (m *BasicMiddleware) MiddlewareSecurity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var method = request.Method
		method = strings.ToUpper(method)
		//////// request body size check
		if Config.Security.Limiter.Body.Active {
			var isNotReadOnlyMethod = method == "POST" || method == "PUT" || method == "DELETE" || method == "PATCH"
			if isNotReadOnlyMethod {
				var requestURI = request.RequestURI
				var exceptSlice = Config.Security.Limiter.Body.Except
				var isBypassed = false
				for _, exceptOne := range exceptSlice {
					if exceptOne == requestURI{
						isBypassed = true
					}
				}
				if !isBypassed {
					request.Body = http.MaxBytesReader(writer, request.Body, Config.Security.Limiter.Body.MaxSize)
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
