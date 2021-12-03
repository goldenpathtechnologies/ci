package utils

import (
	"log"
	"os"
)

func HandleError(err error, logError bool) {
	if err != nil {
		if logError {
			log.Fatal(err)
		}
		if ScreenBufferActive {
			ExitScreenBuffer()
		}
		os.Exit(1)
	}
}

func PrintAndExit(data string) {
	_, err := os.Stdout.WriteString(data)
	HandleError(err, true)
	os.Exit(0)
}