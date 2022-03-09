package base

import "net/http"

// provides functions to get info about authorized user.
type UserPiper interface {
	// get pipe by request context. Use only if you provided token and user middlewares.
	GetByContext(request *http.Request) UserPipe
}

// get info about authorized user.
type UserPipe interface {
	IsAdmin() bool
	GetID() string
	GetUsername() string
	GetPassword() string
}
