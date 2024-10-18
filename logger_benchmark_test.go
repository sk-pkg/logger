package logger

import (
	"context"
	"go.uber.org/zap/zapcore"
	"testing"

	"go.uber.org/zap"
)

func BenchmarkManager_Info(b *testing.B) {
	logger, _ := New()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "benchmark info message")
	}
}

func BenchmarkManager_InfoWithFields(b *testing.B) {
	logger, _ := New()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "benchmark info message", zap.Int("count", i), zap.String("benchmark", "true"))
	}
}

func BenchmarkManager_Error(b *testing.B) {
	logger, _ := New()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Error(ctx, "benchmark error message")
	}
}

func BenchmarkManager_WithTraceID(b *testing.B) {
	logger, _ := New()
	ctx := context.WithValue(context.Background(), TraceIDKey, "benchmark-trace-id")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "benchmark message with trace id")
	}
}

func BenchmarkManager_SetLevel(b *testing.B) {
	logger, _ := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.SetLevel(zapcore.InfoLevel)
	}
}
