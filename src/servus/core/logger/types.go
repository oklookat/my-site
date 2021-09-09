package logger

import (
	"os"
	"time"
)

// logging levels
const (
	SilentLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
)

// console colors
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
	colorWhite  = "\033[97m"
)

// main functions
type loggerI interface {
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(message string)
	Panic(err error)
}

// Config config
type Config struct {
	LogLevel       int
	WriteToConsole bool
	WriteToFile    struct {
		Activated bool
		Dir       string
		MaxLogFiles int
		MaxLogSize int64
	}
}

type Logger struct {
	loggerI
	Config
	fileWriterInfo fileWriterInfo
}

type fileWriterInfo struct {
	fullPath string // path to log file
	fileDate time.Time // date when file created
	file     *os.File // file itself
}

// log to file JSON
type logFile struct {
	Level   string `json:"level"`
	Time    int64  `json:"time"`
	At      string `json:"at"`
	Message string `json:"message"`
}
