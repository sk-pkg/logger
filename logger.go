// Package logger provides a flexible and extensible logging system built on top of Zap.
// It supports various log levels, file rotation, and context-aware logging with trace IDs.
package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log levels
const (
	DebugLevel zapcore.Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel
)

// Constants for default configuration
const (
	defaultDriver          = "stdout"
	defaultLevel           = InfoLevel
	defaultCallerSkip      = 1
	defaultStacktraceLevel = DPanicLevel
	TraceIDKey             = "trace_id"
)

type (
	// Option is a function that configures the logger options
	Option func(*option)

	// option holds the configuration for the logger
	option struct {
		driver          string                // Log driver: "stdout" or "file"
		level           zapcore.Level         // Minimum log level
		logPath         string                // Path for log files (only used when driver is "file")
		encoderConfig   zapcore.EncoderConfig // Encoder configuration for log formatting
		callerSkip      int                   // Number of stack frames to skip when logging caller info
		maxAge          time.Duration         // Maximum age of log files before rotation
		rotationTime    time.Duration         // Time between log file rotations
		useColor        bool                  // Whether to use colored output (only for console encoder)
		stacktraceLevel zapcore.Level         // Minimum log level for stacktrace
	}

	// Manager manages the logger instance and provides logging methods
	Manager struct {
		Zap   *zap.Logger     // Underlying Zap logger instance
		level zap.AtomicLevel // Atomic level for dynamic level changes
	}
)

// DefaultEncoderConfig is the default encoder configuration for log formatting
var DefaultEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "T",
	LevelKey:       "L",
	NameKey:        "N",
	MessageKey:     "M",
	CallerKey:      "C",
	StacktraceKey:  "S",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

// WithDriver sets the logger driver
//
// Parameters:
//   - driver: The driver to use ("stdout" or "file")
//
// Returns:
//   - Option: A function that sets the driver in the option struct
func WithDriver(driver string) Option {
	return func(o *option) {
		o.driver = driver
	}
}

// WithLevel sets the minimum log level
//
// Parameters:
//   - level: The minimum log level to set
//
// Returns:
//   - Option: A function that sets the level in the option struct
func WithLevel(level string) Option {
	return func(o *option) {
		switch level {
		case "debug":
			o.level = DebugLevel
		case "info":
			o.level = InfoLevel
		case "warn":
			o.level = WarnLevel
		case "error":
			o.level = ErrorLevel
		case "dpanic":
			o.level = DPanicLevel
		case "panic":
			o.level = PanicLevel
		case "fatal":
			o.level = FatalLevel
		default:
			panic("invalid log level")
		}
	}
}

// WithLogPath sets the log file path (only used when driver is "file")
//
// Parameters:
//   - path: The path where log files will be stored
//
// Returns:
//   - Option: A function that sets the log path in the option struct
func WithLogPath(path string) Option {
	return func(o *option) {
		o.logPath = path
	}
}

// WithEncoderConfig sets the encoder configuration for log formatting
//
// Parameters:
//   - config: The encoder configuration to use
//
// Returns:
//   - Option: A function that sets the encoder config in the option struct
func WithEncoderConfig(config zapcore.EncoderConfig) Option {
	return func(o *option) {
		o.encoderConfig = config
	}
}

// WithCallerSkip sets the number of callers to skip when logging caller info
//
// Parameters:
//   - skip: The number of callers to skip
//
// Returns:
//   - Option: A function that sets the caller skip in the option struct
func WithCallerSkip(skip int) Option {
	return func(o *option) {
		o.callerSkip = skip
	}
}

// WithMaxAge sets the maximum age for log files before rotation
//
// Parameters:
//   - maxAge: The maximum age for log files
//
// Returns:
//   - Option: A function that sets the max age in the option struct
func WithMaxAge(maxAge time.Duration) Option {
	return func(o *option) {
		o.maxAge = maxAge
	}
}

