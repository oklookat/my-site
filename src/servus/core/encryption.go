package core

import "servus/core/internal/cryptor"

type Encryption struct {
	config *ConfigFile
	AES    *cryptor.AES
	Bcrypt *cryptor.Bcrypt
	Argon  *cryptor.Argon
}

func (e *Encryption) boot() {
	// aes.
	var aes = cryptor.AES{}
	var aesCfg = e.config.Security.Encryption.AES
	aes.Secret = aesCfg.Secret
	e.AES = &aes
	// argon.
	var argon = cryptor.Argon{}
	var argonCfg = e.config.Security.Encryption.Argon
	argon.Memory = argonCfg.Memory * 1024
	argon.Iterations = argonCfg.Iterations
	argon.Parallelism = argonCfg.Parallelism
	argon.SaltLength = argonCfg.SaltLength
	argon.KeyLength = argonCfg.KeyLength
	e.Argon = &argon
	// bcrypt.
	var bcrypt = cryptor.Bcrypt{}
	var bcryptCfg = e.config.Security.Encryption.Bcrypt
	bcrypt.Cost = bcryptCfg.Cost
	e.Bcrypt = &bcrypt
}
