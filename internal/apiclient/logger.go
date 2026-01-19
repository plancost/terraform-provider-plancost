package apiclient

import (
	"github.com/rs/zerolog"
)

// LeveledLogger is a logger that implements the retryablehttp.LeveledLogger interface
type LeveledLogger struct {
	Logger zerolog.Logger
}

// Error logs an error message
func (l *LeveledLogger) Error(msg string, keysAndValues ...interface{}) {
	l.Logger.Error().Fields(keysAndValues).Msg(msg)
}

// Info logs an info message
func (l *LeveledLogger) Info(msg string, keysAndValues ...interface{}) {
	l.Logger.Info().Fields(keysAndValues).Msg(msg)
}

// Debug logs a debug message
func (l *LeveledLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.Logger.Debug().Fields(keysAndValues).Msg(msg)
}

// Warn logs a warning message
func (l *LeveledLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.Logger.Warn().Fields(keysAndValues).Msg(msg)
}
