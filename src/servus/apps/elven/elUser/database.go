package elUser

import (
	"context"
	"github.com/pkg/errors"
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

func dbCreateUser(user modelUser) error {
	var hashedPassword, _ = servus.Utils.HashPassword(user.password)
	query, err := servus.DB.Connection.Query(context.Background(), "INSERT INTO users (role, username, password) VALUES ($1, $2, $3)", user.role, user.username, hashedPassword)
	defer query.Close()
	return dbCheckErr(err)
}

func dbFindUserBy(username string) (modelUser, error) {
	var user = modelUser{}
	row := servus.DB.Connection.QueryRow(context.Background(), "SELECT * FROM users WHERE username=$1 LIMIT 1", username)
	err := row.Scan(&user.id, &user.role, &user.username, &user.password, &user.regIp, &user.regAgent, &user.createdAt, &user.updatedAt)
	if err != nil {
		switch err.Error() {
		case "no rows in result set":
			return user, errors.New("PIPE_USER_NOT_FOUND")
		default:
			servus.Logger.Error(err.Error())
			return user, err
		}
	}
	return user, nil
}

func dbDeleteUserBy(id string) error {
	query, err := servus.DB.Connection.Query(context.Background(), "DELETE FROM users WHERE id=$1", id)
	defer query.Close()
	return dbCheckErr(err)
}

func dbCheckErr(err error) error {
	if err != nil {
		servus.Logger.Error(err.Error())
		return err
	}
	return nil
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

func dbCreateToken(token modelToken) error {
	query, err := servus.DB.Connection.Query(context.Background(), "INSERT INTO tokens (user_id, token) VALUES ($1, $2)", token.userID, token.token)
	defer query.Close()
	return dbCheckErr(err)
}
