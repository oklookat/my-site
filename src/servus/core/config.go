package core

import (
	"encoding/json"
	"os"
	"servus/core/internal/controlTelegram"
	"servus/core/internal/cors"
	"servus/core/internal/cryptor"
	"servus/core/internal/iHTTP"
	"servus/core/internal/logger"
)

// main configuration (config.json).
type ConfigFile struct {
	// debug mode active? Writes to logger debug information etc.
	Debug bool `json:"debug"`
	// Timezone - timezone for database.
	//
	// see TZ Database Name in https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List
	Timezone string `json:"timezone"`
	// app host. Ex: localhost.
	Host string `json:"host"`
	// app port. Ex: 3333.
	Port string `json:"port"`
	// database settings.
	DB *DatabaseConfig `json:"db"`
	// writes messages to console and file.
	Logger *logger.Config `json:"logger"`
	// protect your ass from hackers.
	Security struct {
		// HTTPS connection.
		HTTPS struct {
			Active bool `json:"active"`
			// absolute path to certificate file.
			CertPath string `json:"certPath"`
			// absolute path to key file.
			KeyPath string `json:"keyPath"`
		} `json:"https"`
		// cookie.
		//
		// see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie
		Cookie *iHTTP.ConfigCookie `json:"cookie"`
		// see: https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
		CORS *cors.Config `json:"cors"`
		// limit request things.
		Limiter struct {
			// limit body.
			Body struct {
				Active  bool     `json:"active"`
				MaxSize int64    `json:"maxSize"`
				Except  []string `json:"except"`
			} `json:"body"`
		} `json:"limiter"`
		// data encryption.
		Encryption *cryptor.Config `json:"encryption"`
	}
	// files uploading.
	Uploads struct {
		// files will be saved here.
		To string `json:"to"`
		// temp folder before file saved.
		Temp string `json:"temp"`
	} `json:"uploads"`
	// control server / get notifications via third-party services.
	Control struct {
		Telegram *controlTelegram.Config `json:"telegram"`
	} `json:"control"`
}

// load config file from path.
func (c *ConfigFile) load(path string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	return
}

type DatabaseConfig struct {
	// PostgreSQL settings.
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
