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

// ConfigFile - main configuration (config.json).
type ConfigFile struct {
	// Debug - is debug mode active. Writes to logger debug information etc.
	Debug bool `json:"debug"`
	// Timezone - timezone for database.
	//
	// see TZ Database Name in https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List
	Timezone string `json:"timezone"`
	// Host - app host. Ex: localhost.
	Host string `json:"host"`
	// Port - app port. Ex: 3333.
	Port string `json:"port"`
	// DB - database settings.
	DB *DatabaseConfig `json:"db"`
	// Logger - writes messages to console and file.
	Logger *logger.Config `json:"logger"`
	// Security - protect your ass from hackers.
	Security struct {
		// HTTPS - HTTPS connection.
		HTTPS struct {
			Active bool `json:"active"`
			// CertPath - absolute path of certificate file.
			CertPath string `json:"certPath"`
			// CertPath - absolute path of key file.
			KeyPath string `json:"keyPath"`
		} `json:"https"`
		// Cookie - settings for HTTP.SetCookie.
		//
		// see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie
		Cookie *iHTTP.ConfigCookie `json:"cookie"`
		// CORS - see https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
		CORS *cors.Config `json:"cors"`
		// Limiter - limit request things.
		Limiter struct {
			// Body - limit body.
			Body struct {
				Active  bool     `json:"active"`
				MaxSize int64    `json:"maxSize"`
				Except  []string `json:"except"`
			} `json:"body"`
		} `json:"limiter"`
		// Encryption - data encryption.
		Encryption *cryptor.Config `json:"encryption"`
	}
	// Uploads - files uploading.
	Uploads struct {
		// To - files will be saved here.
		To string `json:"to"`
		// Temp - temp folder before file saved.
		Temp string `json:"temp"`
	} `json:"uploads"`
	// Control - control server / get notifications via third-party services.
	Control struct {
		Telegram *controlTelegram.Config `json:"telegram"`
	} `json:"control"`
}

// load - load config file from path.
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
		// Host - like: localhost.
		Host string `json:"host"`
		// Port - like: 5432.
		Port string `json:"port"`
		// User - like: postgres.
		User string `json:"user"`
		// Password - like: qwerty.
		Password string `json:"password"`
		// DbName - name of database.
		DbName string `json:"database"`
	} `json:"postgres"`
}
