package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	logDir      = "logs" // keep all logs in a folder
	logFilePath = filepath.Join(logDir, "app.log")
	maxLines    = 10000
	lineCount   = 0
	file        *os.File
)

func init() {
	setup()
}

// setup initializes the log file and directory
func setup() {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("failed to create log directory: %v", err)
	}

	var err error
	file, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	log.SetOutput(os.Stdout)
}

// write handles both console + file logging
func write(level string, msg string) {
	lineCount++
	if lineCount >= maxLines {
		rotateFile()
	}

	// Add timestamp
	timestamp := time.Now().Format(time.RFC3339)

	// Get caller info
	_, fileName, line, _ := runtime.Caller(2)
	entry := fmt.Sprintf("[%s] %s [%s:%d] %s",
		level,
		timestamp,
		filepath.Base(fileName),
		line,
		msg,
	)

	// Print to console
	log.Println(entry)

	// Write to current log file
	fmt.Fprintln(file, entry)
}

// rotateFile creates a new log file and keeps the old one
func rotateFile() {
	file.Close()

	// rename old file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupName := filepath.Join(logDir, fmt.Sprintf("app_%s.log", timestamp))

	err := os.Rename(logFilePath, backupName)
	if err != nil {
		log.Printf("failed to rotate log file: %v", err)
	}

	// open new fresh log file
	var openErr error
	file, openErr = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if openErr != nil {
		log.Fatalf("failed to open new log file: %v", openErr)
	}

	lineCount = 0
	log.Printf("[logger] log rotated -> %s", backupName)
}

// ----------------- Public Helpers ------------------

// Info logs general information
func Info(msg string) {
	write("INFO", msg)
}

// Warn logs warning
func Warn(msg string) {
	write("WARN", msg)
}

// Error logs an error
func Error(err error, msg string) {
	if err == nil {
		return
	}
	write("ERROR", fmt.Sprintf("%s | %v", msg, err))
}
