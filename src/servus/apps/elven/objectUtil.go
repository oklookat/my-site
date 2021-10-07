package elven

import (
	"github.com/pkg/errors"
	"net/http"
	"servus/core"
	"servus/core/modules/cryptor"
	"strings"
)

// objectUtil - useful utilities.
type objectUtil struct {
}

// PipeAuth - represents auth status in secured routes.
type PipeAuth struct {
	Access             bool
	UserAndTokenExists bool
	IsAdmin            bool
	User               *ModelUser
	Token              *ModelToken
}

// getEncryptedToken - get encrypted token from auth header or cookies.
func (u *objectUtil) getEncryptedToken(request *http.Request) (string, error) {
	// get from cookie.
	var cookieToken, err = request.Cookie("token")
	if err == nil && len(cookieToken.Value) > 4 {
		return cookieToken.Value, nil
	}
	// get from authorization header.
	// get something like: 'Elven tokenHere'.
	var authHeader = request.Header.Get("Authorization")
	if len(authHeader) < 8 {
		return "", errors.New("authorization token not found (not in cookie, not in authorization header)")
	}
	// remove 6 symbols (Elven and space) to get only token.
	var authToken = authHeader[:len(authHeader)-6]
	return authToken, nil
}

// getUserAndTokenByEncrypted - get ModelUser and ModelToken by encrypted (AES hex) token.
func (u *objectUtil) getUserAndTokenByEncrypted(encryptedToken string) (user *ModelUser, token *ModelToken, err error) {
	user = &ModelUser{}
	token = &ModelToken{}
	// get token id from encrypted token.
	var tokenID, aesErr = cryptor.AESDecrypt(encryptedToken, core.Config.Secret)
	if aesErr.HasErrors {
		return nil, nil, aesErr.AdditionalErr
	}
	// find token by id.
	token, err = eToken.databaseFind(tokenID)
	if err != nil || token == nil {
		return nil, nil, err
	}
	// find user by id in found token.
	user, err = eUser.databaseFind(token.UserID)
	if user == nil {
		return nil, nil, err
	}
	return
}

// getUserAndToken - get ModelUser and ModelToken by request.
func (u *objectUtil) getUserAndToken(request *http.Request) (user *ModelUser, token *ModelToken, err error) {
	user = &ModelUser{}
	token = &ModelToken{}
	tokenString, err := u.getEncryptedToken(request)
	if err != nil {
		return user, token, err
	}
	user, token, err = u.getUserAndTokenByEncrypted(tokenString)
	return user, token, nil
}

// isMethodReadOnly - check is HTTP method readonly.
func (u *objectUtil) isMethodReadOnly(method string) bool {
	return method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions
}

// getIP - get IP by request.
func (u *objectUtil) getIP(request *http.Request) (ip string) {
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

// setLastAgents - writes ip and user agent to ModelToken.
func (u *objectUtil) setLastAgents(request *http.Request, token *ModelToken) (err error) {
	if request == nil {
		return errors.New("setLastAgents: request nil pointer.")
	}
	if token == nil {
		return errors.New("setLastAgents: token nil pointer.")
	}
	var lastAgent = request.UserAgent()
	var lastIP = u.getIP(request)
	token.LastAgent = new(string)
	*token.LastAgent = lastAgent
	token.LastIP = new(string)
	*token.LastIP = lastIP
	err = eToken.databaseUpdate(token)
	return
}

// setAuthAgents - writes last ip and user agent to ModelToken.
func (u *objectUtil) setAuthAgents(request *http.Request, token *ModelToken) {
	token.AuthAgent = new(string)
	*token.AuthAgent = request.UserAgent()
	token.AuthIP = new(string)
	*token.AuthIP = u.getIP(request)
	_ = eToken.databaseUpdate(token)
}

// createPipeAuth - check user permissions by request and accessType and write last agents and ip.
func (u *objectUtil) createPipeAuth(request *http.Request, accessType string) (auth PipeAuth) {
	auth = PipeAuth{}
	var err error
	auth.User, auth.Token, err = u.getUserAndToken(request)
	auth.UserAndTokenExists = auth.User != nil && auth.Token != nil && err == nil
	if auth.UserAndTokenExists {
		_ = u.setLastAgents(request, auth.Token)
		switch auth.User.Role {
		default:
			break
		case "admin":
			auth.IsAdmin = true
			auth.Access = true
			return
		}
	}
	switch accessType {
	default:
		auth.Access = false
		break
	case accessTypeReadOnly:
		auth.Access = u.isMethodReadOnly(request.Method)
		break
	case accessTypeAdminOnly:
		auth.Access = false
		break
	}
	return
}

// getPipeAuth - get authData from request context. If nothing - returns nil.
func (u *objectUtil) getPipeAuth(request *http.Request) *PipeAuth {
	auth, ok := request.Context().Value(CtxAuthData).(PipeAuth)
	if !ok {
		return nil
	}
	return &auth
}
