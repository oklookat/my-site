package middleware

import "net/http"

type mFunc func(http.Handler) http.Handler

type Instance struct {
	cors        mFunc
	limitBody   mFunc
	provideHTTP mFunc
	asJson      mFunc
}

func (m *Instance) New(
	cors mFunc,
	limitBody mFunc,
	httpHelper mFunc,
) {
	if cors == nil || limitBody == nil || httpHelper == nil {
		panic("[middleware]: one of args has nil pointer")
	}
	m.cors = cors
	m.limitBody = limitBody
	m.provideHTTP = httpHelper
	m.asJson = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Add("Content-Type", "application/json; charset=utf-8")
			next.ServeHTTP(writer, request)
		})
	}
}

func (i *Instance) AsJson() func(http.Handler) http.Handler {
	return i.asJson
}

func (i *Instance) CORS() func(http.Handler) http.Handler {
	return i.cors
}

func (i *Instance) LimitBody() func(http.Handler) http.Handler {
	return i.limitBody
}

func (i *Instance) ProvideHTTP() func(http.Handler) http.Handler {
	return i.provideHTTP
}
