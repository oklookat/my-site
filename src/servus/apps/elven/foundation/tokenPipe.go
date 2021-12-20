package foundation

import "net/http"

type TokenPiper interface {
	GetByRequest(request *http.Request) (TokenPipe, error)
	GetByContext(request *http.Request) TokenPipe
}

type TokenPipe interface {
	GetID() string
	GetUserID() string
	GetToken() string
}
