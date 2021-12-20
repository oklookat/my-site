package file

import (
	"database/sql"
	"fmt"
	"math"
	"servus/apps/elven/model"
)

// queryGetAll - validated query params in files GetAll.
type queryGetAll struct {
	page  int
	start string
	by    string
}

// getAll - get files in database by queryFileGetAll.
func (q *queryGetAll) getAll() (files []model.File, totalPages int, err error) {
	// get pages count.
	var queryCount = "SELECT count(*) FROM files"
	totalPages = 1
	err = call.DB.Conn.Get(&totalPages, queryCount)
	err = call.DB.CheckError(err)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, nil
	}
	files = make([]model.File, 0)
	totalPages = int(math.Round(float64(totalPages) / float64(filesPageSize)))
	if q.page > totalPages {
		return
	}
	// get files
	var query = fmt.Sprintf("SELECT * FROM files ORDER BY %v %v, id %v LIMIT $1 OFFSET $2", q.by, q.start, q.start)
	rows, err := call.DB.Conn.Queryx(query, filesPageSize, (q.page-1)*filesPageSize)
	err = call.DB.CheckError(err)
	for rows.Next() {
		file := model.File{}
		err = rows.StructScan(&file)
		if err != nil {
			return
		}
		files = append(files, file)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	return
}
