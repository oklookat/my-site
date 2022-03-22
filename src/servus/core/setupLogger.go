package core

import "servus/core/internal/logger"

func (i *Instance) setupLogger() {
	var log = logger.New(i.Config.Logger)
	i.Logger = log
}
