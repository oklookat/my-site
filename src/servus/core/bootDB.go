package core

import (
	"fmt"
	"servus/core/modules/database"
	"servus/core/modules/logger"
)

// bootDB - connect to database and get connection.
func bootDB(config *ConfigFile, logger *logger.Logger) *database.DB {
	if config == nil {
		panic("bootDB: config nil pointer.")
	}
	if logger == nil {
		panic("bootDB: logger nil pointer.")
	}
	var pgUser = config.DB.Postgres.User
	var pgPassword = config.DB.Postgres.Password
	var pgPort = config.DB.Postgres.Port
	var pgDb = config.DB.Postgres.DbName
	var timeZone = config.Timezone
	var connStr = fmt.Sprintf("user=%v password=%v port=%v dbname=%v sslmode=disable TimeZone=%v", pgUser, pgPassword, pgPort, pgDb, timeZone)
	return database.New(connStr, logger)
}
