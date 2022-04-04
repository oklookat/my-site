package core

import "servus/core/internal/logger"

func (i *Instance) setupLogger() error {
	var log, err = logger.New(i.Config.Logger)
	if err != nil {
		return err
	}
	i.Logger = log
	return nil
}
