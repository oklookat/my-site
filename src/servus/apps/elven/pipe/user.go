package pipe

import (
	"net/http"
	"servus/apps/elven/base"
	"servus/apps/elven/model"
)

type User struct {
	model *model.User
}

// is pipe model not empty?
func (u *User) IsAuthorized() bool {
	return u.model != nil
}

// get pipe by request context.
func (u *User) GetByContext(request *http.Request) base.UserPipe {
	pipe, ok := request.Context().Value(CtxUser).(base.UserPipe)
	if !ok {
		var emptyPipe = &User{}
		return emptyPipe
	}
	return pipe
}

// get pipe by user id.
func (u *User) GetByID(id string) (base.UserPipe, error) {
	var user = &model.User{}
	user.ID = id
	found, err := user.FindByID()
	if !found {
		return nil, err
	}
	var pipe = &User{}
	pipe.model = user
	return pipe, err
}

func (u *User) IsAdmin() bool {
	return u.IsAuthorized() && u.model.Role == userRoleAdmin
}

func (u *User) GetID() string {
	if !u.IsAuthorized() {
		return ""
	}
	return u.model.ID
}

func (u *User) GetUsername() string {
	if !u.IsAuthorized() {
		return ""
	}
	return u.model.Username
}

func (u *User) GetPassword() string {
	if !u.IsAuthorized() {
		return ""
	}
	return u.model.Password
}
