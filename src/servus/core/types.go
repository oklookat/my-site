package core

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
		Limiter struct {
			Body struct {
				Active  bool     `json:"active"`
				MaxSize int64    `json:"maxSize"`
				Except  []string `json:"except"`
			} `json:"body"`
		} `json:"limiter"`
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
	}
	Uploads struct {
		To   string `json:"to"`
		Temp string `json:"temp"`
	} `json:"Uploads"`
}

// BasicUtils - must have utilities.
type BasicUtils struct {
}

// BasicMiddleware - hello I need basic middleware for my API.
type BasicMiddleware struct {
}
