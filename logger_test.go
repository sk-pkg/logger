package logger

import (
	"errors"
	"go.uber.org/zap"
	"testing"
)

func TestLogger(t *testing.T) {
	logger, err := New()

	if err != nil {
		t.Fatal(err)
	}

	logger.Sync()

	err = errors.New("This is an error testing logger ")

	logger = logger.Named("test")

	logger.Info("Info")
	logger.Debug("Debug")
	logger.Error("Error occurs", zap.Error(err))
}
