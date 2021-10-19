package elven

import (
	"database/sql"
	"fmt"
	"math"
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
	page int
	start  string
	by string
}

// getAll - get files in database by queryFileGetAll.
func (q *queryFileGetAll) getAll() (files []ModelFile, totalPages int, err error){
	// get pages count.
	var queryCount = "SELECT count(*) FROM files"
	totalPages = 1
	err = instance.DB.Conn.Get(&totalPages, queryCount)
	err = instance.DB.CheckError(err)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, nil
	}
	files = make([]ModelFile, 0)
	totalPages = int(math.Round(float64(totalPages / filesPageSize)))
	if q.page > totalPages {
		return
	}
	// get files
	var query = fmt.Sprintf("SELECT * FROM files ORDER BY %v %v, id %v LIMIT $1 OFFSET $2", q.by, q.start, q.start)
	rows, err := instance.DB.Conn.Queryx(query, filesPageSize, (q.page - 1) * filesPageSize)
	err = instance.DB.CheckError(err)
	for rows.Next() {
		file := ModelFile{}
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

// create - create file in database.
func (f *ModelFile) create() (err error){
	var query = `INSERT INTO files (user_id, hash, path, name, original_name, extension, size) VALUES (:user_id, :hash, :path, :name, :original_name, :extension, :size) RETURNING *`
	stmt, err := instance.DB.Conn.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()
	err = stmt.Get(f, f)
	err = instance.DB.CheckError(err)
	if err != nil {
		return err
	}
	return
}

// findByID - find one file in database by id field.
func (f *ModelFile) findByID() (found bool, err error){
	var query = "SELECT * FROM files WHERE id=$1 LIMIT 1"
	err = instance.DB.Conn.Get(f, query, f.ID)
	err = instance.DB.CheckError(err)
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
	err = instance.DB.Conn.Get(f, query, f.Hash)
	err = instance.DB.CheckError(err)
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
	_, err = instance.DB.Conn.Exec(query, f.ID)
	err = instance.DB.CheckError(err)
	return
}
