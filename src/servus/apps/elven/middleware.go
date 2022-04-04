package elven

import (
	"context"
	"net/http"
	"servus/apps/elven/pipe"
)

type middleware struct {
}

// only authorized user can access.
func (m *middleware) AuthorizedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var h = call.Http.Get(request)

		// get pipe.
		var user = pipe.User{}
		userPipe := user.GetByContext(request)

		// check rights.
		if !userPipe.IsAuthorized() {
			h.Send(requestErrors.Forbidden(), 401, nil)
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
		var user = pipe.User{}
		userPipe := user.GetByContext(request)

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
		h.Send(requestErrors.Forbidden(), statusCode, nil)
		return

	})
}

// only admin can access.
func (m *middleware) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var h = call.Http.Get(request)

		// check rights.
		var user = pipe.User{}
		userPipe := user.GetByContext(request)
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
		h.Send(requestErrors.Forbidden(), statusCode, nil)
		return

	})
}

// get token if exists and provide token pipe to request context.
func (m *middleware) ProvideTokenPipe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var ctx = request.Context()
		// get token.
		var token = pipe.Token{}
		var tokenPipe, err = token.GetByRequest(request)
		var isHasTokenPipe = !(tokenPipe == nil || err != nil)
		if isHasTokenPipe {
			var tokenCtx = context.WithValue(ctx, pipe.CtxToken, tokenPipe)
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
		var token = pipe.Token{}
		tokenPipe := token.GetByContext(request)
		if !tokenPipe.IsExists() {
			next.ServeHTTP(response, request)
			return
		}

		// get user.
		var user = pipe.User{}
		var userPipe, err = user.GetByID(tokenPipe.GetUserID())
		var isHasUserPipe = !(userPipe == nil || err != nil)
		if isHasUserPipe {
			var userCtx = context.WithValue(ctx, pipe.CtxUser, userPipe)
			*request = *request.WithContext(userCtx)
		}
		next.ServeHTTP(response, request)
	})
}
