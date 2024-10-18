package logger

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			name: "Default options",
			opts: []Option{},
		},
		{
			name: "Custom options",
			opts: []Option{
				WithDriver("stdout"),
				WithLevel("info"),
				WithColor(true),
			},
		},
		{
			name: "Invalid driver",
			opts: []Option{
				WithDriver("invalid"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.opts...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func TestManager_LogLevels(t *testing.T) {
	core, recorded := observer.New(zapcore.WarnLevel)
	logger := &Manager{
		Zap: zap.New(core),
	}

	ctx := context.Background()

	tests := []struct {
		name     string
		logFunc  func(context.Context, string, ...zap.Field)
		message  string
		wantLogs int
	}{
		{"Info", logger.Info, "info message", 0},
		{"Debug", logger.Debug, "debug message", 0},
		{"Warn", logger.Warn, "warn message", 1},
		{"Error", logger.Error, "error message", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logFunc(ctx, tt.message)
			assert.Equal(t, tt.wantLogs, recorded.Len())
			if tt.wantLogs > 0 {
				assert.Equal(t, tt.message, recorded.All()[recorded.Len()-1].Message)
			}
		})
	}
}

func TestManager_WithTraceID(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	logger := &Manager{
		Zap: zap.New(core),
	}

	ctx := context.WithValue(context.Background(), TraceIDKey, "test-trace-id")
	logger.Info(ctx, "message with trace id")

	assert.Equal(t, 1, recorded.Len())
	assert.Equal(t, "test-trace-id", recorded.All()[0].ContextMap()["TraceID"])
}

func TestManager_SetLevel(t *testing.T) {
	logger, err := New()
	assert.NoError(t, err)

	logger.SetLevel(zapcore.DebugLevel)
	assert.Equal(t, zapcore.DebugLevel, logger.level.Level())

	logger.SetLevel(zapcore.ErrorLevel)
	assert.Equal(t, zapcore.ErrorLevel, logger.level.Level())
}

func TestManager_Named(t *testing.T) {
	logger, err := New()
	assert.NoError(t, err)

	named := logger.Named("test")
	assert.NotNil(t, named)
}

func TestManager_With(t *testing.T) {
	logger, err := New()
	assert.NoError(t, err)

	with := logger.With(zap.String("key", "value"))
	assert.NotNil(t, with)
}
