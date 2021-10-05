package elven

import (
	"database/sql"
	"github.com/jackc/pgx/v4"
	"servus/core"
	"servus/core/modules/cryptor"
	"time"
)


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

func dbUserScanRow(row pgx.Row, user *ModelUser)(err error) {
	err = row.Scan(&user.ID, &user.Role, &user.Username, &user.Password, &user.RegIP, &user.RegAgent, &user.CreatedAt, &user.UpdatedAt)
	err = core.Utils.DBCheckError(err)
	return err
}

func dbUserCreate(user ModelUser) (new ModelUser, err error) {
	new = ModelUser{}
	hashedPassword, err := cryptor.BHash(user.Password)
	if err != nil {
		core.Logger.Error(err.Error())
		return new, err
	}
	var query = `INSERT INTO users (role, username, password) VALUES ($1, $2, $3) RETURNING *`
	row := core.Database.QueryRow(query, user.Role, user.Username, hashedPassword)
	err = dbUserScanRow(row, &new)
	return new, err
}

func dbUserFind(id string) (found *ModelUser, err error){
	found = &ModelUser{}
	var query = "SELECT * FROM users WHERE id=$1 LIMIT 1"
	row := core.Database.QueryRow(query, id)
	err = dbUserScanRow(row, found)
	switch err != nil {
	case true:
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	case false:
		return found, err
	}
	return found, err
}

func dbUserFindBy(username string) (found *ModelUser, err error) {
	found = &ModelUser{}
	var query = "SELECT * FROM users WHERE username=$1 LIMIT 1"
	row := core.Database.QueryRowx(query, username)
	err = row.StructScan(found)
	err = core.Utils.DBCheckError(err)
	switch err != nil {
	case true:
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	case false:
		return found, err
	}
	return
}

func dbUserDeleteBy(id string) error {
	var query = "DELETE FROM users WHERE id=$1"
	_, err := core.Database.Exec(query, id)
	err = core.Utils.DBCheckError(err)
	return err
}
