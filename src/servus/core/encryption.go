package core

import "servus/core/internal/cryptor"

type encryptor struct {
	config *EncryptorConfig
	AES    *cryptor.AES
	BCrypt *cryptor.BCrypt
	Argon  *cryptor.Argon
}

type EncryptorConfig struct {
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

func (e *encryptor) new(config *EncryptorConfig) {
	e.config = config
	// aes.
	var aes = cryptor.AES{}
	var aesCfg = e.config.AES
	aes.Secret = aesCfg.Secret
	e.AES = &aes
	// argon.
	var argon = cryptor.Argon{}
	var argonCfg = e.config.Argon
	argon.Memory = argonCfg.Memory * 1024
	argon.Iterations = argonCfg.Iterations
	argon.Parallelism = argonCfg.Parallelism
	argon.SaltLength = argonCfg.SaltLength
	argon.KeyLength = argonCfg.KeyLength
	e.Argon = &argon
	// bcrypt.
	var bcrypt = cryptor.BCrypt{}
	var bcryptCfg = e.config.Bcrypt
	bcrypt.New(bcryptCfg.Cost)
	e.BCrypt = &bcrypt
}
