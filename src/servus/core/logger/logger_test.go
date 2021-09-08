package logger

import (
	"fmt"
	"os"
	"testing"
)


func TestLogger_Log(t *testing.T) {
	var logger = New(Config{LogLevel: InfoLevel, WriteToConsole: true})
	logger.WriteToFile.Activated = true
	logger.WriteToFile.Dir = GetExecuteDir() + "/testLogs/"
	for i := 0; i < 10; i++ {
		logger.Error(fmt.Sprintf("testing: %v", i))
	}
	//err = os.RemoveAll(logger.WriteToFile.Dir)
	//if err != nil {
	//	println(err.Error())
	//}
}
func GetExecuteDir() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path
}