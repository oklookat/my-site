package database

import (
	"github.com/jmoiron/sqlx"
)

type DB struct {
	Connection *sqlx.DB
}
