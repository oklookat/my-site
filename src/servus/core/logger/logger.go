package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)


// New create new logger
func New(config Config) Logger {
	var logger = Logger{Config: config}
	var err = newFileWriter(&logger, time.Now())
	if err != nil {
		panic(err)
	}
	return logger
}


func (l *Logger) Debug(message string) {
	logWriter(l, DebugLevel, message)
}
func (l *Logger) Info(message string) {
	logWriter(l, InfoLevel, message)
}
func (l *Logger) Warn(message string) {
	logWriter(l, WarnLevel, message)
}
func (l *Logger) Error(message string) {
	logWriter(l, ErrorLevel, message)
}
func (l *Logger) Panic(err error) {
	logWriter(l, PanicLevel, err.Error())
	os.Exit(1)
}

func logWriter(l *Logger, calledLevel int, message string) {
	var calledLevelString string
	var errorColor string
	switch calledLevel {
	case DebugLevel:
		calledLevelString = "debug"
		errorColor = colorGray
	case InfoLevel:
		calledLevelString = "info"
		errorColor = colorGray
	case WarnLevel:
		calledLevelString = "warn"
		errorColor = colorYellow
	case ErrorLevel:
		calledLevelString = "error"
		errorColor = colorRed
	case PanicLevel:
		calledLevelString = "panic"
		errorColor = colorRed
	default:
		calledLevelString = "unknown"
		errorColor = colorGray
	}
	var at = "unknown"
	// get caller function
	if _, file, line, ok := runtime.Caller(2); ok {
		at = file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	}
	var logLevel = l.Config.LogLevel
	if logLevel > calledLevel || logLevel == SilentLevel {
		return
	}
	// write to file
	if l.Config.WriteToFile.Activated {
		logFileWriter(l, calledLevelString, at, message)
	}
	// format message for console
	if l.Config.WriteToConsole {
		var curTime = fmt.Sprintf("%v%v%v", colorGray, time.Now().Format("15:04:05"), colorReset)
		var calledLevelStringColor = fmt.Sprintf("%v%v%v", errorColor, calledLevelString, colorReset)
		var atColor = fmt.Sprintf("%v%v%v", colorBlue, at, colorReset)
		var messageColor = fmt.Sprintf("%v%v%v", colorCyan, message, colorReset)
		var formatted = fmt.Sprintf("%v > %v > %v > %v", curTime, calledLevelStringColor, atColor, messageColor)
		fmt.Println(formatted)
	}
}

func logFileWriter(l *Logger, calledLevelString string, at string, message string) {
	logFileWriterOptimizer(l, time.Now())
	// make json
	var jsonStruct = logFile{
		Level:   calledLevelString,
		Time:    time.Now().Unix(),
		At:      at,
		Message: message}
	var _json, err = json.Marshal(jsonStruct)
	if err != nil {
		log.Println(err)
		return
	}
	// write json to log file
	_, err = l.fileWriterInfo.file.WriteString(fmt.Sprintln(string(_json)))
	if err != nil {
		log.Println(err)
		return
	}
}

func logFileWriterOptimizer(l *Logger, currentDate time.Time) (bool, error) {
	// create new log file if new day
	var _, _, fileDay = l.fileWriterInfo.fileDate.Date()
	var _, _, currentDay = currentDate.Date()
	if fileDay != currentDay{
		err := newFileWriter(l, currentDate)
		if err != nil {
			return false, err
		}
	}
	// create new log if current log size very big
	logFileStat, err := l.fileWriterInfo.file.Stat()
	if err != nil {
		log.Println(err)
		return false, err
	}
	// log size > * bytes
	if logFileStat.Size() > l.Config.maxLogSize {
		l.fileWriterInfo.file, err = os.Create(l.fileWriterInfo.fullPath)
		if err != nil {
			log.Println(err)
			return false, err
		}
	}
	return true, nil
}

func newFileWriter(l *Logger, date time.Time) error {
	if !l.Config.WriteToFile.Activated {
		return nil
	}
	if l.fileWriterInfo.file != nil {
		err := l.fileWriterInfo.file.Close()
		if err != nil {
		}
	}
	var dirPath = l.Config.WriteToFile.Dir
	l.fileWriterInfo.fileDate = date
	var currentTime = date.Format("02.01.2006")
	currentTime = strings.Replace(currentTime, ".", "_", -1)
	var fileName = fmt.Sprintf("log_%v.txt", currentTime)
	var filePath = fmt.Sprintf("%v/%v", dirPath, fileName)
	l.fileWriterInfo.fullPath = filePath
	// create dir if no exists
	var err error
	if _, err = os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0700)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	// create log file
	logFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	l.fileWriterInfo.file = logFile
	return nil
}
