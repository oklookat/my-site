package foundation

import "net/http"

type MiddlewareAuthorizedOnly interface {
	AuthorizedOnly(next http.Handler) http.Handler
}

type MiddlewareAdminOnly interface {
	AdminOnly(next http.Handler) http.Handler
}

type MiddlewareSafeMethodsOnly interface {
	SafeMethodsOnly(next http.Handler) http.Handler
}
