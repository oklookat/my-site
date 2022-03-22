package core

// Instance - servus kernel. Provides cool things.
type Instance struct {
	Utils      Utils
	Dirs       Directories
	Config     *Config
	Logger     Logger
	Banhammer  Banhammer
	Middleware Middlewarer
	Encryptor  *Encryptor
	Control    Controller
}

// boot servus.
func (i *Instance) Boot() {
	i.setupUtils()
	i.setupDirectories()
	i.setupConfig()
	i.setupLogger()
	i.setupBanhammer()
	i.setupControl()
	i.setupMiddleware()
	i.setupEncryptor()
	i.setupDatabase()
}
