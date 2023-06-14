package main

import (
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

	loggers.Debug("Debug")
	loggers.Info("Info", zap.String("info", "value"))
	loggers.Error("error", zap.Error(errors.New("debug info")))
}
