package logger

import (
	"fmt"
	"os"
	"testing"
)

// GetExecuteDir - get server execution directory.
func GetExecuteDir() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path
}

func create() Logger {
	loggerConfig := Config{
		LogLevel:       LevelDebug,
		WriteToConsole: true,
	}
	loggerConfig.WriteToFile.Activated = true
	var wtfDir = fmt.Sprintf("%v/logs/", GetExecuteDir())
	var wtfMaxLogFiles = 1
	var wtfMaxLogSize = 10000
	loggerConfig.WriteToFile.Dir = wtfDir
	loggerConfig.WriteToFile.MaxLogFiles = wtfMaxLogFiles
	loggerConfig.WriteToFile.MaxLogSize = int64(wtfMaxLogSize)
	return New(loggerConfig)
}

func TestNew(t *testing.T) {
	create()
}

func TestLogger_Info(t *testing.T) {
	var logger = create()
	for i := 0; i < 100; i++ {
		logger.Info("testing: info. Formatted: %v.", "yes")
	}
}

func TestLogger_Debug(t *testing.T) {
	var logger = create()
	for i := 0; i < 100; i++ {
		logger.Debug("testing: debug. Formatted: %v.", "yes")
	}
}

func TestLogger_Warn(t *testing.T) {
	var logger = create()
	for i := 0; i < 100; i++ {
		logger.Warn("testing: warn. Formatted: %v.", "yes")
	}
}

func TestLogger_Error(t *testing.T) {
	var logger = create()
	for i := 0; i < 100; i++ {
		logger.Error("testing: error. Formatted: %v.", "yes")
	}
}

func BenchmarkLogger_Info(b *testing.B) {
	var logger = create()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark: info. Formatted: %v.", "yes")
	}
}
