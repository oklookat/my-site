package elControllers

import (
	"net/http"
	"servus/core/errorCollector"
)

func Login(response http.ResponseWriter, request *http.Request){
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


	//for key, value := range request.Form{
	//	println("lf")
	//	println(fmt.Sprintf("%v / %v", key, value))
	//}
}


func Logout(w http.ResponseWriter, r *http.Request){

}
