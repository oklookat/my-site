package _core

import (
	"fmt"
	"servus/_core/logger"
)

func bootLogger() logger.Logger {
	// create log file
	var dir = fmt.Sprintf("%v/logs/", servus.Utils.GetExecuteDir())
	loggerConfig := logger.Config{
		LogLevel:       logger.DebugLevel,
		WriteToConsole: true,
		WriteToFile: struct {
			Activated bool
			Dir       string
		}{Activated: true, Dir: dir},
	}
	return logger.New(loggerConfig)
}
