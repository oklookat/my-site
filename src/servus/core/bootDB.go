package core

import (
	"fmt"
	"servus/core/database"
	"servus/core/logger"
)

func bootDB(config ConfigFile, logger *logger.Logger) *database.DB{
	var pgUser = config.DB.Postgres.User
	var pgPassword = config.DB.Postgres.Password
	var pgPort = config.DB.Postgres.Port
	var pgDb = config.DB.Postgres.DbName
	var timeZone = config.Timezone
	var connStr = fmt.Sprintf("user=%v password=%v port=%v dbname=%v sslmode=disable TimeZone=%v", pgUser, pgPassword, pgPort, pgDb, timeZone)
	return database.New(connStr, logger)
}
