package elven

import (
	"net/http"
	"servus/core"
	"servus/core/modules/cryptor"
)

const (
	accessTypeAdminOnly  = "ELVEN_ACCESS_ADMIN_ONLY"
	accessTypeReadOnly   = "ELVEN_ACCESS_READ_ONLY"
	accessTypeAuthorized = "ELVEN_ACCESS_AUTHORIZED"
)

// PipeAuth - represents auth status on secured routes.
type PipeAuth struct {
	Access             bool
	UserAndTokenExists bool
	IsAdmin            bool
	User               *ModelUser
	Token              *ModelToken
}

// get - get PipeAuth from request context.
func (a *PipeAuth) get(request *http.Request) {
	auth, ok := request.Context().Value(CtxAuthData).(PipeAuth)
	if !ok {
		a.UserAndTokenExists = false
		return
	}
	*a = auth
	return
}

// create - get ModelUser, ModelToken, write permissions, agents etc. Used in middleware.
func (a *PipeAuth) create(request *http.Request, accessType string) {
	var err error
	token, found := a.getEncryptedToken(request)
	if found {
		err = a.setUserAndToken(token)
	}
	a.IsAdmin = false
	a.UserAndTokenExists = a.User != nil && a.Token != nil && err == nil
	if a.UserAndTokenExists {
		_ = a.Token.setLastAgents(request)
		switch a.User.Role {
		default:
			break
		case userRoleAdmin:
			a.IsAdmin = true
			a.Access = true
			return
		}
	}
	switch accessType {
	default:
		a.Access = false
		break
	case accessTypeReadOnly:
		a.Access = a.isMethodReadOnly(request)
		break
	case accessTypeAdminOnly:
		a.Access = false
		break
	case accessTypeAuthorized:
		a.Access = a.UserAndTokenExists && a.User.Role == userRoleUser
		break
	}
	return
}

// setUserAndToken - set ModelUser and ModelToken to PipeAuth.
func (a *PipeAuth) setUserAndToken(encryptedToken string) (err error) {
	// get token id from encrypted token.
	a.User = nil
	a.Token = nil
	tokenID, err := cryptor.AESDecrypt(encryptedToken, core.Config.Secret)
	if err != nil {
		return err
	}
	// find token by id.
	var tempToken = ModelToken{ID: tokenID}
	found, err := tempToken.findByID()
	if !found || err != nil {
		return err
	}
	// find user by id in found token.
	var tempUser = ModelUser{ID: tempToken.UserID}
	found, err = tempUser.findByID()
	if !found || err != nil {
		return err
	}
	a.User = &tempUser
	a.Token = &tempToken
	return
}

// getEncryptedToken - get encryptedToken from cookie or authorization header.
func (a *PipeAuth) getEncryptedToken(request *http.Request) (token string, found bool) {
	// get from cookie.
	found = false
	token = ""
	cookieToken, err := request.Cookie("token")
	if err == nil && len(cookieToken.Value) > 4 {
		found = true
		token = cookieToken.Value
		return
	}
	// get from authorization header.
	// get something like: 'Elven tokenHere'.
	var authHeader = request.Header.Get("Authorization")
	if len(authHeader) < 12 {
		return
	}
	// remove 6 symbols (Elven and space) to get only token.
	token = authHeader[:len(authHeader)-6]
	found = true
	return
}

// isMethodReadOnly - check is HTTP method readonly.
func (a *PipeAuth) isMethodReadOnly(request *http.Request) bool {
	var method = request.Method
	return method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions
}
