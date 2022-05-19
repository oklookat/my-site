package user

import (
	"servus/core/external/database"
	"strings"
	"time"

	"errors"
)

var (
	ErrUserUsernameMinMax       = errors.New("validation: username min length 4, max 24")
	ErrUserUsernameAlphanumeric = errors.New("validation: username must be alphanumeric")
	ErrUserPasswordMinMax       = errors.New("validation: password min length 8, max 64")
	ErrUserPasswordWrongSymbols = errors.New("validation: wrong symbols used in password")
)

var userAdapter = database.Adapter[Model]{}

// represents user in database.
type Model struct {
	ID        string    `json:"id" db:"id"`
	Role      string    `json:"role" db:"role"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"password" db:"password"`
	RegIP     *string   `json:"reg_ip" db:"reg_ip"`
	RegAgent  *string   `json:"reg_agent" db:"reg_agent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// create user in database.
func (m *Model) Create() (err error) {
	if err = m.hookBeforeChange(); err != nil {
		call.Logger.Error(err.Error())
		return
	}
	var query = `INSERT INTO users (role, username, password) VALUES ($1, $2, $3) RETURNING *`
	err = userAdapter.Get(m, query, m.Role, m.Username, m.Password)
	return
}

func (m *Model) Update() (err error) {
	if err = m.hookBeforeChange(); err != nil {
		return
	}
	var query = "UPDATE users SET role=$1, username=$2, password=$3 WHERE id=$4 RETURNING *"
	err = userAdapter.Get(m, query, m.Role, m.Username, m.Password, m.ID)
	return
}

// find user in database by id in UserModel.
func (m *Model) FindByID() (found bool, err error) {
	found = false
	var query = "SELECT * FROM users WHERE id=$1 LIMIT 1"
	founded, err := userAdapter.Find(query, m.ID)
	if err != nil {
		return
	}
	if founded != nil {
		found = true
		*m = *founded
	}
	return
}

// find user in database by username in UserModel.
func (m *Model) FindByUsername() (found bool, err error) {
	found = false
	var query = "SELECT * FROM users WHERE username=$1 LIMIT 1"
	founded, err := userAdapter.Find(query, m.Username)
	if err != nil {
		return
	}
	if founded != nil {
		found = true
		*m = *founded
	}
	return
}

// delete user by id in UserModel.
func (m *Model) DeleteByID() (err error) {
	var query = "DELETE FROM users WHERE id=$1"
	_, err = userAdapter.Exec(query, m.ID)
	return
}

// change data before send it to DB.
func (m *Model) hookBeforeChange() (err error) {
	// convert to lower
	m.Username = strings.ToLower(m.Username)
	// check if password not hashed.
	var isHashed = call.Encryptor.Argon.IsHash(m.Password)
	if !isHashed {
		hash, err := call.Encryptor.Argon.Hash(m.Password)
		if err != nil {
			return err
		}
		m.Password = hash
	}
	return
}
