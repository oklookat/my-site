package elven

import (
	"database/sql"
	"time"
)

const (
	userRoleAdmin = "admin"
	userRoleUser  = "user"
)

// ModelUser - represents user in database.
type ModelUser struct {
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
func (u *ModelUser) create() (err error) {
	u.Password, err = instance.Encryption.Argon.Hash(u.Password)
	if err != nil {
		instance.Logger.Error(err.Error())
		return
	}
	var query = `INSERT INTO users (role, username, password) VALUES ($1, $2, $3) RETURNING *`
	err = instance.DB.Conn.Get(u, query, u.Role, u.Username, u.Password)
	err = instance.DB.CheckError(err)
	return
}

// findByID - find user in database by id in ModelUser.
func (u *ModelUser) findByID() (found bool, err error) {
	var query = "SELECT * FROM users WHERE id=$1 LIMIT 1"
	err = instance.DB.Conn.Get(u, query, u.ID)
	err = instance.DB.CheckError(err)
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

// findByUsername - find user in database by username in ModelUser.
func (u *ModelUser) findByUsername() (found bool, err error) {
	var query = "SELECT * FROM users WHERE username=$1 LIMIT 1"
	err = instance.DB.Conn.Get(u, query, u.Username)
	err = instance.DB.CheckError(err)
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

// deleteByID - delete user by id in ModelUser.
func (u *ModelUser) deleteByID() (err error) {
	var query = "DELETE FROM users WHERE id=$1"
	_, err = instance.DB.Conn.Exec(query, u.ID)
	err = instance.DB.CheckError(err)
	if err == sql.ErrNoRows {
		return nil
	}
	return
}
