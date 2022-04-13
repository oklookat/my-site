package auth

import (
	"fmt"
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
	isUserFound, err := user.FindByUsername()
	if err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}
	var isPasswordCorrect = false
	if isUserFound {
		// compare found user password and provided password.
		isPasswordCorrect, _ = call.Encryptor.Argon.Compare(body.Password, user.Password)
	}

	if !isUserFound || !isPasswordCorrect {
		// warn wrong cases to avoid bruteforce
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
		var tokenJSON = fmt.Sprintf(`{"token": "%v"}`, encryptedToken)
		h.Send(tokenJSON, 200, err)
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

	// delete token.
	var token = model.Token{}
	token.ID = pipe.GetID()
	var err error
	if err = token.DeleteByID(); err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}
	_ = h.UnsetCookie("token")
	h.Send("", 200, err)
}
