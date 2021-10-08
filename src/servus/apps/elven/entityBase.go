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
		ec := errorMan.New()
		var auth = oUtil.createPipeAuth(request, accessTypeReadOnly)
		if !auth.Access {
			ec.AddEAuthForbidden([]string{"auth"})
			b.Send(response, ec.GetErrors(), 403)
			return
		}
		var ctx = request.Context()
		*request = *request.WithContext(context.WithValue(ctx, CtxAuthData, auth))
		next.ServeHTTP(response, request)
	})
}

func (b *entityBase) middlewareAdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		ec := errorMan.New()
		var auth = oUtil.createPipeAuth(request, accessTypeAdminOnly)
		if !auth.Access {
			ec.AddEAuthForbidden([]string{"auth"})
			b.Send(response, ec.GetErrors(), 403)
			return
		}
		var ctx = request.Context()
		*request = *request.WithContext(context.WithValue(ctx, CtxAuthData, auth))
		next.ServeHTTP(response, request)
	})
}
