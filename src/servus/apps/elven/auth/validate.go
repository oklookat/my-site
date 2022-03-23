package auth

import (
	"encoding/json"
	"io"
	"servus/core/external/validator"
)

// request body that user should send to login.
type Body struct {
	Username string
	Password string
	Type     string
}

func (a *Body) Validate(body io.ReadCloser) (isValid bool) {
	isValid = false

	// body.
	if err := json.NewDecoder(body).Decode(a); err != nil {
		return
	}

	// username.
	var username = a.Username
	if validator.IsEmpty(&username) {
		return
	}

	// password.
	var password = a.Password
	if validator.IsEmpty(&password) {
		return
	}

	// auth type.
	var authType = a.Type
	if validator.IsEmpty(&authType) {
		return
	}
	var isAuthType = authType == "cookie" || authType == "direct"
	if !isAuthType {
		return
	}

	isValid = true
	return
}
