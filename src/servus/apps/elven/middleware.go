package elven

import (
	"context"
	"net/http"
	"servus/core"
	"servus/core/external/errorMan"
)

type CtxAuthDataPipe string

const (
	CtxAuthData CtxAuthDataPipe = "ELVEN_PIPE_AUTH_DATA"
)

type middleware struct {
}

func (m *middleware) getHTTP(r *http.Request) core.HTTP {
	var h, _ = call.Middleware.GetHTTP(r)
	return h
}

// authorizedOnly - only authorized user can access.
func (m *middleware) authorizedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var h = m.getHTTP(request)
		var pipe = AuthPipe{}
		pipe.create(request, accessTypeAuthorized)
		if !pipe.Access {
			h.Send(errorMan.ThrowForbidden(), 403, nil)
			return
		}
		var ctx = request.Context()
		*request = *request.WithContext(context.WithValue(ctx, CtxAuthData, pipe))
		next.ServeHTTP(response, request)
	})
}

// middlewareSafeMethodsOnly - allow only safe methods for non-privileged users.
//
// https://developer.mozilla.org/en-US/docs/Glossary/Safe/HTTP
func (m *middleware) safeMethodsOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var h = m.getHTTP(request)
		var pipe = AuthPipe{}
		pipe.create(request, accessTypeReadOnly)
		if !pipe.Access {
			h.Send(errorMan.ThrowForbidden(), 403, nil)
			return
		}
		var ctx = request.Context()
		*request = *request.WithContext(context.WithValue(ctx, CtxAuthData, pipe))
		next.ServeHTTP(response, request)
	})
}

// adminOnly - only admin can access.
func (m *middleware) adminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var h = m.getHTTP(request)
		var pipe = AuthPipe{}
		pipe.create(request, accessTypeAdminOnly)
		if !pipe.Access {
			h.Send(errorMan.ThrowForbidden(), 403, nil)
			return
		}
		var ctx = request.Context()
		*request = *request.WithContext(context.WithValue(ctx, CtxAuthData, pipe))
		next.ServeHTTP(response, request)
	})
}
