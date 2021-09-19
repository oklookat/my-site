package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"servus/core/modules/logger"
)

var pLogger *logger.Logger

func New(connectionStr string, _logger *logger.Logger) *DB {
	pLogger = _logger
	var ctx = context.Background()
	db, err := pgx.Connect(ctx, connectionStr)
	if err != nil {
		pLogger.Panic(err)
	}
	err = db.Ping(ctx)
	if err != nil {
		pLogger.Panic(err)
	}
	return &DB{
		Connection: db,
	}
}
