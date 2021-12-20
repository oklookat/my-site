package logger

import (
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

func create() *Instance {
	var cnf = &Config{Level: 1}
	return New(cnf)
}

func TestNew(t *testing.T) {
	create()
}

func TestLogger_Info(t *testing.T) {
	var logger = create()
	for i := 0; i < 100; i++ {
		logger.Info("testing: info")
	}
}

func TestLogger_Debug(t *testing.T) {
	var logger = create()
	for i := 0; i < 100; i++ {
		logger.Debug("testing: debug")
	}
}

func TestLogger_Warn(t *testing.T) {
	var logger = create()
	for i := 0; i < 100; i++ {
		logger.Warn("testing: warn")
	}
}

func TestLogger_Error(t *testing.T) {
	var logger = create()
	for i := 0; i < 100; i++ {
		logger.Error("testing: error")
	}
}

func BenchmarkLogger_Info(b *testing.B) {
	var logger = create()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark: info")
	}
}
