package database

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var con *Connector

// controls database connection.
type Connector struct {
	connection *sqlx.DB
	config     *Config
	logger     Logger
}

// create new instance and connect.
func (c *Connector) New(config *Config, logger Logger) {
	if config == nil {
		var verr = errors.New("config nil pointer")
		verr = c.wrapError(verr)
		panic(verr)
	}
	if logger == nil {
		var verr = errors.New("logger nil pointer")
		verr = c.wrapError(verr)
		panic(verr)
	}
	// set
	c.config = config
	c.logger = logger
	var pgHost = config.Postgres.Host
	var pgUser = config.Postgres.User
	var pgPassword = config.Postgres.Password
	var pgPort = config.Postgres.Port
	var pgDb = config.Postgres.DbName
	var timeZone = config.Timezone
	// connect
	var connectionStr = fmt.Sprintf(`host=%s user=%s password=%s port=%s 
	dbname=%s sslmode=disable TimeZone=%s`, pgHost, pgUser, pgPassword, pgPort, pgDb, timeZone)
	connection, err := sqlx.Connect("pgx", connectionStr)
	if err != nil {
		err = c.wrapError(err)
		c.logger.Panic(err)
		return
	}
	c.connection = connection
	con = c
}

func (c *Connector) wrapError(err error) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, "[database] ")
}

// err == sql.ErrNoRows.
func (c *Connector) isNotFound(err error) bool {
	return err != nil && err == sql.ErrNoRows
}

// database error checking. If error - send err to logger and return err. If no rows - error will not send to logger.
func (c *Connector) checkError(err error) error {
	if err == nil || c.isNotFound(err) {
		return nil
	}
	err = c.wrapError(err)
	c.logger.Error(err.Error())
	return err
}
