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

type Instance struct {
	config *Config
}

type Config struct {
	// log level.
	Level int
}

// create new Logger instance.
func New(config *Config) *Instance {
	if config == nil {
		panic("[logger]: config nil pointer")
	}
	var logger = Instance{config: config}
	return &logger
}

// calls when new log added. Writes log depending on settings.
func (i *Instance) bus(lev leveler) {
	// if silent level - logger not call anything.
	if i.config.Level > lev.getLevel() || i.config.Level == LevelSilent {
		return
	}
	i.writeToConsole(lev)
}

// prints information to console pretty.
func (i *Instance) writeToConsole(lev leveler) {
	var currentTimePretty = colorGray + time.Now().Format("15:04:05") + colorReset
	var levelPretty = lev.getColor() + lev.getLevelWord() + colorReset
	var messagePretty = colorCyan + lev.getMessage() + colorReset
	var pretty = fmt.Sprintf("%v > %v > %v", currentTimePretty, levelPretty, messagePretty)
	fmt.Println(pretty)
}

// print debug message.
func (i *Instance) Debug(message string) {
	var lev = level{number: LevelDebug, word: "debug", color: colorGray, message: message}
	i.bus(&lev)
}

// print info message.
func (i *Instance) Info(message string) {
	var lev = level{number: LevelInfo, word: "info", color: colorGray, message: message}
	i.bus(&lev)
}

// print warn message.
func (i *Instance) Warn(message string) {
	var lev = level{number: LevelWarn, word: "warn", color: colorYellow, message: message}
	i.bus(&lev)
}

// print error.
func (i *Instance) Error(message string) {
	var lev = level{number: LevelError, word: "error", color: colorRed, message: message}
	i.bus(&lev)
}

// throw panic. This func not calls panic(). It is like other levels, but executes os.Exit.
func (i *Instance) Panic(err error) {
	var lev = level{number: LevelPanic, word: "panic", color: colorRed, message: err.Error()}
	i.bus(&lev)
	os.Exit(1)
}
