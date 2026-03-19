package logger

import (
	"io"
	"log"
	"os"
)

// Logger is the central logger for the entire application.
var Logger *log.Logger

func init() {
	writer := io.Writer(os.Stdout)
	Logger = log.New(writer, "spamfilter: ", log.Ltime|log.Lshortfile)
}

