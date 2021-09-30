package elven

import (
	"github.com/pkg/errors"
	"net/http"
	"servus/core"
	"servus/core/modules/errorCollector"
)

var errDatabaseNotFound = errors.New("database: not found")
var errTokenNotPresented = errors.New("authorization token not found (not in cookie, not in authorization header).")

// used when wrong username or password
func errAuthWrongCredentials(response http.ResponseWriter) {
	var theResponse = core.HttpResponse{ResponseWriter: response}
	var ec = errorCollector.New()
	ec.AddEAuthIncorrect([]string{"auth"})
	theResponse.Send(ec.GetErrors(), 401)
	return
}

// used when auth unknown error
func errAuth500(response http.ResponseWriter){
	var theResponse = core.HttpResponse{ResponseWriter: response}
	var ec = errorCollector.New()
	ec.AddEUnknown([]string{"auth"}, "Server error during auth.")
	theResponse.Send(ec.GetErrors(), 500)
	return
}
