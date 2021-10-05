package elven

import (
	"fmt"
	"servus/core"
	"time"
)


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


func dbArticleCreate(article ModelArticle) (new *ModelArticle, err error) {
	new = &ModelArticle{}
	var sql = `INSERT INTO articles (user_id, is_published, title, content, slug) VALUES (:user_id, :is_published, :title, :content, :slug) RETURNING *`
	dbArticleBeforeChangeHook(&article)
	row, err := core.Database.NamedQuery(sql, &article)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return nil, err
	}
	row.Next()
	err = row.StructScan(new)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return nil, err
	}
	_ = dbArticleAfterChangeHook(new)
	return new, err
}

func dbArticleUpdate(article *ModelArticle) error {
	dbArticleBeforeChangeHook(article)
	var sql = "UPDATE articles SET user_id=:user_id, is_published=:is_published, title=:title, content=:content, slug=:slug WHERE id=:id RETURNING *"
	row, err := core.Database.NamedQuery(sql, article)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return err
	}
	row.Next()
	err = row.StructScan(article)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return err
	}
	_ = dbArticleAfterChangeHook(article)
	return err
}

func dbArticleFind(id string) (found *ModelArticle, err error) {
	found = &ModelArticle{}
	var sql = "SELECT * FROM articles WHERE id=$1 LIMIT 1"
	err = core.Database.Get(found, sql, id)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return nil, err
	}
	return found, err
}

func dbArticleDelete(id string) error {
	var sql = "DELETE FROM articles WHERE id=$1"
	_, err := core.Database.Exec(sql, id)
	err = core.Utils.DBCheckError(err)
	return err
}

func dbArticlesGet(validated validatedArticlesGetAll) (articles []ModelArticle, err error) {
	articles = make([]ModelArticle, 0)
	var query string
	switch validated.show {
	case "published":
		// use sprintf to format validated.by and start because with $ database throws syntax error (I don't know why)
		// it not allows sql injection, because start and by checked in validator
		query = fmt.Sprintf("SELECT * FROM articles WHERE id >= $1 AND is_published=true ORDER BY %v %v, id %v LIMIT $2 + 1", validated.by, validated.start, validated.start)
		break
	case "drafts":
		query = fmt.Sprintf("SELECT * FROM articles WHERE id >= $1 AND is_published=false ORDER BY %v %v, id %v LIMIT $2 + 1", validated.by, validated.start, validated.start)
		break
	case "all":
		query = fmt.Sprintf("SELECT * FROM articles WHERE id >= $1 ORDER BY %v %v, id %v LIMIT $2 + 1", validated.by, validated.start, validated.start)
		break
	}
	err = core.Database.Select(&articles, query, validated.cursor, articlesPageSize)
	err = core.Utils.DBCheckError(err)
	return articles, err
}
