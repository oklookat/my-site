package elven

import (
	"context"
	"servus/core"
	"servus/core/modules/cryptor"
)


func dbUserCreate(user modelUser) (modelUser, error) {
	var newUser = modelUser{}
	var hashedPassword, err = cryptor.BHash(user.password)
	if err != nil {
		core.Logger.Error(err.Error())
		return newUser, err
	}
	var sql = `INSERT INTO users (role, username, password) VALUES ($1, $2, $3) RETURNING id, role, username, password, reg_ip, reg_agent, created_at, updated_at`
	query := core.Database.Connection.QueryRow(context.Background(), sql, user.role, user.username, hashedPassword)
	err = query.Scan(&newUser.id, &newUser.role, &newUser.username, &newUser.password, &newUser.regIP, &newUser.regAgent, &newUser.createdAt, &newUser.updatedAt)
	err = core.Utils.DBCheckError(err)
	return newUser, err
}

func dbUserFind(id string) (modelUser, error){
	var user = modelUser{}
	var sql = "SELECT * FROM users WHERE id=$1 LIMIT 1"
	row := core.Database.Connection.QueryRow(context.Background(), sql, id)
	err := row.Scan(&user.id, &user.role, &user.username, &user.password, &user.regIP, &user.regAgent, &user.createdAt, &user.updatedAt)
	err = core.Utils.DBCheckError(err)
	return user, err
}

func dbUserFindBy(username string) (modelUser, error) {
	var user = modelUser{}
	var sql = "SELECT * FROM users WHERE username=$1 LIMIT 1"
	row := core.Database.Connection.QueryRow(context.Background(), sql, username)
	err := row.Scan(&user.id, &user.role, &user.username, &user.password, &user.regIP, &user.regAgent, &user.createdAt, &user.updatedAt)
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
