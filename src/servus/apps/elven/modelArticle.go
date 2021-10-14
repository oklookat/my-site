package elven

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	"math"
	"math/rand"
	"servus/core"
	"time"
)

// ModelArticle - represents article in database.
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

// ArticleContent - represents ModelArticle content in database.
type ArticleContent struct {
	Time   int64 `json:"time"`
	Blocks []struct {
		ID   string      `json:"id"`
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	} `json:"blocks"`
	Version string `json:"version"`
}

// queryArticleGetAll - validated query params in article GetAll.
type queryArticleGetAll struct {
	page    int
	show    string
	by      string
	start   string
	preview bool
}

func (a ArticleContent) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *ArticleContent) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("articleContent: failed convert value to []byte")
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, &a)
}

// getAll - get articles by queryArticleGetAll.
func (q *queryArticleGetAll) getAll() (articles []ModelArticle, pagesCount int, err error) {
	articles = make([]ModelArticle, 0)
	var query string
	var queryCount string
	// get pages count
	var by = q.by
	var start = q.start
	switch q.show {
	case "published":
		// use sprintf to format validated.by and start because with $ database throws syntax error (I don't know why)
		// it not allows sql injection, because start and by checked in validator
		query = fmt.Sprintf("SELECT * FROM articles WHERE is_published=true ORDER BY %v %v, id %v LIMIT $1 OFFSET $2", by, start, start)
		queryCount = "SELECT count(*) FROM articles WHERE is_published=true"
		break
	case "drafts":
		query = fmt.Sprintf("SELECT * FROM articles WHERE is_published=false ORDER BY %v %v, id %v LIMIT $1 OFFSET $2", by, start, start)
		queryCount = "SELECT count(*) FROM articles WHERE is_published=false"
		break
	}
	// get pages count.
	pagesCount = 1
	row := core.Database.QueryRowx(queryCount)
	err = row.Scan(&pagesCount)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, nil
	}
	pagesCount = int(math.Round(float64(pagesCount / articlesPageSize)))
	// get articles.
	err = core.Database.Select(&articles, query, articlesPageSize, (q.page - 1) * articlesPageSize)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	return
}

// create - create article in database.
func (a *ModelArticle) create() (err error) {
	var query = `INSERT INTO articles (user_id, is_published, title, content, slug) VALUES (:user_id, :is_published, :title, :content, :slug) RETURNING *`
	a.hookBeforeChange()
	row, err := core.Database.NamedQuery(query, a)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return err
	}
	row.Next()
	err = row.StructScan(a)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return
	}
	_ = a.hookAfterChange()
	return
}

// update - update article in database.
func (a *ModelArticle) update() (err error) {
	a.hookBeforeChange()
	var query = "UPDATE articles SET user_id=:user_id, is_published=:is_published, title=:title, content=:content, slug=:slug WHERE id=:id RETURNING *"
	row, err := core.Database.NamedQuery(query, a)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return
	}
	row.Next()
	err = row.StructScan(a)
	err = core.Utils.DBCheckError(err)
	if err != nil {
		return
	}
	_ = a.hookAfterChange()
	return
}

// findByID - find article in database by id field.
func (a *ModelArticle) findByID() (found bool, err error) {
	var query = "SELECT * FROM articles WHERE id=$1 LIMIT 1"
	err = core.Database.Get(a, query, a.ID)
	err = core.Utils.DBCheckError(err)
	found = false
	if err != nil {
		if err == sql.ErrNoRows {
			return found, nil
		}
		return
	}
	found = true
	return
}

// deleteByID - delete article from database by id field.
func (a *ModelArticle) deleteByID() (err error) {
	var query = "DELETE FROM articles WHERE id=$1"
	_, err = core.Database.Exec(query, a.ID)
	err = core.Utils.DBCheckError(err)
	return
}

// hookBeforeChange - executes before article create or update.
func (a *ModelArticle) hookBeforeChange() {
	if len(a.Title) == 0 {
		a.Title = "Untitled"
	}
	// create temp slug (ULID)
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	a.Slug = ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

// hookAfterChange - executes after article create or update.
func (a *ModelArticle) hookAfterChange() (err error) {
	// create normal slug
	a.Slug = slug.Make(a.Title) + "-" + a.ID
	var query = "UPDATE articles SET slug=:slug WHERE id=:id RETURNING *"
	row, err := core.Database.NamedQuery(query, a)
	if err != nil {
		return err
	}
	err = row.StructScan(a)
	if err != nil {
		return err
	}
	return nil
}
