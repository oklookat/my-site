package elven

import (
	"context"
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

func dbUserCreate(user ModelUser) (newUser ModelUser, err error) {
	newUser = ModelUser{}
	hashedPassword, err := cryptor.BHash(user.Password)
	if err != nil {
		core.Logger.Error(err.Error())
		return newUser, err
	}
	var sql = `INSERT INTO users (role, username, password) VALUES ($1, $2, $3) RETURNING *`
	query := core.Database.Connection.QueryRow(context.Background(), sql, user.Role, user.Username, hashedPassword)
	err = query.Scan(&newUser.ID, &newUser.Role, &newUser.Username, &newUser.Password, &newUser.RegIP, &newUser.RegAgent, &newUser.CreatedAt, &newUser.UpdatedAt)
	err = core.Utils.DBCheckError(err)
	return newUser, err
}

func dbUserFind(id string) (ModelUser, error){
	var user = ModelUser{}
	var sql = "SELECT * FROM users WHERE id=$1 LIMIT 1"
	row := core.Database.Connection.QueryRow(context.Background(), sql, id)
	err := row.Scan(&user.ID, &user.Role, &user.Username, &user.Password, &user.RegIP, &user.RegAgent, &user.CreatedAt, &user.UpdatedAt)
	err = core.Utils.DBCheckError(err)
	return user, err
}

func dbUserFindBy(username string) (ModelUser, error) {
	var user = ModelUser{}
	var sql = "SELECT * FROM users WHERE username=$1 LIMIT 1"
	row := core.Database.Connection.QueryRow(context.Background(), sql, username)
	err := row.Scan(&user.ID, &user.Role, &user.Username, &user.Password, &user.RegIP, &user.RegAgent, &user.CreatedAt, &user.UpdatedAt)
	err = core.Utils.DBCheckError(err)
	return user, err
}

func dbUserDeleteBy(id string) error {
	var sql = "DELETE FROM users WHERE id=$1"
	query, err := core.Database.Connection.Query(context.Background(), sql, id)
	defer query.Close()
	err = core.Utils.DBCheckError(err)
	return err
}
