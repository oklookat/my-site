package elven

import (
	"fmt"
	"net/http"
	"servus/core"
	"servus/core/modules/cryptor"
)

const articlesPageSize = 2

// controllerLogin -  generate token if username and password are correct.
func (a *entityAuth) controllerLogin(response http.ResponseWriter, request *http.Request) {
	val, em, _ := a.validatorControllerLogin(request)
	if em.HasErrors() {
		a.Send(response, em.GetJSON(), 400)
		return
	}
	user, err := eUser.databaseFindBy(val.Username)
	if err != nil {
		a.err500(response, request, err)
		return
	}
	if user == nil {
		a.err401(response)
		return
	}
	var isPassword = cryptor.BHashCheck(val.Password, user.Password)
	if !isPassword {
		a.err401(response)
		return
	}
	// token generating.
	// first we generate fake token model to get created token ID.
	var tokenModel = ModelToken{UserID: user.ID, Token: "-1"}
	err = eToken.databaseCreate(&tokenModel)
	if err != nil {
		a.err500(response, request, err)
		return
	}
	// then we get newly created token model id and encrypt it. That's we send to user as token.
	encryptedToken, aesErr := cryptor.AESEncrypt(tokenModel.ID, core.Config.Secret)
	if aesErr.HasErrors {
		a.err500(response, request, err)
		return
	}
	// get hash from generated token.
	// user gets encrypted token, but database gets hash. In general, we do the same as with the password.
	encryptedTokenHash, err := cryptor.BHash(encryptedToken)
	if err != nil {
		a.err500(response, request, err)
		return
	}
	// now we replace fake token with real token in database.
	tokenModel.Token = encryptedTokenHash
	oUtils.setAuthAgents(request, &tokenModel)
	if err != nil {
		a.err500(response, request, err)
		return
	}
	// based on auth type we send token.
	switch val.Type {
	case "direct":
		var direct = fmt.Sprintf(`{token: "%v"}`, encryptedToken)
		a.Send(response, direct, 200)
		return
	case "cookie":
		core.Utils.SetCookie(&response, "token", encryptedToken)
		a.Send(response, "", 200)
		return
	default:
		a.err500(response, request, err)
		return
	}
}

// logout - get token from user and delete.
func (a *entityAuth) controllerLogout(response http.ResponseWriter, request *http.Request) {
	// get token from cookie or auth header.
	var token, err = oUtils.getEncryptedToken(request)
	if err != nil {
		a.err401(response)
		return
	}
	// get user and token instances by encrypted token.
	_, tokenModel, err := oUtils.getUserAndTokenByEncrypted(token)
	if err != nil {
		a.err401(response)
		return
	}
	if tokenModel != nil {
		// delete token.
		_ = eToken.databaseDelete(tokenModel.ID)
	}
	a.Send(response, "", 200)
}
