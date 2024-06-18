package main

import (
	"context"
	"errors"
	"github.com/sk-pkg/logger"
	"go.uber.org/zap"
	"log"
)

func main() {
	loggers, err := logger.New(logger.WithLevel("debug"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.WithValue(context.Background(), logger.TraceIDKey, "123456")

	loggers.Debug(ctx, "Debug")
	loggers.Info(ctx, "Info", zap.String("info", "value"))
	loggers.Error(ctx, "error", zap.Error(errors.New("debug info")))
}
