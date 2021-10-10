package elven

import (
	"database/sql"
	"fmt"
	"servus/core"
	"time"
)

type ModelFile struct {
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

// queryFileGetAll - validated query params in files GetAll.
type queryFileGetAll struct {
	cursor string
	start  string
	by string
}

// getAll - get files in database by queryFileGetAll.
func (q *queryFileGetAll) getAll() (files []ModelFile, err error){
	files = make([]ModelFile, 0)
	var query string
	query = fmt.Sprintf("SELECT * FROM files WHERE id >= $1 ORDER BY %v %v, id %v LIMIT $2 + 1", q.by, q.start, q.start)
	err = core.Database.Select(&files, query, q.cursor, filesPageSize)
	err = core.Utils.DBCheckError(err)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

// create - create file in database.
func (f *ModelFile) create() (err error){
	var query = `INSERT INTO files (user_id, hash, path, name, original_name, extension, size) VALUES (:user_id, :hash, :path, :name, :original_name, :extension, :size) RETURNING *`
	row := core.Database.QueryRowx(query, f)
	err = row.StructScan(f)
	err = core.Utils.DBCheckError(err)
	return
}

// findByID - find one file in database by id field.
func (f *ModelFile) findByID() (found bool, err error){
	var query = "SELECT * FROM files WHERE id=$1 LIMIT 1"
	err = core.Database.Get(f, query, f.ID)
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

// findByHash - find file in database by hash field.
func (f *ModelFile) findByHash() (found bool, err error){
	var query = "SELECT * FROM files WHERE hash=$1 LIMIT 1"
	err = core.Database.Get(found, query, f.Hash)
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

// deleteByID - delete file in database by id field.
func (f *ModelFile) deleteByID() (err error) {
	var query = "DELETE FROM files WHERE id=$1"
	_, err = core.Database.Exec(query, f.ID)
	err = core.Utils.DBCheckError(err)
	return
}
