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
	"strconv"
	"strings"
	"time"
)

const articlesPageSize = 2

// entityArticle - manage the articles.
type entityArticle struct {
	*entityBase
	queryControllerGetAll *queryArticleControllerGetAll
	bodyChange            *BodyArticleChange
}

// BodyArticleChange - represents the body of the request that the user should send. Used in create and update methods.
type BodyArticleChange struct {
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
	_ = a.validatorGetAll(request, isAdmin)
	if a.EC.HasErrors() {
		a.Send(response, a.EC.GetErrors(), 400)
		return
	}
	// get articles based on query params.
	articles, err := a.databaseGetAll()
	if err != nil {
		core.Logger.Error(fmt.Sprintf("articles get error: %v", err.Error()))
		a.EC.AddEUnknown([]string{"articles"}, "error while getting articles.")
		a.Send(response, a.EC.GetErrors(), 500)
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
		a.EC.AddEUnknown([]string{"articles"}, "error while getting articles.")
		a.Send(response, a.EC.GetErrors(), 500)
		return
	}
	a.Send(response, string(jsonResponse), 200)
}

// controllerGetOne - GET url/id.
func (a *entityArticle) controllerGetOne(response http.ResponseWriter, request *http.Request) {
	isAdmin := a.isAdmin(request)
	var params = mux.Vars(request)
	var id = params["id"]
	var article, err = a.databaseFind(id)
	if err != nil {
		a.err500(response, request, err)
		return
	}
	if article == nil {
		a.EC.AddENotFound([]string{"article"})
		a.Send(response, a.EC.GetErrors(), 404)
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
	var authData = oUtil.getPipeAuth(request)
	a.bodyChange = &BodyArticleChange{}
	var err = json.NewDecoder(request.Body).Decode(a.bodyChange)
	if err != nil {
		a.EC.AddEValidationAllowed([]string{"article"}, []string{"title", "content"})
		a.Send(response, a.EC.GetErrors(), 400)
		return
	}
	a.validatorBodyChange(a.bodyChange)
	if a.EC.HasErrors() {
		a.Send(response, a.EC.GetErrors(), 400)
		return
	}
	var article = ModelArticle{UserID: authData.User.ID, IsPublished: false, Title: a.bodyChange.Title, Content: a.bodyChange.Content}
	err = a.databaseCreate(&article)
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
	//var params = mux.Vars(request)
	a.bodyChange = &BodyArticleChange{}
	err := json.NewDecoder(request.Body).Decode(a.bodyChange)
	if err != nil {
		a.EC.AddEValidationAllowed([]string{"article"}, []string{"isPublished", "title", "content"})
		a.Send(response, a.EC.GetErrors(), 400)
		return
	}
	a.validatorBodyChange(a.bodyChange)
	if a.EC.HasErrors() {
		a.Send(response, a.EC.GetErrors(), 400)
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
		a.EC.AddENotFound([]string{"article"})
		a.Send(response, a.EC.GetErrors(), 404)
		return
	}
	article.Title = a.bodyChange.Title
	article.Content = a.bodyChange.Content
	if a.bodyChange.IsPublished != nil {
		article.IsPublished = *a.bodyChange.IsPublished
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
		a.EC.AddENotFound([]string{"article"})
		a.Send(response, a.EC.GetErrors(), 404)
		return
	}
	a.Send(response, "", 200)
}

// err500 - write error to logger and send 500 error to user.
func (a *entityArticle) err500(response http.ResponseWriter, request *http.Request, err error) {
	a.Logger.Warn("entityArticle code 500 at: %v. Error: %v", request.URL.Path, err.Error())
	a.EC.AddEUnknown([]string{"articles"}, "server error")
	a.Send(response, a.EC.GetErrors(), 500)
	return
}

// err403 - send an error if the user is not allowed to do something.
func (a *entityArticle) err403(response http.ResponseWriter){
	a.EC.AddEAuthForbidden([]string{"article"})
	a.Send(response, a.EC.GetErrors(), 403)
	return
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

// validatorGetAll - validate query params in request depending to ModelArticle
// if validation error - returns errorCollector JSON (err.Error()).
func (a *entityArticle) validatorGetAll(request *http.Request, isAdmin bool) (err error) {
	a.queryControllerGetAll = &queryArticleControllerGetAll{}
	var queryParams = request.URL.Query()
	// validate "show" param
	var show = queryParams.Get("show")
	if len(show) == 0 {
		show = "published"
	} else {
		strings.ToLower(show)
	}
	var isShowInvalid = show != "published" && show != "drafts" && show != "all"
	switch isShowInvalid {
	case true:
		a.EC.AddEValidationAllowed([]string{"show"}, []string{"published", "drafts", "all"})
		break
	case false:
		if (show == "drafts" || show == "all") && !isAdmin {
			a.EC.AddEAuthForbidden([]string{"show"})
		}
		break
	}
	a.queryControllerGetAll.show = show
	// validate "by" param
	var by = queryParams.Get("by")
	if len(by) == 0 {
		by = "published"
	} else {
		strings.ToLower(by)
	}
	var isByInvalid = by != "created" && by != "published" && by != "updated"
	if isByInvalid {
		a.EC.AddEValidationAllowed([]string{"by"}, []string{"created", "published", "updated"})
	} else if (by == "updated" || by == "created") && !isAdmin{
		a.EC.AddEAuthForbidden([]string{"by"})
	}
	switch by {
	case "created":
		by = "created_at"
		break
	case "updated":
		by = "updated_at"
		break
	case "published":
		by = "published_at"
		break
	}
	a.queryControllerGetAll.by = by
	// validate "start" param
	var start = queryParams.Get("start")
	if len(start) == 0 {
		start = "newest"
	} else {
		strings.ToLower(start)
	}
	var isStartInvalid = start != "newest" && start != "oldest"
	if isStartInvalid {
		a.EC.AddEValidationAllowed([]string{"start"}, []string{"newest", "oldest"})
	} else {
		switch start {
		case "newest":
			start = "DESC"
			break
		case "oldest":
			start = "ASC"
			break
		}
	}
	a.queryControllerGetAll.start = start
	// validate "preview" param
	var preview = queryParams.Get("preview")
	if len(preview) == 0 {
		preview = "true"
	}
	var previewBool bool
	previewBool, err = strconv.ParseBool(preview)
	if err != nil {
		a.EC.AddEValidationAllowed([]string{"preview"}, []string{"bool"})
		previewBool = true
	}
	a.queryControllerGetAll.preview = previewBool
	// validate "cursor" param
	var cursor = queryParams.Get("cursor")
	if len(cursor) == 0 {
		cursor = "0"
	}
	a.queryControllerGetAll.cursor = cursor
	// finally
	return
}

// validatorBodyChange - validate BodyArticleChange. Writes result in errorCollector instance.
func (a *entityArticle) validatorBodyChange(body *BodyArticleChange) {
	if len(body.Title) > 124 {
		a.EC.AddEValidationMinMax([]string{"title"}, 1, 124)
	}
}

// databaseGetAll - get ModelArticles by request and GetAllQuery.
func (a *entityArticle) databaseGetAll() (articles []ModelArticle, err error) {
	if a.queryControllerGetAll == nil {
		return nil, errors.New("databaseGetAll: getAllQuery nil pointer.")
	}
	articles = make([]ModelArticle, 0)
	var query string
	var by = a.queryControllerGetAll.by
	var start = a.queryControllerGetAll.start
	var cursor = a.queryControllerGetAll.cursor
	switch a.queryControllerGetAll.show {
	case "published":
		// use sprintf to format validated.by and start because with $ database throws syntax error (I don't know why)
		// it not allows sql injection, because start and by checked in validator
		query = fmt.Sprintf("SELECT * FROM articles WHERE id >= $1 AND is_published=true ORDER BY %v %v, id %v LIMIT $2 + 1", by, start, start)
		break
	case "drafts":
		query = fmt.Sprintf("SELECT * FROM articles WHERE id >= $1 AND is_published=false ORDER BY %v %v, id %v LIMIT $2 + 1", by, start, start)
		break
	case "all":
		query = fmt.Sprintf("SELECT * FROM articles WHERE id >= $1 ORDER BY %v %v, id %v LIMIT $2 + 1", by, start, start)
		break
	}
	err = core.Database.Select(&articles, query, cursor, articlesPageSize)
	err = core.Utils.DBCheckError(err)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

// databaseCreate - create ModelArticle.
func (a *entityArticle) databaseCreate(article *ModelArticle) (err error) {
	var query = `INSERT INTO articles (user_id, is_published, title, content, slug) VALUES (:user_id, :is_published, :title, :content, :slug) RETURNING *`
	a.databaseBeforeChange(article)
	row, err := core.Database.NamedQuery(query, &article)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return err
	}
	row.Next()
	err = row.StructScan(article)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return
	}
	_ = a.databaseAfterChange(article)
	return
}

// databaseUpdate - update ModelArticle.
func (a *entityArticle) databaseUpdate(article *ModelArticle) (err error) {
	a.databaseBeforeChange(article)
	var query = "UPDATE articles SET user_id=:user_id, is_published=:is_published, title=:title, content=:content, slug=:slug WHERE id=:id RETURNING *"
	row, err := core.Database.NamedQuery(query, article)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return
	}
	row.Next()
	err = row.StructScan(article)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return
	}
	_ = a.databaseAfterChange(article)
	return
}

