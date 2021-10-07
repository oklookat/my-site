package logger

import (
	"fmt"
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
