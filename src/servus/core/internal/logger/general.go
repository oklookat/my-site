package logger

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// New create new Logger instance.
func New(config Config) *Logger {
	var logger = Logger{Config: config}
	return &logger
}

// bus - calls when new log added. Writes log depending on settings.
func (l *Logger) bus(lev leveler) {
	// if silent level - logger not call anything.
	if l.Config.LogLevel > lev.getLevel() || l.Config.LogLevel == LevelSilent {
		return
	}
	if l.Config.WriteToFile.Activated {
		l.writeToFile(lev)
	}
	if l.Config.WriteToConsole {
		l.writeToConsole(lev)
	}
}

// writeToConsole - prints information to console pretty.
func (l *Logger) writeToConsole(lev leveler) {
	var currentTimePretty = colorGray + time.Now().Format("15:04:05") + colorReset
	var levelPretty = lev.getColor() + lev.getLevelWord() + colorReset
	var callerPretty = colorBlue + l.getAt() + colorReset
	var messagePretty = colorCyan + lev.getMessage() + colorReset
	var pretty = fmt.Sprintf("%v > %v > %v > %v", currentTimePretty, levelPretty, callerPretty, messagePretty)
	fmt.Println(pretty)
}

// writeToFile - write log to file.
func (l *Logger) writeToFile(lev leveler) {
	l.newLogFile()
	// write json to the log file.
	var theLog = logFile{
		Level:   lev.getLevelWord(),
		Time:    time.Now().Unix(),
		At:      l.getAt(),
		Message: lev.getMessage(),
	}
	theLogJson, err := json.Marshal(&theLog)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = l.file.instance.WriteString(string(theLogJson) + "\n")
	if err != nil {
		log.Println(err)
		return
	}
}

// newLogFile - create (or not) a log file, depending on settings.
func (l *Logger) newLogFile() {
	//// check and create log file.
	var err error
	var date = time.Now()
	if l.file.instance != nil {
		// create new log file if current log file size very big.
		logFileStat, err := l.file.instance.Stat()
		if err != nil {
			var prettyErr = errors.Wrap(err, "logger: newLogFile failed get log size. Error")
			log.Println(prettyErr)
			return
		}
		// create new log file if new day.
		var _, _, createdAt = l.file.created.Date()
		var _, _, currentDay = date.Date()
		var notContinue = createdAt == currentDay && l.Config.WriteToFile.MaxLogSize > logFileStat.Size()
		if notContinue {
			return
		}
		// close old file.
		err = l.file.instance.Close()
		if err != nil {
			var prettyErr = errors.Wrap(err, "logger: newLogFile failed to close old log file. Error")
			log.Println(prettyErr)
		}
	}
	var dirPath = l.Config.WriteToFile.Dir
	var fileName = "at_" + strconv.FormatInt(date.Unix(), 10) + ".log"
	var pathToLogFile = dirPath + "/" + fileName
	l.file.path = pathToLogFile
	// create log dir if no exists.
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0700)
		if err != nil {
			var prettyErr = errors.Wrap(err, "logger: newLogFile failed to create logs dir. Error")
			log.Println(prettyErr)
		}
	}
	// create log file.
	instance, err := os.Create(pathToLogFile)
	if err != nil {
		var prettyErr = errors.Wrap(err, "logger: newLogFile failed to create a log file. Error")
		log.Println(prettyErr)
		return
	}
	l.file.instance = instance
	l.file.created = date
	//////// delete latest log file if max log files in dir reached.
	files, err := os.ReadDir(dirPath)
	if err != nil {
		var prettyErr = errors.Wrap(err, "logger: newLogFile failed to read logs dir. Error")
		log.Println(prettyErr)
		return
	}
	var maxLogFiles = l.Config.WriteToFile.MaxLogFiles
	// if max log files reached.
	if !(len(files) > maxLogFiles) {
		return
	}
	// get count of files would be deleted.
	var countDifference = len(files) - maxLogFiles
	var attempts = 0
	for index := range files {
		// oldest files must be first in this cycle, because log filename has unix timestamp. And we have three attempts to delete.
		if countDifference <= 0 || attempts > 3 {
			break
		}
		if files[index].IsDir() {
			continue
		}
		// if not log file.
		if !strings.Contains(files[index].Name(), "at_") {
			continue
		}
		var firstLogPath = dirPath + "/" + files[index].Name()
		err := os.Remove(firstLogPath)
		if err != nil {
			var prettyErr = errors.Wrap(err, "logger: newLogFile failed to remove old log file. Error")
			log.Println(prettyErr)
			attempts++
			continue
		}
		countDifference--
	}
}