// databaseFind - find ModelArticle.
func (a *entityArticle) databaseFind(id string) (found *ModelArticle, err error) {
	found = &ModelArticle{}
	var query = "SELECT * FROM articles WHERE id=$1 LIMIT 1"
	err = core.Database.Get(found, query, id)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return nil, err
	}
	return
}

// databaseDelete - delete ModelArticle.
func (a *entityArticle) databaseDelete(id string) (err error) {
	var query = "DELETE FROM articles WHERE id=$1"
	_, err = core.Database.Exec(query, id)
	err = core.Utils.DBCheckError(err)
	return
}

// databaseBeforeChange - executes before article create or update.
func (a *entityArticle) databaseBeforeChange(article *ModelArticle) {
	if len(article.Title) == 0 {
		article.Title = "Untitled"
	}
	// create temp slug (ULID)
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	article.Slug = ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

// databaseBeforeChange - executes after article create or update.
func (a *entityArticle) databaseAfterChange(article *ModelArticle) (err error) {
	// create normal slug
	article.Slug = slug.Make(article.Title) + "-" + article.ID
	var query = "UPDATE articles SET slug=:slug WHERE id=:id RETURNING *"
	row, err := core.Database.NamedQuery(query, &article)
	if err != nil {
		return err
	}
	err = row.StructScan(&article)
	if err != nil {
		return err
	}
	return nil
}

