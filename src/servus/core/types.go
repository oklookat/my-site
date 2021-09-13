package core

import (
	"net/http"
	"servus/core/database"
	"servus/core/logger"
)

// Servus this struct collect all server stuff
type Servus struct {
	DB         *database.DB
	Logger     logger.Logger
	Config     ConfigFile
	Utils      Utils
	Middleware Middleware
}

// ConfigFile config.json struct
type ConfigFile struct {
	Debug    bool   `json:"Debug"`
	Secret   string `json:"Secret"`
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
			AllowOrigin      []string `json:"allowOrigin"`
			AllowMethods     []string `json:"allowMethods"`
			AllowHeaders     []string `json:"allowHeaders"`
			ExposeHeaders    []string `json:"exposeHeaders"`
			AllowCredentials bool     `json:"allowCredentials"`
			MaxAge           int64    `json:"maxAge"`
		} `json:"cors"`
	}
}

// Utils servus utilities
type Utils struct {
	utils
}

// Servus utilities
type utils interface {
	RemoveSpaces(str string) string
	GetExecuteDir() string
	HashPassword(password string) (string, error)
	HashPasswordCheck(password, hash string) bool
	EncryptAES(text string) (encrypted string, err error)
	DecryptAES(encrypted string) (text string, err error)
	SetCookie(response *http.ResponseWriter, name string, value string)
}

// Middleware - hello I need basic middleware for my API
type Middleware struct {
	middleware
}

type middleware interface {
	MiddlewareAsJSON(next http.Handler) http.Handler
	MiddlewareCORS(next http.Handler) http.Handler
}
