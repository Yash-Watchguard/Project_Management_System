package logger

import (
	"os"
	"strings"
	"testing"

	
)

func TestLogger(t *testing.T) {
	// Create temporary log file
	tmpFile, err := os.CreateTemp("", "testlog*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // clean up

	// Set environment variable for logger
	os.Setenv("LOGFILE", tmpFile.Name())

	// Initialize logger
	InitLogger()
	defer CloseLogger()

	// Log messages
	Info("info message")
	Error("error message")

	// Flush writes (CloseLogger closes the file)
	CloseLogger()

	// Read the file
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	logText := string(content)

	if !strings.Contains(logText, "[INFO] info message") {
		t.Errorf("expected info message in log, got: %s", logText)
	}

	if !strings.Contains(logText, "[ERROR] error message") {
		t.Errorf("expected error message in log, got: %s", logText)
	}
}
