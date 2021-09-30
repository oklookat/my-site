package elven

import (
	"context"
	"servus/core"
)

// dbModifyToken - uses for modify data in dbToken methods. Ex: cut long user agent before make update request to database. You don't need this function.
func dbModifyToken(token modelToken) modelToken {
	if len(token.authAgent) > 323 {
		var cut = 323 - len(token.authAgent)
		token.authAgent = token.authAgent[:len(token.authAgent)-cut]
	}
	if len(token.lastAgent) > 323 {
		var cut = 323 - len(token.lastAgent)
		token.lastAgent = token.lastAgent[:len(token.lastAgent)-cut]
	}
	if len(token.authIP) > 53 {
		var cut = 53 - len(token.authIP)
		token.authIP = token.authIP[:len(token.authIP)-cut]
	}
	if len(token.lastIP) > 53 {
		var cut = 53 - len(token.lastIP)
		token.lastIP = token.lastIP[:len(token.lastIP)-cut]
	}
	return token
}

// dbTokenCreate - create token in database. ATTENTION: its function writes only token and user id, other data will be ignored. For write full data see dbTokenUpdate.
func dbTokenCreate(token *modelToken) (created modelToken, err error) {
	created = modelToken{}
	var sql = "INSERT INTO tokens (user_id, token) VALUES ($1, $2) RETURNING *"
	row := core.Database.Connection.QueryRow(context.Background(), sql, &token.userID, &token.token)
	err = row.Scan(&created.id, &created.userID, &created.token, &created.lastIP, &created.lastAgent, &created.authIP, &created.authAgent, &created.createdAt, &created.updatedAt)
	err = core.Utils.DBCheckError(err)
	return created, err
}

// dbTokenUpdate - updates token in database. All fields (except update and created dates) must be filled.
func dbTokenUpdate(token *modelToken) (updated modelToken, err error) {
	updated = modelToken{}
	*token = dbModifyToken(*token)
	var sql = "UPDATE tokens SET user_id=$1, token=$2, last_ip=$3, last_agent=$4, auth_ip=$5, auth_agent=$6 WHERE id=$7 RETURNING *"
	row := core.Database.Connection.QueryRow(context.Background(), sql, &token.userID, &token.token, &token.lastIP, &token.lastAgent, &token.authIP, &token.authAgent, &token.id)
	err = row.Scan(&updated.id, &updated.userID, &updated.token, &updated.lastIP, &updated.lastAgent, &updated.authIP, &updated.authAgent, &updated.createdAt, &updated.updatedAt)
	err = core.Utils.DBCheckError(err)
	return updated, err
}

// dbTokenFind - find token in database.
func dbTokenFind(id string) (found modelToken, err error) {
	found = modelToken{}
	var sql = "SELECT * FROM tokens WHERE id=$1 LIMIT 1"
	row := core.Database.Connection.QueryRow(context.Background(), sql, id)
	err = row.Scan(&found.id, &found.userID, &found.token, &found.lastIP, &found.lastAgent, &found.authIP, &found.authAgent, &found.createdAt, &found.updatedAt)
	err = core.Utils.DBCheckError(err)
	return found, err
}

// dbTokenDelete - delete token from database.
func dbTokenDelete(id string) error {
	var sql = "DELETE FROM tokens WHERE id=$1"
	_, err := core.Database.Connection.Exec(context.Background(), sql, id)
	return err
}
