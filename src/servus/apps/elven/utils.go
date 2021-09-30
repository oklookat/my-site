package elven

import (
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
func getUserAndTokenByToken(tokenHex string) (modelUser, modelToken, error) {
	var user = modelUser{}
	var token = modelToken{}
	// get token id from encrypted token
	var tokenID, aesErr = cryptor.AESDecrypt(tokenHex, core.Config.Secret)
	if aesErr.HasErrors {
		return user, token, aesErr.AdditionalErr
	}
	// find token by id
	token, err := dbTokenFind(tokenID)
	if err != nil {
		return user, token, err
	}
	// find user by id in found token
	user, err = dbUserFind(token.userID)
	if err != nil {
		return user, token, err
	}
	return user, token, nil
}

// getUserAndTokenByRequest - get user and token by request.
func getUserAndTokenByRequest(request *http.Request) (modelUser, modelToken, error) {
	var user = modelUser{}
	var token = modelToken{}
	var tokenString, err = getToken(request)
	if err != nil {
		return user, token, err
	}
	user, token, err = getUserAndTokenByToken(tokenString)
	if err != nil {
		return user, token, err
	}
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
func setLastAgents(request *http.Request, token modelToken) modelToken {
	var lastAgent = request.UserAgent()
	var lastIP = getIP(request)
	token.lastAgent = lastAgent
	token.lastIP = lastIP
	newToken, err := dbTokenUpdate(&token)
	if err != nil {
		return token
	}
	return newToken
}

// createAuthData - check user permissions by request and accessType and write last agents and ip. Returns authData.
func createAuthData(request *http.Request, accessType string) authData {
	var auth = authData{}
	var user, token, err = getUserAndTokenByRequest(request)
	if err == nil {
		auth.user = user
		auth.token = token
		auth.token = setLastAgents(request, token)
		switch user.role {
		default:
			break
		case "admin":
			auth.access = true
			return auth
		}
	}
	switch accessType {
	default:
		auth.access = false
		break
	case accessTypeReadOnly:
		auth.access = isMethodReadOnly(request.Method)
		break
	case accessTypeAdminOnly:
		auth.access = false
		break
	}
	return auth
}

// getAuthData - get authData from request context. If nothing - returns nil.
func getAuthData(request *http.Request) *authData {
	var auth, err = request.Context().Value(ctxAuthData).(authData)
	if err {
		return nil
	}
	return &auth
}
