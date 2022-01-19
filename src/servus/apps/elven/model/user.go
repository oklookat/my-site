package model

import (
	"servus/core/external/database"
	"servus/core/external/validator"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var ErrUserUsernameMinMax = errors.New("validation: username min length 4, max 24")
var ErrUserUsernameAlphanumeric = errors.New("validation: username must be alphanumeric")
var ErrUserPasswordMinMax = errors.New("validation: password min length 8, max 64")
var ErrUserPasswordWrongSymbols = errors.New("validation: wrong symbols used in password")

// represents user in database.
type User struct {
	ID        string    `json:"id" db:"id"`
	Role      string    `json:"role" db:"role"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"password" db:"password"`
	RegIP     *string   `json:"reg_ip" db:"reg_ip"`
	RegAgent  *string   `json:"reg_agent" db:"reg_agent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

var userAdapter = database.Adapter[User]{}

// create user in database.
func (u *User) Create() (err error) {
	err = u.hookBeforeChange()
	if err != nil {
		call.Logger.Error(err.Error())
		return
	}
	var query = `INSERT INTO users (role, username, password) VALUES ($1, $2, $3) RETURNING *`
	err = userAdapter.Get(u, query, u.Role, u.Username, u.Password)
	return
}

func (u *User) Update() (err error) {
	err = u.hookBeforeChange()
	if err != nil {
		return
	}
	var query = "UPDATE users SET role=$1, username=$2, password=$3 WHERE id=$4 RETURNING *"
	err = userAdapter.Get(u, query, u.Role, u.Username, u.Password, u.ID)
	return
}

// find user in database by id in UserModel.
func (u *User) FindByID() (found bool, err error) {
	found = false
	var query = "SELECT * FROM users WHERE id=$1 LIMIT 1"
	founded, err := userAdapter.Find(query, u.ID)
	if err != nil {
		return
	}
	if founded != nil {
		found = true
		*u = *founded
	}
	return
}

// find user in database by username in UserModel.
func (u *User) FindByUsername() (found bool, err error) {
	found = false
	var query = "SELECT * FROM users WHERE username=$1 LIMIT 1"
	founded, err := userAdapter.Find(query, u.Username)
	if err != nil {
		return
	}
	if founded != nil {
		found = true
		*u = *founded
	}
	return
}

// delete user by id in UserModel.
func (u *User) DeleteByID() (err error) {
	var query = "DELETE FROM users WHERE id=$1"
	_, err = userAdapter.Exec(query, u.ID)
	return
}

// validate username from UserModel. Used in cmd create user.
func (u *User) ValidateUsername() error {
	if validator.MinMax(&u.Username, 4, 24) {
		return ErrUserUsernameMinMax
	}
	if !validator.IsAlphanumeric(&u.Username) {
		return ErrUserUsernameAlphanumeric
	}
	return nil
}

// validate UserModel password. Used in cmd create user.
func (u *User) ValidatePassword() error {
	if validator.MinMax(&u.Password, 8, 64) {
		return ErrUserPasswordMinMax
	}
	if !validator.IsAlphanumericWithSymbols(&u.Password) {
		return ErrUserPasswordWrongSymbols
	}
	return nil
}

// change data before send it to DB.
func (u *User) hookBeforeChange() (err error) {
	// convert to lower
	u.Username = strings.ToLower(u.Username)
	// validate.
	err = u.ValidateUsername()
	if err != nil {
		return
	}
	err = u.ValidatePassword()
	if err != nil {
		return
	}
	// check if password not hashed.
	var isHashed = call.Encryptor.Argon.IsHash(u.Password)
	if !isHashed {
		hash, err := call.Encryptor.Argon.Hash(u.Password)
		if err != nil {
			return err
		}
		u.Password = hash
	}
	return
}
