package user

import (
	"encoding/json"
	"net/http"
	"servus/apps/elven/model"
)

// GET
// send some user data by token.
func (u *Instance) getMe(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	pipe := u.pipe.GetByContext(request)
	var resp = Response{}
	resp.IsAdmin = pipe.IsAdmin()
	resp.Username = pipe.GetUsername()
	bytes, err := json.Marshal(resp)
	if err != nil {
		h.Send(u.throw.Server(), 500, err)
		return
	}
	h.Send(string(bytes), 200, err)
}

// POST
// change username or password.
// body: what change (username or password); password to confirm; new value for change.
func (u *Instance) change(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	// validate body.
	var body = Body{}
	validator := body.Validate(request.Body)
	if validator.HasErrors() {
		h.Send(validator.GetJSON(), 400, nil)
		return
	}
	// get pipe.
	pipe := u.pipe.GetByContext(request)
	match, err := call.Encryptor.Argon.Compare(body.Password, pipe.GetPassword())
	if err != nil || !match {
		h.Send(u.throw.NotAuthorized(), 401, err)
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
			h.Send(u.throw.Server(), 500, err)
			return
		}
		if userExists {
			validator.Add("username")
			h.Send(validator.GetJSON(), 409, err)
			return
		}
	} else if body.What == "password" {
		user.Password = body.NewValue
	}
	err = user.Update()
	if err != nil {
		// check if validation error.
		validationErr := err == model.ErrUserUsernameMinMax || err == model.ErrUserUsernameAlphanumeric || err == model.ErrUserPasswordMinMax || err == model.ErrUserPasswordWrongSymbols
		if validationErr {
			validator.Add(body.What)
			h.Send(validator.GetJSON(), 400, err)
			return
		}
		h.Send(u.throw.Server(), 500, err)
		return
	}
	h.Send("", 200, err)
}
