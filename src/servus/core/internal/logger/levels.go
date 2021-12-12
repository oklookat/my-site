package logger

import (
	"os"
	"runtime"
	"strconv"
	"strings"
)

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

// getAt - get place where logger was called. Returns string like limitBody.go:11.
func (l *Logger) getAt() string {
	at := "unknown"
	if _, file, line, ok := runtime.Caller(4); ok {
		at = file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	}
	return at
}

// Debug - print debug message.
func (l *Logger) Debug(message string) {
	var lev = level{at: l.getAt(), number: LevelDebug, word: "debug", color: colorGray, message: message}
	l.bus(&lev)
}

// Info - print info message.
func (l *Logger) Info(message string) {
	var lev = level{at: l.getAt(), number: LevelInfo, word: "info", color: colorGray, message: message}
	l.bus(&lev)
}

// Warn - print warn message.
func (l *Logger) Warn(message string) {
	var lev = level{at: l.getAt(), number: LevelWarn, word: "warn", color: colorYellow, message: message}
	l.bus(&lev)
}

// Error - throw error.
func (l *Logger) Error(message string) {
	var lev = level{at: l.getAt(), number: LevelError, word: "error", color: colorRed, message: message}
	l.bus(&lev)
}

// Panic - throw panic. This func not throw panic(). It is like other levels, but executes os.Exit.
func (l *Logger) Panic(err error) {
	var lev = level{at: l.getAt(), number: LevelPanic, word: "panic", color: colorRed, message: err.Error()}
	l.bus(&lev)
	os.Exit(1)
}
