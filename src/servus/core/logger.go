package core

import (
	"fmt"
	"servus/core/internal/logger"
)

type Logger interface {
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(message string)
	Panic(err error)
}

type LoggerConfig struct {
	// WriteToConsole - write messages to console.
	WriteToConsole bool `json:"writeToConsole"`
	// WriteToConsole - write messages to file.
	WriteToFile struct {
		Active bool `json:"active"`
		// WriteToConsole - max log files in logs folder. After reaching limit, the oldest file will be removed.
		MaxLogFiles int `json:"maxLogFiles"`
		// MaxLogSize - in bytes. After reaching limit new log files be created.
		MaxLogSize int64 `json:"maxLogSize"`
	} `json:"writeToFile"`
}

type Log struct {
	config *LoggerConfig
	ins *logger.Logger
}

// TODO: fix func position in logger
func (l *Log) new(c *LoggerConfig, executionDir string) {
	l.config = c
	// create log file
	var writeToConsole = l.config.WriteToConsole
	loggerConfig := logger.Config{
		LogLevel:       logger.LevelDebug,
		WriteToConsole: writeToConsole,
	}
	// write to file
	var wtfActive = l.config.WriteToFile.Active
	loggerConfig.WriteToFile.Activated = wtfActive
	if wtfActive {
		var wtfDir = fmt.Sprintf("%v/logs/", executionDir)
		var wtfMaxLogFiles = l.config.WriteToFile.MaxLogFiles
		var wtfMaxLogSize = l.config.WriteToFile.MaxLogSize
		loggerConfig.WriteToFile.Dir = wtfDir
		loggerConfig.WriteToFile.MaxLogFiles = wtfMaxLogFiles
		loggerConfig.WriteToFile.MaxLogSize = wtfMaxLogSize
	}
	l.ins = logger.New(loggerConfig)
}

func (l *Log) Debug(message string) {
	l.ins.Debug(message)
}

func (l *Log) Info(message string) {
	l.ins.Info(message)
}

func (l *Log) Warn(message string) {
	l.ins.Warn(message)
}

func (l *Log) Error(message string) {
	l.ins.Error(message)
}

func (l *Log) Panic(err error) {
	l.ins.Panic(err)
}
