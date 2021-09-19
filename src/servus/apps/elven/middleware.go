package elven

import "net/http"

func middlewareReadOnly(next http.Handler) http.Handler{
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

			next.ServeHTTP(writer, request)
		})
}
