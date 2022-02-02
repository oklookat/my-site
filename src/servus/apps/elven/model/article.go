package model

import (
	"fmt"
	"math"
	"time"
)

const ArticlePageSize = 2

// represents article in database.
type Article struct {
	ID         string  `json:"id" db:"id"`
	UserID     string  `json:"user_id" db:"user_id"`
	CategoryID *string `json:"category_id" db:"category_id"`
	// name of category available only when we get article(s)
	CategoryName *string    `json:"category_name" db:"category_name"`
	IsPublished  bool       `json:"is_published" db:"is_published"`
	Title        string     `json:"title" db:"title"`
	Content      string     `json:"content" db:"content"`
	PublishedAt  *time.Time `json:"published_at" db:"published_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// get query to get article(s) with join category name
func (a *Article) getQueryGetter() string {
	return `
	SELECT art.*, cats.name as category_name
	FROM articles as art
	LEFT JOIN article_categories as cats
	ON art.category_id = cats.id

	`
}

// get query to get rows count
func (a *Article) getQueryRowsCount() string {
	return `SELECT count(*) FROM articles `
}

// get paginated.
func (a *Article) GetPaginated(by string, start string, show string, page int, preview bool) (articles map[int]*Article, totalPages int, err error) {
	var isPublished string
	if show == "published" {
		isPublished = "true"
	} else if show == "drafts" {
		isPublished = "false"
	} else {
		return
	}
	var queryCount = a.getQueryRowsCount() + "WHERE is_published=$1"
	// get pages count.
	totalPages = 1
	err = IntAdapter.Get(&totalPages, queryCount, isPublished)
	if err != nil {
		return
	}
	totalPages = int(math.Round(float64(totalPages) / float64(ArticlePageSize)))
	if page > totalPages {
		return
	}
	// get.
	// attention: potential sql injection - check values in validator before use it
	var query = a.getQueryGetter() + fmt.Sprintf(`
		WHERE art.is_published = $1
		ORDER BY %v %v, id %v LIMIT $2 OFFSET $3;
		`, by, start, start)
	articles, err = articleAdapter.GetRows(query, isPublished, ArticlePageSize, (page-1)*ArticlePageSize)
	return
}

// create in database.
func (a *Article) Create() (err error) {
	a.hookBeforeChange()
	var query = `
	INSERT INTO articles (user_id, category_id, 
	is_published, title, content) VALUES ($1, $2, $3, $4, $5) RETURNING *`
	err = articleAdapter.Get(a, query, a.UserID, a.CategoryID, a.IsPublished, a.Title, a.Content)
	if err != nil {
		return
	}
	_ = a.hookAfterChange()
	return
}

// update all article in database.
func (a *Article) Update() (err error) {
	a.hookBeforeChange()
	var query = `UPDATE articles SET user_id=$1, category_id=$2,
	is_published=$3, title=$4, content=$5, published_at=$6 WHERE id=$7 RETURNING *`
	err = articleAdapter.Get(a, query, a.UserID, a.CategoryID, a.IsPublished, a.Title, a.Content, a.PublishedAt, a.ID)
	if err != nil {
		return
	}
	_ = a.hookAfterChange()
	return
}

// find article in database by id field.
func (a *Article) FindByID() (found bool, err error) {
	found = false
	var query = a.getQueryGetter() + "WHERE art.id=$1 LIMIT 1"
	founded, err := articleAdapter.Find(query, a.ID)
	if err != nil {
		return
	}
	if founded != nil {
		found = true
		*a = *founded
	}
	return
}

// delete article from database by id field.
func (a *Article) DeleteByID() (err error) {
	var query = "DELETE FROM articles WHERE id=$1"
	_, err = articleAdapter.Exec(query, a.ID)
	return
}

// executes before article create or update.
func (a *Article) hookBeforeChange() {
	if len(a.Title) == 0 {
		a.Title = "Untitled"
	}
	// check published.
	if a.IsPublished && a.PublishedAt == nil {
		var cur = time.Now()
		a.PublishedAt = &cur
	}
	// check category.
	if a.CategoryID != nil {
		var cat = ArticleCategory{}
		cat.ID = *a.CategoryID
		found, err := cat.FindByID()
		if err == nil && !found {
			// if category not found - reset category
			a.CategoryID = nil
		}
	}
}

// executes after article create or update.
func (a *Article) hookAfterChange() (err error) {
	return
}
