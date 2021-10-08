package elven

import (
	"database/sql"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	"math/rand"
	"servus/core"
	"servus/core/modules/cryptor"
	"time"
)

// databaseGetAll - get ModelArticles by request and GetAllQuery.
func (a *entityArticle) databaseGetAll(validated *queryArticleControllerGetAll) (articles []ModelArticle, err error) {
	if validated == nil {
		return nil, errors.New("databaseGetAll: getAllQuery nil pointer.")
	}
	articles = make([]ModelArticle, 0)
	var query string
	var by = validated.by
	var start = validated.start
	var cursor = validated.cursor
	switch validated.show {
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

func (u *entityUser) databaseCreate(user *ModelUser) (err error) {
	hashedPassword, err := cryptor.BHash(user.Password)
	if err != nil {
		core.Logger.Error(err.Error())
		return
	}
	var query = `INSERT INTO users (role, username, password) VALUES ($1, $2, $3) RETURNING *`
	row := core.Database.QueryRowx(query, user.Role, user.Username, hashedPassword)
	err = row.StructScan(user)
	err = core.Utils.DBCheckError(err)
	if err == sql.ErrNoRows {
		return nil
	}
	return
}

func (u *entityUser) databaseFind(id string) (found *ModelUser, err error) {
	found = &ModelUser{}
	var query = "SELECT * FROM users WHERE id=$1 LIMIT 1"
	row := core.Database.QueryRowx(query, id)
	err = row.StructScan(found)
	err = core.Utils.DBCheckError(err)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

func (u *entityUser) databaseFindBy(username string) (found *ModelUser, err error) {
	found = &ModelUser{}
	var query = "SELECT * FROM users WHERE username=$1 LIMIT 1"
	row := core.Database.QueryRowx(query, username)
	err = row.StructScan(found)
	err = core.Utils.DBCheckError(err)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

func (u *entityUser) databaseDelete(id string) (err error) {
	var query = "DELETE FROM users WHERE id=$1"
	_, err = core.Database.Exec(query, id)
	err = core.Utils.DBCheckError(err)
	if err == sql.ErrNoRows {
		return nil
	}
	return
}
