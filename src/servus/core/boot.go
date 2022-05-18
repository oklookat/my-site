package core

// boot servus.
func (i *Instance) Boot() {
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
	if err := i.setupLogger(); err != nil {
		i.Logger.Panic(err)
		return
	}

	// database.
	i.Logger.Info("servus: setup database")
	// if err := i.setupDatabase(); err != nil {
	// 	i.Logger.Panic(err)
	// 	return
	// }

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

	// cors.
	i.Logger.Info("servus: setup CORS")
	if err := i.setupCors(); err != nil {
		i.Logger.Panic(err)
		return
	}

	// limiter.
	i.Logger.Info("servus: setup limiter")
	if err := i.setupLimiter(); err != nil {
		i.Logger.Panic(err)
		return
	}

	// encryptor.
	i.Logger.Info("servus: setup encryptor")
	if err := i.setupEncryptor(); err != nil {
		i.Logger.Panic(err)
		return
	}

	// http.
	i.Logger.Info("servus: setup HTTP")
	if err := i.setupHTTP(); err != nil {
		i.Logger.Panic(err)
		return
	}

	// control commands.
	i.Logger.Info("servus: setup control commands")
	i.setupControlCommands()

	i.Logger.Info("servus: ready")
}
