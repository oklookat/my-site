package core

import (
	"net/http"
)

type middleware struct {
	cors        func(http.Handler) http.Handler
	limitBody   func(http.Handler) http.Handler
	provideHTTP func(http.Handler) http.Handler
	asJson      func(http.Handler) http.Handler
}

func (m *middleware) new(cors MiddlewareFunc, limitBody MiddlewareFunc, _http MiddlewareFunc) {
	m.cors = cors
	m.limitBody = limitBody
	m.provideHTTP = _http
	m.asJson = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Add("Content-Type", "application/json; charset=utf-8")
			next.ServeHTTP(writer, request)
		})
	}
}

func (m *middleware) AsJson() func(http.Handler) http.Handler {
	return m.asJson
}

func (m *middleware) CORS() func(http.Handler) http.Handler {
	return m.cors
}

func (m *middleware) LimitBody() func(http.Handler) http.Handler {
	return m.limitBody
}

func (m *middleware) ProvideHTTP() func(http.Handler) http.Handler {
	return m.provideHTTP
}
