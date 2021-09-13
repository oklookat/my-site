package core

import "net/http"


// MiddlewareAsJSON this middleware set application/json header
func (m *Middleware) MiddlewareAsJSON(next http.Handler) http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})
}

// MiddlewareCORS this middleware set cors headers depending on config file
func (m *Middleware) MiddlewareCORS(next http.Handler) http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if servus.Config.Security.CORS.Active {
			parseCorsConfig(writer, request)
			// bypass any auth if simple request
			var method = request.Method
			var isSimpleRequest = method == "OPTIONS"
			if isSimpleRequest {
				writer.WriteHeader(200)
				writer.Write([]byte(""))
				return
			}
		}
		next.ServeHTTP(writer, request)
	})
}
