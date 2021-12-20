package article

import (
	"database/sql"
	"fmt"
	"math"
	"servus/apps/elven/model"
)

// ArticleBody - represents the body of the request that the user should send. Used in create and update methods.
type Body struct {
	IsPublished *bool                 `json:"is_published"`
	Title       *string               `json:"title"`
	Content     *model.ArticleContent `json:"content"`
}

// queryGetAll - validated query params in article GetAll.
type queryGetAll struct {
	page    int
	show    string
	by      string
	start   string
	preview bool
}

// getAll - get articles by queryArticleGetAll.
func (q *queryGetAll) getAll() (articles []model.Article, totalPages int, err error) {
	var query string
	var queryCount string
	var by = q.by
	var start = q.start
	switch q.show {
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
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	articles = make([]model.Article, 0)
	totalPages = int(math.Round(float64(totalPages) / float64(pageSize)))
	if q.page > totalPages {
		return
	}
	// get articles.
	rows, err := call.DB.Conn.Queryx(query, pageSize, (q.page-1)*pageSize)
	err = call.DB.CheckError(err)
	for rows.Next() {
		article := model.Article{}
		err = rows.StructScan(&article)
		if err != nil {
			return
		}
		articles = append(articles, article)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	return
}
