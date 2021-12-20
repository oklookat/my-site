package foundation

import "net/http"

type UserPiper interface {
	GetByContext(request *http.Request) UserPipe
	//GetByUsername(username string) (UserPipe, error)
}

type UserPipe interface {
	IsAdmin() bool
	GetID() string
	GetUsername() string
	GetPassword() string
}
