package token

import (
	"net/http"
	"servus/core/external/database"
	"servus/core/external/utils"
	"time"
)

var tokenAdapter = database.Adapter[Model]{}

// represents token in database.
type Model struct {
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

// create Model in database.
func (m *Model) Create() (err error) {
	var query = `INSERT INTO tokens (user_id, token, last_ip, 
		last_agent, auth_ip, auth_agent) VALUES ($1, $2, $3, 
			$4, $5, $6) RETURNING *`
	err = tokenAdapter.Get(m, query, m.UserID, m.Token, m.LastIP, m.LastAgent, m.AuthIP, m.AuthAgent)
	return
}

// updates TokenModel in database. All fields except update and created dates must be filled.
func (m *Model) Update() (err error) {
	m.hookBeforeUpdate()
	var query = `UPDATE tokens SET user_id=$1, token=$2, last_ip=$3, 
	last_agent=$4, auth_ip=$5, auth_agent=$6 WHERE id=$7 RETURNING *`
	err = tokenAdapter.Get(m, query, m.UserID, m.Token, m.LastIP, m.LastAgent, m.AuthIP, m.AuthAgent, m.ID)
	return
}

// find TokenModel in database by id field.
func (m *Model) FindByID() (found bool, err error) {
	found = false
	var query = "SELECT * FROM tokens WHERE id=$1 LIMIT 1"
	founded, err := tokenAdapter.Find(query, m.ID)
	if err != nil {
		return
	}
	if founded != nil {
		found = true
		*m = *founded
	}
	return
}

// delete TokenModel from database by id field.
func (m *Model) DeleteByID() (err error) {
	var query = "DELETE FROM tokens WHERE id=$1"
	_, err = tokenAdapter.Exec(query, m.ID)
	return
}

// executes before token update.
func (m *Model) hookBeforeUpdate() {
	if m.AuthAgent != nil && len(*m.AuthAgent) > 323 {
		var authAgent = *m.AuthAgent
		var cut = 323 - len(authAgent)
		*m.AuthAgent = authAgent[:len(authAgent)-cut]
	}
	if m.LastAgent != nil && len(*m.LastAgent) > 323 {
		var lastAgent = *m.LastAgent
		var cut = 323 - len(lastAgent)
		*m.LastAgent = lastAgent[:len(lastAgent)-cut]
	}
	if m.AuthIP != nil && len(*m.AuthIP) > 53 {
		var authIP = *m.AuthIP
		var cut = 53 - len(authIP)
		*m.AuthIP = authIP[:len(authIP)-cut]
	}
	if m.LastIP != nil && len(*m.LastIP) > 53 {
		var lastIP = *m.LastIP
		var cut = 53 - len(lastIP)
		*m.LastIP = lastIP[:len(lastIP)-cut]
	}
}

// writes last ip and user agent then updating model in database.
func (m *Model) SetAuthAgents(request *http.Request) (err error) {
	m.AuthAgent = new(string)
	*m.AuthAgent = request.UserAgent()
	m.AuthIP = new(string)
	*m.AuthIP = utils.GetIP(request)
	err = m.Update()
	return
}

// writes ip and user agent then updating model in database.
func (m *Model) SetLastAgents(request *http.Request) (err error) {
	var lastAgent = request.UserAgent()
	var lastIP = utils.GetIP(request)
	m.LastAgent = new(string)
	*m.LastAgent = lastAgent
	m.LastIP = new(string)
	*m.LastIP = lastIP
	err = m.Update()
	return
}

// generate token.
//
// returns:
//
// token - token for user
//
// hash - saved in db as TokenModel.Token.
func (m *Model) Generate(userID string) (token string, hash string, err error) {
	// token generating.
	// first we generate fake token model to get created token ID.
	m.UserID = userID
	m.Token = "-1"
	if err = m.Create(); err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = m.DeleteByID()
		}
	}()
	// then we get newly created token model id and encrypt im. That's we send to user as token.
	token, err = call.Encryptor.AES.Encrypt(m.ID)
	if err != nil {
		return
	}
	// get hash from generated token.
	// user gets encrypted token, but database gets hash. In general, we do the same as with the password.
	hash, err = call.Encryptor.Argon.Hash(token)
	if err != nil {
		return
	}
	// now we replace fake token with real token in database.
	m.Token = hash
	err = m.Update()
	return
}
