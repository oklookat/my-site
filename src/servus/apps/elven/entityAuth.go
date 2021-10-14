package elven

import (
	"fmt"
	"net/http"
	"servus/core"
	"servus/core/modules/cryptor"
	"servus/core/modules/errorMan"
)


// entityAuth - manage authorization.
type entityAuth struct {
	*entityBase
}

// controllerLogin -  generate token if username and password are correct.
func (a *entityAuth) controllerLogin(response http.ResponseWriter, request *http.Request) {
	val, em, _ := a.validatorControllerLogin(request)
	if em.HasErrors() {
		a.Send(response, em.GetJSON(), 400)
		return
	}
	var user = ModelUser{Username: val.Username}
	found, err := user.findByUsername()
	if err != nil {
		a.err500(response, request, err)
		return
	}
	if !found {
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
	err = tokenModel.create()
	if err != nil {
		a.err500(response, request, err)
		return
	}
	defer func() {
		if err != nil {
			_ = tokenModel.deleteByID()
		}
	}()
	// then we get newly created token model id and encrypt it. That's we send to user as token.
	encryptedToken, err := cryptor.AESEncrypt(tokenModel.ID, core.Config.Secret)
	if err != nil {
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
	err = tokenModel.setAuthAgents(request)
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
	var auth = PipeAuth{}
	auth.get(request)
	if auth.UserAndTokenExists {
		_ = auth.Token.deleteByID()
	}
	a.Send(response, "", 200)
}

// err500 - like unknown error.
func (a *entityAuth) err500(response http.ResponseWriter, request *http.Request, err error) {
	a.Logger.Warn("entityAuth code 500 at: %v. Error: %v", request.URL.Path, err.Error())
	a.Send(response, errorMan.ThrowServer(), 500)
	return
}

// err401 - like wrong username or password.
func (a *entityAuth) err401(response http.ResponseWriter) {
	a.Send(response, errorMan.ThrowForbidden(), 403)
	return
}
