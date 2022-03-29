package base

import "net/http"

// get info about authorized user token.
type TokenPipe interface {
	// is token model exists?
	IsExists() bool

	// get pipe by request. Use in middleware.
	GetByRequest(request *http.Request) (TokenPipe, error)

	// get pipe by request context. Use only if you provided middleware.
	GetByContext(request *http.Request) TokenPipe

	// get token ID.
	GetID() string

	// get token user ID.
	GetUserID() string

	// get token.
	GetToken() string
}
