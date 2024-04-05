package main

import (
	"log"
	"os"
)

type logger struct {
	errLog *log.Logger
}

func initLogger() *logger {
	file, err := os.OpenFile(loggerFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return &logger{
		errLog: log.New(file, "Error Logger:\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
