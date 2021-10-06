package database

import (
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

type Logger interface {
	Panic(err error)
}