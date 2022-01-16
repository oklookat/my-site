package model

import (
	"fmt"
	"math"
	"time"
)

const FilesPageSize = 2

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

func (f *File) GetPaginated(by string, start string, page int) (files map[int]File, totalPages int, err error) {
	// get pages count.
	var queryCount = "SELECT count(*) FROM files"
	totalPages = 1
	err = call.DB.Conn.Get(&totalPages, queryCount)
	err = call.DB.CheckError(err)
	if err != nil {
		return
	}
	files = make(map[int]File, 0)
	totalPages = int(math.Round(float64(totalPages) / float64(FilesPageSize)))
	if page > totalPages {
		return
	}
	// get.
	var query = fmt.Sprintf("SELECT * FROM files ORDER BY %v %v, id %v LIMIT $1 OFFSET $2", by, start, start)
	rows, err := call.DB.Conn.Queryx(query, FilesPageSize, (page-1)*FilesPageSize)
	defer func() {
		_ = rows.Close()
	}()
	if call.DB.IsNotFound(err) {
		err = nil
		return
	}
	err = call.DB.CheckError(err)
	if err != nil {
		return
	}
	var mapCounter = 0
	for rows.Next() {
		file := File{}
		err = rows.StructScan(&file)
		if err != nil {
			return
		}
		files[mapCounter] = file
		mapCounter++
	}
	return
}

// create file in database.
func (f *File) Create() (err error) {
	var query = `INSERT INTO files (user_id, hash, path, name, original_name, extension, size) VALUES (:user_id, :hash, :path, :name, :original_name, :extension, :size) RETURNING *`
	stmt, err := call.DB.Conn.PrepareNamed(query)
	defer func() {
		_ = stmt.Close()
	}()
	if err != nil {
		return
	}
	err = stmt.Get(f, f)
	err = call.DB.CheckError(err)
	return
}

// find one file in database by id field.
func (f *File) FindByID() (found bool, err error) {
	found = false
	var query = "SELECT * FROM files WHERE id=$1 LIMIT 1"
	err = call.DB.Conn.Get(f, query, f.ID)
	if call.DB.IsNotFound(err) {
		err = nil
		return
	}
	err = call.DB.CheckError(err)
	if err != nil {
		return
	}
	found = true
	return
}

// find file in database by hash field.
func (f *File) FindByHash() (found bool, err error) {
	found = false
	var query = "SELECT * FROM files WHERE hash=$1 LIMIT 1"
	err = call.DB.Conn.Get(f, query, f.Hash)
	if call.DB.IsNotFound(err) {
		err = nil
		return
	}
	err = call.DB.CheckError(err)
	if err != nil {
		return
	}
	found = true
	return
}

// delete file in database by id field.
func (f *File) DeleteByID() (err error) {
	var query = "DELETE FROM files WHERE id=$1"
	_, err = call.DB.Conn.Exec(query, f.ID)
	err = call.DB.CheckError(err)
	return
}
