package auth

import (
	"fmt"
	"net/http"
	"servus/apps/elven/model"
)

// generate token if username and password are correct.
func (a *Instance) login(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	// validate credentials.
	body := Body{}
	validator := body.Validate(request.Body)
	if validator.HasErrors() {
		h.Send(validator.GetJSON(), 400, nil)
		return
	}
	// find user.
	var user = model.User{Username: body.Username}
	found, err := user.FindByUsername()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(a.throw.NotAuthorized(), 401, err)
		return
	}
	// compare found user and provided password.
	var isPassword, _ = call.Encryptor.Argon.Compare(body.Password, user.Password)
	if !isPassword {
		h.Send(a.throw.NotAuthorized(), 401, err)
		return
	}
	// generate token.
	var tokenModel = model.Token{}
	encryptedToken, _, err := tokenModel.Generate(user.ID)
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	err = tokenModel.SetAuthAgents(request)
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	// send token depending by auth type.
	switch body.Type {
	case "direct":
		var direct = fmt.Sprintf(`{token: "%v"}`, encryptedToken)
		h.Send(direct, 200, err)
		return
	case "cookie":
		h.SetCookie("token", encryptedToken)
		h.Send("", 200, err)
		return
	default:
		h.Send(a.throw.Server(), 500, err)
		return
	}
}

// get token from user and delete.
func (a *Instance) logout(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	// get token from cookie or auth header.
	var pipe = a.pipe.GetByContext(request)
	if pipe == nil {
		// not authorized.
		h.Send("", 400, nil)
		return
	}
	// delete token.
	var token = model.Token{}
	token.ID = pipe.GetID()
	_ = token.DeleteByID()
	// send OK.
	h.Send("", 200, nil)
}
