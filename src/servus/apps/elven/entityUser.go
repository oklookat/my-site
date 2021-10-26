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
	var body = struct {
		What     string
		Password string
		NewValue string
	}{}
	em := errorMan.NewValidation()
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		em.Add("body", "wrong value provided.")
		u.Send(response, em.GetJSON(), 400)
		return
	}
	auth := PipeAuth{}
	auth.get(request)
	match, err := instance.Encryption.Argon.Check(body.Password, auth.User.Password)
	if err != nil || !match {
		u.Send(response, errorMan.ThrowNotAuthorized(), 401)
		return
	}
	var changeUsername = false
	switch body.What {
	case "username":
		changeUsername = true
		auth.User.Username = body.NewValue
		err = auth.User.validateUsername()
		break
	case "password":
		auth.User.Password = body.NewValue
		err = auth.User.validatePassword()
		break
	}
	if err != nil {
		em.Add(body.What, "wrong value provided.")
		u.Send(response, em.GetJSON(), 400)
		return
	}
	if changeUsername {
		userExists, err := auth.User.findByUsername()
		if err != nil {
			u.Send(response, errorMan.ThrowServer(), 500)
			return
		}
		if userExists {
			em.Add("username", "already exists")
			u.Send(response, em.GetJSON(), 400)
			return
		}
	}
	err = auth.User.update()
	if err != nil {
		u.Send(response, errorMan.ThrowServer(), 500)
		return
	}
	u.Send(response, "", 200)
	return
}
