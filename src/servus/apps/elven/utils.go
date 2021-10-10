package elven

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"path/filepath"
	"servus/core"
	"servus/core/modules/cryptor"
	"strings"
)

// getEncryptedToken - get encrypted token from auth header or cookies.
func (u *objectUtils) getEncryptedToken(request *http.Request) (string, error) {
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
func (u *objectUtils) getUserAndTokenByEncrypted(encryptedToken string) (user *ModelUser, token *ModelToken, err error) {
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
func (u *objectUtils) getUserAndToken(request *http.Request) (user *ModelUser, token *ModelToken, err error) {
	user = &ModelUser{}
	token = &ModelToken{}
	tokenString, err := u.getEncryptedToken(request)
	if err != nil {
		return user, token, err
	}
	user, token, err = u.getUserAndTokenByEncrypted(tokenString)
	return user, token, nil
}

// isAdmin - check is request by admin.
func (u *objectUtils) isAdmin(request *http.Request) (isAdmin bool) {
	var authData = oUtils.getPipeAuth(request)
	isAdmin = false
	if authData != nil {
		isAdmin = authData.IsAdmin
	}
	return
}

// isMethodReadOnly - check is HTTP method readonly.
func (u *objectUtils) isMethodReadOnly(method string) bool {
	return method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions
}

// getIP - get IP by request.
func (u *objectUtils) getIP(request *http.Request) (ip string) {
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
func (u *objectUtils) setLastAgents(request *http.Request, token *ModelToken) (err error) {
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
func (u *objectUtils) setAuthAgents(request *http.Request, token *ModelToken) {
	token.AuthAgent = new(string)
	*token.AuthAgent = request.UserAgent()
	token.AuthIP = new(string)
	*token.AuthIP = u.getIP(request)
	_ = eToken.databaseUpdate(token)
}

// createPipeAuth - check user permissions by request and accessType and write last agents and ip.
func (u *objectUtils) createPipeAuth(request *http.Request, accessType string) (auth PipeAuth) {
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
func (u *objectUtils) getPipeAuth(request *http.Request) *PipeAuth {
	auth, ok := request.Context().Value(CtxAuthData).(PipeAuth)
	if !ok {
		return nil
	}
	return &auth
}

// generateDirsByHash - generate folders struct from hash like: 1d/2c. Returns error if hash length less than 6 symbols.
func (u *objectUtils) generateDirsByHash(hash string) (result string, err error) {
	if len(hash) < 6 {
		return "", errors.New("generateDirByHash: hash too short")
	}
	var hashFirstTwo1 = hash[0:2]
	var hashFirstTwo2 = hash[2:4]
	result = fmt.Sprintf("%v/%v", hashFirstTwo1, hashFirstTwo2)
	return
}

// deleteEmptyDirsRecursive - collect dirs and delete
// starts from rootPath + relativePath, and goes up to rootPath deleting dirs (if empty) along the way.
// params:
// rootPath - end at: 'D:\Test\'
// relativePath - start from: '123\456\789\music.flac' or '123\456\789\'
// total: if in dirs D:\Test\123\456\789\ has no files, it deletes all dirs up to D:\Test\.
func (u *objectUtils) deleteEmptyDirsRecursive(rootPath string, relativePath string) (err error){
	rootPath = filepath.ToSlash(rootPath)
	relativePath = filepath.ToSlash(relativePath)
	relativePath, _ = filepath.Split(relativePath)
	var pathSlice = strings.Split(relativePath, "/")
	var pathsForDelete []string
	for index := range pathSlice{
		if len(pathSlice[index]) < 1{
			continue
		}
		// make concat and collect paths for recursive. Example:
		// it.1 D:\Test\123\
		// it.2 D:\Test\123\456\
		// it.3 D:\Test\123\456\789\
		var p string
		if (index - 1) > -1 {
			p = rootPath + "/" + pathSlice[index - 1] + "/" + pathSlice[index]
		} else {
			p = rootPath + "/" + pathSlice[index]
		}
		pathsForDelete = append(pathsForDelete, p)
	}
	var deleteDirIfEmpty = func(path string) (err error) {
		entry, err := os.ReadDir(path)
		if len(entry) == 0 {
			err = os.Remove(path)
		}
		return
	}
	// reverse
	for i := len(pathsForDelete) - 1; i >= 0; i-- {
		// it.1: delete D:\Test\123\456\789\
		// it.2: delete D:\Test\123\456\
		// it.3: delete D:\Test\123\
		err = deleteDirIfEmpty(pathsForDelete[i])
		if err != nil {
			break
		}
	}
	return
}
