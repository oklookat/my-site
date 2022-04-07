package article

import (
	"encoding/json"
	"errors"
	"net/http"
	"servus/apps/elven/base"
	"servus/apps/elven/model"
)

// ALL HANDLERS PROTECTED BY SAFE METHODS MIDDLEWARE.

// get paginated articles (GET url/).
func getArticles(response http.ResponseWriter, request *http.Request) {
	var h = call.Http.Get(request)

	var pipe = pipe.GetByContext(request)

	var isAdmin = pipe.IsAdmin()

	// validate.
	validatedBody := &base.ArticleGetParams{}
	ValidateGetParams(validatedBody, request.URL.Query(), isAdmin)

	// get.
	article := model.Article{}
	articles, totalPages, err := article.GetPaginated(validatedBody)
	if err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}

	// fill response.
	var responseContent = base.ResponseContent{}
	responseContent.Meta.CurrentPage = validatedBody.Page
	responseContent.Meta.TotalPages = totalPages
	responseContent.Meta.PerPage = model.ArticlePageSize
	responseContent.Data = articles

	// send.
	jsonResponse, err := json.Marshal(&responseContent)
	if err != nil {
		h.Send(throw.Server(), 500, err)
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

	// find.
	var article = model.Article{ID: id}
	found, err := article.FindByID()
	if err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(throw.NotFound(), 404, err)
		return
	}

	// if not published and user not admin, access denied.
	var pAuth = pipe.GetByContext(request)
	isAdmin = pAuth.IsAdmin()
	if !article.IsPublished && !isAdmin {
		h.Send(throw.Forbidden(), 403, err)
		return
	}

	// send.
	articleJson, err := json.Marshal(article)
	if err != nil {
		h.Send(throw.Server(), 500, err)
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
			h.Send(throw.Server(), 500, err)
			return
		case *base.ValidationError:
			h.Send("bad request", 400, err)
			return
		}
	}

	var pAuth = pipe.GetByContext(request)
	article.UserID = pAuth.GetID()

	// create (get ID after creating).
	if err = article.Create(); err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}

	// get created by id (we need to get category name which only available when we call find() method)
	found, err := article.FindByID()
	if err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(throw.NotFound(), 404, err)
		return
	}

	// send created.
	articleJson, err := json.Marshal(&article)
	if err != nil {
		h.Send(throw.Server(), 500, err)
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
	var article = model.Article{ID: id}
	found, err := article.FindByID()
	if err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(throw.NotFound(), 404, err)
		return
	}

	// validate body.
	var filteredArticle *model.Article
	filteredArticle, err = ValidateBody(request.Method, request.Body, &article)
	if err != nil {
		var valError *base.ValidationError
		if errors.As(err, &valError) {
			h.Send("bad request", 400, err)
			return
		} else {
			h.Send(throw.Server(), 500, err)
			return
		}
	}
	article = *filteredArticle

	// update.
	if err = article.Update(); err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}

	// send updated.
	jsonArticle, err := json.Marshal(article)
	if err != nil {
		h.Send(throw.Server(), 500, err)
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
	var article = model.Article{ID: id}
	found, err := article.FindByID()
	if err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(throw.NotFound(), 404, err)
		return
	}

	// delete.
	if err = article.DeleteByID(); err != nil {
		h.Send(throw.Server(), 500, err)
		return
	}
	h.Send("", 200, err)
}
