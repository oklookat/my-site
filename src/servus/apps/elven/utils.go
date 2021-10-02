package elven

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"servus/core"
	"servus/core/modules/cryptor"
	"strings"
	"time"
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
func getUserAndTokenByToken(tokenHex string) (user ModelUser, token ModelToken, err error) {
	user = ModelUser{}
	token = ModelToken{}
	// get token id from encrypted token
	var tokenID, aesErr = cryptor.AESDecrypt(tokenHex, core.Config.Secret)
	if aesErr.HasErrors {
		return user, token, aesErr.AdditionalErr
	}
	// find token by id
	token, err = dbTokenFind(tokenID)
	if err != nil {
		return user, token, err
	}
	// find user by id in found token
	user, err = dbUserFind(token.UserID)
	if err != nil {
		return user, token, err
	}
	return user, token, nil
}

// getUserAndTokenByRequest - get user and token by request.
func getUserAndTokenByRequest(request *http.Request) (ModelUser, ModelToken, error) {
	var user = ModelUser{}
	var token = ModelToken{}
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
func setLastAgents(request *http.Request, token ModelToken) ModelToken {
	var lastAgent = request.UserAgent()
	var lastIP = getIP(request)
	*token.LastAgent = lastAgent
	*token.LastIP = lastIP
	newToken, err := dbTokenUpdate(&token)
	if err != nil {
		return token
	}
	return newToken
}

// createAuthData - check user permissions by request and accessType and write last agents and ip. Returns authData.
func createAuthData(request *http.Request, accessType string) AuthData {
	var auth = AuthData{}
	var user, token, err = getUserAndTokenByRequest(request)
	if err == nil {
		auth.User = user
		auth.Token = token
		auth.Token = setLastAgents(request, token)
		switch user.Role {
		default:
			break
		case "admin":
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

// paginationEncodeCursor - encode cursor with base 64.
func paginationEncodeCursor(t time.Time, uuid string) string {
	key := fmt.Sprintf("%s,%s", t.Format(time.RFC3339Nano), uuid)
	return base64.StdEncoding.EncodeToString([]byte(key))
}

func paginationDecodeCursor(encodedCursor string) (res time.Time, uuid string, err error) {
	byt, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return
	}
	arrStr := strings.Split(string(byt), ",")
	if len(arrStr) != 2 {
		err = errors.New("cursor is invalid")
		return
	}
	res, err = time.Parse(time.RFC3339Nano, arrStr[0])
	if err != nil {
		return
	}
	uuid = arrStr[1]
	return
}
