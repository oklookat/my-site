package auth

import (
	"fmt"
	"net/http"
	"servus/apps/elven/model"
	"servus/core/external/errorMan"
)

// login - generate token if username and password are correct.
func (a *route) login(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	val, em, _ := a.validate.controllerLogin(request)
	if em.HasErrors() {
		h.Send(em.GetJSON(), 400, nil)
		return
	}
	var user = model.User{Username: val.Username}
	found, err := user.FindByUsername()
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
	var tokenModel = model.Token{}
	encryptedToken, _, err := tokenModel.Generate(user.ID)
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	err = tokenModel.SetAuthAgents(request)
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
func (a *route) logout(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	// get token from cookie or auth header.
	var pipe = a.pipe.GetByContext(request)
	if pipe != nil {
		var token = model.Token{}
		token.ID = pipe.GetID()
		_ = token.DeleteByID()
	}
	h.Send("", 200, nil)
}
