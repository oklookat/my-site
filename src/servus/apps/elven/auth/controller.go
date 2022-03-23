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
	isValid := body.Validate(request.Body)
	if !isValid {
		h.Send("invalid request", 400, nil)
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
		// wrong password.
		var ip = call.Banhammer.GetIpByRequest(request)
		if ip != nil {
			// warn IP.
			call.Banhammer.Warn(ip.String())
		}
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
	if err = tokenModel.SetAuthAgents(request); err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}

	// send token depending on auth type.
	switch body.Type {
	case "direct":
		var tokenJSON = fmt.Sprintf(`{token: "%v"}`, encryptedToken)
		h.Send(tokenJSON, 200, err)
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

// get token from user and delete. | PROTECTED BY AUTHORIZED ONLY MIDDLEWARE.
func (a *Instance) logout(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)

	// get token from cookie or auth header.
	var pipe = a.pipe.GetByContext(request)
	if pipe == nil {
		h.Send("not authorized", 400, nil)
		return
	}

	// delete token.
	var token = model.Token{}
	token.ID = pipe.GetID()
	var err error
	if err = token.DeleteByID(); err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}

	h.Send("", 200, err)
}
