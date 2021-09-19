package elven

import (
	"net/http"
	"servus/core"
	"servus/core/modules/cryptor"
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

func authGetUserAndTokenByToken(tokenHex string) (modelUser, modelToken, error){
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
