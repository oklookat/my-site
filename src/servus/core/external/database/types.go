package database

// writes information.
type Logger interface {
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(message string)
	Panic(err error)
}

type Config struct {
	// see timezones https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List
	Timezone string `json:"timezone"`
	// PostgreSQL.
	Postgres struct {
		// like: localhost.
		Host string `json:"host"`
		// like: 5432.
		Port string `json:"port"`
		// like: postgres.
		User string `json:"user"`
		// like: qwerty.
		Password string `json:"password"`
		// name of database.
		DbName string `json:"database"`
	} `json:"postgres"`
}
