package user

import (
	"encoding/json"
	"net/http"
	"servus/apps/elven/model"
	"servus/core/external/errorMan"
)

// GET
// getMe - send some user data by token.
func (u *route) getMe(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	pipe := u.pipe.GetByContext(request)
	var resp = Response{}
	resp.IsAdmin = pipe.IsAdmin()
	resp.Username = pipe.GetUsername()
	bytes, err := json.Marshal(resp)
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	h.Send(string(bytes), 200, err)
}

// POST
// change - change username or password.
// body: what change (username or password); password to confirm; new value for change.
func (u *route) change(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	// validate body.
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
	// get pipe.
	pipe := u.pipe.GetByContext(request)
	match, err := call.Encryptor.Argon.Compare(body.Password, pipe.GetPassword())
	if err != nil || !match {
		h.Send(errorMan.ThrowNotAuthorized(), 401, err)
		return
	}
	// creating model and copy id from pipe.
	var user = model.User{}
	user.ID = pipe.GetID()
	if body.What == "username" {
		user.Username = body.NewValue
		// check if exists.
		userExists, err := user.FindByUsername()
		if err != nil {
			h.Send(errorMan.ThrowServer(), 500, err)
			return
		}
		if userExists {
			em.Add("username", "already exists")
			h.Send(em.GetJSON(), 400, err)
			return
		}
	} else if body.What == "password" {
		user.Password = body.NewValue
	}
	err = user.Update()
	if err != nil {
		// check if validation error.
		validationErr := err == model.ErrUserUsernameValidation || err == model.ErrUserPasswordValidation
		if validationErr {
			em.Add(body.What, "wrong value provided.")
			h.Send(em.GetJSON(), 400, err)
			return
		}
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	h.Send("", 200, err)
}
