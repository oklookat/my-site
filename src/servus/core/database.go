package core

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Database struct {
	config *ConfigFile
	logger Logger
	Conn   *sqlx.DB
}

func (d *Database) new(config *ConfigFile, logger Logger) (err error) {
	if config == nil {
		return errors.New("[core/database]: config nil pointer")
	}
	if logger == nil {
		return errors.New("[core/database]: logger nil pointer")
	}
	d.config = config
	var pgHost = d.config.DB.Postgres.Host
	var pgUser = d.config.DB.Postgres.User
	var pgPassword = d.config.DB.Postgres.Password
	var pgPort = d.config.DB.Postgres.Port
	var pgDb = d.config.DB.Postgres.DbName
	var timeZone = d.config.Timezone
	var connectionStr = fmt.Sprintf("host=%v user=%v password=%v port=%v dbname=%v sslmode=disable TimeZone=%v", pgHost, pgUser, pgPassword, pgPort, pgDb, timeZone)
	connection, err := sqlx.Connect("pgx", connectionStr)
	if err != nil {
		var errPretty = errors.Wrap(err, "[core/database]")
		return errPretty
	}
	d.logger = logger
	d.Conn = connection
	return
}

// CheckError - database error checking. If error - send err to logger and return err. If no rows - error will not send to logger.
func (d *Database) CheckError(err error) error {
	switch err {
	case nil:
		return nil
	case sql.ErrNoRows:
		return err
	default:
		d.logger.Error(err.Error())
		return err
	}
}
