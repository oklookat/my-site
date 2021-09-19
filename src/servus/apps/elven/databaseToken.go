package elven

import (
	"context"
	"servus/core"
)


func dbTokenCreate(token *modelToken) (modelToken, error) {
	var nToken = modelToken{}
	var sql = "INSERT INTO tokens (user_id, token) VALUES ($1, $2) RETURNING *"
	row := core.Database.Connection.QueryRow(context.Background(), sql, &token.userID, &token.token)
	err := row.Scan(&nToken.id, &nToken.userID, &nToken.token, &nToken.lastIP, &nToken.lastAgent, &nToken.authIP, &nToken.authAgent, &nToken.createdAt, &nToken.updatedAt)
	err = core.Utils.DBCheckError(err)
	return nToken, err
}

func dbTokenUpdate(token *modelToken) (modelToken, error){
	var uToken = modelToken{}
	var sql = "UPDATE tokens SET user_id=$1, token=$2, last_ip=$3, last_agent=$4, auth_ip=$5, auth_agent=$6 WHERE id=$7 RETURNING *"
	row := core.Database.Connection.QueryRow(context.Background(), sql, &token.userID, &token.token, &token.lastIP, &token.lastAgent, &token.authIP, &token.authAgent, &token.id)
	err := row.Scan(&uToken.id, &uToken.userID, &uToken.token, &uToken.lastIP, &uToken.lastAgent, &uToken.authIP, &uToken.authAgent, &uToken.createdAt, &uToken.updatedAt)
	err = core.Utils.DBCheckError(err)
	return uToken, err
}

func dbTokenFind(id string) (modelToken, error){
	var token = modelToken{}
	var sql = "SELECT * FROM tokens WHERE id=$1 LIMIT 1"
	row := core.Database.Connection.QueryRow(context.Background(), sql, id)
	err := row.Scan(&token.id, &token.userID, &token.token, &token.lastIP, &token.lastAgent, &token.authIP, &token.authAgent, &token.createdAt, &token.updatedAt)
	err = core.Utils.DBCheckError(err)
	return token, err
}

func dbTokenDelete(id string) error {
	var sql = "DELETE FROM tokens WHERE id=$1"
	_, err := core.Database.Connection.Exec(context.Background(), sql, id)
	return err
}
