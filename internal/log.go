package internal

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "sub-translate-helper: ", log.Ldate|log.Ltime|log.Lshortfile)

func GetLogger() *log.Logger {
	return logger
}
