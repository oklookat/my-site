package core

import (
	"fmt"
	"servus/core/modules/logger"
)

func bootLogger() *logger.Logger {
	// create log file
	var writeToConsole = Config.Logger.WriteToConsole
	loggerConfig := logger.Config{
		LogLevel:       logger.LevelDebug,
		WriteToConsole: writeToConsole,
	}
	// write to file
	var wtfActive = Config.Logger.WriteToFile.Active
	loggerConfig.WriteToFile.Activated = wtfActive
	if wtfActive {
		var wtfDir = fmt.Sprintf("%v/logs/", Utils.GetExecuteDir())
		var wtfMaxLogFiles = Config.Logger.WriteToFile.MaxLogFiles
		var wtfMaxLogSize = Config.Logger.WriteToFile.MaxLogSize
		loggerConfig.WriteToFile.Dir = wtfDir
		loggerConfig.WriteToFile.MaxLogFiles = wtfMaxLogFiles
		loggerConfig.WriteToFile.MaxLogSize = wtfMaxLogSize
	}
	var theLogger = logger.New(loggerConfig)
	return &theLogger
}
