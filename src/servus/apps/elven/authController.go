package elven

import (
	"fmt"
	"net/http"
	"servus/core/external/errorMan"
)

// login -  generate token if username and password are correct.
func (a *authController) login(response http.ResponseWriter, request *http.Request) {
	val, em, _ := a.validate.controllerLogin(request)
	if em.HasErrors() {
		a.Send(response, em.GetJSON(), 400)
		return
	}
	var user = UserModel{Username: val.Username}
	found, err := user.findByUsername()
	if err != nil {
		a.err500(response, request, err)
		return
	}
	if !found {
		a.err401(response)
		return
	}
	var isPassword, _ = call.Encryption.Argon.Check(val.Password, user.Password)
	if !isPassword {
		a.err401(response)
		return
	}
	// token generating.
	var tokenModel = TokenModel{}
	err, encryptedToken, _ := tokenModel.generate(user.ID)
	if err != nil {
		a.err500(response, request, err)
		return
	}
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
		a.SetCookie(&response, "token", encryptedToken)
		a.Send(response, "", 200)
		return
	default:
		a.err500(response, request, err)
		return
	}
}

// logout - get token from user and delete.
func (a *authController) logout(response http.ResponseWriter, request *http.Request) {
	// get token from cookie or auth header.
	var auth = AuthPipe{}
	auth.get(request)
	if auth.UserAndTokenExists {
		_ = auth.Token.deleteByID()
	}
	a.Send(response, "", 200)
}

// err500 - like unknown error.
func (a *authController) err500(response http.ResponseWriter, request *http.Request, err error) {
	call.Logger.Warn(fmt.Sprintf("entityAuth code 500 at: %v. Error: %v", request.URL.Path, err.Error()))
	a.Send(response, errorMan.ThrowServer(), 500)
	return
}

// err401 - like wrong username or password.
func (a *authController) err401(response http.ResponseWriter) {
	a.Send(response, errorMan.ThrowNotAuthorized(), 401)
	return
}
