package elven

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	"math/rand"
	"net/http"
	"servus/core"
	"servus/core/modules/errorMan"
	"strconv"
	"strings"
	"time"
)

const articlesPageSize = 2

// entityArticle - manage the articles.
type entityArticle struct {
	*entityBase
}

// BodyArticle - represents the body of the request that the user should send. Used in create and update methods.
type BodyArticle struct {
	IsPublished *bool  `json:"is_published"`
	Title       string `json:"title"`
	Content     struct {
		Time   int64 `json:"time"`
		Blocks []struct {
			ID   string      `json:"id"`
			Type string      `json:"type"`
			Data interface{} `json:"data"`
		} `json:"blocks"`
		Version string `json:"version"`
	} `json:"content"`
}

// articleGetAllQuery - parsed query in the article.controllerGetAll.
type queryArticleControllerGetAll struct {
	cursor  string
	show    string
	by      string
	start   string
	preview bool
}

// ModelArticle - represents the article in database.
type ModelArticle struct {
	ID          string         `json:"id" db:"id"`
	UserID      string         `json:"user_id" db:"user_id"`
	IsPublished bool           `json:"is_published" db:"is_published"`
	Title       string         `json:"title" db:"title"`
	Content     ArticleContent `json:"content" db:"content"`
	Slug        string         `json:"slug" db:"slug"`
	PublishedAt *time.Time     `json:"published_at" db:"published_at"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

// controllerGetAll - GET url/
// params:
// cursor = article id
// show = published, drafts, all
// by = created, updated, published
// start = newest (DESC), oldest (ASC)
// preview = true (content < 480 symbols), false (gives you full articles).
func (a *entityArticle) controllerGetAll(response http.ResponseWriter, request *http.Request) {
	var err error
	isAdmin := a.isAdmin(request)
	// validate query params.
	validated, _ := a.validatorControllerGetAll(request, ec, isAdmin)
	if ec.HasErrors() {
		a.Send(response, ec.GetErrors(), 400)
		return
	}
	// get articles based on query params.
	articles, err := a.databaseGetAll(&validated)
	if err != nil {
		core.Logger.Error(fmt.Sprintf("articles get error: %v", err.Error()))
		ec.AddEUnknown([]string{"articles"}, "error while getting articles.")
		a.Send(response, ec.GetErrors(), 500)
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
		ec.AddEUnknown([]string{"articles"}, "error while getting articles.")
		a.Send(response, ec.GetErrors(), 500)
		return
	}
	a.Send(response, string(jsonResponse), 200)
}

// controllerGetOne - GET url/id.
func (a *entityArticle) controllerGetOne(response http.ResponseWriter, request *http.Request) {
	ec := errorMan.New()
	isAdmin := a.isAdmin(request)
	var params = mux.Vars(request)
	var id = params["id"]
	var article, err = a.databaseFind(id)
	if err != nil {
		a.err500(response, request, ec, err)
		return
	}
	if article == nil {
		ec.AddENotFound([]string{"article"})
		a.Send(response, ec.GetErrors(), 404)
		return
	}
	if !article.IsPublished && !isAdmin {
		a.err403(response, ec)
		return
	}
	articleJson, err := json.Marshal(article)
	if err != nil {
		a.err500(response, request, ec, err)
		return
	}
	a.Send(response, string(articleJson), 200)
}

// controllerCreateOne - POST url/.
func (a *entityArticle) controllerCreateOne(response http.ResponseWriter, request *http.Request) {
	ec := errorMan.New()
	var authData = oUtil.getPipeAuth(request)
	bodyChange := &BodyArticle{}
	var err = json.NewDecoder(request.Body).Decode(bodyChange)
	if err != nil {
		ec.AddEValidationAllowed([]string{"article"}, []string{"title", "content"})
		a.Send(response, ec.GetErrors(), 400)
		return
	}
	a.validatorBodyChange(bodyChange, ec)
	if ec.HasErrors() {
		a.Send(response, ec.GetErrors(), 400)
		return
	}
	var article = ModelArticle{UserID: authData.User.ID, IsPublished: false, Title: bodyChange.Title, Content: bodyChange.Content}
	err = a.databaseCreate(&article)
	if err != nil {
		a.err500(response, request, ec, err)
		return
	}
	articleJson, err := json.Marshal(&article)
	if err != nil {
		a.err500(response, request, ec, err)
		return
	}
	a.Send(response, string(articleJson), 200)
}

// controllerUpdateOne - PUT url/id.
func (a *entityArticle) controllerUpdateOne(response http.ResponseWriter, request *http.Request) {
	body, em, err := a.validatorControllerUpdateOne(request)
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

// isAdmin - check is request by admin.
func (a *entityArticle) isAdmin(request *http.Request) (isAdmin bool){
	var authData = oUtil.getPipeAuth(request)
	isAdmin = false
	if authData != nil {
		isAdmin = authData.IsAdmin
	}
	return
}

