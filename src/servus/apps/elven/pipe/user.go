package pipe

import (
	"net/http"
	"servus/apps/elven/foundation"
	"servus/apps/elven/model"
)

type _ctx_user string

const (
	CtxUser       _ctx_user = "ELVEN_USER_PIPE"
	userRoleAdmin           = "admin"
	userRoleUser            = "user"
)

type User struct {
}

type UserPipe struct {
	model *model.User
}

// GetByContext - get pipe by request context. Use only if you provided to request context.
func (u *User) GetByContext(request *http.Request) foundation.UserPipe {
	pipe, ok := request.Context().Value(CtxUser).(UserPipe)
	if !ok {
		return nil
	}
	return &pipe
}

// GetByID - used for ex. providing pipe to request context. Get id from Model.Token.
func (u *User) GetByID(id string) (*UserPipe, error) {
	var md = &model.User{}
	md.ID = id
	found, err := md.FindByID()
	if !found {
		return nil, err
	}
	var pipe = &UserPipe{}
	pipe.model = md
	return pipe, err
}

func (u *UserPipe) IsAdmin() bool {
	return u.model.Role == userRoleAdmin
}

func (u *UserPipe) GetID() string {
	return u.model.ID
}

func (u *UserPipe) GetUsername() string {
	return u.model.Username
}

func (u *UserPipe) GetPassword() string {
	return u.model.Password
}
