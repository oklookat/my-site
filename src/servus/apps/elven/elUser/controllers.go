package elUser

import (
	"fmt"
	"net/http"
	"servus/core/errorCollector"
)

// ControllerAuthLogin generate token if username and pass correct
func ControllerAuthLogin(response http.ResponseWriter, request *http.Request){
	var username = request.FormValue("username")
	var password = request.FormValue("password")
	var authType = request.FormValue("type")
	var isFull = len(username) > 1 && len(password) > 1 && len(authType) > 1
	var ec = errorCollector.New()
	if !isFull {
		ec.AddEValidationEmpty([]string{"username", "password", "type"})
		response.WriteHeader(400)
		response.Write([]byte(ec.GetErrors()))
		return
	}
	// detect auth type
	if authType != "cookie" && authType != "direct"{
		ec.AddEValidationAllowed([]string{"type"}, []string{"cookie", "direct"})
		response.WriteHeader(400)
		response.Write([]byte(ec.GetErrors()))
		return
	}
	var user, err = dbFindUserBy(username)
	// user not found by username
	if err != nil{
		if err.Error() == "PIPE_USER_NOT_FOUND" {
			serviceControllerAuthIncorrect(response)
			return
		}
	}
	// check password
	var isPassword = servus.Utils.HashPasswordCheck(password, user.password)
	// wrong password
	if !isPassword {
		serviceControllerAuthIncorrect(response)
		return
	}
	// generate token
	// send to user token, in database save hashed token
	var encrypted, _ = servus.Utils.EncryptAES(user.id)
	var encryptedHashed, _ = servus.Utils.HashPassword(encrypted)
	var token = modelToken{userID: user.id, token: encryptedHashed}
	dbCreateToken(token)
	response.WriteHeader(200)
	if authType == "direct" {
		var direct = []byte(fmt.Sprintf(`{token: "%v"}`, encrypted))
		response.Write(direct)
		return
	}
	if authType == "cookie" {
		servus.Utils.SetCookie(&response, "token", encrypted)
		response.Write([]byte(""))
		return
	}

}

func ControllerAuthLogout(w http.ResponseWriter, r *http.Request){

}

// used when wrong username or password
func serviceControllerAuthIncorrect(response http.ResponseWriter){
	var ec = errorCollector.New()
	ec.AddEAuthIncorrect([]string{"auth"})
	response.WriteHeader(403)
	response.Write([]byte(ec.GetErrors()))
	return
}