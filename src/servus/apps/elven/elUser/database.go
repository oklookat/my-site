package elUser

import (
	"context"
	"github.com/pkg/errors"
	"servus/core"
	"servus/core/modules/cryptor"
	"time"
)

type modelUser struct {
	id        string
	role      string
	username  string
	password  string
	regIp     string
	regAgent  string
	createdAt time.Time
	updatedAt time.Time
}

type modelToken struct {
	id        string
	userID    string
	token     string
	lastIP    string
	authIP    string
	authAgent string
	createdAt time.Time
	updatedAt time.Time
}

func dbCreateUser(user modelUser) error {
	var hashedPassword, _ = cryptor.BHash(user.password)
	query, err := core.Database.Connection.Query(context.Background(), "INSERT INTO users (role, username, password) VALUES ($1, $2, $3)", user.role, user.username, hashedPassword)
	defer query.Close()
	return dbCheckErr(err)
}

func dbFindUserBy(username string) (modelUser, error) {
	var user = modelUser{}
	row := core.Database.Connection.QueryRow(context.Background(), "SELECT * FROM users WHERE username=$1 LIMIT 1", username)
	err := row.Scan(&user.id, &user.role, &user.username, &user.password, &user.regIp, &user.regAgent, &user.createdAt, &user.updatedAt)
	if err != nil {
		switch err.Error() {
		case "no rows in result set":
			return user, errors.New("PIPE_USER_NOT_FOUND")
		default:
			core.Logger.Error(err.Error())
			return user, err
		}
	}
	return user, nil
}

func dbDeleteUserBy(id string) error {
	query, err := core.Database.Connection.Query(context.Background(), "DELETE FROM users WHERE id=$1", id)
	defer query.Close()
	return dbCheckErr(err)
}

func dbCheckErr(err error) error {
	if err != nil {
		core.Logger.Error(err.Error())
		return err
	}
	return nil
}


func dbCreateToken(token modelToken) error {
	query, err := core.Database.Connection.Query(context.Background(), "INSERT INTO tokens (user_id, token) VALUES ($1, $2)", token.userID, token.token)
	defer query.Close()
	return dbCheckErr(err)
}
