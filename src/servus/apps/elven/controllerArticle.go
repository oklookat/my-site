package elven

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"servus/core"
	"servus/core/modules/errorCollector"
)

const articlesPageSize = 2

type ControllerArticlesPostBody struct {
	Title   string      `json:"title"`
	Content struct {
		Time   int64 `json:"time"`
		Blocks []struct {
			ID   string `json:"id"`
			Type string `json:"type"`
			Data interface{} `json:"data"`
		} `json:"blocks"`
		Version string `json:"version"`
	} `json:"content"`
}


// GET url/
// params:
// cursor = article id
// show = published, drafts, all
// by = created, updated, published
// start = newest (DESC), oldest (ASC)
// preview = true (content < 480 symbols), false (gives you full articles).
func controllerArticlesGetAll(response http.ResponseWriter, request *http.Request) {
	var err error
	var ec = errorCollector.New()
	var theResponse = core.HttpResponse{ResponseWriter: response}
	var authData = getAuthData(request)
	var isAdmin = false
	if authData != nil {
		isAdmin = authData.User.Role == "admin"
	}
	// validate query params
	var validated = validatedArticlesGetAll{}
	validated, err = validatorArticlesGetAll(request, isAdmin)
	if err != nil {
		theResponse.Send(err.Error(), 400)
		return
	}
	// get articles based on query params
	// TODO: fix 500 when not articles found
	articles, err := dbArticlesGetDependingOnValidated(validated)
	if err != nil {
		core.Logger.Error(fmt.Sprintf("articles get error: %v", err.Error()))
		ec.AddEUnknown([]string{"articles"}, "error while getting articles.")
		theResponse.Send(ec.GetErrors(), 500)
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
	// make json
	jsonResponse, err := json.Marshal(&responseContent)
	if err != nil {
		core.Logger.Error(fmt.Sprintf("articles response json marshal error: %v", err.Error()))
		ec.AddEUnknown([]string{"articles"}, "error while getting articles.")
		theResponse.Send(ec.GetErrors(), 500)
		return
	}
	theResponse.Send(string(jsonResponse), 200)
	return
}

// GET url/id
func controllerArticlesGetOne(response http.ResponseWriter, request *http.Request) {
	var theResponse = core.HttpResponse{ResponseWriter: response}
	var ec = errorCollector.New()
	var params = mux.Vars(request)
	var id = params["id"]
	var article, err = dbArticleFind(id)
	if err != nil {
		core.Logger.Error(fmt.Sprintf("error while getting article: %v", err.Error()))
		ec.AddEUnknown([]string{"articles"}, "error while getting article.")
		theResponse.Send(ec.GetErrors(), 500)
		return
	}
	if article.ID == "" {
		ec.AddENotFound([]string{"article"})
		theResponse.Send(ec.GetErrors(), 404)
		return
	}
	articleJson, err := json.Marshal(article)
	if err != nil {
		core.Logger.Error(fmt.Sprintf("error while article marshal: %v", err.Error()))
		ec.AddEUnknown([]string{"articles"}, "error while getting article.")
		theResponse.Send(ec.GetErrors(), 500)
		return
	}
	theResponse.Send(string(articleJson), 200)
	return
}

// POST url/
func controllerArticlesPost(response http.ResponseWriter, request *http.Request) {
	var theResponse = core.HttpResponse{ResponseWriter: response}
	var ec = errorCollector.New()
	var authData = getAuthData(request)
	var postBody ControllerArticlesPostBody
	var err = json.NewDecoder(request.Body).Decode(&postBody)
	if err != nil {
		ec.AddEValidationAllowed([]string{"article"}, []string{"isPublished", "title", "content"})
		theResponse.Send(ec.GetErrors(), 400)
		return
	}
	err = validatorArticlesPost(&postBody)
	if err != nil {
		theResponse.Send(err.Error(), 400)
		return
	}
	contentJson, err := json.Marshal(postBody.Content)
	if err != nil {
		ec.AddEUnknown([]string{"articles"}, "error while creating article")
		theResponse.Send(ec.GetErrors(), 500)
		return
	}
	var article = ModelArticle{UserID: authData.User.ID, IsPublished: false, Title: postBody.Title, Content: JSON(contentJson)}
	newArticle, err := dbArticleCreate(article)
	if err != nil {
		ec.AddEUnknown([]string{"articles"}, "error while creating article")
		theResponse.Send(ec.GetErrors(), 500)
		return
	}
	articleJson, err := json.Marshal(&newArticle)
	if err != nil {
		ec.AddEUnknown([]string{"articles"}, "error while creating article")
		theResponse.Send(ec.GetErrors(), 500)
		return
	}
	theResponse.Send(string(articleJson), 200)
	return
}

// PUT url/id
func controllerArticlesPut(response http.ResponseWriter, request *http.Request) {

}

// DELETE url/id
func controllerArticlesDelete(response http.ResponseWriter, request *http.Request) {
	var theResponse = core.HttpResponse{ResponseWriter: response}
	var ec = errorCollector.New()
	var params = mux.Vars(request)
	var id = params["id"]
	err := dbArticleDelete(id)
	if err != nil {
		ec.AddENotFound([]string{"article"})
		theResponse.Send(ec.GetErrors(), 404)
		return
	}
	theResponse.Send("", 200)
	return
}
