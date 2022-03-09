package pipe

import (
	"net/http"
	"servus/apps/elven/base"
	"servus/apps/elven/model"
)

type _ctx_user string

const CtxUser _ctx_user = "ELVEN_USER_PIPE"
const (
	userRoleAdmin = "admin"
	userRoleUser  = "user"
)

type User struct {
}

type UserPipe struct {
	model *model.User
}

// get pipe by request context.
func (u *User) GetByContext(request *http.Request) base.UserPipe {
	pipe, ok := request.Context().Value(CtxUser).(base.UserPipe)
	if !ok {
		return nil
	}
	return pipe
}

// get pipe by user id.
func (u *User) GetByID(id string) (*UserPipe, error) {
	var user = &model.User{}
	user.ID = id
	found, err := user.FindByID()
	if !found {
		return nil, err
	}
	var pipe = &UserPipe{}
	pipe.model = user
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
