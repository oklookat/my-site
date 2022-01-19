package database

import (
	"database/sql"
)

// provides database functions. You must call Connector.New to init database connection.
type Adapter[T comparable] struct {
}

// execute query with args and put result in dest (1 row).
func (a *Adapter[T]) Get(dest *T, query string, args ...any) (err error) {
	err = con.connection.Get(dest, query, args...)
	err = con.checkError(err)
	return
}

// execute query with args and get rows in array (many rows).
func (a *Adapter[T]) GetRows(query string, args ...any) (scaned map[int]*T, err error) {
	rows, err := con.connection.Queryx(query, args...)
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
	scaned = make(map[int]*T, 0)
	for rows.Next() {
		something := new(T)
		err = rows.StructScan(something)
		if err != nil {
			return
		}
		scaned[mapCounter] = something
		mapCounter++
	}
	return
}

// execute query (without rows).
func (a *Adapter[T]) Exec(query string, args ...any) (res sql.Result, err error) {
	res, err = con.connection.Exec(query, args...)
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
	var empty = *emptySomething == *found
	if empty {
		found = nil
	}
	return
}
