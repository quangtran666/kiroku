package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
)

var (
	// Logger is the global logger instance
	Logger zerolog.Logger

	// logFile holds the log file handle for cleanup
	logFile *os.File
)

// Config holds logging configuration
type Config struct {
	// Level sets the minimum log level (debug, info, warn, error)
	Level string
	// LogToFile enables file logging
	LogToFile bool
	// LogDir is the directory for log files
	LogDir string
	// MaxLogFiles is the number of log files to keep
	MaxLogFiles int
}

// DefaultConfig returns the default logging configuration
func DefaultConfig() Config {
	homeDir, _ := os.UserHomeDir()
	return Config{
		Level:       "info",
		LogToFile:   true,
		LogDir:      filepath.Join(homeDir, ".local", "share", "kiroku", "logs"),
		MaxLogFiles: 5,
	}
}

// Init initializes the global logger
func Init(cfg Config) error {
	// Parse log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure zerolog
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.CallerMarshalFunc = shortCallerMarshalFunc

	var writers []io.Writer

	// File logging
	if cfg.LogToFile {
		if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}

		// Clean up old log files
		cleanOldLogs(cfg.LogDir, cfg.MaxLogFiles)

		// Create new log file with timestamp
		logFileName := fmt.Sprintf("kiroku_%s.log", time.Now().Format("2006-01-02_15-04-05"))
		logPath := filepath.Join(cfg.LogDir, logFileName)

		logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}

		writers = append(writers, logFile)
	}

	// Create multi-writer if we have writers
	var output io.Writer
	if len(writers) > 0 {
		output = zerolog.MultiLevelWriter(writers...)
	} else {
		// Fallback to discard if no writers
		output = io.Discard
	}

	// Create logger with caller info
	Logger = zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()

	Logger.Info().Msg("Logger initialized")

	return nil
}

// Close closes the log file
func Close() {
	if logFile != nil {
		Logger.Info().Msg("Shutting down logger")
		logFile.Close()
	}
}

// RecoverPanic logs panic information and returns the error
func RecoverPanic() {
	if r := recover(); r != nil {
		stack := debug.Stack()
		Logger.Error().
			Interface("panic", r).
			Bytes("stack", stack).
			Msg("Application panic recovered")

		// Also write to stderr for immediate visibility
		fmt.Fprintf(os.Stderr, "\n🔥 Panic occurred: %v\n\nStack trace written to log file.\n", r)
	}
}

// GetLogDir returns the log directory path
func GetLogDir() string {
	cfg := DefaultConfig()
	return cfg.LogDir
}

// shortCallerMarshalFunc formats caller as file:line
func shortCallerMarshalFunc(pc uintptr, file string, line int) string {
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// cleanOldLogs removes old log files, keeping only the most recent ones
func cleanOldLogs(logDir string, maxFiles int) {
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return
	}

	var logFiles []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".log" {
			logFiles = append(logFiles, entry)
		}
	}

	// Remove oldest files if we have too many
	if len(logFiles) >= maxFiles {
		// Sort by name (which includes timestamp, so older files come first)
		for i := 0; i < len(logFiles)-maxFiles+1; i++ {
			oldFile := filepath.Join(logDir, logFiles[i].Name())
			os.Remove(oldFile)
		}
	}
}

// Debug logs a debug message
func Debug() *zerolog.Event {
	return Logger.Debug()
}

// Info logs an info message
func Info() *zerolog.Event {
	return Logger.Info()
}

// Warn logs a warning message
func Warn() *zerolog.Event {
	return Logger.Warn()
}

// Error logs an error message
func Error() *zerolog.Event {
	return Logger.Error()
}

// Fatal logs a fatal message and exits
func Fatal() *zerolog.Event {
	return Logger.Fatal()
}
