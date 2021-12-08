package elven

import (
	"github.com/pkg/errors"
	"servus/core/external/validator"
)



// validateUsername - validate username from UserModel. Used in cmd create user.
func (u *UserModel) validateUsername() error {
	if validator.MinMax(&u.Username, 4, 24) {
		return errors.New("username: min length 4 and max 24")
	}
	if !validator.IsAlphanumeric(&u.Username) {
		return errors.New("username: allowed only alphanumeric")
	}
	return nil
}

// validatePassword - validate UserModel password. Used in cmd create user.
func (u *UserModel) validatePassword() error {
	if len(u.Password) < 8 || len(u.Password) > 64 {
		return errors.New("password: min length 8 and max 64")
	}
	if !validator.IsAlphanumericWithSymbols(&u.Password) {
		return errors.New("password: allowed only alphanumeric and some symbols")
	}
	return nil
}
