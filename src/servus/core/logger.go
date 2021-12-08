package core

import (
	"fmt"
	"servus/core/internal/logger"
)

// bootLogger - boot Logger. Use this after booting the config.
func (c *Core) bootLogger() {
	// create log file
	var writeToConsole = c.Config.Logger.WriteToConsole
	loggerConfig := logger.Config{
		LogLevel:       logger.LevelDebug,
		WriteToConsole: writeToConsole,
	}
	// write to file
	var wtfActive = c.Config.Logger.WriteToFile.Active
	loggerConfig.WriteToFile.Activated = wtfActive
	if wtfActive {
		var wtfDir = fmt.Sprintf("%v/logs/", c.Utils.GetExecutionDir())
		var wtfMaxLogFiles = c.Config.Logger.WriteToFile.MaxLogFiles
		var wtfMaxLogSize = c.Config.Logger.WriteToFile.MaxLogSize
		loggerConfig.WriteToFile.Dir = wtfDir
		loggerConfig.WriteToFile.MaxLogFiles = wtfMaxLogFiles
		loggerConfig.WriteToFile.MaxLogSize = wtfMaxLogSize
	}
	var theLogger = logger.New(loggerConfig)
	c.Logger = &theLogger
}
