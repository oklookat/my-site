package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"servus/apps/elven/model"
)

// send some user data by token (GET).
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

// change username or password (POST).
func (u *Instance) change(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)

	// validate body.
	var body = Body{}
	isValid := body.Validate(request.Body)
	if !isValid {
		h.Send("invalid request", 400, nil)
		return
	}

	// get pipe.
	pipe := u.pipe.GetByContext(request)

	// compare confirm password from body and original password from pipe.
	match, err := call.Encryptor.Argon.Compare(body.Password, pipe.GetPassword())
	if err != nil || !match {
		h.Send(u.throw.NotAuthorized(), 401, err)
		return
	}

	var sendValidationErr = func() {
		h.Send("invalid request", 400, err)
		return
	}

	// configure model & update.
	var user = model.User{}
	user.ID = pipe.GetID()
	switch body.What {
	case "username":
		// validate username.
		isValid = ValidateUsername(body.NewValue)
		if !isValid {
			sendValidationErr()
			return
		}
		user.Username = body.NewValue

		// check is username in use.
		isUsernameTaken, err := user.FindByUsername()
		if err != nil {
			h.Send(u.throw.Server(), 500, err)
			return
		}
		if isUsernameTaken {
			h.Send("username already in use", 409, err)
			return
		}
	case "password":
		// validate password.
		isValid = ValidatePassword(body.NewValue)
		if !isValid {
			sendValidationErr()
			return
		}
		user.Password = body.NewValue
	}

	// update.
	err = user.Update()
	if err != nil {
		h.Send(u.throw.Server(), 500, err)
		return
	}
	h.Send("", 200, err)
}

// create new user. THIS IS NOT A ROUTE.
func (u *Instance) Create(username string, password string, isAdmin bool) (err error) {

	// validate.
	isValid := ValidateUsername(username)
	if !isValid {
		return errors.New("invalid username")
	}
	isValid = ValidatePassword(password)
	if !isValid {
		return errors.New("invalid password")
	}

	// fill.
	var user = model.User{}
	var role string
	if isAdmin {
		role = "admin"
	} else {
		role = "user"
	}
	user.Role = role
	user.Username = username
	user.Password = password

	// user with this username exists?
	found, err := user.FindByUsername()
	if err != nil {
		return
	}
	if found {
		return errors.New("user with this username already exists")
	}

	// create
	err = user.Create()

	return
}

func (u *Instance) DeleteByUsername(username string) (err error) {
	var user = model.User{}
	user.Username = username
	found, err := user.FindByUsername()
	if !found {
		return errors.New("user not found")
	}
	if err != nil {
		return
	}
	err = user.DeleteByID()
	return
}
