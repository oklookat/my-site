package core

import (
	"net/http"

	"github.com/pkg/errors"
)

type _ctxHTTP string

const ctxHTTP _ctxHTTP = "CORE_HTTP_PIPE"

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

func (m *middleware) GetHTTP(request *http.Request) (HTTP, error) {
	var ctx = request.Context()
	var h, ok = ctx.Value(ctxHTTP).(HTTP)
	if !ok {
		var err = errors.New("[core/middleware]: error while get HTTP struct from request context. Are you provided ProvideHTTP limiter?")
		return nil, err
	}
	return h, nil
}
