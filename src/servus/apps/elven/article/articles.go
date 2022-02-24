package article

import (
	"encoding/json"
	"net/http"
	"servus/apps/elven/base"
	"servus/apps/elven/model"
)

// ALL HANDLERS PROTECTED BY SAFE METHODS MIDDLEWARE.

// get paginated articles (GET url/).
func (a *Instance) getArticles(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	var err error
	var isAdmin = false
	var pipe = a.pipe.GetByContext(request)
	isAdmin = pipe != nil && pipe.IsAdmin()
	// validate.
	validatedBody := &base.ArticleGetParams{}
	validator := ValidateGetParams(validatedBody, request.URL.Query(), isAdmin)
	if validator.HasErrors() {
		h.Send(validator.GetJSON(), 400, err)
		return
	}
	// get.
	article := model.Article{}
	articles, totalPages, err := article.GetPaginated(validatedBody)
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	// fill pagination.
	var responseContent = base.ResponseContent{}
	responseContent.Meta.CurrentPage = validatedBody.Page
	responseContent.Meta.TotalPages = totalPages
	responseContent.Meta.PerPage = model.ArticlePageSize
	responseContent.Data = articles
	// make & send json.
	jsonResponse, err := json.Marshal(&responseContent)
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	h.Send(string(jsonResponse), 200, err)
}

// get one article (GET url/id).
func (a *Instance) getArticle(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)
	var isAdmin = false
	// get id from params.
	var id = h.GetRouteArgs()["id"]
	// find.
	var article = model.Article{ID: id}
	found, err := article.FindByID()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(a.throw.NotFound(), 404, err)
		return
	}
	// if not published and user not admin, access denied.
	var pAuth = a.pipe.GetByContext(request)
	isAdmin = pAuth.IsAdmin()
	if !article.IsPublished && !isAdmin {
		h.Send(a.throw.Forbidden(), 403, err)
		return
	}
	// send.
	articleJson, err := json.Marshal(article)
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	h.Send(string(articleJson), 200, err)
}

// сreate new article (POST url/).
func (a *Instance) createArticle(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)

	// validate.
	var article = ValidateBody(request.Method, request.Body, nil)
	if article == nil {
		h.Send("invalid", 400, nil)
		return
	}

	var pAuth = a.pipe.GetByContext(request)
	article.UserID = pAuth.GetID()

	// create (get ID after creating).
	err := article.Create()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}

	// get created by id (we need to get category name which only available when we call find() method)
	found, err := article.FindByID()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(a.throw.NotFound(), 404, err)
		return
	}

	// send created.
	articleJson, err := json.Marshal(&article)
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	h.Send(string(articleJson), 200, err)
}

// update all fields (PUT url/id)
//
// update specific fields (PATCH url/id).
func (a *Instance) updateArticle(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)

	// get id from params.
	var id = h.GetRouteArgs()["id"]

	// find.
	var article = model.Article{ID: id}
	found, err := article.FindByID()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(a.throw.NotFound(), 404, err)
		return
	}

	// validate body.
	var filteredArticle = ValidateBody(request.Method, request.Body, &article)
	if filteredArticle == nil {
		h.Send("invalid", 400, err)
		return
	}
	article = *filteredArticle

	// update.
	err = article.Update()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}

	// send updated.
	jsonArticle, err := json.Marshal(article)
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	h.Send(string(jsonArticle), 200, err)
}

// DELETE url/id. Deletes one article.
func (a *Instance) deleteArticle(response http.ResponseWriter, request *http.Request) {
	var h = call.Utils.GetHTTP(request)

	// get id from params.
	var id = h.GetRouteArgs()["id"]

	// find.
	var article = model.Article{ID: id}
	found, err := article.FindByID()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	if !found {
		h.Send(a.throw.NotFound(), 404, err)
		return
	}

	// delete.
	err = article.DeleteByID()
	if err != nil {
		h.Send(a.throw.Server(), 500, err)
		return
	}
	h.Send("", 200, err)
}
