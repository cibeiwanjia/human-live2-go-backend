// Package logger provides structured logging using Zap
package logger

import (
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	sugarLogger *zap.SugaredLogger
	loggerOnce  sync.Once
)

// Config holds logger configuration
type LogConfig struct {
	Level      string // debug, info, warn, error
	Format     string // json, console
	OutputPath string // stdout, stderr, or file path
}

// Init initializes the global logger
func Init(level, format string) error {
	var err error
	loggerOnce.Do(func() {
		err = initLogger(level, format)
	})
	return err
}

// InitWithConfig initializes logger with config struct
func InitWithConfig(cfg LogConfig) error {
	return Init(cfg.Level, cfg.Format)
}

func initLogger(level, format string) error {
	// Parse log level
	var zapLevel zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn", "warning":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create encoder based on format
	var encoder zapcore.Encoder
	if strings.ToLower(format) == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		// Console format with colors
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Create core with stdout writer
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)

	// Create logger with caller and stacktrace
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	sugarLogger = logger.Sugar()
	return nil
}

// L returns the global sugared logger instance
// If not initialized, returns a development logger
func L() *zap.SugaredLogger {
	if sugarLogger == nil {
		// Return a nop logger if not initialized
		return zap.NewNop().Sugar()
	}
	return sugarLogger
}

// Sync flushes any buffered log entries
func Sync() error {
	if sugarLogger != nil {
		return sugarLogger.Sync()
	}
	return nil
}

// Debug logs a message at DebugLevel
func Debug(args ...interface{}) {
	L().Debug(args...)
}

// Debugf logs a formatted message at DebugLevel
func Debugf(template string, args ...interface{}) {
	L().Debugf(template, args...)
}

// Info logs a message at InfoLevel
func Info(args ...interface{}) {
	L().Info(args...)
}

// Infof logs a formatted message at InfoLevel
func Infof(template string, args ...interface{}) {
	L().Infof(template, args...)
}

// Warn logs a message at WarnLevel
func Warn(args ...interface{}) {
	L().Warn(args...)
}

// Warnf logs a formatted message at WarnLevel
func Warnf(template string, args ...interface{}) {
	L().Warnf(template, args...)
}

// Error logs a message at ErrorLevel
func Error(args ...interface{}) {
	L().Error(args...)
}

// Errorf logs a formatted message at ErrorLevel
func Errorf(template string, args ...interface{}) {
	L().Errorf(template, args...)
}

// Fatal logs a message at FatalLevel and calls os.Exit(1)
func Fatal(args ...interface{}) {
	L().Fatal(args...)
}

// Fatalf logs a formatted message at FatalLevel and calls os.Exit(1)
func Fatalf(template string, args ...interface{}) {
	L().Fatalf(template, args...)
}

// With creates a child logger with additional fields
func With(fields ...interface{}) *zap.SugaredLogger {
	return L().With(fields...)
}

// Named adds a sub-logger name
func Named(name string) *zap.SugaredLogger {
	return L().Named(name)
}
