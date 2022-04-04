package core

import (
	"encoding/json"
	"os"
	"servus/core/external/database"
	"servus/core/internal/controlTelegram"
	"servus/core/internal/iHTTP"
	"servus/core/internal/limiter"
	"servus/core/internal/logger"

	"github.com/oklookat/argument"
	"github.com/oklookat/cryptor"
	"github.com/oklookat/gocors"
)

func (i *Instance) setupConfig() error {
	var config = Config{}

	// get from path.
	var get = func(path string) error {
		return config.load(path)
	}

	// get from data dir.
	var getFromDir = func() error {
		var dataDir, err = i.Dirs.GetData()
		if err != nil {
			return err
		}
		var path = dataDir + "/config.json"
		return get(path)
	}

	var err error
	var configPath = ""

	var argumentus = argument.New()
	argumentus.Add("config", "c", func(values []string) {
		if values == nil || len(values) < 1 {
			return
		}
		configPath = values[0]
	})
	argumentus.Start()

	if len(configPath) > 0 {
		// get by path in args.
		err = get(configPath)
	} else {
		// get from current dir.
		err = getFromDir()
	}

	i.Config = &config
	return err
}

// main configuration (config.json).
type Config struct {

	// debug mode active? Writes to logger debug information etc.
	Debug bool `json:"debug"`

	// app host. Ex: localhost.
	Host string `json:"host"`

	// app port. Ex: 3333.
	Port string `json:"port"`

	// database settings.
	DB *database.Config `json:"db"`

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

		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie
		Cookie *iHTTP.ConfigCookie `json:"cookie"`

		// https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
		CORS *gocors.Config `json:"cors"`

		// limit request things.
		Limiter struct {
			Body *limiter.BodyConfig `json:"body"`
		} `json:"limiter"`

		// data encryption.
		Encryption struct {
			AES    *cryptor.AES    `json:"aes"`
			BCrypt *cryptor.BCrypt `json:"bcrypt"`
			Argon  *cryptor.Argon  `json:"argon"`
		} `json:"encryption"`
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
func (c *Config) load(path string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	return
}
