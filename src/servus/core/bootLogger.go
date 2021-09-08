package core

import (
	"fmt"
	"servus/core/logger"
)

func bootLogger() logger.Logger {
	// create log file
	var dir = fmt.Sprintf("%v/logs/", servus.Utils.GetExecuteDir())
	var writeToConsole = servus.Config.Logger.WriteToConsole
	var writeToFile = servus.Config.Logger.WriteToFile
	loggerConfig := logger.Config{
		LogLevel:       logger.DebugLevel,
		WriteToConsole: writeToConsole,
		WriteToFile: struct {
			Activated bool
			Dir       string
		}{Activated: writeToFile, Dir: dir},
	}
	return logger.New(loggerConfig)
}
