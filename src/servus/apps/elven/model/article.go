package model

import (
	"fmt"
	"math"
	"servus/apps/elven/base"
	"time"
)

const ArticlePageSize = 2

// represents article in database.
type Article struct {
	ID         string  `json:"id" db:"id"`
	UserID     string  `json:"user_id" db:"user_id"`
	CategoryID *string `json:"category_id" db:"category_id"`
	CoverID    *string `json:"cover_id" db:"cover_id"`
	// name of category available only when we get article(s).
	CategoryName *string    `json:"category_name" db:"category_name"`
	IsPublished  bool       `json:"is_published" db:"is_published"`
	Title        string     `json:"title" db:"title"`
	Content      string     `json:"content" db:"content"`
	PublishedAt  *time.Time `json:"published_at" db:"published_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// get query what gets articles + categories names (join)
func (a *Article) getQueryGetterWithCatsNames() string {
	return `
	SELECT art.*, cats.name as category_name
	FROM articles as art
	LEFT JOIN article_categories as cats
	ON art.category_id = cats.id
	`
}

// get query to get article(s) with join category name
func (a *Article) getQueryGetter() string {
	return "SELECT * FROM (" + a.getQueryGetterWithCatsNames() + ") as tentacles\n"
}

// get query to get rows count
func (a *Article) getQueryRowsCount() string {
	return "SELECT count(*) FROM (" + a.getQueryGetterWithCatsNames() + ") as tentacles\n"
}

// get paginated.
func (a *Article) GetPaginated(params *base.ArticleGetParams) (articles map[int]*Article, totalPages int, err error) {
	//
	var isPublished bool
	if params.Show == "published" {
		isPublished = true
	} else if params.Show == "drafts" {
		isPublished = false
	} else {
		return
	}
	var categoryNameExists = params.CategoryName != nil

	// queries. //
	var queryWhereIsPublished = "WHERE is_published = $1 "
	var queryAndCategoryName = "AND category_name = $2 "
	var queryAndCategoryIDnull = "AND category_id IS NULL "

	// get pages count depend on category name. //
	var queryGetCount = a.getQueryRowsCount() + queryWhereIsPublished
	var getPagesCountArgsArr = []any{isPublished}
	if params.WithoutCategory {
		queryGetCount += queryAndCategoryIDnull
	} else if categoryNameExists {
		getPagesCountArgsArr = append(getPagesCountArgsArr, *params.CategoryName)
		queryGetCount += queryAndCategoryName
	}
	totalPages = 1
	err = IntAdapter.Get(&totalPages, queryGetCount, getPagesCountArgsArr...)
	if err != nil {
		return
	}
	totalPages = int(math.Round(float64(totalPages) / float64(ArticlePageSize)))
	if params.Page > totalPages {
		return
	}

	// get articles depend on category name (prepare). //
	var queryGetArticles = a.getQueryGetter() + queryWhereIsPublished
	var getArticlesArgsArr = []any{isPublished}
	var getArticlesArgsLimitOffset = [2]string{"$2", "$3"}
	if params.WithoutCategory {
		queryGetArticles += queryAndCategoryIDnull
	} else if categoryNameExists {
		queryGetArticles += queryAndCategoryName
		getArticlesArgsLimitOffset[0] = "$3"
		getArticlesArgsLimitOffset[1] = "$4"
		getArticlesArgsArr = append(getArticlesArgsArr, *params.CategoryName)
	}

	// get articles & paginate. //
	// attention: potential sql injection - check values in validator before use it
	queryGetArticles += fmt.Sprintf("ORDER BY %v %v, id %v LIMIT %v OFFSET %v", params.By, params.Start, params.Start,
		getArticlesArgsLimitOffset[0], getArticlesArgsLimitOffset[1])
	getArticlesArgsArr = append(getArticlesArgsArr, ArticlePageSize, (params.Page-1)*ArticlePageSize)
	articles, err = articleAdapter.GetRows(queryGetArticles, getArticlesArgsArr...)
	return
}

// create in database. AFTER CREATING RETURNS ONLY ID.
func (a *Article) Create() (err error) {
	a.hookBeforeChange()
	var query = `
	INSERT INTO articles (user_id, category_id, 
	is_published, title, content) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = articleAdapter.Get(a, query, a.UserID, a.CategoryID, a.IsPublished, a.Title, a.Content)
	if err != nil {
		return
	}
	return
}

// update all article in database.
func (a *Article) Update() (err error) {
	a.hookBeforeChange()
	var query = `UPDATE articles SET user_id=$1, category_id=$2, cover_id=$3,
	is_published=$4, title=$5, content=$6, published_at=$7 WHERE id=$8 RETURNING *`
	err = articleAdapter.Get(a, query, a.UserID, a.CategoryID, a.CoverID, a.IsPublished, a.Title, a.Content, a.PublishedAt, a.ID)
	if err != nil {
		return
	}
	return
}

// find article in database by id field.
func (a *Article) FindByID() (found bool, err error) {
	found = false
	var query = a.getQueryGetter() + "WHERE id=$1 LIMIT 1"
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
	// article published and no published date? wtf lets fix that.
	if a.IsPublished && a.PublishedAt == nil {
		var cur = time.Now()
		a.PublishedAt = &cur
	}
}
