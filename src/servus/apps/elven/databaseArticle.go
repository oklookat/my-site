package elven

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/oklog/ulid/v2"
	"math/rand"
	"servus/core"
	"time"
)

type ModelArticle struct {
	ID          string `json:"id" db:"id"`
	UserID      string `json:"userID" db:"user_id"`
	IsPublished bool   `json:"isPublished" db:"is_published"`
	Title       string `json:"title" db:"title"`
	Content     string `json:"content" db:"content"`
	Slug        string     `json:"slug" db:"slug"`
	PublishedAt *time.Time `json:"publishedAt" db:"published_at"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
}

func dbArticleBeforeChangeHook(article *ModelArticle) {
	// create temp slug (ULID)
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	article.Slug = ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

func dbArticleAfterChangeHook(article *ModelArticle) (err error) {
	// create normal slug
	article.Slug = slug.Make(article.Title) + "-" + article.ID
	var sql = "UPDATE articles SET slug=:slug WHERE id=:id RETURNING *"
	row, err := core.Database.Connection.NamedQuery(sql, &article)
	if err != nil {
		return err
	}
	err = row.StructScan(&article)
	if err != nil {
		return err
	}
	return nil
}

func dbArticleCreate(article ModelArticle) (new *ModelArticle, err error) {
	new = &ModelArticle{}
	var sql = `INSERT INTO articles (user_id, is_published, title, content, slug) VALUES (:user_id, :is_published, :title, :content, :slug) RETURNING *`
	dbArticleBeforeChangeHook(&article)
	row, err := core.Database.Connection.NamedQuery(sql, &article)
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

func dbArticleUpdate(article ModelArticle) (updated *ModelArticle, err error) {
	updated = &ModelArticle{}
	dbArticleBeforeChangeHook(&article)
	var sql = "UPDATE articles SET user_id=:user_id, is_published=:is_published, title=:title, content=:content, slug=:slug WHERE id=:id RETURNING *"
	row, err := core.Database.Connection.NamedQuery(sql, article)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return nil, err
	}
	row.Next()
	err = row.StructScan(updated)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return nil, err
	}
	_ = dbArticleAfterChangeHook(updated)
	return updated, err
}

func dbArticleFind(id string) (found *ModelArticle, err error) {
	found = &ModelArticle{}
	var sql = "SELECT * FROM articles WHERE id=$1 LIMIT 1"
	err = core.Database.Connection.Select(&found, sql, id)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return nil, err
	}
	return found, err
}

func dbArticleDelete(id string) error {
	var sql = "DELETE FROM articles WHERE id=$1"
	_, err := core.Database.Connection.Exec(sql, id)
	err = core.Utils.DBCheckError(err)
	return err
}

func dbArticlesGetDependingOnValidated(validated validatedArticlesGetAll) (articles []ModelArticle, err error){
	articles = []ModelArticle{}
	switch validated.show {
	case "published":
		// use sprintf to format validated.by and start because with $ database throws syntax error (I don't know why)
		// it not allows sql injection, because start and by checked in validator
		var sql = fmt.Sprintf("SELECT * FROM articles WHERE id >= $1 AND is_published=true ORDER BY %v %v, id %v LIMIT $2 + 1", validated.by, validated.start, validated.start)
		err = core.Database.Connection.Select(&articles, sql, validated.cursor, articlesPageSize)
		break
	case "drafts":
		var sql = fmt.Sprintf("SELECT * FROM articles WHERE id >= $1 AND is_published=false ORDER BY %v %v, id %v LIMIT $2 + 1", validated.by, validated.start, validated.start)
		err = core.Database.Connection.Select(&articles, sql, validated.cursor, articlesPageSize)
		break
	case "all":
		var sql = fmt.Sprintf("SELECT * FROM articles WHERE id >= $1 ORDER BY %v %v, id %v LIMIT $2 + 1", validated.by, validated.start, validated.start)
		err = core.Database.Connection.Select(&articles, sql, validated.cursor, articlesPageSize)
		break
	}
	err = core.Utils.DBCheckError(err)
	return articles, err
}