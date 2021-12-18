package core

import "servus/core/internal/cryptor"

type encryptor struct {
	config *EncryptorConfig
	AES    *cryptor.AES
	BCrypt *cryptor.BCrypt
	Argon  *cryptor.Argon
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
