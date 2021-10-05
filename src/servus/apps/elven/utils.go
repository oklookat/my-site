package elven

import (
	"github.com/pkg/errors"
	"net/http"
	"servus/core"
	"servus/core/modules/cryptor"
	"strings"
)

// getToken - get token from auth header or cookies.
func getToken(request *http.Request) (string, error) {
	// get from cookie
	var cookieToken, err = request.Cookie("token")
	if err == nil && len(cookieToken.Value) > 4 {
		return cookieToken.Value, nil
	}
	// get from authorization header
	var authHeader = request.Header.Get("Authorization") // get something like: 'Elven tokenHere'
	if len(authHeader) < 8 {
		return "", errTokenNotPresented
	}
	var authToken = authHeader[:len(authHeader)-6] // remove 6 symbols (Elven and space) to get only token
	return authToken, nil
}

// getUserAndTokenByToken - get user and token model by encrypted token.
func getUserAndTokenByToken(tokenHex string) (user *ModelUser, token *ModelToken, err error) {
	user = &ModelUser{}
	token = &ModelToken{}
	// get token id from encrypted token
	var tokenID, aesErr = cryptor.AESDecrypt(tokenHex, core.Config.Secret)
	if aesErr.HasErrors {
		return nil, nil, aesErr.AdditionalErr
	}
	// find token by id
	token, err = dbTokenFind(tokenID)
	if err != nil {
		return user, token, err
	}
	if token == nil {
		return nil, nil, err
	}
	// find user by id in found token
	user, err = dbUserFind(token.UserID)
	return user, token, err
}

// getUserAndTokenByRequest - get user and token by request.
func getUserAndTokenByRequest(request *http.Request) (user *ModelUser, token *ModelToken, err error) {
	user = &ModelUser{}
	token = &ModelToken{}
	tokenString, err := getToken(request)
	if err != nil {
		return user, token, err
	}
	user, token, err = getUserAndTokenByToken(tokenString)
	return user, token, nil
}

// isMethodReadOnly - check is method readonly.
func isMethodReadOnly(method string) bool {
	return method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions
}

// getIP - get IP by request.
func getIP(request *http.Request) string {
	var ip = ""
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
	return ip
}

// setLastAgents - writes ip and user agent to token model.
func setLastAgents(request *http.Request, token *ModelToken) error {
	if request == nil {
		return errors.New("setLastAgents: request nil pointer.")
	}
	if token == nil {
		return errors.New("setLastAgents: token nil pointer.")
	}
	var lastAgent = request.UserAgent()
	var lastIP = getIP(request)
	token.LastAgent = new(string)
	*token.LastAgent = lastAgent
	token.LastIP = new(string)
	*token.LastIP = lastIP
	err := dbTokenUpdate(token)
	return err
}

func setAuthAgents(request *http.Request, token *ModelToken) {
	token.AuthAgent = new(string)
	*token.AuthAgent = request.UserAgent()
	token.AuthIP = new(string)
	*token.AuthIP = getIP(request)
	_ = dbTokenUpdate(token)
}

// createAuthData - check user permissions by request and accessType and write last agents and ip. Returns authData.
func createAuthData(request *http.Request, accessType string) AuthData {
	var auth = AuthData{}
	var err error
	auth.User, auth.Token, err = getUserAndTokenByRequest(request)
	auth.UserAndTokenExists = auth.User != nil && auth.Token != nil && err == nil
	if auth.UserAndTokenExists {
		_ = setLastAgents(request, auth.Token)
		switch auth.User.Role {
		default:
			break
		case "admin":
			auth.IsAdmin = true
			auth.Access = true
			return auth
		}
	}
	switch accessType {
	default:
		auth.Access = false
		break
	case accessTypeReadOnly:
		auth.Access = isMethodReadOnly(request.Method)
		break
	case accessTypeAdminOnly:
		auth.Access = false
		break
	}
	return auth
}

// getAuthData - get authData from request context. If nothing - returns nil.
func getAuthData(request *http.Request) *AuthData {
	var auth, ok = request.Context().Value(CtxAuthData).(AuthData)
	if !ok {
		return nil
	}
	return &auth
}
