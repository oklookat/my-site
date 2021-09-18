package elven

import (
	"encoding/json"
	"fmt"
	"net/http"
	"servus/core"
	"servus/core/modules/cryptor"
	"servus/core/modules/errorCollector"
)


// ControllerAuthLogin generate token if username and pass correct
func controllerAuthLogin(response http.ResponseWriter, request *http.Request) {
	var ec = errorCollector.New()
	// get and convert request body
	var loginBody controllerAuthLoginBody
	var err = json.NewDecoder(request.Body).Decode(&loginBody)
	if err != nil {
		ec.AddEValidationAllowed([]string{"request body"}, []string{"username", "password", "type"})
		response.WriteHeader(400)
		response.Write([]byte(ec.GetErrors()))
		return
	}
	// get user credentials and other data
	var username = loginBody.Username
	var password = loginBody.Password
	var authType = loginBody.Type
	// validate
	var ecErrors = validatorAuth(username, password, authType)
	if ecErrors != nil {
		response.WriteHeader(400)
		response.Write([]byte(ecErrors.Error()))
		return
	}
	// find username in database
	user, err := dbUserFindBy(username)
	if err != nil {
		// user not found by username
		if err == errDatabaseNotFound {
			errAuthWrongCredentials(response)
			return
		} else {
			// server token saving error
			errAuth500(response)
			return
		}
	}
	var isPassword = cryptor.BHashCheck(password, user.password)
	// wrong password
	if !isPassword {
		errAuthWrongCredentials(response)
		return
	}
	// token generating
	// first we generate fake token to get created token ID
	var token = modelToken{userID: user.id, token: "-1"}
	token, err = dbTokenCreate(&token)
	if err != nil {
		errAuth500(response)
		return
	}
	// then we encrypt this token id
	encrypted, aesErr := cryptor.AESEncrypt(token.id, core.Config.Secret)
	if aesErr.HasErrors {
		core.Logger.Error(fmt.Sprintf("%v / %v", aesErr.AdditionalErr.Error(), aesErr.OriginalErr.Error()))
		errAuth500(response)
		return
	}
	// get hash from generated token
	// user gets token, but database gets hash. In general, we do the same as with the password.
	encryptedHashed, err := cryptor.BHash(encrypted)
	if err != nil {
		core.Logger.Error(err.Error())
		errAuth500(response)
		return
	}
	// store hashed token to model
	token.token = encryptedHashed
	// update token model
	_, err = dbTokenUpdate(&token)
	if err != nil {
		errAuth500(response)
		return
	}
	switch authType {
	case "direct":
		var direct = []byte(fmt.Sprintf(`{token: "%v"}`, encrypted))
		response.Write(direct)
		return
	case "cookie":
		core.Utils.SetCookie(&response, "token", encrypted)
		response.Write([]byte(""))
		return
	default:
		core.Logger.Error("authController: this is authType switch default action. Check your code for strange stuff.")
		errAuth500(response)
		return
	}
}

//func controllerAuthLogout(response http.ResponseWriter, request *http.Request) {
//	var theResponse = core.HttpResponse{ResponseWriter: response}
//	var token, err = authGrabToken(request)
//	if err != nil {
//		var ec = errorCollector.New()
//		ec.AddEAuthIncorrect([]string{"logout"})
//		theResponse.Send(ec.GetErrors(), 401)
//		return
//	}
//
//}
