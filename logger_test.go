package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/sk-pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNew(t *testing.T) {
	// Test creating a new logger manager with default options
	mgr, err := logger.New()
	assert.NoError(t, err)
	assert.NotNil(t, mgr)

	// Test creating a new logger manager with custom options
	mgr, err = logger.New(
		logger.WithDriver("file"),
		logger.WithLogPath("/tmp/test-log-"),
		logger.WithLevel("debug"),
	)
	assert.NoError(t, err)
	assert.NotNil(t, mgr)
}

func TestLoggingWithTraceID(t *testing.T) {
	var buf bytes.Buffer
	writer := zapcore.AddSync(&buf)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writer,
		zapcore.DebugLevel,
	)

	zapLogger := zap.New(core)
	mgr := &logger.Manager{Zap: zapLogger}

	ctx := context.WithValue(context.Background(), logger.TraceIDKey, "test-trace-id")

	// Test Info level logging
	mgr.Info(ctx, "This is an info message")

	// Flush and sync the buffer
	mgr.Sync()

	var loggedMessage map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &loggedMessage)
	assert.NoError(t, err)

	assert.Equal(t, "INFO", loggedMessage["L"])
	assert.Equal(t, "This is an info message", loggedMessage["M"])
	assert.Equal(t, "test-trace-id", loggedMessage["TraceID"])
}

func TestLoggingLevels(t *testing.T) {
	var buf bytes.Buffer
	writer := zapcore.AddSync(&buf)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		writer,
		zapcore.DebugLevel,
	)

	zapLogger := zap.New(core)
	mgr := &logger.Manager{Zap: zapLogger}

	ctx := context.WithValue(context.Background(), logger.TraceIDKey, "test-trace-id")

	// Test all level logging
	mgr.Info(ctx, "This is an info message")
	mgr.Error(ctx, "This is an error message")
	mgr.Debug(ctx, "This is a debug message")
	mgr.Warn(ctx, "This is a warn message")

	// Capture and assert the logs
	mgr.Sync()

	logs := buf.String()
	assert.Contains(t, logs, "INFO")
	assert.Contains(t, logs, "ERROR")
	assert.Contains(t, logs, "DEBUG")
	assert.Contains(t, logs, "WARN")
}
