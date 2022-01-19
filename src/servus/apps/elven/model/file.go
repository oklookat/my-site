package model

import (
	"fmt"
	"math"
	"servus/core/external/database"
	"time"
)

const FilePageSize = 2

type File struct {
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

var fileAdapter = database.Adapter[File]{}

func (f *File) GetPaginated(by string, start string, page int) (files map[int]*File, totalPages int, err error) {
	// get pages count.
	var queryCount = "SELECT count(*) FROM files"
	totalPages = 1
	err = IntAdapter.Get(&totalPages, queryCount)
	if err != nil {
		return
	}
	totalPages = int(math.Round(float64(totalPages) / float64(FilePageSize)))
	if page > totalPages {
		return
	}
	// get.
	var query = fmt.Sprintf("SELECT * FROM files ORDER BY %s %s, id %s LIMIT $1 OFFSET $2", by, start, start)
	files, err = fileAdapter.GetRows(query, FilePageSize, (page-1)*FilePageSize)
	return
}

// create file in database.
func (f *File) Create() (err error) {
	var query = `INSERT INTO files (user_id, hash, path, name, 
		original_name, extension, size) VALUES ($1, $2, $3, 
			$4, $5, $6, $7) 
		RETURNING *`
	err = fileAdapter.Get(f, query, f.UserID, f.Hash, f.Path, f.Name, f.OriginalName, f.Extension, f.Size)
	return
}

// find one file in database by id field.
func (f *File) FindByID() (found bool, err error) {
	found = false
	var query = "SELECT * FROM files WHERE id=$1 LIMIT 1"
	founded, err := fileAdapter.Find(query, f.ID)
	if err != nil {
		return
	}
	if founded != nil {
		found = true
		*f = *founded
	}
	return
}

// find file in database by hash field.
func (f *File) FindByHash() (found bool, err error) {
	found = false
	var query = "SELECT * FROM files WHERE hash=$1 LIMIT 1"
	founded, err := fileAdapter.Find(query, f.ID)
	if err != nil {
		return
	}
	if founded != nil {
		found = true
		*f = *founded
	}
	return
}

// delete file in database by id field.
func (f *File) DeleteByID() (err error) {
	var query = "DELETE FROM files WHERE id=$1"
	_, err = fileAdapter.Exec(query, f.ID)
	return
}
