package log

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// NewLogger return a zap.Logger instance.
func GetLogger() *zap.Logger {
	if logger != nil {
		return logger
	}

	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.DebugLevel)

	config := zap.Config{
		Level:            level,
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		// Fatal error
		fmt.Fprintf(os.Stdout, "Failed to create logger.")
		os.Exit(1)
	}

	return logger
}
