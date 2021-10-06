package logger

import (
	"os"
	"time"
)

type Config struct {
	LogLevel       int
	WriteToConsole bool
	WriteToFile    struct {
		Activated   bool
		// Dir - for logs.
		Dir         string
		// MaxLogFiles - in log dir.
		MaxLogFiles int
		// MaxLogSize - in bytes.
		MaxLogSize  int64
	}
}

type Logger struct {
	Config Config
	file file
}

type file struct {
	// log file in system.
	instance *os.File
	// path to log file.
	path string
	// date when file created.
	created time.Time
}

type logFile struct {
	Level   string `json:"level"`
	Time    int64  `json:"time"`
	At      string `json:"at"`
	Message string `json:"message"`
}
