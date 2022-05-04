package cmd

import (
	"errors"
	"servus/apps/elven/user"
	"strconv"
)

type UserOps interface {
	DeleteByUsername(username string) error
	Create(username string, password string, isAdmin bool) error
}

func createUser(isAdmin bool, values []string) (err error) {

	// get args.
	if len(values) < 3 {
		return errors.New("not all args provided")
	}
	var username = values[0]
	var password = values[1]
	var deleteIfExists bool
	deleteIfExists, err = strconv.ParseBool(values[2])
	if err != nil {
		return errors.New("delete_if_exists must be a boolean")
	}

	err = user.Create(username, password, isAdmin)
	if err == nil {
		return
	}

	var isUsernameExists = err.Error() == "user with this username already exists"

	if isUsernameExists {
		if !deleteIfExists {
			// delete if exists disabled, go back.
			return
		}

		// delete old user.
		if err = user.DeleteByUsername(username); err != nil {
			return
		}

		// create new user.
		if err = user.Create(username, password, isAdmin); err != nil {
			return
		}

		return
	}

	return
}
