package database

import (
	"database/sql"
)

// provides database functions. You must call Connector.New to init database connection (once, in core).
type Adapter[T comparable] struct {
}

// execute query with args and put result in dest (1 row).
func (a *Adapter[T]) Get(dest *T, query string, args ...any) (err error) {
	if dest == nil {
		dest = new(T)
	}
	// if result empty, keep dest in original state, no overwrite
	var destCopy T
	destCopy = *dest
	err = con.Connection.Get(&destCopy, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = con.checkError(err)
		}
		return
	}
	*dest = destCopy
	return
}

// execute query with args and get rows in array (many rows).
func (a *Adapter[T]) GetRows(query string, args ...any) (result map[int]*T, err error) {
	rows, err := con.Connection.Queryx(query, args...)
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()
	err = con.checkError(err)
	if err != nil {
		return
	}
	var mapCounter = 0
	result = make(map[int]*T, 0)
	for rows.Next() {
		newType := new(T)
		err = rows.StructScan(newType)
		if err != nil {
			return
		}
		result[mapCounter] = newType
		mapCounter++
	}
	return
}

// execute query (without rows).
func (a *Adapter[T]) Exec(query string, args ...any) (res sql.Result, err error) {
	res, err = con.Connection.Exec(query, args...)
	err = con.checkError(err)
	return
}

// find and put in found. found = nil if error or not found.
func (a *Adapter[T]) Find(query string, args ...any) (found *T, err error) {
	found = new(T)
	err = a.Get(found, query, args...)
	if err != nil {
		found = nil
		return
	}
	var emptySomething = new(T)
	var isEmpty = *emptySomething == *found
	if isEmpty {
		found = nil
	}
	return
}
