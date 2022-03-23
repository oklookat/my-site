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
	// utils.
	println("servus: setup utils")
	i.setupUtils()

	// dirs.
	println("servus: setup directories")
	if err := i.setupDirectories(); err != nil {
		panic(err)
	}

	// config.
	println("servus: setup config")
	if err := i.setupConfig(); err != nil {
		panic(err)
	}

	// logger.
	println("servus: setup logger")
	i.setupLogger()

	// control.
	i.Logger.Info("servus: setup control")
	if err := i.setupControl(); err != nil {
		i.Logger.Panic(err)
		return
	}

	// banhammer.
	i.Logger.Info("servus: setup banhammer")
	if err := i.setupBanhammer(); err != nil {
		i.Logger.Panic(err)
		return
	}

	// middleware.
	i.Logger.Info("servus: setup middleware")
	i.setupMiddleware()

	// encryptor.
	i.Logger.Info("servus: setup encryptor")
	i.setupEncryptor()

	// database.
	i.Logger.Info("servus: setup database")
	if err := i.setupDatabase(); err != nil {
		i.Logger.Panic(err)
		return
	}

	// control commands.
	i.Logger.Info("servus: setup control commands")
	i.setupControlCommands()

	i.Logger.Info("servus: ready")
}