// WithRotationTime sets the time between log file rotations
//
// Parameters:
//   - rotationTime: The time between log file rotations
//
// Returns:
//   - Option: A function that sets the rotation time in the option struct
func WithRotationTime(rotationTime time.Duration) Option {
	return func(o *option) {
		o.rotationTime = rotationTime
	}
}

// WithColor enables or disables colored output (only for console encoder)
//
// Parameters:
//   - useColor: Whether to use colored output
//
// Returns:
//   - Option: A function that sets the color usage in the option struct
func WithColor(useColor bool) Option {
	return func(o *option) {
		o.useColor = useColor
	}
}

// WithStacktraceLevel sets the minimum log level for stacktrace
//
// Parameters:
//   - level: The minimum log level for stacktrace
//
// Returns:
//   - Option: A function that sets the stacktrace level in the option struct
func WithStacktraceLevel(level string) Option {
	return func(o *option) {
		switch level {
		case "debug":
			o.stacktraceLevel = DebugLevel
		case "info":
			o.stacktraceLevel = InfoLevel
		case "warn":
			o.stacktraceLevel = WarnLevel
		case "error":
			o.stacktraceLevel = ErrorLevel
		case "dpanic":
			o.stacktraceLevel = DPanicLevel
		case "panic":
			o.stacktraceLevel = PanicLevel
		case "fatal":
			o.stacktraceLevel = FatalLevel
		default:
			panic("invalid log level")
		}
	}
}

// New creates a new logger manager with the given options
//
// Parameters:
//   - opts: A variadic list of Option functions to configure the logger
//
// Returns:
//   - *Manager: A new Manager instance
//   - error: An error if the logger creation fails
//
// Example:
//
//	logger, err := New(
//	    WithDriver("file"),
//	    WithLogPath("/var/log/myapp/"),
//	    WithLevel(InfoLevel),
//	    WithColor(true),
//	)
//	if err != nil {
//	    // Handle error
//	}
func New(opts ...Option) (*Manager, error) {
	// Initialize default options
	opt := &option{
		driver:          defaultDriver,
		level:           defaultLevel,
		encoderConfig:   DefaultEncoderConfig,
		callerSkip:      defaultCallerSkip,
		maxAge:          7 * 24 * time.Hour,
		rotationTime:    24 * time.Hour,
		stacktraceLevel: defaultStacktraceLevel,
	}

	// Apply provided options
	for _, f := range opts {
		f(opt)
	}

	// Create atomic level for dynamic level changes
	level := zap.NewAtomicLevelAt(opt.level)

	// Create encoder based on color option
	var encoder zapcore.Encoder
	if opt.useColor {
		encoder = zapcore.NewConsoleEncoder(opt.encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(opt.encoderConfig)
	}

	var core zapcore.Core
	var err error

	// Create core based on driver
	switch opt.driver {
	case "stdout":
		core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
	case "file":
		core, err = newFileCore(opt, encoder, level)
		if err != nil {
			return nil, fmt.Errorf("failed to create file core: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown driver: %s", opt.driver)
	}

	// Create Zap logger
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(opt.callerSkip),
		zap.ErrorOutput(zapcore.AddSync(os.Stderr)),
		zap.AddStacktrace(opt.stacktraceLevel),
	)

	// Return new Manager instance
	return &Manager{
		Zap:   logger,
		level: level,
	}, nil
}

// newFileCore creates a new zapcore.Core for file-based logging
//
// Parameters:
//   - opt: The option struct containing configuration
//   - encoder: The zapcore.Encoder to use
//   - level: The zap.AtomicLevel for dynamic level changes
//
// Returns:
//   - zapcore.Core: A new Core for file-based logging
//   - error: An error if the file core creation fails
func newFileCore(opt *option, encoder zapcore.Encoder, level zap.AtomicLevel) (zapcore.Core, error) {
	// Create rotatelogs hook
	hook, err := rotatelogs.New(
		opt.logPath+"%Y-%m-%d.log",
		rotatelogs.WithMaxAge(opt.maxAge),
		rotatelogs.WithRotationTime(opt.rotationTime),
	)
	if err != nil {
		return nil, err
	}

	// Create and return new Core
	return zapcore.NewCore(encoder, zapcore.AddSync(hook), level), nil
}

// getTraceIDFromContext extracts the TraceID from the context
//
// Parameters:
//   - ctx: The context.Context to extract the TraceID from
//
// Returns:
//   - string: The extracted TraceID, or an empty string if not found
func getTraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// getLoggerWithTraceID returns a logger with the TraceID field added if present in the context
//
// Parameters:
//   - ctx: The context.Context to extract the TraceID from
//
// Returns:
//   - *zap.Logger: A logger with the TraceID field added if present
func (m *Manager) getLoggerWithTraceID(ctx context.Context) *zap.Logger {
	traceID := getTraceIDFromContext(ctx)
	if traceID == "" {
		return m.Zap
	}

	return m.Zap.With(zap.String("TraceID", traceID))
}

// SetLevel dynamically changes the log level
//
// Parameters:
//   - level: The new zapcore.Level to set
func (m *Manager) SetLevel(level zapcore.Level) {
	m.level.SetLevel(level)
}

// Info logs a message at InfoLevel
//
// Parameters:
//   - ctx: The context.Context for this log entry
//   - msg: The message to log
//   - fields: Optional fields to add to the log entry
func (m *Manager) Info(ctx context.Context, msg string, fields ...zap.Field) {
	logger := m.getLoggerWithTraceID(ctx)
	logger.Info(msg, fields...)
}

// Error logs a message at ErrorLevel with a stack trace
//
// Parameters:
//   - ctx: The context.Context for this log entry
//   - msg: The message to log
//   - fields: Optional fields to add to the log entry
func (m *Manager) Error(ctx context.Context, msg string, fields ...zap.Field) {
	logger := m.getLoggerWithTraceID(ctx)
	logger.Error(msg, fields...)
}

// Debug logs a message at DebugLevel
//
// Parameters:
//   - ctx: The context.Context for this log entry
//   - msg: The message to log
//   - fields: Optional fields to add to the log entry
func (m *Manager) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	logger := m.getLoggerWithTraceID(ctx)
	logger.Debug(msg, fields...)
}

