package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"servus/apps/elven/model"
)

var (
	ErrUserExists   = errors.New("username already exists")
	ErrUserNotFound = errors.New("user not found")
)

// send some user data by token (GET).
func getMe(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)
	var pipe = pipe.GetByContext(request)
	var resp = Response{}
	resp.IsAdmin = pipe.IsAdmin()
	resp.Username = pipe.GetUsername()
	bytes, err := json.Marshal(resp)
	if err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}
	h.Send(string(bytes), 200, err)
}

// change username or password (POST).
func change(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)

	// validate body.
	var body = Body{}
	if err := body.Validate(request.Body); err != nil {
		h.Send("invalid request", 400, nil)
		return
	}

	// get pipe.
	var pipe = pipe.GetByContext(request)

	// compare confirm password from body and original password from pipe.
	match, err := call.Encryptor.Argon.Compare(body.Password, pipe.GetPassword())
	if err != nil || !match {
		h.Send(throw.NotAuthorized(), 401, err)
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
		if err = ValidateUsername(body.NewValue); err != nil {
			sendValidationErr()
			return
		}
		user.Username = body.NewValue

		// check is username in use.
		isUsernameTaken, err := user.FindByUsername()
		if err != nil {
			h.Send(throw.Server(), 500, err)
			return
		}
		if isUsernameTaken {
			h.Send("username already in use", 409, err)
			return
		}
	case "password":
		// validate password.
		if err = ValidatePassword(body.NewValue); err != nil {
			sendValidationErr()
			return
		}
		user.Password = body.NewValue
	}

	// update.
	if err = user.Update(); err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}
	h.Send("", 200, err)
}

// create new user. THIS IS NOT A ROUTE.
func Create(username string, password string, isAdmin bool) (err error) {

	// validate.
	if err = ValidateUsername(username); err != nil {
		return errors.New("invalid username")
	}
	if err = ValidatePassword(password); err != nil {
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
		return ErrUserExists
	}

	// create
	err = user.Create()

	return
}

// delete user. THIS IS NOT A ROUTE.
func DeleteByUsername(username string) (err error) {
	var user = model.User{}
	user.Username = username
	found, err := user.FindByUsername()
	if !found {
		return ErrUserNotFound
	}
	if err != nil {
		return
	}
	err = user.DeleteByID()
	return
}
