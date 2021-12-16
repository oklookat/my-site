package core

import (
	"servus/core/internal/logger"
)

type LoggerConfig struct {
	// Level - log level.
	Level int `json:"level"`
}

type Log struct {
	config *LoggerConfig
	ins    *logger.Logger
}

// TODO: fix func position in logger
func (l *Log) new(c *LoggerConfig, executionDir string) {
	if c == nil {
		panic("[core/logger]: config nil pointer")
	}
	l.config = c
	l.ins = logger.New(c.Level)
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
