package logger

import (
	"os"
	"time"
)


// logging levels.
const (
	// LevelSilent - no messages.
	LevelSilent = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
)

// console colors.
const (
	// colorReset - no color.
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
	//colorWhite  = "\033[97m"
)

type Config struct {
	LogLevel       int
	WriteToConsole bool
	WriteToFile    struct {
		Activated bool
		// Dir - for logs.
		Dir string
		// MaxLogFiles - in log dir.
		MaxLogFiles int
		// MaxLogSize - in bytes.
		MaxLogSize int64
	}
}

type Logger struct {
	Config Config
	file   file
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

// leveler - represents methods for level.
type leveler interface {
	// getLevel - get error level number. Used for internal comparison.
	getLevel() int
	// getLevelWord - example: debug level must return "debug". Or debug level can return whatever. But better, level must return self name. Used for write to file and to console.
	getLevelWord() string
	// getColor - get console color for level.
	getColor() string
	// getMessage - get user message for level.
	getMessage() string
}

// level - represents level of log.
type level struct {
	at      string
	number  int
	word    string
	color   string
	message string
}
