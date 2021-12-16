package logger

import (
	"fmt"
	"os"
	"time"
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

type Logger struct {
	Level int
}

// New create new Logger instance.
func New(level int) *Logger {
	var logger = Logger{Level: level}
	return &logger
}

// bus - calls when new log added. Writes log depending on settings.
func (l *Logger) bus(lev leveler) {
	// if silent level - logger not call anything.
	if l.Level > lev.getLevel() || l.Level == LevelSilent {
		return
	}
	l.writeToConsole(lev)
}

// writeToConsole - prints information to console pretty.
func (l *Logger) writeToConsole(lev leveler) {
	var currentTimePretty = colorGray + time.Now().Format("15:04:05") + colorReset
	var levelPretty = lev.getColor() + lev.getLevelWord() + colorReset
	var messagePretty = colorCyan + lev.getMessage() + colorReset
	var pretty = fmt.Sprintf("%v > %v > %v", currentTimePretty, levelPretty, messagePretty)
	fmt.Println(pretty)
}

// Debug - print debug message.
func (l *Logger) Debug(message string) {
	var lev = level{number: LevelDebug, word: "debug", color: colorGray, message: message}
	l.bus(&lev)
}

// Info - print info message.
func (l *Logger) Info(message string) {
	var lev = level{number: LevelInfo, word: "info", color: colorGray, message: message}
	l.bus(&lev)
}

// Warn - print warn message.
func (l *Logger) Warn(message string) {
	var lev = level{number: LevelWarn, word: "warn", color: colorYellow, message: message}
	l.bus(&lev)
}

// Error - throw error.
func (l *Logger) Error(message string) {
	var lev = level{number: LevelError, word: "error", color: colorRed, message: message}
	l.bus(&lev)
}

// Panic - throw panic. This func not throw panic(). It is like other levels, but executes os.Exit.
func (l *Logger) Panic(err error) {
	var lev = level{number: LevelPanic, word: "panic", color: colorRed, message: err.Error()}
	l.bus(&lev)
	os.Exit(1)
}
