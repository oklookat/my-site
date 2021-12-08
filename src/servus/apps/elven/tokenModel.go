package elven

import (
	"database/sql"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

// TokenModel - represents token in database.
type TokenModel struct {
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

// create - create TokenModel in database.
func (t *TokenModel) create() (err error) {
	var query = "INSERT INTO tokens (user_id, token, last_ip, last_agent, auth_ip, auth_agent) VALUES (:user_id, :token, :last_ip, :last_agent, :auth_ip, :auth_agent) RETURNING *"
	stmt, err := call.DB.Conn.PrepareNamed(query)
	if err != nil {
		return
	}
	defer func() {
		_ = stmt.Close()
	}()
	err = stmt.Get(t, t)
	err = call.DB.CheckError(err)
	return
}

// update - updates TokenModel in database. All fields except update and created dates must be filled.
func (t *TokenModel) update() (err error) {
	t.hookBeforeUpdate()
	var query = "UPDATE tokens SET user_id=:user_id, token=:token, last_ip=:last_ip, last_agent=:last_agent, auth_ip=:auth_ip, auth_agent=:auth_agent WHERE id=:id RETURNING *"
	stmt, err := call.DB.Conn.PrepareNamed(query)
	if err != nil {
		return
	}
	defer func() {
		_ = stmt.Close()
	}()
	err = stmt.Get(t, t)
	err = call.DB.CheckError(err)
	return
}

// databaseFind - find TokenModel in database by id field.
func (t *TokenModel) findByID() (found bool, err error) {
	var query = "SELECT * FROM tokens WHERE id=$1 LIMIT 1"
	err = call.DB.Conn.Get(t, query, t.ID)
	err = call.DB.CheckError(err)
	found = false
	if err != nil {
		if err == sql.ErrNoRows {
			return found, nil
		}
		return
	}
	found = true
	return
}

// deleteByID - delete TokenModel from database by id field.
func (t *TokenModel) deleteByID() (err error) {
	var query = "DELETE FROM tokens WHERE id=$1"
	_, err = call.DB.Conn.Exec(query, t.ID)
	if err == sql.ErrNoRows {
		return nil
	}
	return
}

// hookBeforeUpdate - executes before token update.
func (t *TokenModel) hookBeforeUpdate() {
	if t.AuthAgent != nil && len(*t.AuthAgent) > 323 {
		var authAgent = *t.AuthAgent
		var cut = 323 - len(authAgent)
		*t.AuthAgent = authAgent[:len(authAgent)-cut]
	}
	if t.LastAgent != nil && len(*t.LastAgent) > 323 {
		var lastAgent = *t.LastAgent
		var cut = 323 - len(lastAgent)
		*t.LastAgent = lastAgent[:len(lastAgent)-cut]
	}
	if t.AuthIP != nil && len(*t.AuthIP) > 53 {
		var authIP = *t.AuthIP
		var cut = 53 - len(authIP)
		*t.AuthIP = authIP[:len(authIP)-cut]
	}
	if t.LastIP != nil && len(*t.LastIP) > 53 {
		var lastIP = *t.LastIP
		var cut = 53 - len(lastIP)
		*t.LastIP = lastIP[:len(lastIP)-cut]
	}
}

// setAuthAgents - writes last ip and user agent then updating model in database.
func (t *TokenModel) setAuthAgents(request *http.Request) (err error) {
	if request == nil {
		return errors.New("setLastAgents: request nil pointer.")
	}
	if t == nil {
		return errors.New("setLastAgents: token nil pointer.")
	}
	t.AuthAgent = new(string)
	*t.AuthAgent = request.UserAgent()
	t.AuthIP = new(string)
	*t.AuthIP = getIP(request)
	err = t.update()
	return
}

// setLastAgents - writes ip and user agent then updating model in database.
func (t *TokenModel) setLastAgents(request *http.Request) (err error) {
	if request == nil {
		return errors.New("setLastAgents: request nil pointer.")
	}
	if t == nil {
		return errors.New("setLastAgents: token nil pointer.")
	}
	var lastAgent = request.UserAgent()
	var lastIP = getIP(request)
	t.LastAgent = new(string)
	*t.LastAgent = lastAgent
	t.LastIP = new(string)
	*t.LastIP = lastIP
	err = t.update()
	return
}

// generate - generate token.
// returns: token - token for user, hash - saved in db as TokenModel.Token.
func (t *TokenModel) generate(userID string) (err error, token string, hash string) {
	// token generating.
	// first we generate fake token model to get created token ID.
	t.UserID = userID
	t.Token = "-1"
	err = t.create()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = t.deleteByID()
		}
	}()
	// then we get newly created token model id and encrypt it. That's we send to user as token.
	token, err = call.Encryption.AES.Encrypt(t.ID)
	if err != nil {
		return
	}
	// get hash from generated token.
	// user gets encrypted token, but database gets hash. In general, we do the same as with the password.
	hash, err = call.Encryption.Argon.Hash(token)
	if err != nil {
		return
	}
	// now we replace fake token with real token in database.
	t.Token = hash
	err = t.update()
	return
}
