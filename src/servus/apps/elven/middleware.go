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
		var h = call.Utils.GetHTTP(request)
		// check rights.
		var user = pipe.User{}
		userPipe := user.GetByContext(request)
		if userPipe == nil {
			h.Send(requestErrors.Forbidden(), 403, nil)
			return
		}
		next.ServeHTTP(response, request)
	})
}

// allow only safe methods for non-privileged users.
//
// https://developer.mozilla.org/en-US/docs/Glossary/Safe/HTTP
func (m *middleware) SafeMethodsOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var h = call.Utils.GetHTTP(request)
		// check method.
		var method = request.Method
		var safeMethod = method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions
		// check rights.
		var user = pipe.User{}
		userPipe := user.GetByContext(request)
		if !safeMethod {
			if userPipe == nil || !userPipe.IsAdmin() {
				h.Send(requestErrors.Forbidden(), 403, nil)
				return
			}
		}
		next.ServeHTTP(response, request)
	})
}

// only admin can access.
func (m *middleware) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var h = call.Utils.GetHTTP(request)
		// check rights.
		var user = pipe.User{}
		userPipe := user.GetByContext(request)
		if userPipe == nil || !userPipe.IsAdmin() {
			h.Send(requestErrors.Forbidden(), 403, nil)
			return
		}
		next.ServeHTTP(response, request)
	})
}

// gets token if exists and provides token pipe to request context.
func (m *middleware) ProvideTokenPipe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var ctx = request.Context()
		var token = pipe.Token{}
		tokenPipe, err := token.GetByRequest(request)
		hasToken := !(tokenPipe == nil || err != nil)
		if hasToken {
			var tokenCtx = context.WithValue(ctx, pipe.CtxToken, tokenPipe)
			*request = *request.WithContext(tokenCtx)
		}
		next.ServeHTTP(response, request)
	})
}

// gets user by token pipe and provides user pipe to request context. Use only after ProvideTokenPipe.
func (m *middleware) ProvideUserPipe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var ctx = request.Context()
		// get token.
		var token = pipe.Token{}
		tokenPipe := token.GetByContext(request)
		if tokenPipe != nil {
			// get user.
			var user = pipe.User{}
			userPipe, err := user.GetByID(tokenPipe.GetUserID())
			hasUser := !(userPipe == nil || err != nil)
			if hasUser {
				var userCtx = context.WithValue(ctx, pipe.CtxUser, userPipe)
				*request = *request.WithContext(userCtx)
			}
		}
		next.ServeHTTP(response, request)
	})
}
