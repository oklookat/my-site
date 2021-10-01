package elven

import (
	"net/http"
	"servus/core"
)

// GET url/
// params:
// page = number
// show = published, drafts, all
// by = created, updated, published
// start = newest (DESC), oldest (ASC)
// preview = true (content < 480 symbols), false (gives you full articles)
func controllerArticlesGet(response http.ResponseWriter, request *http.Request) {
	var theResponse = core.HttpResponse{ResponseWriter: response}
	var authData = getAuthData(request)
	var isAdmin = false
	if authData != nil {
		isAdmin = authData.User.Role == "admin"
	}
	var validated, err = validatorArticleQueryParams(request, isAdmin)
	if err != nil {
		theResponse.Send(err.Error(), 400)
		return
	}

}

// POST url/
func controllerArticlesPost(response http.ResponseWriter, request *http.Request) {

}

// PUT /url/id
func controllerArticlesPut(response http.ResponseWriter, request *http.Request) {

}

// DELETE /url/id
func controllerArticlesDelete(response http.ResponseWriter, request *http.Request) {

}
