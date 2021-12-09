package core

type Core struct {
	Utils      *Utils
	Config     *ConfigFile
	Logger     Logger
	HTTP       *HTTP
	Middleware *Middleware
	Encryption *Encryption
	DB         *Database
	Control     *Control
}

// Boot - boot Core.
func (c *Core) Boot() {
	c.bootUtils()
	c.bootConfig()
	c.bootLogger()
	c.bootMiddleware()
	c.bootEncryption()
	c.bootDatabase()
	c.Control = &Control{}
	c.Control.boot(c)
}
