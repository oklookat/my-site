package article

import (
	"fmt"
	"math"
	"servus/core/external/database"
	"strconv"
	"strings"
	"time"
)

const pageSize = 2

var articleAdapter = database.Adapter[Model]{}

// represents article in database.
type Model struct {
	ID          string     `json:"id" db:"id"`
	UserID      string     `json:"user_id" db:"user_id"`
	CoverID     *string    `json:"cover_id" db:"cover_id"`
	IsPublished bool       `json:"is_published" db:"is_published"`
	Title       string     `json:"title" db:"title"`
	Content     string     `json:"content" db:"content"`
	PublishedAt *time.Time `json:"published_at" db:"published_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`

	// available only when we get article(s) (JOIN).
	CoverPath      *string `json:"cover_path" db:"cover_path"`
	CoverExtension *string `json:"cover_extension" db:"cover_extension"`
}

// get query what gets articles + article cover (JOIN).
func (m *Model) queryGetWithCats(withSensitive bool, withContent bool) string {
	var query = `SELECT art.id, art.cover_id, art.is_published, art.title, art.published_at, `

	if withContent {
		query += "art.content, "
	}

	if withSensitive {
		query += "art.user_id, art.created_at, art.updated_at, "
	}

	query += `
	file.path as cover_path, file.extension as cover_extension
	FROM articles as art

	LEFT JOIN files as file
	ON art.cover_id = file.id
	`
	return query
}

// get query to get article(s) with join additional fields.
func (m *Model) queryGetSelectAll(withSensitive bool, withContent bool) string {
	return "SELECT * FROM (" + m.queryGetWithCats(withSensitive, withContent) + ") as tentacles\n"
}

// get paginated.
func (m *Model) GetPaginated(params *GetParams, isAdmin bool) (articles map[int]*Model, totalPages int, err error) {
	totalPages = 1

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

	var query = m.queryGetSelectAll(isAdmin, false)

	// is published.
	query += "WHERE is_published = " + addGetAllArg(!params.Drafts) + " "

	// title.
	if params.Title != nil {
		*params.Title = strings.ToLower(*params.Title)
		var dollar = addGetAllArg(*params.Title)
		query += "AND LOWER(title) LIKE '%'||" + dollar + "||'%' "
	}

	// get pages count.
	var queryCount = "SELECT count(*) FROM (" + query + ") as tentacles"
	if err = database.IntAdapter.Get(&totalPages, queryCount, getAllArgs...); err != nil {
		return
	}

	totalPages = int(math.Round(float64(totalPages) / float64(pageSize)))
	if totalPages < 1 {
		totalPages = 1
		return
	}
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
	getAllArgs = append(getAllArgs, pageSize, (params.Page-1)*pageSize)

	// get all.
	articles, err = articleAdapter.GetRows(query, getAllArgs...)
	return
}

// create in database.
func (m *Model) Create() (err error) {
	m.hookBeforeChange()
	var query = `
	INSERT INTO articles 
	(user_id, cover_id, is_published, title, content) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING *`
	err = articleAdapter.Get(m, query, m.UserID, m.CoverID, m.IsPublished, m.Title, m.Content)
	return
}

// update all article in database.
func (m *Model) Update() (err error) {
	m.hookBeforeChange()
	var query = `UPDATE articles SET 
	user_id=$1, cover_id=$2,
	is_published=$3, title=$4, content=$5, 
	published_at=$6 
	WHERE id=$7 RETURNING *`
	if err = articleAdapter.Get(m, query, m.UserID,
		m.CoverID, m.IsPublished, m.Title,
		m.Content, m.PublishedAt, m.ID); err != nil {
		return
	}
	return
}

// find article in database by id field.
func (m *Model) FindByID(isAdmin bool) (found bool, err error) {
	found = false
	var query = m.queryGetSelectAll(isAdmin, true) + "WHERE id=$1 LIMIT 1"
	founded, err := articleAdapter.Find(query, m.ID)
	if err != nil {
		return
	}
	if founded != nil {
		found = true
		*m = *founded
	}
	return
}

// delete article from database by id field.
func (m *Model) DeleteByID() (err error) {
	var query = "DELETE FROM articles WHERE id=$1"
	_, err = articleAdapter.Exec(query, m.ID)
	return
}

// executes before article create or update.
func (m *Model) hookBeforeChange() {

	// title.
	if len(strings.TrimSpace(m.Title)) < 1 {
		m.Title = "Untitled"
	}

	// article published but published date not exists? Fix that.
	if m.IsPublished && m.PublishedAt == nil {
		var currentTime = time.Now()
		m.PublishedAt = &currentTime
	}
}
