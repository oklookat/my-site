package elven

import (
	"encoding/json"
	"fmt"
	"net/http"
	"servus/core"
	"servus/core/modules/cryptor"
	"servus/core/modules/errorMan"
	"servus/core/modules/validator"
)

// entityAuth - manage authorization.
type entityAuth struct {
	*entityBase
	bodyLogin *bodyAuthLogin
}

// authLoginBody - parsed request body in auth.controllerLogin.
type bodyAuthLogin struct {
	Username string
	Password string
	Type     string
}

// controllerLogin -  generate token if username and password are correct.
func (a *entityAuth) controllerLogin(response http.ResponseWriter, request *http.Request) {
	em, _ := a.validatorControllerLogin(request)
	if em.HasErrors() {
		a.Send(response, em.GetJSON(), 400)
		return
	}
	user, err := eUser.databaseFindBy(a.bodyLogin.Username)
	if err != nil {
		a.err500(response, request, err)
		return
	}
	if user == nil {
		a.errWrongCredentials(response)
		return
	}
	var isPassword = cryptor.BHashCheck(a.bodyLogin.Password, user.Password)
	if !isPassword {
		a.errWrongCredentials(response)
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
	oUtil.setAuthAgents(request, &tokenModel)
	if err != nil {
		a.err500(response, request, err)
		return
	}
	// based on auth type we send token.
	switch a.bodyLogin.Type {
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
	var token, err = oUtil.getEncryptedToken(request)
	if err != nil {
		a.errWrongCredentials(response)
		return
	}
	// get user and token instances by encrypted token.
	_, tokenModel, err := oUtil.getUserAndTokenByEncrypted(token)
	if err != nil {
		a.errWrongCredentials(response)
		return
	}
	if tokenModel != nil {
		// delete token.
		_ = eToken.databaseDelete(tokenModel.ID)
	}
	a.Send(response, "", 200)
}

// validatorControllerLogin - validate request body. Writes result in errorMan instance.
func (a *entityAuth) validatorControllerLogin(request *http.Request) (em *errorMan.EValidation, err error) {
	em = errorMan.NewValidation()
	a.bodyLogin = &bodyAuthLogin{}
	err = json.NewDecoder(request.Body).Decode(&a.bodyLogin)
	if err != nil {
		em.Add("body", "provide username, password, type fields.")
	} else {
		// get user credentials and other data.
		var username = a.bodyLogin.Username
		var password = a.bodyLogin.Password
		var authType = a.bodyLogin.Type
		if validator.IsEmpty(&username) {
			em.Add("username", "username cannot be empty.")
		}
		if validator.IsEmpty(&password) {
			em.Add("password", "password cannot be empty.")
		}
		if validator.IsEmpty(&authType) {
			em.Add("type", "type cannot be empty.")
		} else {
			var isAuthType = authType == "cookie" || authType == "direct"
			if !isAuthType {
				em.Add("type", "wrong type.")
			}
		}
	}
	return
}

// errWrongCredentials - like wrong username or password.
func (a *entityAuth) errWrongCredentials(response http.ResponseWriter) {
	a.Send(response, errorMan.ThrowForbidden(), 403)
	return
}
