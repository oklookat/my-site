package user

import (
	"net/http"
	"servus/apps/elven/base"
)

type _string string

const (
	CTX                  _string = "ELVEN_USER_PIPE"
	AccessTypeAdminOnly          = "ELVEN_ACCESS_ADMIN_ONLY"
	AccessTypeReadOnly           = "ELVEN_ACCESS_READ_ONLY"
	AccessTypeAuthorized         = "ELVEN_ACCESS_AUTHORIZED"
	userRoleAdmin                = "admin"
	userRoleUser                 = "user"
)

type Pipe struct {
	model *Model
}

// is pipe model not empty?
func (p *Pipe) IsAuthorized() bool {
	return p.model != nil
}

// get pipe by request context.
func (p *Pipe) GetByContext(request *http.Request) base.UserPipe {
	pipe, ok := request.Context().Value(CTX).(base.UserPipe)
	if !ok {
		var emptyPipe = &Pipe{}
		return emptyPipe
	}
	return pipe
}

// get pipe by user id.
func (p *Pipe) GetByID(id string) (base.UserPipe, error) {
	var user = &Model{}
	user.ID = id
	found, err := user.FindByID()
	if !found {
		return nil, err
	}
	var pipe = &Pipe{}
	pipe.model = user
	return pipe, err
}

func (p *Pipe) IsAdmin() bool {
	return p.IsAuthorized() && p.model.Role == userRoleAdmin
}

func (p *Pipe) GetID() string {
	if !p.IsAuthorized() {
		return ""
	}
	return p.model.ID
}

func (p *Pipe) GetUsername() string {
	if !p.IsAuthorized() {
		return ""
	}
	return p.model.Username
}

func (p *Pipe) GetPassword() string {
	if !p.IsAuthorized() {
		return ""
	}
	return p.model.Password
}
