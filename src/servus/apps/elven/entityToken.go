package elven

import (
	"database/sql"
	"servus/core"
	"time"
)

// entityToken - manage token.
type entityToken struct {
}

// ModelToken - represents token in database.
type ModelToken struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	LastIP    *string   `json:"last_ip" db:"last_ip"`
	LastAgent *string   `json:"last_agent" db:"last_agent"`
	AuthIP    *string   `json:"auth_ip" db:"auth_ip"`
	AuthAgent *string   `json:"auth_agent" db:"auth_agent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// databaseCreate - create ModelToken in database. ATTENTION: its function writes only Token and user id, other data will be ignored. For write full data see dbTokenUpdate.
func (t *entityToken) databaseCreate(token *ModelToken) (err error) {
	var query = "INSERT INTO tokens (user_id, token) VALUES ($1, $2) RETURNING *"
	row := core.Database.QueryRowx(query, &token.UserID, &token.Token)
	err = row.StructScan(token)
	err = core.Utils.DBCheckError(err)
	return
}

// databaseUpdate - updates ModelToken in database. All fields (except update and created dates) must be filled.
func (t *entityToken) databaseUpdate(token *ModelToken) (err error) {
	t.databaseBeforeUpdate(token)
	var query = "UPDATE tokens SET user_id=$1, token=$2, last_ip=$3, last_agent=$4, auth_ip=$5, auth_agent=$6 WHERE id=$7 RETURNING *"
	row := core.Database.QueryRowx(query, &token.UserID, &token.Token, &token.LastIP, &token.LastAgent, &token.AuthIP, &token.AuthAgent, &token.ID)
	err = row.StructScan(token)
	err = core.Utils.DBCheckError(err)
	return
}

// databaseFind - find ModelToken in database.
func (t *entityToken) databaseFind(id string) (found *ModelToken, err error) {
	found = &ModelToken{}
	var query = "SELECT * FROM tokens WHERE id=$1 LIMIT 1"
	row := core.Database.QueryRowx(query, id)
	err = row.StructScan(found)
	err = core.Utils.DBCheckError(err)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

// databaseDelete - delete ModelToken from database.
func (t *entityToken) databaseDelete(id string) (err error) {
	var query = "DELETE FROM tokens WHERE id=$1"
	_, err = core.Database.Exec(query, id)
	if err == sql.ErrNoRows {
		return nil
	}
	return
}

// databaseBeforeUpdate - executes before token update.
func (t *entityToken) databaseBeforeUpdate(token *ModelToken) {
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
}
