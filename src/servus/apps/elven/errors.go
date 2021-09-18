package elven

import (
	"github.com/pkg/errors"
	"net/http"
	"servus/core/modules/errorCollector"
)

var errDatabaseNotFound = errors.New("database: not found")
var errTokenNotPresented = errors.New("authorization: token not found not in cookie, not in authorization header")

// used when wrong username or password
func errAuthWrongCredentials(response http.ResponseWriter) {
	var ec = errorCollector.New()
	ec.AddEAuthIncorrect([]string{"auth"})
	response.WriteHeader(401)
	response.Write([]byte(ec.GetErrors()))
	return
}

// used when auth unknown error
func errAuth500(response http.ResponseWriter){
	var ec = errorCollector.New()
	ec.AddEUnknown([]string{"auth"}, "Server error during auth.")
	response.WriteHeader(500)
	response.Write([]byte(ec.GetErrors()))
	return
}
