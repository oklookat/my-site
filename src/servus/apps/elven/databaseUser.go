package elven

import (
	"github.com/jackc/pgx/v4"
	"servus/core"
	"servus/core/modules/cryptor"
	"time"
)


type ModelUser struct {
	ID        string
	Role      string
	Username  string
	Password  string
	RegIP     *string
	RegAgent  *string
	CreatedAt time.Time
	UpdatedAt time.Time
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
	var sql = `INSERT INTO users (role, username, password) VALUES ($1, $2, $3) RETURNING *`
	row := core.Database.Connection.QueryRow(sql, user.Role, user.Username, hashedPassword)
	err = dbUserScanRow(row, &new)
	return new, err
}

func dbUserFind(id string) (found ModelUser, err error){
	found = ModelUser{}
	var sql = "SELECT * FROM users WHERE id=$1 LIMIT 1"
	row := core.Database.Connection.QueryRow(sql, id)
	err = dbUserScanRow(row, &found)
	return found, err
}

func dbUserFindBy(username string) (found ModelUser, err error) {
	found = ModelUser{}
	var sql = "SELECT * FROM users WHERE username=$1 LIMIT 1"
	row := core.Database.Connection.QueryRow(sql, username)
	err = dbUserScanRow(row, &found)
	return found, err
}

func dbUserDeleteBy(id string) error {
	var sql = "DELETE FROM users WHERE id=$1"
	_, err := core.Database.Connection.Exec(sql, id)
	err = core.Utils.DBCheckError(err)
	return err
}
