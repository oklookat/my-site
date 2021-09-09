package core

import (
	"fmt"
	"servus/core/logger"
)

func bootLogger() logger.Logger {
	// create log file
	var writeToConsole = servus.Config.Logger.WriteToConsole
	loggerConfig := logger.Config{
		LogLevel:       logger.DebugLevel,
		WriteToConsole: writeToConsole,
	}
	// write to file
	var wtfActive = servus.Config.Logger.WriteToFile.Active
	loggerConfig.WriteToFile.Activated = wtfActive
	if wtfActive {
		var wtfDir = fmt.Sprintf("%v/logs/", servus.Utils.GetExecuteDir())
		var wtfMaxLogFiles = servus.Config.Logger.WriteToFile.MaxLogFiles
		var wtfMaxLogSize = servus.Config.Logger.WriteToFile.MaxLogSize
		loggerConfig.WriteToFile.Dir = wtfDir
		loggerConfig.WriteToFile.MaxLogFiles = wtfMaxLogFiles
		loggerConfig.WriteToFile.MaxLogSize = wtfMaxLogSize
	}
	return logger.New(loggerConfig)
}
