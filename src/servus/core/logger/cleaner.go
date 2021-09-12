package logger

import (
	"fmt"
	"os"
	"strings"
)

func cleanerMaxLogFiles(l *Logger, dirPath string) {
	// delete the oldest file if max log files reached
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return
	}
	var maxLogFiles = l.Config.WriteToFile.MaxLogFiles
	var isMaxLogFilesNotReached = !(len(files) > maxLogFiles)
	if isMaxLogFilesNotReached {
		return
	}
	// get count of file would be deleted
	var countDifference = len(files) - maxLogFiles
	var errCounter = 0
	for _, file := range files {
		// must oldest files in this cycle be first, because log filename unix timestamp. If im wrong, plz tell me.
		if countDifference <= 0 || errCounter > 3 {
			break
		}
		if file.IsDir() {
			continue
		}
		// if not log file
		if !strings.Contains(file.Name(), "at_") {
			continue
		}
		// if its file and its log file
		var firstLogPath = fmt.Sprintf("%v/%v", dirPath, file.Name())
		err := os.Remove(firstLogPath)
		if err != nil {
			errCounter++
			continue
		}
		countDifference--
	}
}
