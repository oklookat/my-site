package logger

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
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
	at string
	number  int
	word    string
	color   string
	message string
}

func (l *level) getLevel() int {
	return l.number
}

func (l *level) getLevelWord() string {
	return l.word
}

func (l *level) getColor() string {
	return l.color
}

func (l *level) getMessage() string {
	return l.message
}

// getAt - get place where logger was called. Returns string like main.go:11.
func (l *Logger) getAt() string {
	at := "unknown"
	if _, file, line, ok := runtime.Caller(4); ok {
		at = file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	}
	return at
}

// Debug - print debug message.
func (l *Logger) Debug(format string, args ...interface{}) {
	var lev = level{at: l.getAt(), number: LevelDebug, word: "debug", color: colorGray, message: fmt.Sprintf(format, args...)}
	l.bus(&lev)
}

// Info - print info message.
func (l *Logger) Info(format string, args ...interface{}) {
	var lev = level{at: l.getAt(), number: LevelInfo, word: "info", color: colorGray, message: fmt.Sprintf(format, args...)}
	l.bus(&lev)
}

// Warn - print warn message.
func (l *Logger) Warn(format string, args ...interface{}) {
	var lev = level{at: l.getAt(), number: LevelWarn, word: "warn", color: colorYellow, message: fmt.Sprintf(format, args...)}
	l.bus(&lev)
}

// Error - throw error.
func (l *Logger) Error(format string, args ...interface{}) {
	var lev = level{at: l.getAt(), number: LevelError, word: "error", color: colorRed, message: fmt.Sprintf(format, args...)}
	l.bus(&lev)
}

// Panic - throw panic. This func not throw panic(). It is like other levels, but executes os.Exit.
func (l *Logger) Panic(err error) {
	var lev = level{at: l.getAt(), number: LevelPanic, word: "panic", color: colorRed, message: err.Error()}
	l.bus(&lev)
	os.Exit(1)
}
