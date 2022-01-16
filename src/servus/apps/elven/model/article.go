package model

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/gosimple/slug"
	"github.com/oklog/ulid/v2"
)

const ArticlePageSize = 2

// represents article in database.
type Article struct {
	ID           string     `json:"id" db:"id"`
	UserID       string     `json:"user_id" db:"user_id"`
	CategoryName *string    `json:"category_name" db:"category_name"`
	IsPublished  bool       `json:"is_published" db:"is_published"`
	Title        string     `json:"title" db:"title"`
	Content      string     `json:"content" db:"content"`
	Slug         string     `json:"slug" db:"slug"`
	PublishedAt  *time.Time `json:"published_at" db:"published_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// get paginated.
func (a *Article) GetPaginated(by string, start string, show string, page int, preview bool) (articles map[int]Article, totalPages int, err error) {
	var query string
	var queryCount string
	switch show {
	case "published":
		// use sprintf to format validated.by and start because with $ database throws syntax error (I don't know why)
		// it not allows sql injection, because start and by checked in validator
		query = fmt.Sprintf("SELECT * FROM articles WHERE is_published=true ORDER BY %v %v, id %v LIMIT $1 OFFSET $2", by, start, start)
		queryCount = "SELECT count(*) FROM articles WHERE is_published=true"
	case "drafts":
		query = fmt.Sprintf("SELECT * FROM articles WHERE is_published=false ORDER BY %v %v, id %v LIMIT $1 OFFSET $2", by, start, start)
		queryCount = "SELECT count(*) FROM articles WHERE is_published=false"
	}
	// get pages count.
	totalPages = 1
	err = call.DB.Conn.Get(&totalPages, queryCount)
	err = call.DB.CheckError(err)
	if err != nil {
		return
	}
	totalPages = int(math.Round(float64(totalPages) / float64(ArticlePageSize)))
	if page > totalPages {
		return
	}
	// get.
	rows, err := call.DB.Conn.Queryx(query, ArticlePageSize, (page-1)*ArticlePageSize)
	defer func() {
		_ = rows.Close()
	}()
	if call.DB.IsNotFound(err) {
		err = nil
		return
	}
	err = call.DB.CheckError(err)
	if err != nil {
		return
	}
	var mapCounter = 0
	articles = make(map[int]Article, 0)
	for rows.Next() {
		article := Article{}
		err = rows.StructScan(&article)
		if err != nil {
			return
		}
		if preview {
			article.Content = ""
		}
		articles[mapCounter] = article
		mapCounter++
	}
	return
}

// create in database.
func (a *Article) Create() (err error) {
	a.hookBeforeChange()
	var query = `INSERT INTO articles (user_id, category_name, 
		is_published, title, content, slug) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`
	err = call.DB.Conn.Get(a, query, a.UserID, a.CategoryName, a.IsPublished, a.Title, a.Content, a.Slug)
	err = call.DB.CheckError(err)
	if err != nil {
		return
	}
	_ = a.hookAfterChange()
	return
}

// update all article in database.
func (a *Article) Update() (err error) {
	a.hookBeforeChange()
	var query = `UPDATE articles SET user_id=$1, category_name=$2,
	is_published=$3, title=$4, content=$5, slug=$6, published_at=$7 WHERE id=$8 RETURNING *`
	err = call.DB.Conn.Get(a, query, a.UserID, a.CategoryName, a.IsPublished, a.Title, a.Content, a.Slug, a.PublishedAt, a.ID)
	err = call.DB.CheckError(err)
	if err != nil {
		return
	}
	_ = a.hookAfterChange()
	return
}

// find article in database by id field.
func (a *Article) FindByID() (found bool, err error) {
	found = false
	var query = "SELECT * FROM articles WHERE id=$1 LIMIT 1"
	err = call.DB.Conn.Get(a, query, a.ID)
	if call.DB.IsNotFound(err) {
		err = nil
		return
	}
	err = call.DB.CheckError(err)
	if err != nil {
		return
	}
	found = true
	return
}

// delete article from database by id field.
func (a *Article) DeleteByID() (err error) {
	var query = "DELETE FROM articles WHERE id=$1"
	_, err = call.DB.Conn.Exec(query, a.ID)
	err = call.DB.CheckError(err)
	return
}

// executes before article create or update.
func (a *Article) hookBeforeChange() {
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
	// check category.
	if a.CategoryName != nil {
		var cat = ArticleCategory{}
		cat.ID = *a.CategoryName
		found, err := cat.FindByName()
		if err == nil && !found {
			// if category not found - reset category
			a.CategoryName = nil
		}
	}
}

// executes after article create or update.
func (a *Article) hookAfterChange() (err error) {
	// create normal slug.
	a.Slug = slug.Make(a.Title) + "-" + a.ID
	var query = "UPDATE articles SET slug=$1 WHERE id=$2 RETURNING *"
	row := call.DB.Conn.QueryRowx(query, a.Slug, a.ID)
	err = row.StructScan(a)
	return
}
