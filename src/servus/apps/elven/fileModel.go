package elven

import (
	"database/sql"
	"time"
)

type FileModel struct {
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


// create - create file in database.
func (f *FileModel) create() (err error){
	var query = `INSERT INTO files (user_id, hash, path, name, original_name, extension, size) VALUES (:user_id, :hash, :path, :name, :original_name, :extension, :size) RETURNING *`
	stmt, err := call.DB.Conn.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()
	err = stmt.Get(f, f)
	err = call.DB.CheckError(err)
	if err != nil {
		return err
	}
	return
}

// findByID - find one file in database by id field.
func (f *FileModel) findByID() (found bool, err error){
	var query = "SELECT * FROM files WHERE id=$1 LIMIT 1"
	err = call.DB.Conn.Get(f, query, f.ID)
	err = call.DB.CheckError(err)
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
func (f *FileModel) findByHash() (found bool, err error){
	var query = "SELECT * FROM files WHERE hash=$1 LIMIT 1"
	err = call.DB.Conn.Get(f, query, f.Hash)
	err = call.DB.CheckError(err)
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
func (f *FileModel) deleteByID() (err error) {
	var query = "DELETE FROM files WHERE id=$1"
	_, err = call.DB.Conn.Exec(query, f.ID)
	err = call.DB.CheckError(err)
	return
}
