package core

import "servus/core/internal/cryptor"

func (i *Instance) setupEncryptor() error {
	var err error

	// create.
	var cr = &cryptor.Instance{}
	if err = cr.New(i.Config.Security.Encryption); err != nil {
		return err
	}

	// set.
	var en = &Encryptor{}
	en.AES = cr.AES
	en.Argon = cr.Argon
	en.BCrypt = cr.BCrypt
	i.Encryptor = en
	return err
}
