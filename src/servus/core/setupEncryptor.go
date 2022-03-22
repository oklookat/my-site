package core

import "servus/core/internal/cryptor"

func (i *Instance) setupEncryptor() {
	var cr = &cryptor.Instance{}
	cr.New(i.Config.Security.Encryption)
	//
	var en = &Encryptor{}
	en.AES = cr.AES
	en.Argon = cr.Argon
	en.BCrypt = cr.BCrypt
	i.Encryptor = en
}
