package logger

import (
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
	Debug *log.Logger
)

func Init() {
	Info = log.New(os.Stdout, "[INFO] ", log.LstdFlags)
	Error = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)
	Debug = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags)
}
