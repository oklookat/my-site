package core

import (
	"net/http"
	"servus/core/internal/corsParse"
	"strings"
)

type _ctxHTTP string
const ctxHTTP _ctxHTTP = "CORE_HTTP_PIPE"

// Middleware - basic middleware.
type Middleware struct {
	core *Core
	cors cors
}

func (c *Core) bootMiddleware() {
	var corsConfig = corsParse.Config{
		AllowCredentials: c.Config.Security.CORS.AllowCredentials,
		AllowOrigin:      c.Config.Security.CORS.AllowOrigin,
		AllowMethods:     c.Config.Security.CORS.AllowMethods,
		AllowHeaders:     c.Config.Security.CORS.AllowHeaders,
		ExposeHeaders:    c.Config.Security.CORS.ExposeHeaders,
		MaxAge:           c.Config.Security.CORS.MaxAge,
	}
	var corsInstance = corsParse.New(corsConfig)
	c.Middleware = &Middleware{core: c, cors: &corsInstance}
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
		if m.core.Config.Security.CORS.Active {
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
		if m.core.Config.Security.Limiter.Body.Active {
			var isNotReadOnlyMethod = !(method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions)
			if isNotReadOnlyMethod {
				var requestURI = request.RequestURI
				var exceptSlice = m.core.Config.Security.Limiter.Body.Except
				var isBypassed = false
				for _, exceptOne := range exceptSlice {
					if exceptOne == requestURI {
						isBypassed = true
					}
				}
				if !isBypassed {
					// TODO: fix max body length
					request.Body = http.MaxBytesReader(writer, request.Body, m.core.Config.Security.Limiter.Body.MaxSize)
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
