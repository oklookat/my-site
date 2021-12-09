package elven

import (
	"encoding/json"
	"net/http"
	"servus/core/external/errorMan"
)

// GET
// getMe - send some user data by token.
func (u *userRoute) getMe(response http.ResponseWriter, request *http.Request) {
	var h = u.middleware.getHTTP(request)
	pipe := AuthPipe{}
	pipe.get(request)
	var resp = ResponseUser{}
	resp.IsAdmin = pipe.IsAdmin
	resp.Username = pipe.User.Username
	bytes, err := json.Marshal(resp)
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	h.Send(string(bytes), 200, err)
	return
}

// POST
// change - change username or password.
// body: what change (username or password); password to confirm; new value for change.
func (u *userRoute) change(response http.ResponseWriter, request *http.Request) {
	var h = u.middleware.getHTTP(request)
	var body = struct {
		What     string
		Password string
		NewValue string
	}{}
	em := errorMan.NewValidation()
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		em.Add("body", "wrong value provided.")
		h.Send(em.GetJSON(), 400, err)
		return
	}
	pipe := AuthPipe{}
	pipe.get(request)
	match, err := call.Encryption.Argon.Check(body.Password, pipe.User.Password)
	if err != nil || !match {
		h.Send(errorMan.ThrowNotAuthorized(), 401, err)
		return
	}
	var changeUsername = false
	switch body.What {
	case "username":
		changeUsername = true
		pipe.User.Username = body.NewValue
		err = pipe.User.validateUsername()
		break
	case "password":
		pipe.User.Password = body.NewValue
		err = pipe.User.validatePassword()
		break
	}
	if err != nil {
		em.Add(body.What, "wrong value provided.")
		h.Send(em.GetJSON(), 400, err)
		return
	}
	if changeUsername {
		userExists, err := pipe.User.findByUsername()
		if err != nil {
			h.Send(errorMan.ThrowServer(), 500, err)
			return
		}
		if userExists {
			em.Add("username", "already exists")
			h.Send(em.GetJSON(), 400, err)
			return
		}
	}
	err = pipe.User.update()
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	h.Send("", 200, err)
	return
}
