package model

import (
	"fmt"
	"math"
	"servus/apps/elven/base"
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

func (f *File) queryGetSelectAll() string {
	return "SELECT * FROM files "
}

func (f *File) queryGetCount() string {
	return "SELECT count(*) FROM files "
}

func (f *File) GetPaginated(params *base.FileGetParams) (files map[int]*File, totalPages int, err error) {

	// preapare query.
	var query = f.queryGetSelectAll()
	var limitOffsetDollars = [2]string{"$1", "$2"}

	// args to get paginated files.
	var getAllArgs = []any{}
	if params.Extension != nil {
		limitOffsetDollars[0] = "$2"
		limitOffsetDollars[1] = "$3"
		// add where arg.
		getAllArgs = append(getAllArgs, *params.Extension)
		query += "WHERE extension = $1 "
	} else if params.ExtensionType != nil {
		var typd = *params.ExtensionType
		switch typd {
		case "image":
			query += "WHERE extension IN ('jpeg', 'jpg', 'gif', 'png', 'svg', 'bmp', 'webp') "
		case "audio":
			query += "WHERE extension IN ('mp3', 'flac', 'wav', 'ogg') "
		case "video":
			query += "WHERE extension IN ('mpg', 'mpeg', 'webm', 'mp4') "
		}
	}

	// get pages count.
	var queryCount = "SELECT count(*) FROM (" + query + ") as tentacles"
	totalPages = 1
	err = IntAdapter.Get(&totalPages, queryCount, getAllArgs...)
	if err != nil {
		return
	}
	totalPages = int(math.Round(float64(totalPages) / float64(FilePageSize)))
	if params.Page > totalPages {
		return
	}

	// sort. WARNING: potential SQL injection, be careful and validate this params.
	query += fmt.Sprintf("ORDER BY %s %s, id %s ", params.By, params.Start, params.Start)

	// add limit offset args.
	query += fmt.Sprintf("LIMIT %s OFFSET %s ", limitOffsetDollars[0], limitOffsetDollars[1])
	getAllArgs = append(getAllArgs, FilePageSize, (params.Page-1)*FilePageSize)

	// get all.
	files, err = fileAdapter.GetRows(query, getAllArgs...)
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
	founded, err := fileAdapter.Find(query, f.Hash)
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
