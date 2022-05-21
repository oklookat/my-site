package core

import "github.com/oklookat/cryptor"

func (i *Instance) setupEncryptor() error {
	var err error

	var module = &Encryptor{}

	var config = i.Config.Security.Encryption

	// AES.
	aes, err := cryptor.New_AES(config.AES.Secret)
	if err != nil {
		return err
	}
	module.AES = aes

	// BCrypt.
	var bcrypt = cryptor.New_BCrypt(config.BCrypt.Cost)
	module.BCrypt = bcrypt

	// Argon.
	var argon = cryptor.New_Argon(config.Argon)
	module.Argon = argon

	// set.
	i.Encryptor = module
	return err
}
