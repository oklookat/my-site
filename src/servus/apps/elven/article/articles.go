package article

import (
	"encoding/json"
	"errors"
	"net/http"
	"servus/apps/elven/base"
)

// ALL HANDLERS PROTECTED BY SAFE METHODS MIDDLEWARE.

// get paginated articles (GET url/).
func getArticles(response http.ResponseWriter, request *http.Request) {
	var err error
	var h = call.Http.Get(request)
	var pipe = pipe.GetByContext(request)
	var isAdmin = pipe.IsAdmin()

	// validate.
	validatedBody := &GetParams{}
	ValidateGetParams(validatedBody, request.URL.Query(), isAdmin)

	// get.
	article := Model{}
	articles, totalPages, err := article.GetPaginated(validatedBody, isAdmin)
	if err != nil {
		h.Send("", 500, err)
		return
	}

	// fill response.
	var responseContent = base.ResponseContent{}
	responseContent.Meta.CurrentPage = validatedBody.Page
	responseContent.Meta.TotalPages = totalPages
	responseContent.Meta.PerPage = pageSize
	responseContent.Data = articles

	// send.
	jsonResponse, err := json.Marshal(&responseContent)
	if err != nil {
		h.Send("", 500, err)
		return
	}
	h.Send(string(jsonResponse), 200, err)
}

// get one article (GET url/id).
func getArticle(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)
	var isAdmin = false

	// get id from params.
	var id = h.GetRouteArgs()["id"]

	var pAuth = pipe.GetByContext(request)
	isAdmin = pAuth.IsAdmin()

	// find.
	var article = Model{ID: id}
	found, err := article.FindByID(isAdmin)
	if err != nil {
		h.Send("", 500, err)
		return
	}
	if !found {
		h.Send("", 404, err)
		return
	}

	// if not published and user not admin, access denied.
	if !article.IsPublished && !isAdmin {
		h.Send("", 403, err)
		return
	}

	// send.
	articleJson, err := json.Marshal(article)
	if err != nil {
		h.Send("", 500, err)
		return
	}
	h.Send(string(articleJson), 200, err)
}

// сreate new article (POST url/).
func createArticle(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)

	// validate.
	var article, err = ValidateBody(request.Method, request.Body, nil)
	if err != nil {
		switch err.(type) {
		default:
			h.Send("", 500, err)
			return
		case *base.ValidationError:
			h.Send("", 400, err)
			return
		}
	}

	var pAuth = pipe.GetByContext(request)
	article.UserID = pAuth.GetID()

	// create (get ID after creating).
	if err = article.Create(); err != nil {
		h.Send("", 500, err)
		return
	}

	// send created.
	articleJson, err := json.Marshal(&article)
	if err != nil {
		h.Send("", 500, err)
		return
	}
	h.Send(string(articleJson), 200, err)
}

// update all fields (PUT url/id)
//
// update specific fields (PATCH url/id).
func updateArticle(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)

	// get id from params.
	var id = h.GetRouteArgs()["id"]

	// find.
	var article = Model{ID: id}

	// isAdmin = true because this handler protected by admin-only
	found, err := article.FindByID(true)
	if err != nil {
		h.Send("", 500, err)
		return
	}
	if !found {
		h.Send("", 404, err)
		return
	}

	// validate body.
	var filteredArticle *Model
	filteredArticle, err = ValidateBody(request.Method, request.Body, &article)
	if err != nil {
		var valError *base.ValidationError
		if errors.As(err, &valError) {
			h.Send("", 400, err)
			return
		} else {
			h.Send("", 500, err)
			return
		}
	}
	article = *filteredArticle

	// update.
	if err = article.Update(); err != nil {
		h.Send("", 500, err)
		return
	}

	// send updated.
	jsonArticle, err := json.Marshal(article)
	if err != nil {
		h.Send("", 500, err)
		return
	}
	h.Send(string(jsonArticle), 200, err)
}

// DELETE url/id. Deletes one article.
func deleteArticle(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)

	// get id from params.
	var id = h.GetRouteArgs()["id"]

	// find.
	var article = Model{ID: id}

	// isAdmin = true because this handler protected by admin-only
	found, err := article.FindByID(true)
	if err != nil {
		h.Send("", 500, err)
		return
	}
	if !found {
		h.Send("", 404, err)
		return
	}

	// delete.
	if err = article.DeleteByID(); err != nil {
		h.Send("", 500, err)
		return
	}
	h.Send("", 200, err)
}
