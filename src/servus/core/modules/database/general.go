package database

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"servus/core/modules/logger"
)

var pLogger *logger.Logger

func New(connectionStr string, _logger *logger.Logger) *DB {
	pLogger = _logger
	db, err := sqlx.Connect("pgx", connectionStr)
	if err != nil {
		pLogger.Panic(err)
	}
	return &DB{
		Connection: db,
	}
}
