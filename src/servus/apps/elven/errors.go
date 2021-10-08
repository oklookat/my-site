package elven

import (
	"net/http"
	"servus/core/modules/errorMan"
)

// err500 - write error to logger and send 500 error to user.
func (a *entityArticle) err500(response http.ResponseWriter, request *http.Request, err error) {
	a.Logger.Warn("entityArticle code 500 at: %v. Error: %v", request.URL.Path, err.Error())
	a.Send(response, errorMan.ThrowServer(), 500)
	return
}

// err500 - like unknown error.
func (a *entityAuth) err500(response http.ResponseWriter, request *http.Request, err error) {
	a.Logger.Warn("entityAuth code 500 at: %v. Error: %v", request.URL.Path, err.Error())
	a.Send(response, errorMan.ThrowServer(), 500)
	return
}

// err403 - send an error if the user is not allowed to do something.
func (a *entityArticle) err403(response http.ResponseWriter){
	a.Send(response, errorMan.ThrowForbidden(), 403)
	return
}
