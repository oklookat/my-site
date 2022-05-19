package user

import (
	"encoding/json"
	"io"
	"servus/apps/elven/base"
	"servus/core/external/validator"
)

type Body struct {
	What     string
	Password string
	NewValue string
}

func (u *Body) Validate(body io.ReadCloser) (err error) {
	err = json.NewDecoder(body).Decode(u)
	return
}

func ValidateUsername(username string) (err error) {
	var valErr = base.ValidationError{}
	if validator.MinMax(&username, 4, 24) {
		valErr.New("username")("invalid length: min 4, max 24")
		return &valErr
	}
	if !validator.IsAlphanumeric(&username) {
		valErr.New("username")("allowed only alphanumeric")
		return &valErr
	}
	return
}

func ValidatePassword(password string) (err error) {
	var valErr = base.ValidationError{}
	if validator.MinMax(&password, 8, 64) {
		valErr.New("password")("invalid length: min 8, max 64")
		return &valErr
	}
	if !validator.IsAlphanumericWithSymbols(&password) {
		valErr.New("password")("allowed only alphanumeric with symbols")
		return &valErr
	}
	return
}
