package way

import "net/http"

// Middleware - making things before something.
type Middleware func(next http.Handler) http.Handler

// middleware - internal representation of Middleware.
type middleware struct {
	// executed - all middleware executed = true, one of middleware send response = false.
	executed bool
	// chain - middlewares here.
	chain http.Handler
}

// add - make one big Middleware from middlewares. At the end, call ServeHTTP.
func (m *middleware) add(middlewares []Middleware) {
	m.chain = m.makeChain(middlewares, m)
}

// ServeHTTP - calls when all middlewares executed.
func (m *middleware) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	m.executed = true
}

// makeChain - make one big middleware from middlewares. At the end, call next.
func (m *middleware) makeChain(middlewares []Middleware, next http.Handler) http.Handler {
	// https://gist.github.com/husobee/fd23681261a39699ee37#gistcomment-3111569
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i](next)
	}
	return next
}
