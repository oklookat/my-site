package elven

import (
	"fmt"
	"net/http"
	"servus/core/external/errorMan"
)

// login - generate token if username and password are correct.
func (a *authRoute) login(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	val, em, _ := a.validate.controllerLogin(request)
	if em.HasErrors() {
		h.Send(em.GetJSON(), 400, nil)
		return
	}
	var user = UserModel{Username: val.Username}
	found, err := user.findByUsername()
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	if !found {
		h.Send(errorMan.ThrowNotAuthorized(), 401, err)
		return
	}
	var isPassword, _ = call.Encryptor.Argon.Compare(val.Password, user.Password)
	if !isPassword {
		h.Send(errorMan.ThrowNotAuthorized(), 401, err)
		return
	}
	// token generating.
	var tokenModel = TokenModel{}
	err, encryptedToken, _ := tokenModel.generate(user.ID)
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	err = tokenModel.setAuthAgents(request)
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	// based on auth type we send token.
	switch val.Type {
	case "direct":
		var direct = fmt.Sprintf(`{token: "%v"}`, encryptedToken)
		h.Send(direct, 200, err)
		return
	case "cookie":
		h.SetCookie("token", encryptedToken)
		h.Send("", 200, err)
		return
	default:
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
}

// logout - get token from user and delete.
func (a *authRoute) logout(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	// get token from cookie or auth header.
	var pipe = AuthPipe{}
	pipe.get(request)
	if pipe.UserAndTokenExists {
		_ = pipe.Token.deleteByID()
	}
	h.Send("", 200, nil)
}
