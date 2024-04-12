package log

import (
	"log"
	"os"
)

const (
	loggerFileName = "golim.log"
)

type Logger struct {
	errLog *log.Logger
}

func InitLogger() *Logger {
	file, err := os.OpenFile(loggerFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return &Logger{
		errLog: log.New(file, "Error Logger:\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
