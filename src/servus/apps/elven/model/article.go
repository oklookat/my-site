package model

import (
	"fmt"
	"math"
	"servus/apps/elven/base"
	"strconv"
	"strings"
	"time"
)

const ArticlePageSize = 2

// represents article in database.
type Article struct {
	ID          string     `json:"id" db:"id"`
	UserID      string     `json:"user_id" db:"user_id"`
	CategoryID  *string    `json:"category_id" db:"category_id"`
	CoverID     *string    `json:"cover_id" db:"cover_id"`
	IsPublished bool       `json:"is_published" db:"is_published"`
	Title       string     `json:"title" db:"title"`
	Content     string     `json:"content" db:"content"`
	PublishedAt *time.Time `json:"published_at" db:"published_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`

	// available only when we get article(s) (JOIN).
	CategoryName   *string `json:"category_name" db:"category_name"`
	CoverPath      *string `json:"cover_path" db:"cover_path"`
	CoverExtension *string `json:"cover_extension" db:"cover_extension"`
}

// get query what gets articles + categories names (JOIN)
func (a *Article) queryGetWithCats() string {
	var withCats = `
	SELECT art.*, cats.name as category_name
	FROM articles as art
	LEFT JOIN article_categories as cats
	ON art.category_id = cats.id
	`
	var withFile = `
	SELECT with_cat.*, 
	file.path as cover_path, file.extension as cover_extension
	FROM (` + withCats + `) as with_cat
	LEFT JOIN files as file
	ON with_cat.cover_id = file.id
	`
	return withFile
}

// get query to get article(s) with join category name
func (a *Article) queryGetSelectAll() string {
	return "SELECT * FROM (" + a.queryGetWithCats() + ") as tentacles\n"
}

// get query to get rows count
func (a *Article) queryGetCount() string {
	return "SELECT count(*) FROM (" + a.queryGetWithCats() + ") as tentacles\n"
}

// get paginated.
func (a *Article) GetPaginated(params *base.ArticleGetParams) (articles map[int]*Article, totalPages int, err error) {
	// preapare.
	var getAllDollars = make([]string, 0)
	var getAllArgs = make([]any, 0)
	var addGetAllArg = func(arg any) (insertedDollar string) {
		// add dollar.
		var dollar = "$" + strconv.Itoa(len(getAllDollars)+1)
		getAllDollars = append(getAllDollars, dollar)

		// add arg.
		getAllArgs = append(getAllArgs, arg)
		return getAllDollars[len(getAllDollars)-1]
	}

	var query = a.queryGetSelectAll()

	// is published.
	query += "WHERE is_published = " + addGetAllArg(params.Published) + " "

	// category name.
	if params.CategoryName != nil {
		query += "AND category_name = " + addGetAllArg(params.CategoryName) + " "
	} else if params.WithoutCategory {
		query += "AND category_id IS NULL "
	}

	// title.
	if params.Title != nil {
		*params.Title = strings.ToLower(*params.Title)
		var dollar = addGetAllArg(*params.Title)
		query += "AND LOWER(title) LIKE '%'||" + dollar + "||'%' "
	}

	// get pages count.
	var queryCount = "SELECT count(*) FROM (" + query + ") as tentacles"
	totalPages = 1
	if err = IntAdapter.Get(&totalPages, queryCount, getAllArgs...); err != nil {
		return
	}
	totalPages = int(math.Round(float64(totalPages) / float64(FilePageSize)))
	if params.Page > totalPages {
		return
	}

	// WARNING: potential 'ORDER BY' SQL injection, be careful and validate 'params.By'.
	var start = "DESC"
	if !params.Newest {
		start = "ASC"
	}
	query += fmt.Sprintf(`ORDER BY %s %s, id %s `, params.By, start, start)

	// add limit offset args (paginate).
	var limitOffsetDollars = [2]int{1, 2}
	limitOffsetDollars[0] = len(getAllDollars) + 1
	limitOffsetDollars[1] = len(getAllDollars) + 2
	query += fmt.Sprintf("LIMIT $%v OFFSET $%v ", limitOffsetDollars[0], limitOffsetDollars[1])
	getAllArgs = append(getAllArgs, ArticlePageSize, (params.Page-1)*ArticlePageSize)

	// get all.
	articles, err = articleAdapter.GetRows(query, getAllArgs...)
	return
}

// create in database. AFTER CREATING RETURNS ONLY ID.
func (a *Article) Create() (err error) {
	a.hookBeforeChange()
	var query = `
	INSERT INTO articles 
	(user_id, category_id, cover_id, is_published, title, content) 
	VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING id`
	err = articleAdapter.Get(a, query, a.UserID, a.CategoryID, a.CoverID, a.IsPublished, a.Title, a.Content)
	if err != nil {
		return
	}
	return
}

// update all article in database.
func (a *Article) Update() (err error) {
	a.hookBeforeChange()
	var query = `UPDATE articles SET 
	user_id=$1, category_id=$2, cover_id=$3,
	is_published=$4, title=$5, content=$6, 
	published_at=$7 
	WHERE id=$8 RETURNING *`
	if err = articleAdapter.Get(a, query, a.UserID, a.CategoryID,
		a.CoverID, a.IsPublished, a.Title,
		a.Content, a.PublishedAt, a.ID); err != nil {
		return
	}
	return
}

// find article in database by id field.
func (a *Article) FindByID() (found bool, err error) {
	found = false
	var query = a.queryGetSelectAll() + "WHERE id=$1 LIMIT 1"
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

	// title.
	if len(strings.TrimSpace(a.Title)) < 1 {
		a.Title = "Untitled"
	}

	// article published and no published date? wtf lets fix that.
	if a.IsPublished && a.PublishedAt == nil {
		var cur = time.Now()
		a.PublishedAt = &cur
	}
}
