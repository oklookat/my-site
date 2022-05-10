package elven

import (
	"context"
	"net/http"
	"servus/apps/elven/token"
	"servus/apps/elven/user"
)

type middleware struct {
}

// only authorized user can access.
func (m *middleware) AuthorizedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var h = call.Http.Get(request)

		// get pipe.
		var userGetter = user.Pipe{}
		var userPipe = userGetter.GetByContext(request)

		// check rights.
		if !userPipe.IsAuthorized() {
			h.Send("", 401, nil)
			return
		}

		next.ServeHTTP(response, request)
	})
}

// allow only safe methods if user not admin.
//
// https://developer.mozilla.org/en-US/docs/Glossary/Safe/HTTP
func (m *middleware) SafeMethodsOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var h = call.Http.Get(request)

		// check method.
		var method = request.Method
		var isSafeMethod = method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions
		if isSafeMethod {
			next.ServeHTTP(response, request)
			return
		}

		// check rights.
		var userGetter = user.Pipe{}
		var userPipe = userGetter.GetByContext(request)

		// allow if admin.
		var isAdmin = userPipe.IsAdmin()
		if isAdmin {
			next.ServeHTTP(response, request)
			return
		}

		// access denied. Set status code.
		var statusCode = 403 // 403 = authorized, but not admin.
		if !userPipe.IsAuthorized() {
			statusCode = 401 // 401 = not authorized
		}
		h.Send("", statusCode, nil)
	})
}

// only admin can access.
func (m *middleware) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var h = call.Http.Get(request)

		// check rights.
		var userGetter = user.Pipe{}
		var userPipe = userGetter.GetByContext(request)
		var isAdmin = userPipe.IsAdmin()
		if isAdmin {
			next.ServeHTTP(response, request)
			return
		}

		// access denied. Set status code.
		var statusCode = 403 // 403 = authorized, but not admin.
		if !userPipe.IsAuthorized() {
			statusCode = 401 // 401 = not authorized
		}
		h.Send("", statusCode, nil)
	})
}

// get token if exists and provide token pipe to request context.
func (m *middleware) ProvideTokenPipe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var ctx = request.Context()
		// get token.
		var tokenGetter = token.Pipe{}
		var tokenPipe, err = tokenGetter.GetByRequest(request)
		var isHasTokenPipe = !(tokenPipe == nil || err != nil)
		if isHasTokenPipe {
			var tokenCtx = context.WithValue(ctx, token.CTX, tokenPipe)
			*request = *request.WithContext(tokenCtx)
		}
		next.ServeHTTP(response, request)
	})
}

// get user by token pipe and provide user pipe to request context. Use only after ProvideTokenPipe.
func (m *middleware) ProvideUserPipe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var ctx = request.Context()

		// get token.
		var tokenGetter = token.Pipe{}
		tokenPipe := tokenGetter.GetByContext(request)
		if !tokenPipe.IsExists() {
			next.ServeHTTP(response, request)
			return
		}

		// get user.
		var userGetter = user.Pipe{}
		var userPipe, err = userGetter.GetByID(tokenPipe.GetUserID())
		var isHasUserPipe = !(userPipe == nil || err != nil)
		if isHasUserPipe {
			var userCtx = context.WithValue(ctx, user.CTX, userPipe)
			*request = *request.WithContext(userCtx)
		}
		next.ServeHTTP(response, request)
	})
}
