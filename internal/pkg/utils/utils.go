// Package utils contains helpful tools that assist in program operation.
package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// InitFileLogging Initializes logging to a file and returns the function that closes that file
func InitFileLogging() func() {
	var (
		exe  string
		err  error
		file *os.File
	)

	if exe, err = os.Executable(); err != nil {
		log.Fatal(err)
	}

	exeDir := filepath.Dir(exe)
	logFile := fmt.Sprintf("%v/.log", exeDir)
	if file, err = os.OpenFile(logFile, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0644); err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	return func() {
		if err := file.Close(); err != nil {
			log.SetOutput(os.Stdout)
			log.Fatal(err)
		}
	}
}
