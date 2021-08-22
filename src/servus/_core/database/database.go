package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"servus/_core/logger"
)

var pLogger *logger.Logger
var pConnection *pgx.Conn

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
	pConnection = db
	return &DB{
		Connection: pConnection,
		User:       UserObject{},
	}
}

type DB struct {
	Connection *pgx.Conn
	User       UserObject
}

// Service start
func structEmptyErr() error {
	return errors.New("STRUCT_EMPTY")
}


func errorHandler(err error) error {
	var errReadable error = nil
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			defer debugPrinter(err)
			switch pgErr.Code {
			case pgerrcode.ConnectionException:
				pLogger.Panic(err)
			case pgerrcode.UniqueViolation: // for ex. username already exists
				return errors.New("E_EXISTS")
			case pgerrcode.NotNullViolation: // null value provided for NOT NULL
				return errors.New("E_NOT_NULL")
			case pgerrcode.CheckViolation: // for ex. username has min length 4
				return errors.New("E_CHECK")
			default:
				return errors.New("E_UNKNOWN")
			}
		}
		pLogger.Error(fmt.Sprintf("%v\n", err))
	}
	return errReadable
}

func debugPrinter(err error){
	pLogger.Debug(fmt.Sprintf("%v\n", err))
}
// Service end
