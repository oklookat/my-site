package elven

import (
	"database/sql"
	"servus/core"
	"servus/core/modules/cryptor"
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
	hashedPassword, err := cryptor.BHash(u.Password)
	if err != nil {
		core.Logger.Error(err.Error())
		return
	}
	var query = `INSERT INTO users (role, username, password) VALUES ($1, $2, $3) RETURNING *`
	row := core.Database.QueryRowx(query, u.Role, u.Username, hashedPassword)
	err = row.StructScan(u)
	err = core.Utils.DBCheckError(err)
	if err == sql.ErrNoRows {
		return nil
	}
	return
}

// findByID - find user in database by id in ModelUser.
func (u *ModelUser) findByID() (found bool, err error) {
	var query = "SELECT * FROM users WHERE id=$1 LIMIT 1"
	row := core.Database.QueryRowx(query, u.ID)
	err = row.StructScan(u)
	err = core.Utils.DBCheckError(err)
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
	row := core.Database.QueryRowx(query, u.Username)
	err = row.StructScan(u)
	err = core.Utils.DBCheckError(err)
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
	_, err = core.Database.Exec(query, u.ID)
	err = core.Utils.DBCheckError(err)
	if err == sql.ErrNoRows {
		return nil
	}
	return
}
