package base

import "net/http"

// template for response.
type ResponseContent struct {
	Meta struct {
		PerPage     int `json:"per_page"`
		TotalPages  int `json:"total_pages"`
		CurrentPage int `json:"current_page"`
	} `json:"meta"`
	Data interface{} `json:"data"`
}

type MiddlewareAuthorizedOnly interface {
	AuthorizedOnly(next http.Handler) http.Handler
}

type MiddlewareAdminOnly interface {
	AdminOnly(next http.Handler) http.Handler
}

type MiddlewareSafeMethodsOnly interface {
	SafeMethodsOnly(next http.Handler) http.Handler
}
