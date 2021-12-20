package cryptor

type Config struct {
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
}

type Instance struct {
	config *Config
	AES    *AES
	BCrypt *BCrypt
	Argon  *Argon
}

func (e *Instance) New(config *Config) {
	e.config = config
	// aes.
	var aes = AES{}
	var aesCfg = e.config.AES
	aes.Secret = aesCfg.Secret
	e.AES = &aes
	// argon.
	var argon = Argon{}
	var argonCfg = e.config.Argon
	argon.Memory = argonCfg.Memory * 1024
	argon.Iterations = argonCfg.Iterations
	argon.Parallelism = argonCfg.Parallelism
	argon.SaltLength = argonCfg.SaltLength
	argon.KeyLength = argonCfg.KeyLength
	e.Argon = &argon
	// bcrypt.
	var bcrypt = BCrypt{}
	var bcryptCfg = e.config.Bcrypt
	bcrypt.New(bcryptCfg.Cost)
	e.BCrypt = &bcrypt
}
