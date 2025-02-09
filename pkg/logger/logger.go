package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger handles logging for informational and error messages
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

// NewLogger initializes a new Logger with output to stdout and stderr
func NewLogger() *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "", 0),
		errorLogger: log.New(os.Stderr, "", 0),
	}
}

// formatLog formats the log message with timestamp and level
func formatLog(level, message string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05,000")
	return fmt.Sprintf("[%s] %s %s", timestamp, level, message)
}

// Info logs informational messages with formatting support
func (l *Logger) Info(message string, args ...interface{}) {
	if len(args) > 0 {
		l.infoLogger.Println(formatLog("INFO", fmt.Sprintf(message, args...)))
	} else {
		l.infoLogger.Println(formatLog("INFO", message))
	}
}

// Error logs error messages with optional formatting
func (l *Logger) Error(err error, message string, args ...interface{}) {
	if message != "" {
		if len(args) > 0 {
			l.errorLogger.Println(formatLog("ERROR", fmt.Sprintf(message+" - %v", append(args, err)...)))
		} else {
			l.errorLogger.Println(formatLog("ERROR", fmt.Sprintf("%s - %v", message, err)))
		}
	} else {
		l.errorLogger.Println(formatLog("ERROR", err.Error()))
	}
}

// Fatal logs fatal errors and exits the application
func (l *Logger) Fatal(err error) {
	l.errorLogger.Fatalf("[%s] FATAL FATAL: %v", time.Now().Format("2006-01-02 15:04:05,000"), err)
}
