package elven

import (
	"net/http"
)

func authGrabToken(request *http.Request) (string, error){
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

//func authGetUserAndTokenByToken(token string) (modelUser, modelToken){
//	var decrypted, err = cryptor.AESDecrypt(token, core.Config.Secret)
//
//	var hashedToken, err = cryptor.BHash(token)
//}
