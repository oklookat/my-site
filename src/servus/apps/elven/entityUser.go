package elven

import (
	"encoding/json"
	"net/http"
	"servus/core/external/errorMan"
)

// entityUser - manage users.
type entityUser struct {
	*entityBase
}

type ResponseUser struct {
	IsAdmin  bool   `json:"is_admin"`
	Username string `json:"username"`
}

// GET
// controllerGetMe - send some user data by token.
func (u *entityUser) controllerGetMe(response http.ResponseWriter, request *http.Request) {
	auth := PipeAuth{}
	auth.get(request)
	var resp = ResponseUser{}
	resp.IsAdmin = auth.IsAdmin
	resp.Username = auth.User.Username
	bytes, err := json.Marshal(resp)
	if err != nil {
		u.Send(response, errorMan.ThrowServer(), 500)
		return
	}
	u.Send(response, string(bytes), 200)
	return
}

// POST
// controllerChange - change username or password.
// body: what change (username or password); password to confirm; new value for change.
func (u *entityUser) controllerChange(response http.ResponseWriter, request *http.Request) {
	var data = struct {
		what     string
		password string
		newValue string
	}{}
	em := errorMan.NewValidation()
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		em.Add("body", "wrong value provided.")
		u.Send(response, em.GetJSON(), 400)
		return
	}
	auth := PipeAuth{}
	auth.get(request)
	match, err := instance.Encryption.Argon.Check(data.password, auth.User.Password)
	if err != nil || !match {
		u.Send(response, errorMan.ThrowNotAuthorized(), 401)
		return
	}
	switch data.what {
	case "username":
		auth.User.Username = data.newValue
		err = auth.User.validateUsername()
		break
	case "password":
		auth.User.Password = data.password
		err = auth.User.validatePassword()
		break
	}
	if err != nil {
		em.Add(data.what, "wrong value provided.")
		u.Send(response, em.GetJSON(), 400)
		return
	}
	err = auth.User.update()
	if err != nil {
		u.Send(response, errorMan.ThrowServer(), 500)
		return
	}
	u.Send(response, "", 200)
	return
}
