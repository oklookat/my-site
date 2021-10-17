package elven

import (
	"github.com/pkg/errors"
	"servus/core/external/validator"
)

// entityUser - manage users.
type entityUser struct {
}


// validatorUsername - validate username from ModelUser. Used in cmd create user.
func (u *entityUser) validatorUsername(username string) error {
	if validator.MinMax(&username, 4, 24) {
		return errors.New("username: min length 4 and max 24")
	}
	if !validator.IsAlphanumeric(&username) {
		return errors.New("username: allowed only alphanumeric")
	}
	return nil
}

// validatorPassword - validate ModelUser password. Used in cmd create user.
func (u *entityUser) validatorPassword(password string) error {
	if len(password) < 8 || len(password) > 64 {
		return errors.New("password: min length 8 and max 64")
	}
	if !validator.IsAlphanumericWithSymbols(&password) {
		return errors.New("password: allowed only alphanumeric and some symbols")
	}
	return nil
}
