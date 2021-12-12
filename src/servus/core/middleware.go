package core

import (
	"github.com/pkg/errors"
	"net/http"
	"servus/core/internal/iHTTP"
)

type _ctxHTTP string

const ctxHTTP _ctxHTTP = "CORE_HTTP_PIPE"

type Middleware struct {
	CORS        func(http.Handler) http.Handler
	LimitBody   func(http.Handler) http.Handler
	ProvideHTTP func(http.Handler) http.Handler
	AsJson      func(http.Handler) http.Handler
}

func (m *Middleware) new(cors MiddlewareFunc, limitBody MiddlewareFunc, _http MiddlewareFunc) {
	m.CORS = cors
	m.LimitBody = limitBody
	m.ProvideHTTP = _http
	m.AsJson = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Add("Content-Type", "application/json; charset=utf-8")
			next.ServeHTTP(writer, request)
		})
	}
}

// GetHTTP - get HTTP from request context.
func (m *Middleware) GetHTTP(request *http.Request) (HTTP, error) {
	var ctx = request.Context()
	var h, ok = ctx.Value(ctxHTTP).(*iHTTP.Instance)
	if !ok {
		var err = errors.New("[core/middleware]: error while get HTTP struct from request context. Are you provided ProvideHTTP limiter?")
		return nil, err
	}
	return h, nil
}
