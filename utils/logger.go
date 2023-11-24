package utils

import (
	"log"
	"os"
)

func SetupLogging() *os.File {
	logFile, err := os.OpenFile("mylog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Fail to open log file: ", err)
	}

	log.SetOutput(logFile)
	return logFile
}
