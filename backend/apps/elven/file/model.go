package file

import (
	"fmt"
	"math"
	"servus/core/external/database"
	"strconv"
	"strings"
	"time"
)

const pageSize = 2

var fileAdapter = database.Adapter[Model]{}

type Model struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	Hash         string    `json:"hash" db:"hash"`
	Path         string    `json:"path" db:"path"`
	Name         string    `json:"name" db:"name"`
	OriginalName string    `json:"original_name" db:"original_name"`
	Extension    string    `json:"extension" db:"extension"`
	Size         int64     `json:"size" db:"size"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func (m *Model) queryGetSelectAll() string {
	return "SELECT * FROM files "
}

func (m *Model) GetPaginated(params *GetParams) (files map[int]*Model, totalPages int, err error) {
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

	// args to get paginated files.

	var query = m.queryGetSelectAll()

	// extensions.
	var isExtensionsExists = params.Extensions != nil && len(params.Extensions) > 0
	if isExtensionsExists {
		// add extension dollars & args
		var extensionsLen = len(params.Extensions)
		for i := 0; i < extensionsLen; i++ {
			addGetAllArg(params.Extensions[i])
		}
		query += "WHERE extension IN (" + strings.Join(getAllDollars, ",") + ") "
	}

	// filename. DEPENDS ON PREVIOUS PARAM! IF YOU NEED TO CHANGE PREVIOUS PARAMS, RECHECK THIS CODE.
	if params.Filename != nil {
		*params.Filename = strings.ToLower(*params.Filename)
		var dollar = addGetAllArg(*params.Filename)
		var whereQuery = "LOWER(original_name) LIKE '%'||" + dollar + "||'%' "
		// > 1 because except previous appending.
		if len(getAllArgs) > 1 {
			query += "AND " + whereQuery
		} else {
			query += "WHERE " + whereQuery
		}
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

	// sort. WARNING: potential SQL injection, be careful and validate this params.
	query += fmt.Sprintf(`ORDER BY %s %s, id %s `, params.By, params.Start, params.Start)

	// add limit offset args (paginate).
	var limitOffsetDollars = [2]int{1, 2}
	limitOffsetDollars[0] = len(getAllDollars) + 1
	limitOffsetDollars[1] = len(getAllDollars) + 2
	query += fmt.Sprintf("LIMIT $%v OFFSET $%v ", limitOffsetDollars[0], limitOffsetDollars[1])
	getAllArgs = append(getAllArgs, pageSize, (params.Page-1)*pageSize)

	// get all.
	files, err = fileAdapter.GetRows(query, getAllArgs...)
	return
}

// create file in database.
func (m *Model) Create() (err error) {
	var query = `INSERT INTO files 
	(user_id, hash, path, name, original_name, extension, size) 
	VALUES 
	($1, $2, $3, $4, $5, $6, $7) 
	RETURNING *`
	err = fileAdapter.Get(m, query, m.UserID, m.Hash, m.Path, m.Name, m.OriginalName, m.Extension, m.Size)
	return
}

// find one file in database by id field.
func (m *Model) FindByID() (found bool, err error) {
	found = false
	var query = "SELECT * FROM files WHERE id=$1 LIMIT 1"
	founded, err := fileAdapter.Find(query, m.ID)
	if err != nil {
		return
	}
	if founded != nil {
		found = true
		*m = *founded
	}
	return
}

// find file in database by hash field.
func (m *Model) FindByHash() (found bool, err error) {
	found = false
	var query = "SELECT * FROM files WHERE hash=$1 LIMIT 1"
	founded, err := fileAdapter.Find(query, m.Hash)
	if err != nil {
		return
	}
	if founded != nil {
		found = true
		*m = *founded
	}
	return
}

// delete file in database by id field.
func (m *Model) DeleteByID() (err error) {
	var query = "DELETE FROM files WHERE id=$1"
	_, err = fileAdapter.Exec(query, m.ID)
	return
}
