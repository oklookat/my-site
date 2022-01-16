package base

import "net/http"

// provides functions to get info about authorized user token.
type TokenPiper interface {
	// get pipe by request. Use in middleware.
	GetByRequest(request *http.Request) (TokenPipe, error)
	// get pipe by request context. Use only if you provided middleware.
	GetByContext(request *http.Request) TokenPipe
}

// get info about authorized user token.
type TokenPipe interface {
	GetID() string
	GetUserID() string
	GetToken() string
}
