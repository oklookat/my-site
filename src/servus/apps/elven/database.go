package elven

import (
	"context"
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

func dbCreateUser(user modelUser) error{
	_, err := servus.DB.Connection.Query(context.Background(), "INSERT INTO users (role, username, password) VALUES ($1, $2, $3)", user.role, user.username, user.password)
	if err != nil {
		servus.Logger.Error(err.Error())
		return err
	}
	return nil
}

func dbSearchUserBy(username string) (modelUser, error){
	// TODO: fix conn busy err when user recreating
	var user = modelUser{}
	row := servus.DB.Connection.QueryRow(context.Background(), "SELECT * FROM users WHERE username=$1 LIMIT 1", username)
	err := row.Scan(&user.id, &user.role, &user.username, &user.password, &user.regIp, &user.regAgent, &user.createdAt, &user.updatedAt)
	if err != nil {
		switch err.Error() {
		case "no rows in result set":
			break
		default:
			servus.Logger.Error(err.Error())
			return user, err
		}
	}
	return user, nil
}

func dbDeleteUserBy(id string) error{
	_, err := servus.DB.Connection.Query(context.Background(), "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		servus.Logger.Error(err.Error())
		return err
	}
	return nil
}
