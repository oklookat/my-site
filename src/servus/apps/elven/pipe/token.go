package pipe

import (
	"net/http"
	"servus/apps/elven/base"
	"servus/apps/elven/model"
)

type _ctx_token string

const (
	CtxToken _ctx_token = "ELVEN_TOKEN_PIPE"
)

type Token struct {
}

type TokenPipe struct {
	model *model.Token
}

// get pipe by request context. Use only if you provided to request context.
func (t *Token) GetByContext(request *http.Request) base.TokenPipe {
	pipe, ok := request.Context().Value(CtxToken).(base.TokenPipe)
	if !ok {
		return nil
	}
	return pipe
}

// used for ex. providing pipe to request context.
func (t *Token) GetByRequest(request *http.Request) (base.TokenPipe, error) {

	// get encrypted.
	encrypted, found := t.getEncryptedByRequest(request)
	if !found {
		return nil, nil
	}

	// decrypt to get ID.
	id, err := call.Encryptor.AES.Decrypt(encrypted)
	if err != nil {
		return nil, err
	}

	// search by ID.
	var md = &model.Token{}
	md.ID = id
	found, err = md.FindByID()
	if !found || err != nil {
		return nil, err
	}

	_ = md.SetLastAgents(request)

	// create pipe.
	var pipe = &TokenPipe{}
	pipe.model = md
	return pipe, err
}

// get encrypted token from request cookie or headers.
func (t *Token) getEncryptedByRequest(request *http.Request) (encrypted string, found bool) {
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
	encrypted = authHeader[:len(authHeader)-6]

	found = true
	return
}

func (t *TokenPipe) GetID() string {
	return t.model.ID
}

func (t *TokenPipe) GetUserID() string {
	return t.model.UserID
}

func (t *TokenPipe) GetToken() string {
	return t.model.Token
}
