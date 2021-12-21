package base

import "net/http"

// TokenPiper - provides functions to get info about authorized user token.
type TokenPiper interface {
	// GetByRequest - get pipe by request. Use in middleware.
	GetByRequest(request *http.Request) (TokenPipe, error)
	// GetByContext - get pipe by request context. Use only if you provided middleware.
	GetByContext(request *http.Request) TokenPipe
}

// TokenPipe - get info about authorized user token.
type TokenPipe interface {
	GetID() string
	GetUserID() string
	GetToken() string
}
