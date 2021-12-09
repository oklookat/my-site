package core

import (
	"encoding/json"
	"fmt"
	"os"
)

// ConfigFile - main configuration (config.json).
type ConfigFile struct {
	// Debug - is debug mode active. Writes to logger debug information etc.
	Debug    bool   `json:"debug"`
	// Timezone - timezone for database.
	//
	// see TZ Database Name in https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List
	Timezone string `json:"timezone"`
	// Host - app host. Ex: localhost.
	Host     string `json:"host"`
	// Port - app port. Ex: 3333.
	Port     string `json:"port"`
	// DB - database settings.
	DB       struct {
		// Driver - not working now.
		Driver   string `json:"driver"`
		// PostgreSQL settings.
		Postgres struct {
			// Host - like: localhost.
			Host     string `json:"host"`
			// Port - like: 5432.
			Port     string `json:"port"`
			// User - like: postgres.
			User     string `json:"user"`
			// Password - like: qwerty.
			Password string `json:"password"`
			// DbName - name of database.
			DbName   string `json:"database"`
		} `json:"postgres"`
	} `json:"db"`
	// Logger - writes messages to console and file.
	Logger struct {
		// WriteToConsole - write messages to console.
		WriteToConsole bool `json:"writeToConsole"`
		// WriteToConsole - write messages to file.
		WriteToFile    struct {
			Active      bool   `json:"active"`
			// WriteToConsole - max log files in logs folder. After reaching limit, the oldest file will be removed.
			MaxLogFiles int   `json:"maxLogFiles"`
			// MaxLogSize - in bytes. After reaching limit new log files be created.
			MaxLogSize  int64 `json:"maxLogSize"`
		} `json:"writeToFile"`
	} `json:"logger"`
	// Security - protect your ass from hackers.
	Security struct {
		// HTTPS - HTTPS connection.
		HTTPS struct {
			Active   bool   `json:"active"`
			// CertPath - absolute path of certificate file.
			CertPath string `json:"certPath"`
			// CertPath - absolute path of key file.
			KeyPath  string `json:"keyPath"`
		} `json:"https"`
		// Cookie - settings for HTTP.SetCookie.
		//
		// see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie
		Cookie struct {
			Domain   string `json:"domain"`
			Path     string `json:"path"`
			MaxAge   string `json:"maxAge"`
			HttpOnly bool   `json:"httpOnly"`
			Secure   bool   `json:"secure"`
			SameSite string `json:"sameSite"`
		} `json:"cookie"`
		// CORS - see https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
		CORS struct {
			Active           bool     `json:"active"`
			AllowCredentials bool     `json:"allowCredentials"`
			AllowOrigin      []string `json:"allowOrigin"`
			AllowMethods     []string `json:"allowMethods"`
			AllowHeaders     []string `json:"allowHeaders"`
			ExposeHeaders    []string `json:"exposeHeaders"`
			MaxAge           int64   `json:"maxAge"`
		} `json:"cors"`
		// Limiter - limit request things.
		Limiter struct {
			// Body - limit body.
			Body struct {
				Active  bool     `json:"active"`
				MaxSize int64   `json:"maxSize"`
				Except  []string `json:"except"`
			} `json:"body"`
		} `json:"limiter"`
		// Encryption - data encryption.
		Encryption struct {
			AES struct {
				// Secret - 32 bytes length.
				Secret string `json:"secret"`
			} `json:"aes"`
			Bcrypt struct {
				Cost int `json:"cost"`
			} `json:"bcrypt"`
			// Argon - see: https://github.com/alexedwards/argon2id#changing-the-parameters
			Argon struct {
				Memory      uint32 `json:"memory"`
				Iterations  uint32 `json:"iterations"`
				Parallelism uint8  `json:"parallelism"`
				SaltLength  uint32 `json:"saltLength"`
				KeyLength   uint32 `json:"keyLength"`
			} `json:"argon"`
		} `json:"encryption"`
	}
	// Uploads - files uploading.
	Uploads struct {
		// To - files will be saved here.
		To   string `json:"to"`
		// Temp - temp folder before file saved.
		Temp string `json:"temp"`
	} `json:"uploads"`
	// Control - control server / get notifications via third-party services.
	Control struct {
		// Telegram - cool messenger.
		Telegram struct {
			Enabled      bool     `json:"enabled"`
			// Token - bot token.
			Token        string   `json:"token"`
			// AllowedUsers - bot accepts messages only from these users (id).
			AllowedUsers []int64 `json:"allowedUsers"`
			// AllowedChats - bot sends messages only to these chats (id).
			AllowedChats []int64 `json:"allowedChats"`
			// HttpLevel - HTTP errors logging level.
			HttpLevel    uint     `json:"httpLevel"`
		} `json:"telegram"`
	} `json:"control"`
}

// bootConfig - boot config file.
func (c *Core) bootConfig() {
	var config = ConfigFile{}
	var path = fmt.Sprintf("%v/settings/config.json", c.Utils.GetExecutionDir())
	err := config.boot(path)
	if err != nil {
		panic(err)
	}
	c.Config = &config
}

// boot - load config file from path.
func (c *ConfigFile) boot(path string) (err error) {
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
