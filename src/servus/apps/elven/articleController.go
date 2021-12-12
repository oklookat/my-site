package elven

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"servus/core/external/errorMan"
)

const articlesPageSize = 2

// getAll - GET url/. Get paginated articles.
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
func (a *articleRoute) getAll(response http.ResponseWriter, request *http.Request) {
	var h = a.middleware.getHTTP(request)
	var err error
	var isAdmin = false
	var pipe = AuthPipe{}
	pipe.get(request)
	isAdmin = pipe.UserAndTokenExists && pipe.IsAdmin
	// validate query params.
	val, em, _ := a.validate.controllerGetAll(request, isAdmin)
	if em.HasErrors() {
		h.Send(em.GetJSON(), 400, err)
		return
	}
	// get articles by query params.
	articles, pages, err := val.getAll()
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
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
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	h.Send(string(jsonResponse), 200, err)
}

// getOne - GET url/id. Get one article.
func (a *articleRoute) getOne(response http.ResponseWriter, request *http.Request) {
	var h = a.middleware.getHTTP(request)
	var isAdmin = false
	var params = mux.Vars(request)
	var id = params["id"]
	var article = ArticleModel{ID: id}
	found, err := article.findByID()
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	if !found {
		h.Send(errorMan.ThrowNotFound(), 404, err)
		return
	}
	var pAuth = AuthPipe{}
	pAuth.get(request)
	isAdmin = pAuth.IsAdmin
	if !article.IsPublished && !isAdmin {
		h.Send(errorMan.ThrowForbidden(), 403, err)
		return
	}
	articleJson, err := json.Marshal(article)
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	h.Send(string(articleJson), 200, err)
}

// create - POST url/. Creates new article.
func (a *articleRoute) create(response http.ResponseWriter, request *http.Request) {
	var h = a.middleware.getHTTP(request)
	val, em, _ := a.validate.body(request)
	if em.HasErrors() {
		h.Send(em.GetJSON(), 400, nil)
		return
	}
	var pAuth = AuthPipe{}
	pAuth.get(request)
	var article = ArticleModel{UserID: pAuth.User.ID, IsPublished: false, Title: *val.Title, Content: *val.Content}
	err := article.create()
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	articleJson, err := json.Marshal(&article)
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	h.Send(string(articleJson), 200, err)
}

// update - PUT (update all available fields) or PATCH (update specific field) url/id.
func (a *articleRoute) update(response http.ResponseWriter, request *http.Request) {
	var h = a.middleware.getHTTP(request)
	var params = mux.Vars(request)
	var id = params["id"]
	var article = ArticleModel{ID: id}
	found, err := article.findByID()
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	if !found {
		h.Send(errorMan.ThrowNotFound(), 404, err)
		return
	}
	body, em, err := a.validate.body(request)
	if em.HasErrors() {
		h.Send(em.GetJSON(), 400, err)
		return
	}
	var isTitle = body.Title != nil
	var isContent = body.Content != nil
	var isPublished = body.IsPublished != nil
	// if PUT method we need full article to update. If PATCH - we need at least one field.
	if request.Method == http.MethodPut {
		if !isTitle || !isContent || !isPublished {
			em.Add("body", "provide all values.")
			h.Send(em.GetJSON(), 400, err)
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
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	jsonArticle, err := json.Marshal(article)
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	h.Send(string(jsonArticle), 200, err)
}

// delete - DELETE url/id. Deletes one article.
func (a *articleRoute) delete(response http.ResponseWriter, request *http.Request) {
	var h = a.middleware.getHTTP(request)
	var params = mux.Vars(request)
	var id = params["id"]
	var article = ArticleModel{ID: id}
	found, err := article.findByID()
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	if !found {
		h.Send(errorMan.ThrowNotFound(), 404, err)
		return
	}
	err = article.deleteByID()
	if err != nil {
		h.Send(errorMan.ThrowServer(), 500, err)
		return
	}
	h.Send("", 200, err)
	return
}
