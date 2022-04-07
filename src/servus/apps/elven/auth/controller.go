package auth

import (
	"encoding/json"
	"net/http"
	"servus/apps/elven/model"
)

// generate token if username and password are correct.
func login(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)

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
		h.Send(throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(throw.NotAuthorized(), 401, err)
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
		h.Send(throw.NotAuthorized(), 401, err)
		return
	}

	// generate token.
	var tokenModel = model.Token{}
	encryptedToken, _, err := tokenModel.Generate(user.ID)
	if err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}
	if err = tokenModel.SetAuthAgents(request); err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}

	// send token depending on auth type.
	switch body.Type {
	case "direct":
		type tokend struct {
			Token string `json:"token"`
		}
		var tokened = tokend{Token: encryptedToken}
		var tokenBytes, err = json.Marshal(tokened)
		if err != nil {
			h.Send(throw.Server(), 500, err)
			return
		}
		h.Send(string(tokenBytes), 200, err)
		return
	case "cookie":
		h.SetCookie("token", encryptedToken)
		h.Send("", 200, err)
		return
	default:
		h.Send(throw.Server(), 500, err)
		return
	}
}

// get token from user and delete. | PROTECTED BY AUTHORIZED ONLY MIDDLEWARE.
func logout(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)

	// get token from cookie or auth header.
	var pipe = pipe.GetByContext(request)
	if !pipe.IsExists() {
		h.Send("not authorized", 400, nil)
		return
	}

	// delete token.
	var token = model.Token{}
	token.ID = pipe.GetID()
	var err error
	if err = token.DeleteByID(); err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}

	h.Send("", 200, err)
}
