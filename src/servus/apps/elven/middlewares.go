package elven

import (
	"context"
	"net/http"
	"servus/core"
	"servus/core/modules/errorCollector"
)

type ctxAuthDataPipe string

const (
	accessTypeAdminOnly = "ACCESS_ADMIN_ONLY"
	accessTypeReadOnly  = "ACCESS_READ_ONLY"
	ctxAuthData         = "ELVEN_PIPE_AUTH_DATA"
)

type authData struct {
	access bool
	user   modelUser
	token  modelToken
}

func middlewareReadOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var auth = createAuthData(request, accessTypeReadOnly)
		if !auth.access {
			var ec = errorCollector.New()
			ec.AddEAuthForbidden([]string{"AUTH"})
			var res = core.HttpResponse{ResponseWriter: response}
			res.Send(ec.GetErrors(), 403)
			return
		}
		var ctx = request.Context()
		request = request.WithContext(context.WithValue(ctx, ctxAuthData, auth))
		next.ServeHTTP(response, request)
	})
}

func middlewareAdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var auth = createAuthData(request, accessTypeAdminOnly)
		if !auth.access {
			var ec = errorCollector.New()
			ec.AddEAuthForbidden([]string{"AUTH"})
			var res = core.HttpResponse{ResponseWriter: response}
			res.Send(ec.GetErrors(), 403)
			return
		}
		var ctx = request.Context()
		request = request.WithContext(context.WithValue(ctx, ctxAuthData, auth))
		next.ServeHTTP(response, request)
	})
}
