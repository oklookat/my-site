package elUser

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
func ControllerAuthLogin(response http.ResponseWriter, request *http.Request) {
	var loginBody controllerAuthLoginBody
	var err = json.NewDecoder(request.Body).Decode(&loginBody)
	if err != nil {
		http.Error(response, "auth error", http.StatusBadRequest)
		return
	}
	var username = loginBody.Username
	var password = loginBody.Password
	var authType = loginBody.Type
	var isFull = len(username) > 1 && len(password) > 1 && len(authType) > 1
	var ec = errorCollector.New()
	if !isFull {
		ec.AddEValidationEmpty([]string{"username", "password", "type"})
		response.WriteHeader(400)
		response.Write([]byte(ec.GetErrors()))
		return
	}
	// detect auth type
	if authType != "cookie" && authType != "direct" {
		ec.AddEValidationAllowed([]string{"type"}, []string{"cookie", "direct"})
		response.WriteHeader(400)
		response.Write([]byte(ec.GetErrors()))
		return
	}
	user, err := dbFindUserBy(username)
	// user not found by username
	if err != nil {
		if err.Error() == "PIPE_USER_NOT_FOUND" {
			serviceControllerAuthIncorrect(response)
			return
		}
	}
	// check password
	var isPassword = cryptor.BHashCheck(password, user.password)
	// wrong password
	if !isPassword {
		serviceControllerAuthIncorrect(response)
		return
	}
	// generate token
	// send to user token, in database save hashed token
	var encrypted, _ = cryptor.AESEncrypt(user.id, core.Config.Secret)
	var encryptedHashed, _ = cryptor.BHash(encrypted)
	var token = modelToken{userID: user.id, token: encryptedHashed}
	err = dbCreateToken(token)
	if err != nil {
		ec.AddEUnknown([]string{"auth"}, "Server error during auth.")
		response.WriteHeader(500)
		response.Write([]byte(ec.GetErrors()))
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
	}
}

func ControllerAuthLogout(w http.ResponseWriter, r *http.Request) {

}

// used when wrong username or password
func serviceControllerAuthIncorrect(response http.ResponseWriter) {
	var ec = errorCollector.New()
	ec.AddEAuthIncorrect([]string{"auth"})
	response.WriteHeader(403)
	response.Write([]byte(ec.GetErrors()))
	return
}
