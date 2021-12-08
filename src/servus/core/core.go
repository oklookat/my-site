package core

type Core struct {
	Utils      *Utils
	Config     *ConfigFile
	Logger     Logger
	HTTP       *HTTP
	Encryption *Encryption
	DB         *Database
}

// Boot - boot Core.
func (c *Core) Boot() {
	c.bootUtils()
	c.bootConfig()
	c.bootLogger()
	c.bootHTTP()
	c.bootEncryption()
	c.bootDatabase()
}

