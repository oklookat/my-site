package elven

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"servus/core"
	"strings"
	"time"
)

type ModelArticle struct {
	ID string `json:"id"`
	UserID string `json:"userID"`
	IsPublished bool `json:"isPublished"`
	Title string `json:"title"`
	Content string `json:"content"`
	Slug string `json:"slug"`
	PublishedAt *time.Time `json:"publishedAt"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// dbArticleScanRow - scan row or rows to article model.
func dbArticleScanRow(row pgx.Row, article *ModelArticle) (err error){
	err = row.Scan(&article.ID, &article.UserID, &article.IsPublished, &article.Title, &article.Content, &article.Slug, &article.PublishedAt, &article.CreatedAt, &article.UpdatedAt)
	err = core.Utils.DBCheckError(err)
	return err
}

func dbArticleCreate(article ModelArticle) (new ModelArticle, err error){
	new = ModelArticle{}
	var sql = `INSERT INTO articles (user_id, is_published, title, content, slug) VALUES ($1, $2, $3, $4, $5) RETURNING *`
	row := core.Database.Connection.QueryRow(context.Background(), sql, article.UserID, article.IsPublished, article.Title, article.Content, article.Slug)
	err = dbArticleScanRow(row, &new)
	return new, err
}

func dbArticleUpdate(article ModelArticle) (updated ModelArticle, err error){
	updated = ModelArticle{}
	var sql = "UPDATE articles SET user_id=$1, is_published=$2, title=$3, content=$4, slug=$5 WHERE id=$6 RETURNING *"
	row := core.Database.Connection.QueryRow(context.Background(), sql, &article.UserID, &article.IsPublished, &article.Title, &article.Content, &article.Slug)
	err = dbArticleScanRow(row, &updated)
	return updated, err
}

func dbArticleFind(id string) (found ModelArticle, err error){
	found = ModelArticle{}
	var sql = "SELECT * FROM articles WHERE id=$1 LIMIT 1"
	row := core.Database.Connection.QueryRow(context.Background(), sql, id)
	err = dbArticleScanRow(row, &found)
	return found, err
}

// dbArticlesGet - get articles by values in params argument.
func dbArticlesGet(params *ModelArticle) (found ModelArticle, err error){
	var wheres []string
	var wheresValues []interface{}
	var counter = 0
	// append field with counter value
	// like [title=$1, content=$2].
	var appendCounted = func(fieldName string, fieldValue interface{}) {
		counter++
		formatted := fmt.Sprintf("%v=$%v", fieldName, counter)
		wheres = append(wheres, formatted)
		wheresValues = append(wheresValues, fieldValue)
	}
	if params.Title != "" {
		appendCounted("title", params.Title)
	}
	if params.Content != "" {
		appendCounted("content", params.Content)
	}
	if params.Slug != "" {
		appendCounted("slug", params.Slug)
	}
	found = ModelArticle{}
	var sql = fmt.Sprintf("SELECT * FROM articles WHERE %v LIMIT 1", strings.Join(wheres, " "))
	row := core.Database.Connection.QueryRow(context.Background(), sql, wheresValues...)
	err = dbArticleScanRow(row, &found)
	return found, err
}

func dbArticleDelete(id string) error {
	var sql = "DELETE FROM articles WHERE id=$1"
	query, err := core.Database.Connection.Query(context.Background(), sql, id)
	defer query.Close()
	err = core.Utils.DBCheckError(err)
	return err
}
