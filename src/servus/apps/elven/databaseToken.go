package elven

import (
	"context"
	"servus/core"
	"time"
)

type ModelToken struct {
	ID        string
	UserID    string
	Token     string
	LastIP    *string
	LastAgent *string
	AuthIP    *string
	AuthAgent *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// dbModifyToken - uses for modify data in dbToken methods. Ex: cut long user agent before make update request to database. You don't need this function.
func dbModifyToken(token ModelToken) ModelToken {
	if token.AuthAgent != nil && len(*token.AuthAgent) > 323 {
		var authAgent = *token.AuthAgent
		var cut = 323 - len(authAgent)
		*token.AuthAgent = authAgent[:len(authAgent)-cut]
	}
	if token.LastAgent != nil && len(*token.LastAgent) > 323 {
		var lastAgent = *token.LastAgent
		var cut = 323 - len(lastAgent)
		*token.LastAgent = lastAgent[:len(lastAgent)-cut]
	}
	if token.AuthIP != nil && len(*token.AuthIP) > 53 {
		var authIP = *token.AuthIP
		var cut = 53 - len(authIP)
		*token.AuthIP = authIP[:len(authIP)-cut]
	}
	if token.LastIP != nil && len(*token.LastIP) > 53 {
		var lastIP = *token.LastIP
		var cut = 53 - len(lastIP)
		*token.LastIP = lastIP[:len(lastIP)-cut]
	}
	return token
}

// dbTokenCreate - create Token in database. ATTENTION: its function writes only Token and user id, other data will be ignored. For write full data see dbTokenUpdate.
func dbTokenCreate(token *ModelToken) (created ModelToken, err error) {
	created = ModelToken{}
	var sql = "INSERT INTO tokens (user_id, token) VALUES ($1, $2) RETURNING *"
	row := core.Database.Connection.QueryRow(context.Background(), sql, &token.UserID, &token.Token)
	err = row.Scan(&created.ID, &created.UserID, &created.Token, &created.LastIP, &created.LastAgent, &created.AuthIP, &created.AuthAgent, &created.CreatedAt, &created.UpdatedAt)
	err = core.Utils.DBCheckError(err)
	return created, err
}

// dbTokenUpdate - updates Token in database. All fields (except update and created dates) must be filled.
func dbTokenUpdate(token *ModelToken) (updated ModelToken, err error) {
	updated = ModelToken{}
	*token = dbModifyToken(*token)
	var sql = "UPDATE tokens SET user_id=$1, token=$2, last_ip=$3, last_agent=$4, auth_ip=$5, auth_agent=$6 WHERE id=$7 RETURNING *"
	row := core.Database.Connection.QueryRow(context.Background(), sql, &token.UserID, &token.Token, &token.LastIP, &token.LastAgent, &token.AuthIP, &token.AuthAgent, &token.ID)
	err = row.Scan(&updated.ID, &updated.UserID, &updated.Token, &updated.LastIP, &updated.LastAgent, &updated.AuthIP, &updated.AuthAgent, &updated.CreatedAt, &updated.UpdatedAt)
	err = core.Utils.DBCheckError(err)
	return updated, err
}

// dbTokenFind - find Token in database.
func dbTokenFind(id string) (found ModelToken, err error) {
	found = ModelToken{}
	var sql = "SELECT * FROM tokens WHERE id=$1 LIMIT 1"
	row := core.Database.Connection.QueryRow(context.Background(), sql, id)
	err = row.Scan(&found.ID, &found.UserID, &found.Token, &found.LastIP, &found.LastAgent, &found.AuthIP, &found.AuthAgent, &found.CreatedAt, &found.UpdatedAt)
	err = core.Utils.DBCheckError(err)
	return found, err
}

// dbTokenDelete - delete Token from database.
func dbTokenDelete(id string) error {
	var sql = "DELETE FROM tokens WHERE id=$1"
	_, err := core.Database.Connection.Exec(context.Background(), sql, id)
	return err
}
