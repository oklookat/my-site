package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
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

type DB struct {
	Connection *pgx.Conn
}

// service

func errorHandler(err error) error {
	var errReadable error = nil
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			defer warnPrinter(err)
			switch pgErr.Code {
			case pgerrcode.ConnectionException:
				pLogger.Panic(err)
			case pgerrcode.UniqueViolation:
				// for ex. username already exists
				return errors.New("DB_E_EXISTS")
			case pgerrcode.NotNullViolation:
				// null value provided for NOT NULL
				return errors.New("DB_E_NOT_NULL")
			case pgerrcode.CheckViolation:
				// for ex. username has min length 4
				return errors.New("DB_E_CHECK")
			default:
				return errors.New("DB_E_UNKNOWN")
			}
		}
		pLogger.Error(fmt.Sprintf("%v\n", err))
	}
	return errReadable
}

func warnPrinter(err error) {
	pLogger.Warn(fmt.Sprintf("%v", err))
}

func DerefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
