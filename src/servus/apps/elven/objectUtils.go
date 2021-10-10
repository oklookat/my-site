package elven

import (
	"fmt"
	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"servus/core"
	"servus/core/modules/cryptor"
	"strings"
	"time"
)

// objectUtils - useful utilities.
type objectUtils struct {
}

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
	token.ID = tokenID
	found, err := token.findByID()
	if !found || err != nil {
		return nil, nil, err
	}
	// find user by id in found token.
	user.ID = token.UserID
	found, err = user.findByID()
	if !found || err != nil {
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

// createPipeAuth - check user permissions by request and accessType and write last agents and ip.
func (u *objectUtils) createPipeAuth(request *http.Request, accessType string) (auth PipeAuth) {
	auth = PipeAuth{}
	var err error
	auth.User, auth.Token, err = u.getUserAndToken(request)
	auth.UserAndTokenExists = auth.User != nil && auth.Token != nil && err == nil
	if auth.UserAndTokenExists {
		_ = auth.Token.setAuthAgents(request)
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
// relativePath - start from: '123\456\789\'
// total: if in dirs D:\Test\123\456\789\ has no files, it deletes all dirs up to D:\Test\.
func (u *objectUtils) deleteEmptyDirsRecursive(rootPath string, relativePath string) (err error) {
	rootPath = filepath.ToSlash(rootPath)
	relativePath = filepath.ToSlash(relativePath)
	relativePath, _ = filepath.Split(relativePath)
	var deleteDirIfEmpty = func(path string) (err error) {
		entry, err := os.ReadDir(path)
		if len(entry) == 0 {
			err = os.Remove(path)
		}
		return
	}
	var deletePath = rootPath + "/" + relativePath
	for !strings.EqualFold(rootPath, deletePath) {
		// delete dirs. Example:
		// it.1: delete D:\Test\123\456\789\ if empty
		// it.2: delete D:\Test\123\456\ if empty
		// it.3: delete D:\Test\123\ if empty
		// it.4: delete D:\Test\ if empty
		deletePath = path.Dir(deletePath)
		err = deleteDirIfEmpty(deletePath)
		if err != nil {
			break
		}
	}
	return
}

func (u *objectUtils) generateULID() (ul string, err error) {
	current := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(current.UnixNano())), 0)
	ulType, err := ulid.New(ulid.Timestamp(current), entropy)
	if err != nil {
		return "", err
	}
	ul = ulType.String()
	return
}
