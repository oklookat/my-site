package core

import (
	"encoding/json"
	"os"
)

// ConfigFile - configuration file.
type ConfigFile struct {
	Debug    bool   `json:"Debug"`
	Timezone string `json:"Timezone"`
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	DB       struct {
		Driver   string `json:"driver"`
		Postgres struct {
			Host     string `json:"host"`
			Port     string `json:"port"`
			User     string `json:"user"`
			Password string `json:"password"`
			DbName   string `json:"database"`
		} `json:"postgres"`
	} `json:"DB"`
	Logger struct {
		WriteToConsole bool `json:"writeToConsole"`
		WriteToFile    struct {
			Active      bool  `json:"active"`
			MaxLogFiles int   `json:"maxLogFiles"`
			MaxLogSize  int64 `json:"maxLogSize"`
		} `json:"writeToFile"`
	} `json:"Logger"`
	Security struct {
		HTTPS struct {
			Active   bool   `json:"active"`
			CertPath string `json:"certPath"`
			KeyPath  string `json:"keyPath"`
		} `json:"https"`
		Cookie struct {
			Domain   string `json:"domain"`
			Path     string `json:"path"`
			MaxAge   string `json:"maxAge"`
			HttpOnly bool   `json:"httpOnly"`
			Secure   bool   `json:"secure"`
			SameSite string `json:"sameSite"`
		} `json:"cookie"`
		CORS struct {
			Active           bool     `json:"active"`
			AllowCredentials bool     `json:"allowCredentials"`
			AllowOrigin      []string `json:"allowOrigin"`
			AllowMethods     []string `json:"allowMethods"`
			AllowHeaders     []string `json:"allowHeaders"`
			ExposeHeaders    []string `json:"exposeHeaders"`
			MaxAge           int64    `json:"maxAge"`
		} `json:"cors"`
		Limiter struct {
			Body struct {
				Active  bool     `json:"active"`
				MaxSize int64    `json:"maxSize"`
				Except  []string `json:"except"`
			} `json:"body"`
		} `json:"limiter"`
		Encryption struct {
			AES struct {
				Secret string `json:"secret"`
			} `json:"aes"`
			Argon  struct {
				Memory      uint32 `json:"memory"`
				Iterations  uint32 `json:"iterations"`
				Parallelism uint8  `json:"parallelism"`
				SaltLength  uint32 `json:"saltLength"`
				KeyLength   uint32 `json:"keyLength"`
			} `json:"argon"`
			Bcrypt struct {
				Cost int `json:"cost"`
			} `json:"bcrypt"`
		} `json:"encryption"`
	}
	Uploads struct {
		To   string `json:"to"`
		Temp string `json:"temp"`
	} `json:"Uploads"`
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
