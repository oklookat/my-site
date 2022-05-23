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

	// find by username.
	var usr = user.Model{
		Username: username,
	}
	var isExists bool
	isExists, err = usr.FindByUsername()
	if err != nil {
		return
	}

	// delete/not delete if exists.
	if isExists {
		if !deleteIfExists {
			return
		}
		err = usr.DeleteByID()
	}

	// create.
	err = user.Create(username, password, isAdmin)
	if err == nil {
		return
	}

	return
}
