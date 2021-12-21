package auth

import (
	"encoding/json"
	"io"
	"servus/apps/elven/base"
	coreValidator "servus/core/external/validator"
)

// Body - represents the body of the request that the user should send to login.
type Body struct {
	Username string
	Password string
	Type     string
}

func (a *Body) Validate(body io.ReadCloser) base.Validator {
	var val = validate.Create()
	// body.
	err := json.NewDecoder(body).Decode(a)
	if err != nil {
		val.Add("body")
		return val
	}
	// username.
	var username = a.Username
	if coreValidator.IsEmpty(&username) {
		val.Add("username")
	}
	// password.
	var password = a.Password
	if coreValidator.IsEmpty(&password) {
		val.Add("password")
	}
	// auth type.
	var authType = a.Type
	if coreValidator.IsEmpty(&authType) {
		val.Add("type")
	} else {
		var isAuthType = authType == "cookie" || authType == "direct"
		if !isAuthType {
			val.Add("type")
		}
	}
	return val
}
