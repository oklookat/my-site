package elven

import (
	"database/sql"
	"servus/core"
	"servus/core/modules/cryptor"
	"time"
)

// user - user entity.
type entityUser struct {
}

// ModelUser - represents user in database.
type ModelUser struct {
	ID        string `json:"id" db:"id"`
	Role      string `json:"role" db:"role"`
	Username  string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
	RegIP     *string `json:"reg_ip" db:"reg_ip"`
	RegAgent  *string `json:"reg_agent" db:"reg_agent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (u *entityUser) databaseCreate(user *ModelUser) (err error) {
	hashedPassword, err := cryptor.BHash(user.Password)
	if err != nil {
		core.Logger.Error(err.Error())
		return
	}
	var query = `INSERT INTO users (role, username, password) VALUES ($1, $2, $3) RETURNING *`
	row := core.Database.QueryRowx(query, user.Role, user.Username, hashedPassword)
	err = row.StructScan(user)
	err = core.Utils.DBCheckError(err)
	if err == sql.ErrNoRows {
		return nil
	}
	return
}

func (u *entityUser) databaseFind(id string) (found *ModelUser, err error) {
	found = &ModelUser{}
	var query = "SELECT * FROM users WHERE id=$1 LIMIT 1"
	row := core.Database.QueryRowx(query, id)
	err = row.StructScan(found)
	err = core.Utils.DBCheckError(err)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

func (u *entityUser) databaseFindBy(username string) (found *ModelUser, err error){
	found = &ModelUser{}
	var query = "SELECT * FROM users WHERE username=$1 LIMIT 1"
	row := core.Database.QueryRowx(query, username)
	err = row.StructScan(found)
	err = core.Utils.DBCheckError(err)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

func (u *entityUser) databaseDelete(id string) (err error) {
	var query = "DELETE FROM users WHERE id=$1"
	_, err = core.Database.Exec(query, id)
	err = core.Utils.DBCheckError(err)
	if err == sql.ErrNoRows {
		return nil
	}
	return
}
