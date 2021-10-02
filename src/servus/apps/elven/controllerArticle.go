package elven

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"net/http"
	"servus/core"
	"servus/core/modules/errorCollector"
	"servus/core/modules/routerica"
)

// GET url/
// params:
// page = number
// show = published, drafts, all
// by = created, updated, published
// start = newest (DESC), oldest (ASC)
// preview = true (content < 480 symbols), false (gives you full articles)

const articlesPageSize = 2

func controllerArticlesGetAll(response http.ResponseWriter, request *http.Request) {
	var err error
	var ec = errorCollector.New()
	var theResponse = core.HttpResponse{ResponseWriter: response}
	var authData = getAuthData(request)
	var isAdmin = false
	if authData != nil {
		isAdmin = authData.User.Role == "admin"
	}
	var validated = validatedArticleQuery{}
	validated, err = validatorArticleQueryParams(request, isAdmin)
	if err != nil {
		theResponse.Send(err.Error(), 400)
		return
	}
	var rows pgx.Rows
	var whereID = ">="
	if validated.cursor != "0" && validated.start == "DESC"{
		whereID = "<="
	}
	switch validated.show {
	case "published":
		// use sprintf to format validated.by and start because with $ database throws syntax error (I don't know why)
		// it not allows sql injection, because start and by checked in validator
		var sql = fmt.Sprintf("SELECT * FROM articles WHERE id %v $1 AND is_published=true ORDER BY %v %v, id %v LIMIT $2 + 1", whereID, validated.by, validated.start, validated.start)
		rows, err = core.Database.Connection.Query(context.Background(), sql, validated.cursor, articlesPageSize)
		break
	case "drafts":
		var sql = fmt.Sprintf("SELECT * FROM articles WHERE id %v $1 AND is_published=false ORDER BY %v %v, id %v LIMIT $2 + 1", whereID, validated.by, validated.start, validated.start)
		rows, err = core.Database.Connection.Query(context.Background(), sql, validated.cursor, articlesPageSize)
		break
	case "all":
		var sql = fmt.Sprintf("SELECT * FROM articles WHERE id %v $1 ORDER BY %v %v, id %v LIMIT $2 + 1", whereID, validated.by, validated.start, validated.start)
		rows, err = core.Database.Connection.Query(context.Background(), sql, validated.cursor, articlesPageSize)
		break
	}
	defer rows.Close()
	var articles []ModelArticle
	var rowCounter = 0
	var responseContent = ResponseContent{Content: articles}
	responseContent.Cursor.PerPage = articlesPageSize
	for rows.Next(){
		article := ModelArticle{}
		err := dbArticleScanRow(rows, &article)
		if err != nil {
			core.Logger.Error(fmt.Sprintf("article scan rows error: %v", err.Error()))
			continue
		}
		if rowCounter >= articlesPageSize {
			responseContent.Cursor.Next = article.ID
			break
		}
		articles = append(articles, article)
		rowCounter++
	}
	responseContent.Content = articles
	jsonResponse, err := json.Marshal(responseContent)
	if err != nil {
		core.Logger.Error(fmt.Sprintf("articles response json marshal error: %v", err.Error()))
		ec.AddEUnknown([]string{"articles"}, "error while getting articles.")
		theResponse.Send(ec.GetErrors(), 500)
		return
	}
	theResponse.Send(string(jsonResponse), 200)
	return
}

func controllerArticlesGetOne(response http.ResponseWriter, request *http.Request){
	var theResponse = core.HttpResponse{ResponseWriter: response}
	var ec = errorCollector.New()
	var params = routerica.GetParams(request)
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

}

// PUT /url/id
func controllerArticlesPut(response http.ResponseWriter, request *http.Request) {

}

// DELETE /url/id
func controllerArticlesDelete(response http.ResponseWriter, request *http.Request) {

}
