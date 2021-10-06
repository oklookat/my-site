package elven

import (
	"encoding/json"
	"fmt"
	"net/http"
	"servus/core"
	"servus/core/modules/cryptor"
	"servus/core/modules/errorCollector"
)


type controllerAuthLoginBody struct {
	Username string
	Password string
	Type     string
}

// ControllerAuthLogin generate token if username and pass correct
func (c *controllerAuth) login(response http.ResponseWriter, request *http.Request) {
	// TODO: reformat controller and add errorCollector to baseController
	var theResponse = core.HttpResponse{ResponseWriter: response}
	var ec = errorCollector.New()
	// get and convert request body
	var loginBody controllerAuthLoginBody
	var err = json.NewDecoder(request.Body).Decode(&loginBody)
	if err != nil {
		ec.AddEValidationAllowed([]string{"you"}, []string{"username", "password", "type"})
		c.Send(response, ec.GetErrors(), 400)
		return
	}
	// get user credentials and other data
	var username = loginBody.Username
	var password = loginBody.Password
	var authType = loginBody.Type
	// validate
	var ecErrors = validatorAuth(username, password, authType)
	if ecErrors != nil {
		// response error collector JSON
		theResponse.Send(ecErrors.Error(), 400)
		return
	}
	// find username in database
	user, err := dbUserFindBy(username)
	// server user find error
	if err != nil {
		errAuth500(response)
		return
	}
	// user not found by username
	if len(user.ID) < 1 {
		errAuthWrongCredentials(response)
		return
	}
	var isPassword = cryptor.BHashCheck(password, user.Password)
	// wrong password
	if !isPassword {
		errAuthWrongCredentials(response)
		return
	}
	// token generating
	// first we generate fake token model to get created token ID
	var tokenModel = ModelToken{UserID: user.ID, Token: "-1"}
	tokenModel, err = dbTokenCreate(&tokenModel)
	if err != nil {
		errAuth500(response)
		return
	}
	// then we get newly created token model id and encrypt it. That's we send to user as token.
	encryptedToken, aesErr := cryptor.AESEncrypt(tokenModel.ID, core.Config.Secret)
	if aesErr.HasErrors {
		core.Logger.Error(fmt.Sprintf("%v / %v", aesErr.AdditionalErr.Error(), aesErr.OriginalErr.Error()))
		errAuth500(response)
		return
	}
	// get hash from generated token
	// user gets encrypted token, but database gets hash. In general, we do the same as with the password.
	encryptedTokenHash, err := cryptor.BHash(encryptedToken)
	if err != nil {
		core.Logger.Error(err.Error())
		errAuth500(response)
		return
	}
	// now we replace fake token with real token in database
	tokenModel.Token = encryptedTokenHash
	setAuthAgents(request, &tokenModel)
	if err != nil {
		errAuth500(response)
		return
	}
	// based on auth type we send token
	switch authType {
	case "direct":
		var direct = fmt.Sprintf(`{token: "%v"}`, encryptedToken)
		theResponse.Send(direct, 200)
		return
	case "cookie":
		core.Utils.SetCookie(&response, "token", encryptedToken)
		theResponse.Send("", 200)
		return
	default:
		core.Logger.Error("authController: this is authType switch default action. Check your code for oddities.")
		errAuth500(response)
		return
	}
}

func (c *controllerAuth) logout(response http.ResponseWriter, request *http.Request) {
	var theResponse = core.HttpResponse{ResponseWriter: response}
	var ec = errorCollector.New()
	// get token from cookie or auth header
	var token, err = getToken(request)
	if err != nil {
		ec.AddEAuthIncorrect([]string{"logout"})
		theResponse.Send(ec.GetErrors(), 401)
		return
	}
	// get user and token instances by encrypted token
	_, tokenModel, err := getUserAndTokenByToken(token)
	if err != nil {
		ec.AddEAuthIncorrect([]string{"logout"})
		theResponse.Send(ec.GetErrors(), 401)
		return
	}
	// delete token
	_ = dbTokenDelete(tokenModel.ID)
	theResponse.Send("", 200)
}
