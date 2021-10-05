package elven

import (
	"context"
	"net/http"
	"servus/core"
	"servus/core/modules/errorCollector"
)

type CtxAuthDataPipe string

const (
	accessTypeAdminOnly                 = "ACCESS_ADMIN_ONLY"
	accessTypeReadOnly                  = "ACCESS_READ_ONLY"
	CtxAuthData CtxAuthDataPipe = "ELVEN_PIPE_AUTH_DATA"
)

type AuthData struct {
	Access bool
	UserAndTokenExists bool
	IsAdmin bool
	User   *ModelUser
	Token  *ModelToken
}

func middlewareReadOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var auth = createAuthData(request, accessTypeReadOnly)
		if !auth.Access {
			var ec = errorCollector.New()
			ec.AddEAuthForbidden([]string{"AUTH"})
			var res = core.HttpResponse{ResponseWriter: response}
			res.Send(ec.GetErrors(), 403)
			return
		}
		var ctx = request.Context()
		*request = *request.WithContext(context.WithValue(ctx, CtxAuthData, auth))
		next.ServeHTTP(response, request)
	})
}

func middlewareAdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var auth = createAuthData(request, accessTypeAdminOnly)
		if !auth.Access {
			var ec = errorCollector.New()
			ec.AddEAuthForbidden([]string{"AUTH"})
			var res = core.HttpResponse{ResponseWriter: response}
			res.Send(ec.GetErrors(), 403)
			return
		}
		var ctx = request.Context()
		*request = *request.WithContext(context.WithValue(ctx, CtxAuthData, auth))
		next.ServeHTTP(response, request)
	})
}
