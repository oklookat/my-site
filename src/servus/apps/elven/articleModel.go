package elven

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/gosimple/slug"
	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	"math/rand"
	"time"
)

// ArticleModel - represents article in database.
type ArticleModel struct {
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
		ID   *string     `json:"id"`
		Type string      `json:"type"`
		Data interface{} `json:"data"`
		//Tunes *[]struct {
		//	Name interface{} `json:"name"`
		//} `json:"tunes"`
	} `json:"blocks"`
	Version *string `json:"version"`
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

// create - create article in database.
func (a *ArticleModel) create() (err error) {
	a.hookBeforeChange()
	var query = `INSERT INTO articles (user_id, is_published, title, content, slug) VALUES ($1, $2, $3, $4, $5) RETURNING *`
	err = call.DB.Conn.Get(a, query, a.UserID, a.IsPublished, a.Title, a.Content, a.Slug)
	err = call.DB.CheckError(err)
	if err != nil {
		return err
	}
	_ = a.hookAfterChange()
	return
}

// update - update article in database.
func (a *ArticleModel) update() (err error) {
	a.hookBeforeChange()
	var query = "UPDATE articles SET user_id=$1, is_published=$2, title=$3, content=$4, slug=$5, published_at=$6 WHERE id=$7 RETURNING *"
	err = call.DB.Conn.Get(a, query, a.UserID, a.IsPublished, a.Title, a.Content, a.Slug, a.PublishedAt, a.ID)
	err = call.DB.CheckError(err)
	if err != nil {
		return
	}
	_ = a.hookAfterChange()
	return
}

// findByID - find article in database by id field.
func (a *ArticleModel) findByID() (found bool, err error) {
	var query = "SELECT * FROM articles WHERE id=$1 LIMIT 1"
	err = call.DB.Conn.Get(a, query, a.ID)
	err = call.DB.CheckError(err)
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
func (a *ArticleModel) deleteByID() (err error) {
	var query = "DELETE FROM articles WHERE id=$1"
	_, err = call.DB.Conn.Exec(query, a.ID)
	err = call.DB.CheckError(err)
	return
}

// hookBeforeChange - executes before article create or update.
func (a *ArticleModel) hookBeforeChange() {
	if len(a.Title) == 0 {
		a.Title = "Untitled"
	}
	// create temp slug (ULID).
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	a.Slug = ulid.MustNew(ulid.Timestamp(t), entropy).String()
	// check published.
	if a.IsPublished && a.PublishedAt == nil {
		var cur = time.Now()
		a.PublishedAt = &cur
	}
}

// hookAfterChange - executes after article create or update.
func (a *ArticleModel) hookAfterChange() (err error) {
	// create normal slug.
	a.Slug = slug.Make(a.Title) + "-" + a.ID
	var query = "UPDATE articles SET slug=$1 WHERE id=$2 RETURNING *"
	row := call.DB.Conn.QueryRowx(query, a.Slug, a.ID)
	err = row.StructScan(a)
	if err != nil {
		return err
	}
	return nil
}
