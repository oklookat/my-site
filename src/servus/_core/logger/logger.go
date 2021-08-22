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

// levels
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

// main funcs
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
	}
}

type fileWriterInfo struct {
	fullPath string
	file     *os.File
}

// log to file JSON
type logFile struct {
	Level   string `json:"level"`
	Time    int64  `json:"time"`
	At      string `json:"at"`
	Message string `json:"message"`
}

// New create new logger
func New(config Config) Logger {
	var fileWriterInfo = fileWriterInfo{}
	var err = initFileWriter(&config, &fileWriterInfo)
	if err != nil {
		panic(err)
	}
	var logger = Logger{Config: config, fileWriterInfo: fileWriterInfo}
	return logger
}

type Logger struct {
	loggerI
	Config
	fileWriterInfo fileWriterInfo
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

func initFileWriter(config *Config, info *fileWriterInfo) error {
	if config.WriteToFile.Activated {
		var dirPath = config.WriteToFile.Dir
		var currentTime = time.Now().Format("02.01.2006")
		currentTime = strings.Replace(currentTime, ".", "_", -1)
		var fileName = fmt.Sprintf("log_%v.txt", currentTime)
		var filePath = fmt.Sprintf("%v/%v", dirPath, fileName)
		info.fullPath = filePath
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
		info.file = logFile
		return nil
	}
	return nil
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
		println("да")
		return
	}
	if l.Config.WriteToFile.Activated {
		logFileWriter(l, calledLevelString, at, message)
	}
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
	logFileWriterOptimizer(l)
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

func logFileWriterOptimizer(l *Logger) {
	logFileStat, err := l.fileWriterInfo.file.Stat()
	if err != nil {
		log.Println(err)
		return
	}
	if logFileStat.Size() > 10000000 { // log size > 10 MB
		l.fileWriterInfo.file, err = os.Create(l.fileWriterInfo.fullPath)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
