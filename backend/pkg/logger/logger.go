package logger

import (
	"log"
	"os"
)

var defaultLogger *log.Logger

// InitLogger initializes the global logger based on the environment.
func InitLogger(env string) {
	// Use standard logger with useful flags. Keep it simple and portable.
	defaultLogger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
}

// Info logs an informational message.
func Info(msg string, args ...interface{}) {
	logWithLevel("INFO", msg, args...)
}

// Warn logs a warning message.
func Warn(msg string, args ...interface{}) {
	logWithLevel("WARN", msg, args...)
}

// Error logs an error message.
func Error(msg string, args ...interface{}) {
	logWithLevel("ERROR", msg, args...)
}

// Fatal logs a fatal error message and exits the application.
func Fatal(msg string, args ...interface{}) {
	logWithLevel("FATAL", msg, args...)
	os.Exit(1)
}

func logWithLevel(level, msg string, args ...interface{}) {
	if defaultLogger == nil {
		// Fallback if logger not initialized
		log.Printf("[%s] %s %v\n", level, msg, args)
		return
	}

	// In a real application, you'd use a structured logging library like zerolog or zap
	// For simplicity, using tint for dev and basic json for prod here.
	if level == "INFO" || level == "WARN" || level == "ERROR" || level == "FATAL" {
		defaultLogger.Printf("{\"level\":\"%s\", \"message\":\"%s\", \"fields\":%v}\n", level, msg, args)
	} else {
		defaultLogger.Printf("{\"level\":\"%s\", \"message\":\"%s\", \"fields\":%v}\n", level, msg, args)
	}
}
