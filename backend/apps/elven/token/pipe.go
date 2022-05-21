package token

import (
	"net/http"
	"servus/apps/elven/base"
)

type _string string

const (
	CTX _string = "ELVEN_TOKEN_PIPE"
)

type Pipe struct {
	model *Model
}

// is pipe model exists?
func (p *Pipe) IsExists() bool {
	return p.model != nil
}

// get pipe by request context. Use only if you provided to request context.
func (p *Pipe) GetByContext(request *http.Request) base.TokenPipe {
	pipe, ok := request.Context().Value(CTX).(base.TokenPipe)
	if !ok {
		var emptyPipe = &Pipe{}
		return emptyPipe
	}
	return pipe
}

// used for ex. providing pipe to request context.
func (p *Pipe) GetByRequest(request *http.Request) (base.TokenPipe, error) {

	// get encrypted.
	encrypted, found := p.getEncryptedByRequest(request)
	if !found {
		return nil, nil
	}

	// decrypt to get ID.
	id, err := call.Encryptor.AES.Decrypt(encrypted)
	if err != nil {
		return nil, err
	}

	// search by ID.
	var md = &Model{}
	md.ID = id
	found, err = md.FindByID()
	if !found || err != nil {
		return nil, err
	}

	_ = md.SetLastAgents(request)

	// create pipe.
	var pipe = &Pipe{}
	pipe.model = md
	return pipe, err
}

// get encrypted token from request cookie or headers.
func (p *Pipe) getEncryptedByRequest(request *http.Request) (encrypted string, found bool) {
	found = false

	// get from cookie.
	encrypted = ""
	cookieToken, err := request.Cookie("token")
	if err == nil && len(cookieToken.Value) > 4 {
		found = true
		encrypted = cookieToken.Value
		return
	}

	// get from authorization header
	// get string like: 'Elven tokenHere'.
	var authHeader = request.Header.Get("Authorization")
	if len(authHeader) < 12 {
		return
	}

	// remove 6 symbols (Elven and space) to get only token.
	encrypted = authHeader[6:]

	found = true
	return
}

func (p *Pipe) GetID() string {
	if !p.IsExists() {
		return ""
	}
	return p.model.ID
}

func (p *Pipe) GetUserID() string {
	if !p.IsExists() {
		return ""
	}
	return p.model.UserID
}

func (p *Pipe) GetToken() string {
	if !p.IsExists() {
		return ""
	}
	return p.model.Token
}
