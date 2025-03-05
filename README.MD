# Logger

## Introduction

logger is a comprehensive logging package built on top of Zap, offering enhanced functionality and ease of use for Go applications. It supports multiple log levels, file rotation, and context-aware logging with trace IDs.

## Features

- Supports both standard output and file-based logging
- Includes TraceID support for improved log tracing
- Configurable log levels, drivers, and encoder settings
- Easy integration with Go's context for passing TraceID
- Log rotation and retention policies
- Colored output option for console logging

## Installation

```bash
go get -u "github.com/sk-pkg/logger"
```

## Quick Start

Here's a basic example:

```go
package main

import (
    "context"
    "github.com/sk-pkg/logger"
    "log"
)

func main() {
    loggerManager, err := logger.New(logger.WithLevel(logger.DebugLevel))
    if err != nil {
        log.Fatal(err)
    }
    defer loggerManager.Sync()

    ctx := context.WithValue(context.Background(), logger.TraceIDKey, "123456")

    loggerManager.Info(ctx, "This is an info message")
    loggerManager.Debug(ctx, "This is a debug message")
}
```

## Configuration Options

### Log Driver

```go
loggerManager, err := logger.New(
    logger.WithDriver("file"), // or "stdout"
)
```

### Log Level

```go
loggerManager, err := logger.New(
    logger.WithLevel(logger.DebugLevel), // or InfoLevel, WarnLevel, ErrorLevel, DPanicLevel, PanicLevel, FatalLevel
)
```

Or using string:

```go
loggerManager, err := logger.New(
    logger.WithLevel("debug"), // or "info", "warn", "error", "dpanic", "panic", "fatal"
)
```

### Log File Path (only when driver is "file")

```go
loggerManager, err := logger.New(
    logger.WithDriver("file"),
    logger.WithLogPath("/data/logs/"),
)
```

### Custom Encoder Configuration

```go
customEncoderConfig := zapcore.EncoderConfig{
    TimeKey:        "timestamp",
    LevelKey:       "level",
    NameKey:        "logger",
    CallerKey:      "caller",
    MessageKey:     "message",
    StacktraceKey:  "stacktrace",
    LineEnding:     zapcore.DefaultLineEnding,
    EncodeLevel:    zapcore.CapitalLevelEncoder,
    EncodeTime:     zapcore.ISO8601TimeEncoder,
    EncodeDuration: zapcore.SecondsDurationEncoder,
    EncodeCaller:   zapcore.ShortCallerEncoder,
}

loggerManager, err := logger.New(
    logger.WithEncoderConfig(customEncoderConfig),
)
```

### Maximum Age for Log Files

```go
loggerManager, err := logger.New(
    logger.WithDriver("file"),
    logger.WithMaxAge(30 * 24 * time.Hour) // Keep logs for 30 days
)
```

### Log Rotation Time

```go
loggerManager, err := logger.New(
    logger.WithDriver("file"),
    logger.WithRotationTime(24 * time.Hour) // Rotate logs daily
)
```

### Colored Output (for console only)

```go
loggerManager, err := logger.New(
    logger.WithColor(true),
)
```

### Stacktrace Level

```go
loggerManager, err := logger.New(
    logger.WithStacktraceLevel("error"), // or "debug", "info", "warn", "dpanic", "panic", "fatal"
)
```

### Caller Skip Setting

```go
loggerManager, err := logger.New(
    logger.WithCallerSkip(1), // default is 1
)
```

## Logging Methods

The `LoggerManager` provides the following logging methods:

- `Debug(ctx context.Context, msg string, fields ...zap.Field)`
- `Info(ctx context.Context, msg string, fields ...zap.Field)`
- `Warn(ctx context.Context, msg string, fields ...zap.Field)`
- `Error(ctx context.Context, msg string, fields ...zap.Field)`
- `Fatal(ctx context.Context, msg string, fields ...zap.Field)`
- `Panic(ctx context.Context, msg string, fields ...zap.Field)`

Each method accepts a context (for TraceID), a message string, and optional zap.Field values for additional structured logging.

## Additional Methods

### Dynamic Level Change

```go
loggerManager.SetLevel(zapcore.DebugLevel)
```

### Set Caller Skip

```go
loggerManager.SetCallerSkip(2)
```

### Named Logger

```go
namedLogger := loggerManager.Named(ctx, "subsystem")
namedLogger.Info("This is a message from a named logger")
```

### Add Structured Context

```go
withLogger := loggerManager.With(ctx, zap.String("user", "john"))
withLogger.Info("This is a message with additional fields")
```

### Sync Buffered Log Entries

```go
err := loggerManager.Sync()
```

## Complete Example

Here's a comprehensive example showcasing all features:

```go
package main

import (
    "context"
    "errors"
    "github.com/sk-pkg/logger"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "log"
    "time"
)

func main() {
    customEncoderConfig := zapcore.EncoderConfig{
        TimeKey:        "timestamp",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "caller",
        MessageKey:     "message",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.CapitalLevelEncoder,
        EncodeTime:     zapcore.ISO8601TimeEncoder,
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
    }

    loggerManager, err := logger.New(
        logger.WithDriver("file"),
        logger.WithLogPath("/data/logs/"),
        logger.WithLevel(logger.DebugLevel),
        logger.WithEncoderConfig(customEncoderConfig),
        logger.WithMaxAge(7 * 24 * time.Hour),    // Keep logs for 7 days
        logger.WithRotationTime(24 * time.Hour),  // Rotate logs daily
        logger.WithColor(true),                   // Enable colored output (for stdout)
    )
    if err != nil {
        log.Fatal(err)
    }
    defer loggerManager.Sync()

    ctx := context.WithValue(context.Background(), logger.TraceIDKey, "example-trace-id")

    loggerManager.Debug(ctx, "This is a debug message", zap.Int("debugCode", 100))
    loggerManager.Info(ctx, "This is an info message", zap.String("user", "John Doe"))
    loggerManager.Warn(ctx, "This is a warning message", zap.Float64("temperature", 38.5))
    loggerManager.Error(ctx, "This is an error message", zap.Error(errors.New("sample error")))

    // Fatal would typically exit the program
    // loggerManager.Fatal(ctx, "This is a fatal message")
}
```

## Best Practices

1. Always provide a context to logging methods, even if it's `context.Background()`.
2. Use structured logging with `zap.Field` for better log parsing and analysis.
3. Set appropriate log levels for different environments (e.g., debug for development, info for production).
4. Use log rotation and retention policies to manage log file growth.
5. Enable colored output for console logging to improve readability during development.
6. Remember to call `Sync()` on the logger manager before your program exits.

## Performance Considerations

The logging package has been benchmarked to ensure high performance. Consider the following when logging large volumes:

1. Set appropriate log levels, avoiding Debug level in production.
2. Use structured fields instead of string concatenation.
3. Consider using `zap.Field` instead of direct string formatting in high-performance scenarios.

## Reference Documentation

For more detailed information about the underlying zap logger, refer to the [Zap Native Documentation](https://pkg.go.dev/go.uber.org/zap#section-documentation-index).