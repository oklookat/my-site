package model

import (
	"net/http"
	"strings"
	"time"
)

// represents token in database.
type Token struct {
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
func (t *Token) Create() (err error) {
	var query = `INSERT INTO tokens (user_id, token, last_ip, 
		last_agent, auth_ip, auth_agent) VALUES ($1, $2, $3, 
			$4, $5, $6) RETURNING *`
	err = tokenAdapter.Get(t, query, t.UserID, t.Token, t.LastIP, t.LastAgent, t.AuthIP, t.AuthAgent)
	return
}

// updates TokenModel in database. All fields except update and created dates must be filled.
func (t *Token) Update() (err error) {
	t.hookBeforeUpdate()
	var query = `UPDATE tokens SET user_id=$1, token=$2, last_ip=$3, 
	last_agent=$4, auth_ip=$5, auth_agent=$6 WHERE id=$7 RETURNING *`
	err = tokenAdapter.Get(t, query, t.UserID, t.Token, t.LastIP, t.LastAgent, t.AuthIP, t.AuthAgent, t.ID)
	return
}

// find TokenModel in database by id field.
func (t *Token) FindByID() (found bool, err error) {
	found = false
	var query = "SELECT * FROM tokens WHERE id=$1 LIMIT 1"
	founded, err := tokenAdapter.Find(query, t.ID)
	if err != nil {
		return
	}
	if founded != nil {
		found = true
		*t = *founded
	}
	return
}

// delete TokenModel from database by id field.
func (t *Token) DeleteByID() (err error) {
	var query = "DELETE FROM tokens WHERE id=$1"
	_, err = tokenAdapter.Exec(query, t.ID)
	return
}

// executes before token update.
func (t *Token) hookBeforeUpdate() {
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

// writes last ip and user agent then updating model in database.
func (t *Token) SetAuthAgents(request *http.Request) (err error) {
	t.AuthAgent = new(string)
	*t.AuthAgent = request.UserAgent()
	t.AuthIP = new(string)
	*t.AuthIP = t.getIP(request)
	err = t.Update()
	return
}

// writes ip and user agent then updating model in database.
func (t *Token) SetLastAgents(request *http.Request) (err error) {
	var lastAgent = request.UserAgent()
	var lastIP = t.getIP(request)
	t.LastAgent = new(string)
	*t.LastAgent = lastAgent
	t.LastIP = new(string)
	*t.LastIP = lastIP
	err = t.Update()
	return
}

// generate token.
//
// returns:
//
// token - token for user
//
// hash - saved in db as TokenModel.Token.
func (t *Token) Generate(userID string) (token string, hash string, err error) {
	// token generating.
	// first we generate fake token model to get created token ID.
	t.UserID = userID
	t.Token = "-1"
	if err = t.Create(); err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = t.DeleteByID()
		}
	}()
	// then we get newly created token model id and encrypt it. That's we send to user as token.
	token, err = call.Encryptor.AES.Encrypt(t.ID)
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
	t.Token = hash
	err = t.Update()
	return
}

func (t *Token) getIP(request *http.Request) (ip string) {
	ip = ""
	var ips = strings.Split(request.Header.Get("X-FORWARDED-FOR"), ", ")
	for _, theIP := range ips {
		if theIP != "" {
			ip = theIP
			break
		}
	}
	if ip == "" {
		ip = request.RemoteAddr
	}
	return
}