// Warn logs a message at WarnLevel
//
// Parameters:
//   - ctx: The context.Context for this log entry
//   - msg: The message to log
//   - fields: Optional fields to add to the log entry
func (m *Manager) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	logger := m.getLoggerWithTraceID(ctx)
	logger.Warn(msg, fields...)
}

// Fatal logs a message at FatalLevel with a stack trace, then calls os.Exit(1)
//
// Parameters:
//   - ctx: The context.Context for this log entry
//   - msg: The message to log
//   - fields: Optional fields to add to the log entry
func (m *Manager) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	logger := m.getLoggerWithTraceID(ctx)
	logger.Fatal(msg, fields...)
}

// Panic logs a message at PanicLevel with a stack trace, then panics
//
// Parameters:
//   - ctx: The context.Context for this log entry
//   - msg: The message to log
//   - fields: Optional fields to add to the log entry
func (m *Manager) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	logger := m.getLoggerWithTraceID(ctx)
	logger.Panic(msg, fields...)
}

// Sync flushes any buffered log entries
//
// Returns:
//   - error: An error if the sync operation fails
func (m *Manager) Sync() error {
	return m.Zap.Sync()
}

// Named adds a sub-scope to the logger's name
//
// Parameters:
//   - name: The name to add to the logger
//
// Returns:
//   - *zap.Logger: A new logger with the given name added
func (m *Manager) Named(name string) *zap.Logger {
	return m.Zap.Named(name)
}

// With creates a child logger and adds structured context to it
//
// Parameters:
//   - fields: The fields to add to the logger
//
// Returns:
//   - *zap.Logger: A new logger with the given fields added
func (m *Manager) With(fields ...zap.Field) *zap.Logger {
	return m.Zap.With(fields...)
}
