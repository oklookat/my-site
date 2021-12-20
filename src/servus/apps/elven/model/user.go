package model

import (
	"database/sql"
	"servus/core/external/validator"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	userRoleAdmin = "admin"
	userRoleUser  = "user"
)

var ErrUserUsernameValidation = errors.New("username validation failed")
var ErrUserPasswordValidation = errors.New("password validation failed")

// User - represents user in database.
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

// Create - create user in database.
func (u *User) Create() (err error) {
	err = u.hookBeforeChange()
	if err != nil {
		call.Logger.Error(err.Error())
		return
	}
	var query = `INSERT INTO users (role, username, password) VALUES ($1, $2, $3) RETURNING *`
	err = call.DB.Conn.Get(u, query, u.Role, u.Username, u.Password)
	err = call.DB.CheckError(err)
	return
}

func (u *User) Update() (err error) {
	err = u.hookBeforeChange()
	if err != nil {
		return
	}
	var query = "UPDATE users SET role=$1, username=$2, password=$3 WHERE id=$4 RETURNING *"
	err = call.DB.Conn.Get(u, query, u.Role, u.Username, u.Password, u.ID)
	err = call.DB.CheckError(err)
	return
}

// FindByID - find user in database by id in UserModel.
func (u *User) FindByID() (found bool, err error) {
	var query = "SELECT * FROM users WHERE id=$1 LIMIT 1"
	err = call.DB.Conn.Get(u, query, u.ID)
	err = call.DB.CheckError(err)
	found = false
	if err != nil {
		if err == sql.ErrNoRows {
			return found, nil
		}
		return
	}
	found = true
	return
}

// FindByUsername - find user in database by username in UserModel.
func (u *User) FindByUsername() (found bool, err error) {
	var query = "SELECT * FROM users WHERE username=$1 LIMIT 1"
	err = call.DB.Conn.Get(u, query, u.Username)
	err = call.DB.CheckError(err)
	found = false
	if err != nil {
		if err == sql.ErrNoRows {
			return found, nil
		}
		return
	}
	found = true
	return
}

// DeleteByID - delete user by id in UserModel.
func (u *User) DeleteByID() (err error) {
	var query = "DELETE FROM users WHERE id=$1"
	_, err = call.DB.Conn.Exec(query, u.ID)
	err = call.DB.CheckError(err)
	if err == sql.ErrNoRows {
		return nil
	}
	return
}

// validateUsername - validate username from UserModel. Used in cmd create user.
func (u *User) ValidateUsername() error {
	if validator.MinMax(&u.Username, 4, 24) {
		return ErrUserUsernameValidation
	}
	if !validator.IsAlphanumeric(&u.Username) {
		return ErrUserUsernameValidation
	}
	return nil
}

// validatePassword - validate UserModel password. Used in cmd create user.
func (u *User) ValidatePassword() error {
	if len(u.Password) < 8 || len(u.Password) > 64 {
		return ErrUserPasswordValidation
	}
	if !validator.IsAlphanumericWithSymbols(&u.Password) {
		return ErrUserPasswordValidation
	}
	return nil
}

// hookBeforeChange - change data before send it to DB.
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
