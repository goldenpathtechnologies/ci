package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// InitFileLogging Initializes logging to a file and returns the function that closes that file
func InitFileLogging() func() {
	exe, err := os.Executable()
	HandleError(err, false)

	exeDir := filepath.Dir(exe)
	logFile := fmt.Sprintf("%v/.log", exeDir)
	file, err := os.OpenFile(logFile, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0644)
	HandleError(err, false)

	log.SetOutput(file)

	return func() {
		if err := file.Close(); err != nil {
			log.SetOutput(os.Stdout)
			HandleError(err, true)
		}
	}
}
