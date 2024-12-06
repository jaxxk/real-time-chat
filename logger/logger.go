package logger

import (
	"log"

	"github.com/natefinch/lumberjack" // For optional log rotation
)

// Loggers for different levels
var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// Initialize the logger
func Init(logFile string) {
	// Use lumberjack for log rotation
	logFileWriter := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // Max megabytes before rotation
		MaxBackups: 3,  // Max old log files to keep
		MaxAge:     28, // Max number of days to retain old log files
		Compress:   true,
	}

	Info = log.New(logFileWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(logFileWriter, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(logFileWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
