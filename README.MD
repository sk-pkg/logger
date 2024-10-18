# Logger

Logger is a comprehensive logging package built on top of the zap logger, offering enhanced functionality and ease of use for Go applications.

## Features
- Supports both standard output and file-based logging
- Includes TraceID support for improved log tracing
- Configurable log levels, drivers, and encoder settings
- Easy integration with Go's context for passing TraceID
- Log rotation and retention policies
- Colored output option for console logging

## Installation
```shell
go get -u "github.com/sk-pkg/logger"
```

## Quick Start
Here's a basic example to get you started:

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

The `New` function accepts various options to customize the logger:

### Log Driver
```go
loggerManager, err := logger.New(logger.WithDriver("file"))
```
Options: `"stdout"` (default) or `"file"`

### Log Path
When using the `file` driver:
```go
loggerManager, err := logger.New(
    logger.WithDriver("file"),
    logger.WithLogPath("/data/logs/")
)
```

### Log Level
```go
loggerManager, err := logger.New(logger.WithLevel(logger.DebugLevel))
```
Options: `DebugLevel`, `InfoLevel`, `WarnLevel`, `ErrorLevel`, `FatalLevel`

### Zap Encoder Configuration
```go
customEncoderConfig := zapcore.EncoderConfig{
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

loggerManager, err := logger.New(logger.WithEncoderConfig(customEncoderConfig))
```

### Maximum Age for Log Files
```go
loggerManager, err := logger.New(
    logger.WithDriver("file"),
    logger.WithMaxAge(30 * 24 * time.Hour) // Keep logs for 30 days
)
```

### Log File Rotation
```go
loggerManager, err := logger.New(
    logger.WithDriver("file"),
    logger.WithRotationTime(24 * time.Hour) // Rotate logs daily
)
```

### Colored Output
```go
loggerManager, err := logger.New(
    logger.WithDriver("stdout"),
    logger.WithColor(true) // Enable colored output
)
```

## Logger Methods

The `LoggerManager` provides the following logging methods:

- `Debug(ctx context.Context, msg string, fields ...zap.Field)`
- `Info(ctx context.Context, msg string, fields ...zap.Field)`
- `Warn(ctx context.Context, msg string, fields ...zap.Field)`
- `Error(ctx context.Context, msg string, fields ...zap.Field)`
- `Fatal(ctx context.Context, msg string, fields ...zap.Field)`

Each method accepts a context (for TraceID), a message string, and optional zap.Field values for additional structured logging.

## TraceID Integration

Include TraceID in log entries by setting it in the context:

```go
ctx := context.WithValue(context.Background(), logger.TraceIDKey, "unique-trace-id")
loggerManager.Info(ctx, "Log message with TraceID")
```

## Full Example

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

## Reference Documentation
For more detailed information about the underlying zap logger, refer to the [Zap Native Documentation](https://pkg.go.dev/go.uber.org/zap#section-documentation-index).
