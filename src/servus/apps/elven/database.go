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

// databaseCreate - create user.
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

// databaseFind - find user.
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

// databaseFindBy - find user by username.
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

// databaseDelete - delete user.
func (u *entityUser) databaseDelete(id string) (err error) {
	var query = "DELETE FROM users WHERE id=$1"
	_, err = core.Database.Exec(query, id)
	err = core.Utils.DBCheckError(err)
	if err == sql.ErrNoRows {
		return nil
	}
	return
}

// databaseCreate - create ModelToken in database. ATTENTION: its function writes only Token and user id, other data will be ignored. For write full data see dbTokenUpdate.
func (t *entityToken) databaseCreate(token *ModelToken) (err error) {
	var query = "INSERT INTO tokens (user_id, token) VALUES ($1, $2) RETURNING *"
	row := core.Database.QueryRowx(query, &token.UserID, &token.Token)
	err = row.StructScan(token)
	err = core.Utils.DBCheckError(err)
	return
}

// databaseUpdate - updates ModelToken in database. All fields (except update and created dates) must be filled.
func (t *entityToken) databaseUpdate(token *ModelToken) (err error) {
	t.databaseBeforeUpdate(token)
	var query = "UPDATE tokens SET user_id=$1, token=$2, last_ip=$3, last_agent=$4, auth_ip=$5, auth_agent=$6 WHERE id=$7 RETURNING *"
	row := core.Database.QueryRowx(query, &token.UserID, &token.Token, &token.LastIP, &token.LastAgent, &token.AuthIP, &token.AuthAgent, &token.ID)
	err = row.StructScan(token)
	err = core.Utils.DBCheckError(err)
	return
}

// databaseFind - find ModelToken in database.
func (t *entityToken) databaseFind(id string) (found *ModelToken, err error) {
	found = &ModelToken{}
	var query = "SELECT * FROM tokens WHERE id=$1 LIMIT 1"
	row := core.Database.QueryRowx(query, id)
	err = row.StructScan(found)
	err = core.Utils.DBCheckError(err)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

// databaseDelete - delete ModelToken from database.
func (t *entityToken) databaseDelete(id string) (err error) {
	var query = "DELETE FROM tokens WHERE id=$1"
	_, err = core.Database.Exec(query, id)
	if err == sql.ErrNoRows {
		return nil
	}
	return
}

// databaseBeforeUpdate - executes before token update.
func (t *entityToken) databaseBeforeUpdate(token *ModelToken) {
	if token.AuthAgent != nil && len(*token.AuthAgent) > 323 {
		var authAgent = *token.AuthAgent
		var cut = 323 - len(authAgent)
		*token.AuthAgent = authAgent[:len(authAgent)-cut]
	}
	if token.LastAgent != nil && len(*token.LastAgent) > 323 {
		var lastAgent = *token.LastAgent
		var cut = 323 - len(lastAgent)
		*token.LastAgent = lastAgent[:len(lastAgent)-cut]
	}
	if token.AuthIP != nil && len(*token.AuthIP) > 53 {
		var authIP = *token.AuthIP
		var cut = 53 - len(authIP)
		*token.AuthIP = authIP[:len(authIP)-cut]
	}
	if token.LastIP != nil && len(*token.LastIP) > 53 {
		var lastIP = *token.LastIP
		var cut = 53 - len(lastIP)
		*token.LastIP = lastIP[:len(lastIP)-cut]
	}
}

// databaseGetAll - get articles by queryArticleGetAll.
func (a *entityArticle) databaseGetAll(valQuery *queryArticleGetAll) (articles []ModelArticle, err error) {
	if valQuery == nil {
		return nil, errors.New("databaseGetAll: getAllQuery nil pointer.")
	}
	articles = make([]ModelArticle, 0)
	var query string
	var by = valQuery.by
	var start = valQuery.start
	var cursor = valQuery.cursor
	switch valQuery.show {
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

// databaseCreate - create article.
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

// databaseUpdate - update article.
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

// databaseFind - find article.
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

// databaseDelete - delete article.
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

// databaseGetAll - get files by queryFileGetAll.
func (f *entityFile) databaseGetAll(valQuery *queryFileGetAll) (files []ModelFile, err error){
	if valQuery == nil {
		return nil, errors.New("databaseGetAll: getAllQuery nil pointer.")
	}
	files = make([]ModelFile, 0)
	var query string
	var by = valQuery.by
	var cursor = valQuery.cursor
	var start = valQuery.start
	query = fmt.Sprintf("SELECT * FROM files WHERE id >= $1 ORDER BY %v %v, id %v LIMIT $2 + 1", by, start, start)
	err = core.Database.Select(&files, query, cursor, filesPageSize)
	err = core.Utils.DBCheckError(err)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

// databaseCreate - create file.
func (f *entityFile) databaseCreate(file *ModelFile) (err error){
	var query = `INSERT INTO files (user_id, hash, path, name, original_name, extension, size) VALUES (:user_id, :hash, :path, :name, :original_name, :extension, :size) RETURNING *`
	row, err := core.Database.NamedQuery(query, &file)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return err
	}
	row.Next()
	err = row.StructScan(file)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return
	}
	return
}

// databaseFind - find one file in database by id.
func (f *entityFile) databaseFind(id string) (found *ModelFile, err error){
	found = &ModelFile{}
	var query = "SELECT * FROM files WHERE id=$1 LIMIT 1"
	err = core.Database.Get(found, query, id)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return nil, err
	}
	return
}

// databaseFindBy - find one file in database. WARNING: manually control 'field' param and don't allow user input on it, because sql injection can be executed.
func (f *entityFile) databaseFindBy(field string, equals string) (found *ModelFile, err error){
	found = &ModelFile{}
	var query = fmt.Sprintf("SELECT * FROM files WHERE %v=$1 LIMIT 1", field)
	err = core.Database.Get(found, query, equals)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return nil, err
	}
	return
}

// databaseDelete - delete file.
func (f *entityFile) databaseDelete(id string) (err error) {
	var query = "DELETE FROM files WHERE id=$1"
	_, err = core.Database.Exec(query, id)
	err = core.Utils.DBCheckError(err)
	return
}