package user

import (
	"encoding/json"
	"io"
	"servus/core/external/validator"
)

type Body struct {
	What     string
	Password string
	NewValue string
}

func (u *Body) Validate(body io.ReadCloser) (isValid bool) {
	isValid = false
	err := json.NewDecoder(body).Decode(u)
	if err != nil {
		return
	}
	isValid = true
	return
}

func ValidateUsername(username string) (isValid bool) {
	isValid = false
	if validator.MinMax(&username, 4, 24) {
		return
	}
	if !validator.IsAlphanumeric(&username) {
		return
	}
	isValid = true
	return
}

func ValidatePassword(password string) (isValid bool) {
	isValid = false
	if validator.MinMax(&password, 8, 64) {
		return
	}
	if !validator.IsAlphanumericWithSymbols(&password) {
		return
	}
	isValid = true
	return
}
