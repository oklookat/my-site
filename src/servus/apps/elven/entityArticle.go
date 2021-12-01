package elven

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"servus/core/external/errorMan"
)

const articlesPageSize = 2

// entityArticle - manage articles.
type entityArticle struct {
	*entityBase
}

// getAll - GET url/
//
// request params:
//
// page: *number of page*
//
// show: published, drafts
//
// by: created; updated; published
//
// start: newest (DESC); oldest (ASC)
//
// preview: true (content < 480 symbols); false (gives you full articles).
func (a *entityArticle) getAll(response http.ResponseWriter, request *http.Request) {
	var err error
	var isAdmin = false
	var auth = PipeAuth{}
	auth.get(request)
	isAdmin = auth.UserAndTokenExists && auth.IsAdmin
	// validate query params.
	val, em, _ := a.validatorControllerGetAll(request, isAdmin)
	if em.HasErrors() {
		a.Send(response, em.GetJSON(), 400)
		return
	}
	// get articles by query params.
	articles, pages, err := val.getAll()
	if err != nil {
		instance.Logger.Error(fmt.Sprintf("articles get error: %v", err.Error()))
		a.Send(response, errorMan.ThrowServer(), 500)
		return
	}
	// generate response with pagination
	var responseContent = ResponseContent{}
	responseContent.Meta.CurrentPage = val.page
	responseContent.Meta.TotalPages = pages
	responseContent.Meta.PerPage = articlesPageSize
	responseContent.Data = articles
	// make json.
	jsonResponse, err := json.Marshal(&responseContent)
	if err != nil {
		instance.Logger.Error(fmt.Sprintf("articles response json marshal error: %v", err.Error()))
		a.Send(response, errorMan.ThrowServer(), 500)
		return
	}
	a.Send(response, string(jsonResponse), 200)
}

// getOne - GET url/id.
func (a *entityArticle) getOne(response http.ResponseWriter, request *http.Request) {
	var isAdmin = false
	var auth = PipeAuth{}
	auth.get(request)
	isAdmin = auth.UserAndTokenExists && auth.IsAdmin
	var params = mux.Vars(request)
	var id = params["id"]
	var article = ModelArticle{ID: id}
	found, err := article.findByID()
	if err != nil {
		a.err500(response, request, err)
		return
	}
	if !found {
		a.Send(response, errorMan.ThrowNotFound(), 404)
		return
	}
	if !article.IsPublished && !isAdmin {
		a.err403(response)
		return
	}
	articleJson, err := json.Marshal(article)
	if err != nil {
		a.err500(response, request, err)
		return
	}
	a.Send(response, string(articleJson), 200)
}

// create - POST url/.
func (a *entityArticle) create(response http.ResponseWriter, request *http.Request) {
	var pAuth = PipeAuth{}
	pAuth.get(request)
	val, em, _ := a.validatorBody(request)
	if em.HasErrors() {
		a.Send(response, em.GetJSON(), 400)
		return
	}
	var article = ModelArticle{UserID: pAuth.User.ID, IsPublished: false, Title: *val.Title, Content: *val.Content}
	err := article.create()
	if err != nil {
		a.err500(response, request, err)
		return
	}
	articleJson, err := json.Marshal(&article)
	if err != nil {
		a.err500(response, request, err)
		return
	}
	a.Send(response, string(articleJson), 200)
}

// update - PUT (update all available fields) or PATCH (update specific field) url/id.
func (a *entityArticle) update(response http.ResponseWriter, request *http.Request) {
	var params = mux.Vars(request)
	var id = params["id"]
	var article = ModelArticle{ID: id}
	found, err := article.findByID()
	if err != nil {
		a.err500(response, request, err)
		return
	}
	if !found {
		a.Send(response, errorMan.ThrowNotFound(), 404)
		return
	}
	body, em, err := a.validatorBody(request)
	if em.HasErrors() {
		a.Send(response, em.GetJSON(), 400)
		return
	}
	var isTitle = body.Title != nil
	var isContent = body.Content != nil
	var isPublished = body.IsPublished != nil
	// if PUT method we need full article to update. If PATCH - we need at least one field.
	if request.Method == http.MethodPut {
		if !isTitle || !isContent || !isPublished {
			em.Add("body", "provide all values.")
			a.Send(response, em.GetJSON(), 400)
			return
		} else {
			article.Title = *body.Title
			article.Content = *body.Content
			article.IsPublished = *body.IsPublished
		}
	} else if request.Method == http.MethodPatch {
		if isTitle {
			article.Title = *body.Title
		}
		if isContent {
			article.Content = *body.Content
		}
		if isPublished {
			article.IsPublished = *body.IsPublished
		}
	}
	err = article.update()
	if err != nil {
		a.err500(response, request, err)
		return
	}
	jsonArticle, err := json.Marshal(article)
	if err != nil {
		a.err500(response, request, err)
		return
	}
	a.Send(response, string(jsonArticle), 200)
}

// delete - DELETE url/id.
func (a *entityArticle) delete(response http.ResponseWriter, request *http.Request) {
	var params = mux.Vars(request)
	var id = params["id"]
	var article = ModelArticle{ID: id}
	found, err := article.findByID()
	if err != nil {
		a.err500(response, request, err)
		return
	}
	if !found {
		a.Send(response, errorMan.ThrowNotFound(), 404)
		return
	}
	err = article.deleteByID()
	if err != nil {
		a.err500(response, request, err)
		return
	}
	a.Send(response, "", 200)
	return
}

// err500 - write error to logger and send 500 error to user.
func (a *entityArticle) err500(response http.ResponseWriter, request *http.Request, err error) {
	instance.Logger.Warn(fmt.Sprintf("entityArticle code 500 at: %v. Error: %v", request.URL.Path, err.Error()))
	a.Send(response, errorMan.ThrowServer(), 500)
	return
}

// err403 - send an error if the user is not allowed to do something.
func (a *entityArticle) err403(response http.ResponseWriter) {
	a.Send(response, errorMan.ThrowForbidden(), 403)
	return
}
