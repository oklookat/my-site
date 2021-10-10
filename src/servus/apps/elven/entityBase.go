package elven

import (
	"context"
	"net/http"
	"servus/core"
	"servus/core/modules/errorMan"
)

type CtxAuthDataPipe string

const (
	accessTypeAdminOnly                 = "ELVEN_ACCESS_ADMIN_ONLY"
	accessTypeReadOnly                  = "ELVEN_ACCESS_READ_ONLY"
	CtxAuthData         CtxAuthDataPipe = "ELVEN_PIPE_AUTH_DATA"
)

type entityBase struct {
	*core.BaseController
}

func (b *entityBase) middlewareReadOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var auth = oUtils.createPipeAuth(request, accessTypeReadOnly)
		if !auth.Access {
			b.Send(response, errorMan.ThrowForbidden(), 403)
			return
		}
		var ctx = request.Context()
		*request = *request.WithContext(context.WithValue(ctx, CtxAuthData, auth))
		next.ServeHTTP(response, request)
	})
}

func (b *entityBase) middlewareAdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var auth = oUtils.createPipeAuth(request, accessTypeAdminOnly)
		if !auth.Access {
			b.Send(response, errorMan.ThrowForbidden(), 403)
			return
		}
		var ctx = request.Context()
		*request = *request.WithContext(context.WithValue(ctx, CtxAuthData, auth))
		next.ServeHTTP(response, request)
	})
}
