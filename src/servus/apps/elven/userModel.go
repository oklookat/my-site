package elven

import (
	"database/sql"
	"strings"
	"time"
)

const (
	userRoleAdmin = "admin"
	userRoleUser  = "user"
)

// UserModel - represents user in database.
type UserModel struct {
	ID        string    `json:"id" db:"id"`
	Role      string    `json:"role" db:"role"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"password" db:"password"`
	RegIP     *string   `json:"reg_ip" db:"reg_ip"`
	RegAgent  *string   `json:"reg_agent" db:"reg_agent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// create - create user in database.
func (u *UserModel) create() (err error) {
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

func (u *UserModel) update() (err error) {
	err = u.hookBeforeChange()
	if err != nil {
		return
	}
	var query = "UPDATE users SET role=$1, username=$2, password=$3 WHERE id=$4 RETURNING *"
	err = call.DB.Conn.Get(u, query, u.Role, u.Username, u.Password, u.ID)
	err = call.DB.CheckError(err)
	return
}

// findByID - find user in database by id in UserModel.
func (u *UserModel) findByID() (found bool, err error) {
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

// findByUsername - find user in database by username in UserModel.
func (u *UserModel) findByUsername() (found bool, err error) {
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

// deleteByID - delete user by id in UserModel.
func (u *UserModel) deleteByID() (err error) {
	var query = "DELETE FROM users WHERE id=$1"
	_, err = call.DB.Conn.Exec(query, u.ID)
	err = call.DB.CheckError(err)
	if err == sql.ErrNoRows {
		return nil
	}
	return
}

// hookBeforeChange - change data before send it to DB.
func (u *UserModel) hookBeforeChange() (err error){
	// convert to lower
	u.Username = strings.ToLower(u.Username)
	// check if password not hashed.
	conf, _, _, _ := call.Encryption.Argon.ParseHash(u.Password)
	if conf == nil {
		// not hashed
		hash, err := call.Encryption.Argon.Hash(u.Password)
		if err != nil {
			return err
		}
		u.Password = hash
	}
	return
}