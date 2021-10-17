package elven

import (
	"database/sql"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

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

// create - create ModelToken in database. WARNING: its function writes only token and user id, other data will be ignored. For write full data see dbTokenUpdate.
func (t *ModelToken) create() (err error) {
	var query = "INSERT INTO tokens (user_id, token) VALUES ($1, $2) RETURNING *"
	row := instance.DB.Conn.QueryRowx(query, &t.UserID, &t.Token)
	err = row.StructScan(t)
	err = instance.DB.CheckError(err)
	return
}

// update - updates ModelToken in database. All fields except update and created dates must be filled.
func (t *ModelToken) update() (err error) {
	t.hookBeforeUpdate()
	var query = "UPDATE tokens SET user_id=:user_id, token=:token, last_ip=:last_ip, last_agent=:last_agent, auth_ip=:auth_ip, auth_agent=:auth_agent WHERE id=:id RETURNING *"
	stmt, err := instance.DB.Conn.PrepareNamed(query)
	if err != nil {
		return
	}
	defer func() {
		_ = stmt.Close()
	}()
	err = stmt.Get(t, t)
	err = instance.DB.CheckError(err)
	return
}

// databaseFind - find ModelToken in database by id field.
func (t *ModelToken) findByID() (found bool, err error) {
	var query = "SELECT * FROM tokens WHERE id=$1 LIMIT 1"
	err = instance.DB.Conn.Get(t, query, t.ID)
	err = instance.DB.CheckError(err)
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

// deleteByID - delete ModelToken from database by id field.
func (t *ModelToken) deleteByID() (err error) {
	var query = "DELETE FROM tokens WHERE id=$1"
	_, err = instance.DB.Conn.Exec(query, t.ID)
	if err == sql.ErrNoRows {
		return nil
	}
	return
}

// hookBeforeUpdate - executes before token update.
func (t *ModelToken) hookBeforeUpdate() {
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
func (t *ModelToken) setAuthAgents(request *http.Request) (err error) {
	if request == nil {
		return errors.New("setLastAgents: request nil pointer.")
	}
	if t == nil {
		return errors.New("setLastAgents: token nil pointer.")
	}
	t.AuthAgent = new(string)
	*t.AuthAgent = request.UserAgent()
	t.AuthIP = new(string)
	*t.AuthIP = oUtils.getIP(request)
	err = t.update()
	return
}

// setLastAgents - writes ip and user agent then updating model in database.
func (t *ModelToken) setLastAgents(request *http.Request) (err error) {
	if request == nil {
		return errors.New("setLastAgents: request nil pointer.")
	}
	if t == nil {
		return errors.New("setLastAgents: token nil pointer.")
	}
	var lastAgent = request.UserAgent()
	var lastIP = oUtils.getIP(request)
	t.LastAgent = new(string)
	*t.LastAgent = lastAgent
	t.LastIP = new(string)
	*t.LastIP = lastIP
	err = t.update()
	return
}
