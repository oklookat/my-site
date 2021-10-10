package elven

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"servus/core"
	"servus/core/modules/errorMan"
)


// controllerGetAll - GET url/
// params:
// cursor = article id
// show = published, drafts, all
// by = created, updated, published
// start = newest (DESC), oldest (ASC)
// preview = true (content < 480 symbols), false - gives you full articles.
func (a *entityArticle) controllerGetAll(response http.ResponseWriter, request *http.Request) {
	var err error
	isAdmin := oUtils.isAdmin(request)
	// validate query params.
	val, em, _ := a.validatorControllerGetAll(request, isAdmin)
	if em.HasErrors() {
		a.Send(response, em.GetJSON(), 400)
		return
	}
	// get articles based on query params.
	articles, err := a.databaseGetAll(&val)
	if err != nil {
		core.Logger.Error(fmt.Sprintf("articles get error: %v", err.Error()))
		a.Send(response, errorMan.ThrowServer(), 500)
		return
	}
	// generate response with pagination
	var responseContent = ResponseContent{}
	responseContent.Meta.PerPage = articlesPageSize
	if len(articles) >= articlesPageSize {
		var lastElement = len(articles) - 1
		responseContent.Meta.Next = articles[lastElement].ID
		articles = articles[:lastElement]
	}
	responseContent.Data = articles
	// make json.
	jsonResponse, err := json.Marshal(&responseContent)
	if err != nil {
		core.Logger.Error(fmt.Sprintf("articles response json marshal error: %v", err.Error()))
		a.Send(response, errorMan.ThrowServer(), 500)
		return
	}
	a.Send(response, string(jsonResponse), 200)
}

// controllerGetOne - GET url/id.
func (a *entityArticle) controllerGetOne(response http.ResponseWriter, request *http.Request) {
	isAdmin := oUtils.isAdmin(request)
	var params = mux.Vars(request)
	var id = params["id"]
	var article, err = a.databaseFind(id)
	if err != nil {
		a.err500(response, request, err)
		return
	}
	if article == nil {
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

// controllerCreateOne - POST url/.
func (a *entityArticle) controllerCreateOne(response http.ResponseWriter, request *http.Request) {
	var authData = oUtils.getPipeAuth(request)
	val, em, _ := a.validatorBody(request)
	if em.HasErrors() {
		a.Send(response, em.GetJSON(), 400)
		return
	}
	var article = ModelArticle{UserID: authData.User.ID, IsPublished: false, Title: val.Title, Content: val.Content}
	err := a.databaseCreate(&article)
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

// controllerUpdateOne - PUT url/id.
func (a *entityArticle) controllerUpdateOne(response http.ResponseWriter, request *http.Request) {
	body, em, err := a.validatorBody(request)
	if em.HasErrors() {
		a.Send(response, em.GetJSON(), 400)
		return
	}
	var params = mux.Vars(request)
	var id = params["id"]
	article, err := a.databaseFind(id)
	if err != nil {
		a.err500(response, request, err)
		return
	}
	if article == nil {
		a.Send(response, errorMan.ThrowNotFound(), 404)
		return
	}
	article.Title = body.Title
	article.Content = body.Content
	if body.IsPublished != nil {
		article.IsPublished = *body.IsPublished
	}
	err = a.databaseUpdate(article)
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

// controllerDeleteOne - DELETE url/id.
func (a *entityArticle) controllerDeleteOne(response http.ResponseWriter, request *http.Request) {
	var params = mux.Vars(request)
	var id = params["id"]
	err := a.databaseDelete(id)
	if err != nil {
		a.Send(response, errorMan.ThrowNotFound(), 404)
		return
	}
	a.Send(response, "", 200)
}

