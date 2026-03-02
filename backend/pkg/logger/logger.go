package logger

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

func Init(level string) {
	infoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(format string, v ...interface{}) {
	if infoLogger != nil {
		infoLogger.Printf(format, v...)
	} else {
		log.Printf("[INFO] "+format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if errorLogger != nil {
		errorLogger.Printf(format, v...)
	} else {
		log.Printf("[ERROR] "+format, v...)
	}
}
