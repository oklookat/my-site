package database

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var pLogger Logger

func New(connectionStr string, logger Logger) *DB {
	pLogger = logger
	db, err := sqlx.Connect("pgx", connectionStr)
	if err != nil {
		var errPretty = errors.Wrap(err, "database connection failed")
		pLogger.Panic(errPretty)
	}
	return &DB{
		db,
	}
}
