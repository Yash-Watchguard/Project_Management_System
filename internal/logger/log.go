package logger

import (
	"log"
	"os"
)

var logFile *os.File // global file handle

// InitLogger initializes logging to the file specified in the environment variable LOGFILE
func InitLogger() {
	// Get log file path from environment variable
	logPath := os.Getenv("LOGFILE")
	if logPath == "" {
		log.Fatalf("LOGFILE environment variable not set")
	}

	// Open (or create) the log file
	var err error
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Redirect log package output to the file
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // optional: adds date, time, file:line
}

// CloseLogger closes the log file when program ends
func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

// Info logs an informational message
func Info(msg string) {
	log.Println("[INFO]", msg)
}

// Error logs an error message
func Error(msg string) {
	log.Println("[ERROR]", msg)
}
