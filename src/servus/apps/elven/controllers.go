package elven

import (
	"net/http"
	"servus/core/errorCollector"
)

// auth
func controllerAuthLogin(response http.ResponseWriter, request *http.Request){
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
	var isPassword = servus.Utils.HashPasswordCheck(password, user.password)
	// wrong password
	if !isPassword {
		serviceControllerAuthIncorrect(response)
		return
	}
	// TODO: replace 'hello' to user id
	// TODO: handle errors in crypt and decrypt methods
	var hashed, encrypted, _ = servus.Utils.EncryptAES("hello")
	var token = modelToken{userID: user.id, token: encrypted}
	dbCreateToken(token)
	println(hashed)
}

func controllerAuthLogout(w http.ResponseWriter, r *http.Request){

}

// used when wrong username or password
func serviceControllerAuthIncorrect(response http.ResponseWriter){
	var ec = errorCollector.New()
	ec.AddEAuthIncorrect([]string{"auth"})
	response.WriteHeader(403)
	response.Write([]byte(ec.GetErrors()))
	return
}